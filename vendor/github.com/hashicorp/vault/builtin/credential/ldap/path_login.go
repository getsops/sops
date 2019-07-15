package ldap

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login/(?P<username>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "DN (distinguished name) to be used for login.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: username,
			},
		},
	}, nil
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("auth method not configured"), nil
	}

	// Check for a CIDR match.
	if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, cfg.TokenBoundCIDRs) {
		return nil, logical.ErrPermissionDenied
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	policies, resp, groupNames, err := b.Login(ctx, req, username, password)
	// Handle an internal error
	if err != nil {
		return nil, err
	}
	if resp != nil {
		// Handle a logical error
		if resp.IsError() {
			return resp, nil
		}
	} else {
		resp = &logical.Response{}
	}

	auth := &logical.Auth{
		Metadata: map[string]string{
			"username": username,
		},
		InternalData: map[string]interface{}{
			"password": password,
		},
		DisplayName: username,
		Alias: &logical.Alias{
			Name: username,
		},
	}

	cfg.PopulateTokenAuth(auth)

	// Add in configured policies from mappings
	if len(policies) > 0 {
		auth.Policies = append(auth.Policies, policies...)
	}

	resp.Auth = auth

	for _, groupName := range groupNames {
		if groupName == "" {
			continue
		}
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: groupName,
		})
	}
	return resp, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("auth method not configured"), nil
	}

	username := req.Auth.Metadata["username"]
	password := req.Auth.InternalData["password"].(string)

	loginPolicies, resp, groupNames, err := b.Login(ctx, req, username, password)
	if len(loginPolicies) == 0 {
		return resp, err
	}
	finalPolicies := cfg.TokenPolicies
	if len(loginPolicies) > 0 {
		finalPolicies = append(finalPolicies, loginPolicies...)
	}

	if !policyutil.EquivalentPolicies(finalPolicies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	resp.Auth = req.Auth
	resp.Auth.Period = cfg.TokenPeriod
	resp.Auth.TTL = cfg.TokenTTL
	resp.Auth.MaxTTL = cfg.TokenMaxTTL

	// Remove old aliases
	resp.Auth.GroupAliases = nil

	for _, groupName := range groupNames {
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: groupName,
		})
	}

	return resp, nil
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password. Please be sure to
read the note on escaping from the path-help for the 'config' endpoint.
`
