package kubeauth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/briankassouf/jose/crypto"
	"github.com/briankassouf/jose/jws"
	"github.com/briankassouf/jose/jwt"
	"github.com/hashicorp/errwrap"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

var (
	// expectedJWTIssuer is used to verify the iss header on the JWT.
	expectedJWTIssuer = "kubernetes/serviceaccount"

	uidJWTClaimKey = "kubernetes.io/serviceaccount/service-account.uid"

	// errMismatchedSigningMethod is used if the certificate doesn't match the
	// JWT's expected signing method.
	errMismatchedSigningMethod = errors.New("invalid signing method")
)

// pathLogin returns the path configurations for login endpoints
func pathLogin(b *kubeAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: `Name of the role against which the login is being attempted. This field is required`,
			},
			"jwt": {
				Type:        framework.TypeString,
				Description: `A signed JWT for authenticating a service account. This field is required.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin(),
			logical.AliasLookaheadOperation: b.aliasLookahead(),
		},

		HelpSynopsis:    pathLoginHelpSyn,
		HelpDescription: pathLoginHelpDesc,
	}
}

// pathLogin is used to authenticate to this backend
func (b *kubeAuthBackend) pathLogin() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		roleName := data.Get("role").(string)
		if len(roleName) == 0 {
			return logical.ErrorResponse("missing role"), nil
		}

		jwtStr := data.Get("jwt").(string)
		if len(jwtStr) == 0 {
			return logical.ErrorResponse("missing jwt"), nil
		}

		b.l.RLock()
		defer b.l.RUnlock()

		role, err := b.role(ctx, req.Storage, roleName)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid role name \"%s\"", roleName)), nil
		}

		// Check for a CIDR match.
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, role.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}

		config, err := b.config(ctx, req.Storage)
		if err != nil {
			return nil, err
		}
		if config == nil {
			return nil, errors.New("could not load backend configuration")
		}

		serviceAccount, err := b.parseAndValidateJWT(jwtStr, role, config)
		if err != nil {
			return nil, err
		}

		// look up the JWT token in the kubernetes API
		err = serviceAccount.lookup(jwtStr, b.reviewFactory(config))
		if err != nil {
			return nil, err
		}

		auth := &logical.Auth{
			Alias: &logical.Alias{
				Name: serviceAccount.uid(),
				Metadata: map[string]string{
					"service_account_uid":         serviceAccount.uid(),
					"service_account_name":        serviceAccount.name(),
					"service_account_namespace":   serviceAccount.namespace(),
					"service_account_secret_name": serviceAccount.SecretName,
				},
			},
			InternalData: map[string]interface{}{
				"role": roleName,
			},
			Metadata: map[string]string{
				"service_account_uid":         serviceAccount.uid(),
				"service_account_name":        serviceAccount.name(),
				"service_account_namespace":   serviceAccount.namespace(),
				"service_account_secret_name": serviceAccount.SecretName,
				"role":                        roleName,
			},
			DisplayName: fmt.Sprintf("%s-%s", serviceAccount.namespace(), serviceAccount.name()),
		}

		role.PopulateTokenAuth(auth)

		return &logical.Response{
			Auth: auth,
		}, nil
	}
}

// aliasLookahead returns the alias object with the SA UID from the JWT
// Claims.
func (b *kubeAuthBackend) aliasLookahead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		jwtStr := data.Get("jwt").(string)
		if len(jwtStr) == 0 {
			return logical.ErrorResponse("missing jwt"), nil
		}

		// Parse into JWT
		parsedJWT, err := jws.ParseJWT([]byte(jwtStr))
		if err != nil {
			return nil, err
		}

		saUID, ok := parsedJWT.Claims().Get(uidJWTClaimKey).(string)
		if !ok || saUID == "" {
			return nil, errors.New("could not parse UID from claims")
		}

		return &logical.Response{
			Auth: &logical.Auth{
				Alias: &logical.Alias{
					Name: saUID,
				},
			},
		}, nil
	}
}

// parseAndValidateJWT is used to parse, validate and lookup the JWT token.
func (b *kubeAuthBackend) parseAndValidateJWT(jwtStr string, role *roleStorageEntry, config *kubeConfig) (*serviceAccount, error) {
	// Parse into JWT
	parsedJWT, err := jws.ParseJWT([]byte(jwtStr))
	if err != nil {
		return nil, err
	}

	sa := &serviceAccount{}
	validator := &jwt.Validator{
		Expected: jwt.Claims{
			"iss": expectedJWTIssuer,
		},
		Fn: func(c jwt.Claims) error {
			// Decode claims into a service account object
			err := mapstructure.Decode(c, sa)
			if err != nil {
				return err
			}

			// verify the namespace is allowed
			if len(role.ServiceAccountNamespaces) > 1 || role.ServiceAccountNamespaces[0] != "*" {
				if !strutil.StrListContainsGlob(role.ServiceAccountNamespaces, sa.namespace()) {
					return errors.New("namespace not authorized")
				}
			}

			// verify the service account name is allowed
			if len(role.ServiceAccountNames) > 1 || role.ServiceAccountNames[0] != "*" {
				if !strutil.StrListContainsGlob(role.ServiceAccountNames, sa.name()) {
					return errors.New("service account name not authorized")
				}
			}

			return nil
		},
	}

	if err := validator.Validate(parsedJWT); err != nil {
		return nil, err
	}

	// If we don't have any public keys to verify, return the sa and end early.
	if len(config.PublicKeys) == 0 {
		return sa, nil
	}

	// verifyFunc is called for each certificate that is configured in the
	// backend until one of the certificates succeeds.
	verifyFunc := func(cert interface{}) error {
		// Parse Headers and verify the signing method matches the public key type
		// configured. This is done in its own scope since we don't need most of
		// these variables later.
		var signingMethod crypto.SigningMethod
		{
			parsedJWS, err := jws.Parse([]byte(jwtStr))
			if err != nil {
				return err
			}
			headers := parsedJWS.Protected()

			var algStr string
			if headers.Has("alg") {
				algStr = headers.Get("alg").(string)
			} else {
				return errors.New("provided JWT must have 'alg' header value")
			}

			signingMethod = jws.GetSigningMethod(algStr)
			switch signingMethod.(type) {
			case *crypto.SigningMethodECDSA:
				if _, ok := cert.(*ecdsa.PublicKey); !ok {
					return errMismatchedSigningMethod
				}
			case *crypto.SigningMethodRSA:
				if _, ok := cert.(*rsa.PublicKey); !ok {
					return errMismatchedSigningMethod
				}
			default:
				return errors.New("unsupported JWT signing method")
			}
		}

		// validates the signature and then runs the claim validation
		if err := parsedJWT.Validate(cert, signingMethod); err != nil {
			return err
		}

		return nil
	}

	var validationErr error
	// for each configured certificate run the verifyFunc
	for _, cert := range config.PublicKeys {
		err := verifyFunc(cert)
		switch err {
		case nil:
			return sa, nil
		case rsa.ErrVerification, crypto.ErrECDSAVerification, errMismatchedSigningMethod:
			// if the error is a failure to verify or a signing method mismatch
			// continue onto the next cert, storing the error to be returned if
			// this is the last cert.
			validationErr = multierror.Append(validationErr, errwrap.Wrapf("failed to validate JWT: {{err}}", err))
			continue
		default:
			return nil, err
		}
	}

	return nil, validationErr
}

// serviceAccount holds the metadata from the JWT token and is used to lookup
// the JWT in the kubernetes API and compare the results.
type serviceAccount struct {
	Name       string   `mapstructure:"kubernetes.io/serviceaccount/service-account.name"`
	UID        string   `mapstructure:"kubernetes.io/serviceaccount/service-account.uid"`
	SecretName string   `mapstructure:"kubernetes.io/serviceaccount/secret.name"`
	Namespace  string   `mapstructure:"kubernetes.io/serviceaccount/namespace"`
	Aud        []string `mapstructure:"aud"`

	// the JSON returned from reviewing a Projected Service account has a
	// different structure, where the information is in a sub-structure instead of
	// at the top level
	Kubernetes *projectedServiceToken `mapstructure:"kubernetes.io"`
	Expiration int64                  `mapstructure:"exp"`
	IssuedAt   int64                  `mapstructure:"iat"`
}

// uid returns the UID for the service account, preferring the projected service
// account value if found
func (s *serviceAccount) uid() string {
	if s.Kubernetes != nil && s.Kubernetes.ServiceAccount != nil {
		return s.Kubernetes.ServiceAccount.UID
	}
	return s.UID
}

// name returns the name for the service account, preferring the projected
// service account value if found. This is "default" for projected service
// accounts
func (s *serviceAccount) name() string {
	if s.Kubernetes != nil && s.Kubernetes.ServiceAccount != nil {
		return s.Kubernetes.ServiceAccount.Name
	}
	return s.Name
}

// namespace returns the namespace for the service account, preferring the
// projected service account value if found
func (s *serviceAccount) namespace() string {
	if s.Kubernetes != nil {
		return s.Kubernetes.Namespace
	}
	return s.Namespace
}

type projectedServiceToken struct {
	Namespace      string                      `mapstructure:"namespace"`
	Pod            *projectedServiceAccountPod `mapstructure:"pod"`
	ServiceAccount *projectedServiceAccountPod `mapstructure:"serviceaccount"`
}

type projectedServiceAccountPod struct {
	Name string `mapstructure:"name"`
	UID  string `mapstructure:"uid"`
}

// lookup calls the TokenReview API in kubernetes to verify the token and secret
// still exist.
func (s *serviceAccount) lookup(jwtStr string, tr tokenReviewer) error {
	r, err := tr.Review(jwtStr)
	if err != nil {
		return err
	}

	// Verify the returned metadata matches the expected data from the service
	// account.
	if s.name() != r.Name {
		return errors.New("JWT names did not match")
	}
	if s.uid() != r.UID {
		return errors.New("JWT UIDs did not match")
	}
	if s.namespace() != r.Namespace {
		return errors.New("JWT namepaces did not match")
	}

	return nil
}

// Invoked when the token issued by this backend is attempting a renewal.
func (b *kubeAuthBackend) pathLoginRenew() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		roleName := req.Auth.InternalData["role"].(string)
		if roleName == "" {
			return nil, fmt.Errorf("failed to fetch role_name during renewal")
		}

		b.l.RLock()
		defer b.l.RUnlock()

		// Ensure that the Role still exists.
		role, err := b.role(ctx, req.Storage, roleName)
		if err != nil {
			return nil, fmt.Errorf("failed to validate role %s during renewal:%s", roleName, err)
		}
		if role == nil {
			return nil, fmt.Errorf("role %s does not exist during renewal", roleName)
		}

		resp := &logical.Response{Auth: req.Auth}
		resp.Auth.TTL = role.TokenTTL
		resp.Auth.MaxTTL = role.TokenMaxTTL
		resp.Auth.Period = role.TokenPeriod
		return resp, nil
	}
}

const pathLoginHelpSyn = `Authenticates Kubernetes service accounts with Vault.`
const pathLoginHelpDesc = `
Authenticate Kubernetes service accounts.
`
