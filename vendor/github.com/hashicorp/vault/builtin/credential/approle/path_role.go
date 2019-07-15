package approle

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// roleStorageEntry stores all the options that are set on an role
type roleStorageEntry struct {
	tokenutil.TokenParams

	// Name of the role. This field is not persisted on disk. After the role is
	// read out of disk, the sanitized version of name is set in this field for
	// subsequent use of role name elsewhere.
	name string

	// UUID that uniquely represents this role. This serves as a credential
	// to perform login using this role.
	RoleID string `json:"role_id" mapstructure:"role_id"`

	// UUID that serves as the HMAC key for the hashing the 'secret_id's
	// of the role
	HMACKey string `json:"hmac_key" mapstructure:"hmac_key"`

	// Policies that are to be required by the token to access this role. Deprecated.
	Policies []string `json:"policies" mapstructure:"policies"`

	// Number of times the SecretID generated against this role can be
	// used to perform login operation
	SecretIDNumUses int `json:"secret_id_num_uses" mapstructure:"secret_id_num_uses"`

	// Duration (less than the backend mount's max TTL) after which a
	// SecretID generated against the role will expire
	SecretIDTTL time.Duration `json:"secret_id_ttl" mapstructure:"secret_id_ttl"`

	// A constraint, if set, requires 'secret_id' credential to be presented during login
	BindSecretID bool `json:"bind_secret_id" mapstructure:"bind_secret_id"`

	// Deprecated: A constraint, if set, specifies the CIDR blocks from which logins should be allowed,
	// please use SecretIDBoundCIDRs instead.
	BoundCIDRListOld string `json:"bound_cidr_list,omitempty"`

	// Deprecated: A constraint, if set, specifies the CIDR blocks from which logins should be allowed,
	// please use SecretIDBoundCIDRs instead.
	BoundCIDRList []string `json:"bound_cidr_list_list" mapstructure:"bound_cidr_list"`

	// A constraint, if set, specifies the CIDR blocks from which logins should be allowed
	SecretIDBoundCIDRs []string `json:"secret_id_bound_cidrs" mapstructure:"secret_id_bound_cidrs"`

	// Period, if set, indicates that the token generated using this role
	// should never expire. The token should be renewed within the duration
	// specified by this value. The renewal duration will be fixed if the value
	// is not modified on the role. If the `Period` in the role is modified, a
	// token will pick up the new value during its next renewal. Deprecated.
	Period time.Duration `json:"period" mapstructure:"period"`

	// LowerCaseRoleName enforces the lower casing of role names for all the
	// roles that get created since this field was introduced.
	LowerCaseRoleName bool `json:"lower_case_role_name" mapstructure:"lower_case_role_name"`

	// SecretIDPrefix is the storage prefix for persisting secret IDs. This
	// differs based on whether the secret IDs are cluster local or not.
	SecretIDPrefix string `json:"secret_id_prefix" mapstructure:"secret_id_prefix"`
}

// roleIDStorageEntry represents the reverse mapping from RoleID to Role
type roleIDStorageEntry struct {
	Name string `json:"name" mapstructure:"name"`
}

// rolePaths creates all the paths that are used to register and manage an role.
//
// Paths returned:
// role/ - For listing all the registered roles
// role/<role_name> - For registering an role
// role/<role_name>/policies - For updating the param
// role/<role_name>/secret-id-num-uses - For updating the param
// role/<role_name>/secret-id-ttl - For updating the param
// role/<role_name>/token-ttl - For updating the param
// role/<role_name>/token-max-ttl - For updating the param
// role/<role_name>/token-num-uses - For updating the param
// role/<role_name>/bind-secret-id - For updating the param
// role/<role_name>/bound-cidr-list - For updating the param
// role/<role_name>/period - For updating the param
// role/<role_name>/role-id - For fetching the role_id of an role
// role/<role_name>/secret-id - For issuing a secret_id against an role, also to list the secret_id_accessors
// role/<role_name>/custom-secret-id - For assigning a custom SecretID against an role
// role/<role_name>/secret-id/lookup - For reading the properties of a secret_id
// role/<role_name>/secret-id/destroy - For deleting a secret_id
// role/<role_name>/secret-id-accessor/lookup - For reading secret_id using accessor
// role/<role_name>/secret-id-accessor/destroy - For deleting secret_id using accessor
func rolePaths(b *backend) []*framework.Path {
	defTokenFields := tokenutil.TokenFields()

	p := &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("role_name"),
		Fields: map[string]*framework.FieldSchema{
			"role_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"bind_secret_id": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     true,
				Description: "Impose secret_id to be presented when logging in using this role. Defaults to 'true'.",
			},

			"bound_cidr_list": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: `Use "secret_id_bound_cidrs" instead.`,
				Deprecated:  true,
			},

			"secret_id_bound_cidrs": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `Comma separated string or list of CIDR blocks. If set, specifies the blocks of
IP addresses which can perform the login operation.`,
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_policies"),
				Deprecated:  true,
			},

			"secret_id_num_uses": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `Number of times a SecretID can access the role, after which the SecretID
will expire. Defaults to 0 meaning that the the secret_id is of unlimited use.`,
			},

			"secret_id_ttl": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `Duration in seconds after which the issued SecretID should expire. Defaults
to 0, meaning no expiration.`,
			},

			"period": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_period"),
				Deprecated:  true,
			},

			"role_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Identifier of the role. Defaults to a UUID.",
			},

			"local_secret_ids": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `If set, the secret IDs generated using this role will be cluster local. This
can only be set during role creation and once set, it can't be reset later.`,
			},
		},
		ExistenceCheck: b.pathRoleExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathRoleCreateUpdate,
			logical.UpdateOperation: b.pathRoleCreateUpdate,
			logical.ReadOperation:   b.pathRoleRead,
			logical.DeleteOperation: b.pathRoleDelete,
		},
		HelpSynopsis:    strings.TrimSpace(roleHelp["role"][0]),
		HelpDescription: strings.TrimSpace(roleHelp["role"][1]),
	}

	tokenutil.AddTokenFields(p.Fields)

	return []*framework.Path{
		p,
		&framework.Path{
			Pattern: "role/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathRoleList,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-list"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-list"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/local-secret-ids$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathRoleLocalSecretIDsRead,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-local-secret-ids"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-local-secret-ids"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/policies$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"policies": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: tokenutil.DeprecationText("token_policies"),
					Deprecated:  true,
				},
				"token_policies": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: defTokenFields["token_policies"].Description,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRolePoliciesUpdate,
				logical.ReadOperation:   b.pathRolePoliciesRead,
				logical.DeleteOperation: b.pathRolePoliciesDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-policies"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-policies"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/bound-cidr-list$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"bound_cidr_list": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Deprecated: Please use "secret_id_bound_cidrs" instead. Comma separated string or list 
of CIDR blocks. If set, specifies the blocks of IP addresses which can perform the login operation.`,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleBoundCIDRUpdate,
				logical.ReadOperation:   b.pathRoleBoundCIDRListRead,
				logical.DeleteOperation: b.pathRoleBoundCIDRListDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-bound-cidr-list"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-bound-cidr-list"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/secret-id-bound-cidrs$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"secret_id_bound_cidrs": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma separated string or list of CIDR blocks. If set, specifies the blocks of
IP addresses which can perform the login operation.`,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleSecretIDBoundCIDRUpdate,
				logical.ReadOperation:   b.pathRoleSecretIDBoundCIDRRead,
				logical.DeleteOperation: b.pathRoleSecretIDBoundCIDRDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["secret-id-bound-cidrs"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["secret-id-bound-cidrs"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/token-bound-cidrs$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"token_bound_cidrs": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: defTokenFields["token_bound_cidrs"].Description,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleTokenBoundCIDRUpdate,
				logical.ReadOperation:   b.pathRoleTokenBoundCIDRRead,
				logical.DeleteOperation: b.pathRoleTokenBoundCIDRDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["token-bound-cidrs"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["token-bound-cidrs"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/bind-secret-id$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"bind_secret_id": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Default:     true,
					Description: "Impose secret_id to be presented when logging in using this role.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleBindSecretIDUpdate,
				logical.ReadOperation:   b.pathRoleBindSecretIDRead,
				logical.DeleteOperation: b.pathRoleBindSecretIDDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-bind-secret-id"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-bind-secret-id"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/secret-id-num-uses$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"secret_id_num_uses": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Number of times a SecretID can access the role, after which the SecretID will expire.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleSecretIDNumUsesUpdate,
				logical.ReadOperation:   b.pathRoleSecretIDNumUsesRead,
				logical.DeleteOperation: b.pathRoleSecretIDNumUsesDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-secret-id-num-uses"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-secret-id-num-uses"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/secret-id-ttl$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"secret_id_ttl": &framework.FieldSchema{
					Type: framework.TypeDurationSecond,
					Description: `Duration in seconds after which the issued SecretID should expire. Defaults
to 0, meaning no expiration.`,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleSecretIDTTLUpdate,
				logical.ReadOperation:   b.pathRoleSecretIDTTLRead,
				logical.DeleteOperation: b.pathRoleSecretIDTTLDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-secret-id-ttl"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-secret-id-ttl"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/period$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"period": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: tokenutil.DeprecationText("token_period"),
					Deprecated:  true,
				},
				"token_period": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: defTokenFields["token_period"].Description,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRolePeriodUpdate,
				logical.ReadOperation:   b.pathRolePeriodRead,
				logical.DeleteOperation: b.pathRolePeriodDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-period"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-period"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/token-num-uses$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"token_num_uses": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: defTokenFields["token_num_uses"].Description,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleTokenNumUsesUpdate,
				logical.ReadOperation:   b.pathRoleTokenNumUsesRead,
				logical.DeleteOperation: b.pathRoleTokenNumUsesDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-token-num-uses"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-token-num-uses"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/token-ttl$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"token_ttl": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: defTokenFields["token_ttl"].Description,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleTokenTTLUpdate,
				logical.ReadOperation:   b.pathRoleTokenTTLRead,
				logical.DeleteOperation: b.pathRoleTokenTTLDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-token-ttl"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-token-ttl"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/token-max-ttl$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"token_max_ttl": &framework.FieldSchema{
					Type:        framework.TypeDurationSecond,
					Description: defTokenFields["token_max_ttl"].Description,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleTokenMaxTTLUpdate,
				logical.ReadOperation:   b.pathRoleTokenMaxTTLRead,
				logical.DeleteOperation: b.pathRoleTokenMaxTTLDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-token-max-ttl"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-token-max-ttl"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/role-id$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"role_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Identifier of the role. Defaults to a UUID.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   b.pathRoleRoleIDRead,
				logical.UpdateOperation: b.pathRoleRoleIDUpdate,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-id"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-id"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/secret-id/?$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"metadata": &framework.FieldSchema{
					Type: framework.TypeString,
					Description: `Metadata to be tied to the SecretID. This should be a JSON
formatted string containing the metadata in key value pairs.`,
				},
				"cidr_list": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma separated string or list of CIDR blocks enforcing secret IDs to be used from
specific set of IP addresses. If 'bound_cidr_list' is set on the role, then the
list of CIDR blocks listed here should be a subset of the CIDR blocks listed on
the role.`,
				},
				"token_bound_cidrs": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: defTokenFields["token_bound_cidrs"].Description,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleSecretIDUpdate,
				logical.ListOperation:   b.pathRoleSecretIDList,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-secret-id"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-secret-id"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/secret-id/lookup/?$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"secret_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "SecretID attached to the role.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleSecretIDLookupUpdate,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-secret-id-lookup"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-secret-id-lookup"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/secret-id/destroy/?$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"secret_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "SecretID attached to the role.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleSecretIDDestroyUpdateDelete,
				logical.DeleteOperation: b.pathRoleSecretIDDestroyUpdateDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-secret-id-destroy"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-secret-id-destroy"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/secret-id-accessor/lookup/?$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"secret_id_accessor": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Accessor of the SecretID",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleSecretIDAccessorLookupUpdate,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-secret-id-accessor"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-secret-id-accessor"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/secret-id-accessor/destroy/?$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"secret_id_accessor": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Accessor of the SecretID",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleSecretIDAccessorDestroyUpdateDelete,
				logical.DeleteOperation: b.pathRoleSecretIDAccessorDestroyUpdateDelete,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-secret-id-accessor"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-secret-id-accessor"][1]),
		},
		&framework.Path{
			Pattern: "role/" + framework.GenericNameRegex("role_name") + "/custom-secret-id$",
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
				"secret_id": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "SecretID to be attached to the role.",
				},
				"metadata": &framework.FieldSchema{
					Type: framework.TypeString,
					Description: `Metadata to be tied to the SecretID. This should be a JSON
formatted string containing metadata in key value pairs.`,
				},
				"cidr_list": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma separated string or list of CIDR blocks enforcing secret IDs to be used from
specific set of IP addresses. If 'bound_cidr_list' is set on the role, then the
list of CIDR blocks listed here should be a subset of the CIDR blocks listed on
the role.`,
				},
				"token_bound_cidrs": &framework.FieldSchema{
					Type: framework.TypeCommaStringSlice,
					Description: `Comma separated string or list of CIDR blocks. If set, specifies the blocks of
IP addresses which can use the returned token. Should be a subset of the token CIDR blocks listed on the role, if any.`,
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: b.pathRoleCustomSecretIDUpdate,
			},
			HelpSynopsis:    strings.TrimSpace(roleHelp["role-custom-secret-id"][0]),
			HelpDescription: strings.TrimSpace(roleHelp["role-custom-secret-id"][1]),
		},
	}
}

// pathRoleExistenceCheck returns whether the role with the given name exists or not.
func (b *backend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return false, fmt.Errorf("missing role_name")
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return false, err
	}

	return role != nil, nil
}

// pathRoleList is used to list all the Roles registered with the backend.
func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	lock := b.roleLock("")

	lock.RLock()
	defer lock.RUnlock()

	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

// pathRoleSecretIDList is used to list all the 'secret_id_accessor's issued against the role.
func (b *backend) pathRoleSecretIDList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	// Get the role entry
	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role %q does not exist", roleName)), nil
	}

	// Guard the list operation with an outer lock
	b.secretIDListingLock.RLock()
	defer b.secretIDListingLock.RUnlock()

	roleNameHMAC, err := createHMAC(role.HMACKey, role.name)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create HMAC of role_name: {{err}}", err)
	}

	// Listing works one level at a time. Get the first level of data
	// which could then be used to get the actual SecretID storage entries.
	secretIDHMACs, err := req.Storage.List(ctx, fmt.Sprintf("%s%s/", role.SecretIDPrefix, roleNameHMAC))
	if err != nil {
		return nil, err
	}

	var listItems []string
	for _, secretIDHMAC := range secretIDHMACs {
		// For sanity
		if secretIDHMAC == "" {
			continue
		}

		// Prepare the full index of the SecretIDs.
		entryIndex := fmt.Sprintf("%s%s/%s", role.SecretIDPrefix, roleNameHMAC, secretIDHMAC)

		// SecretID locks are not indexed by SecretIDs itself.
		// This is because SecretIDs are not stored in plaintext
		// form anywhere in the backend, and hence accessing its
		// corresponding lock many times using SecretIDs is not
		// possible. Also, indexing it everywhere using secretIDHMACs
		// makes listing operation easier.
		secretIDLock := b.secretIDLock(secretIDHMAC)

		secretIDLock.RLock()

		result := secretIDStorageEntry{}
		if entry, err := req.Storage.Get(ctx, entryIndex); err != nil {
			secretIDLock.RUnlock()
			return nil, err
		} else if entry == nil {
			secretIDLock.RUnlock()
			return nil, fmt.Errorf("storage entry for SecretID is present but no content found at the index")
		} else if err := entry.DecodeJSON(&result); err != nil {
			secretIDLock.RUnlock()
			return nil, err
		}
		listItems = append(listItems, result.SecretIDAccessor)
		secretIDLock.RUnlock()
	}

	return logical.ListResponse(listItems), nil
}

// validateRoleConstraints checks if the role has at least one constraint
// enabled.
func validateRoleConstraints(role *roleStorageEntry) error {
	if role == nil {
		return fmt.Errorf("nil role")
	}

	// At least one constraint should be enabled on the role
	switch {
	case role.BindSecretID:
	case len(role.BoundCIDRList) != 0:
	case len(role.SecretIDBoundCIDRs) != 0:
	case len(role.TokenBoundCIDRs) != 0:
	default:
		return fmt.Errorf("at least one constraint should be enabled on the role")
	}

	return nil
}

// setRoleEntry persists the role and creates an index from roleID to role
// name.
func (b *backend) setRoleEntry(ctx context.Context, s logical.Storage, roleName string, role *roleStorageEntry, previousRoleID string) error {
	if roleName == "" {
		return fmt.Errorf("missing role name")
	}

	if role == nil {
		return fmt.Errorf("nil role")
	}

	// Check if role constraints are properly set
	if err := validateRoleConstraints(role); err != nil {
		return err
	}

	// Create a storage entry for the role
	entry, err := logical.StorageEntryJSON("role/"+strings.ToLower(roleName), role)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("failed to create storage entry for role %q", roleName)
	}

	// Check if the index from the role_id to role already exists
	roleIDIndex, err := b.roleIDEntry(ctx, s, role.RoleID)
	if err != nil {
		return errwrap.Wrapf("failed to read role_id index: {{err}}", err)
	}

	// If the entry exists, make sure that it belongs to the current role
	if roleIDIndex != nil && roleIDIndex.Name != roleName {
		return fmt.Errorf("role_id already in use")
	}

	// When role_id is getting updated, delete the old index before
	// a new one is created
	if previousRoleID != "" && previousRoleID != role.RoleID {
		if err = b.roleIDEntryDelete(ctx, s, previousRoleID); err != nil {
			return errwrap.Wrapf("failed to delete previous role ID index: {{err}}", err)
		}
	}

	// Save the role entry only after all the validations
	if err = s.Put(ctx, entry); err != nil {
		return err
	}

	// If previousRoleID is still intact, don't create another one
	if previousRoleID != "" && previousRoleID == role.RoleID {
		return nil
	}

	// Create a storage entry for reverse mapping of RoleID to role.
	// Note that secondary index is created when the roleLock is held.
	return b.setRoleIDEntry(ctx, s, role.RoleID, &roleIDStorageEntry{
		Name: roleName,
	})
}

// roleEntry reads the role from storage
func (b *backend) roleEntry(ctx context.Context, s logical.Storage, roleName string) (*roleStorageEntry, error) {
	if roleName == "" {
		return nil, fmt.Errorf("missing role_name")
	}

	var role roleStorageEntry

	if entry, err := s.Get(ctx, "role/"+strings.ToLower(roleName)); err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	} else if err := entry.DecodeJSON(&role); err != nil {
		return nil, err
	}

	needsUpgrade := false

	if role.BoundCIDRListOld != "" {
		role.SecretIDBoundCIDRs = strutil.ParseDedupAndSortStrings(role.BoundCIDRListOld, ",")
		role.BoundCIDRListOld = ""
		needsUpgrade = true
	}

	if len(role.BoundCIDRList) != 0 {
		role.SecretIDBoundCIDRs = role.BoundCIDRList
		role.BoundCIDRList = nil
		needsUpgrade = true
	}

	if role.SecretIDPrefix == "" {
		role.SecretIDPrefix = secretIDPrefix
		needsUpgrade = true
	}

	if role.TokenPeriod == 0 && role.Period > 0 {
		role.TokenPeriod = role.Period
	}

	if len(role.TokenPolicies) == 0 && len(role.Policies) > 0 {
		role.TokenPolicies = role.Policies
	}

	if needsUpgrade && (b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby)) {
		entry, err := logical.StorageEntryJSON("role/"+strings.ToLower(roleName), &role)
		if err != nil {
			return nil, err
		}
		if err := s.Put(ctx, entry); err != nil {
			// Only perform upgrades on replication primary
			if !strings.Contains(err.Error(), logical.ErrReadOnly.Error()) {
				return nil, err
			}
		}
	}

	role.name = roleName
	if role.LowerCaseRoleName {
		role.name = strings.ToLower(roleName)
	}

	return &role, nil
}

// pathRoleCreateUpdate registers a new role with the backend or updates the options
// of an existing role
func (b *backend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	// Check if the role already exists
	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	// Create a new entry object if this is a CreateOperation
	switch {
	case role == nil && req.Operation == logical.CreateOperation:
		hmacKey, err := uuid.GenerateUUID()
		if err != nil {
			return nil, errwrap.Wrapf("failed to create role_id: {{err}}", err)
		}
		role = &roleStorageEntry{
			name:              strings.ToLower(roleName),
			HMACKey:           hmacKey,
			LowerCaseRoleName: true,
		}
	case role == nil:
		return logical.ErrorResponse(fmt.Sprintf("role name %q doesn't exist", roleName)), logical.ErrUnsupportedPath
	}

	var resp *logical.Response

	// Handle a backwards compat case
	if tokenTypeRaw, ok := data.Raw["token_type"]; ok {
		switch tokenTypeRaw.(string) {
		case "default-service":
			data.Raw["token_type"] = "service"
			resp.AddWarning("default-service has no useful meaning; adjusting to service")
		case "default-batch":
			data.Raw["token_type"] = "batch"
			resp.AddWarning("default-batch has no useful meaning; adjusting to batch")
		}
	}

	if err := role.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	localSecretIDsRaw, ok := data.GetOk("local_secret_ids")
	if ok {
		switch {
		case req.Operation == logical.CreateOperation:
			localSecretIDs := localSecretIDsRaw.(bool)
			if localSecretIDs {
				role.SecretIDPrefix = secretIDLocalPrefix
			}
		default:
			return logical.ErrorResponse("local_secret_ids can only be modified during role creation"), nil
		}
	}

	previousRoleID := role.RoleID
	if roleIDRaw, ok := data.GetOk("role_id"); ok {
		role.RoleID = roleIDRaw.(string)
	} else if req.Operation == logical.CreateOperation {
		roleID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate role_id: {{err}}", err)
		}
		role.RoleID = roleID
	}
	if role.RoleID == "" {
		return logical.ErrorResponse("invalid role_id supplied, or failed to generate a role_id"), nil
	}

	if bindSecretIDRaw, ok := data.GetOk("bind_secret_id"); ok {
		role.BindSecretID = bindSecretIDRaw.(bool)
	} else if req.Operation == logical.CreateOperation {
		role.BindSecretID = data.Get("bind_secret_id").(bool)
	}

	if boundCIDRListRaw, ok := data.GetFirst("secret_id_bound_cidrs", "bound_cidr_list"); ok {
		role.SecretIDBoundCIDRs = boundCIDRListRaw.([]string)
	}

	if len(role.SecretIDBoundCIDRs) != 0 {
		valid, err := cidrutil.ValidateCIDRListSlice(role.SecretIDBoundCIDRs)
		if err != nil {
			return nil, errwrap.Wrapf("failed to validate CIDR blocks: {{err}}", err)
		}
		if !valid {
			return logical.ErrorResponse("invalid CIDR blocks"), nil
		}
	}

	if secretIDNumUsesRaw, ok := data.GetOk("secret_id_num_uses"); ok {
		role.SecretIDNumUses = secretIDNumUsesRaw.(int)
	} else if req.Operation == logical.CreateOperation {
		role.SecretIDNumUses = data.Get("secret_id_num_uses").(int)
	}
	if role.SecretIDNumUses < 0 {
		return logical.ErrorResponse("secret_id_num_uses cannot be negative"), nil
	}

	if secretIDTTLRaw, ok := data.GetOk("secret_id_ttl"); ok {
		role.SecretIDTTL = time.Second * time.Duration(secretIDTTLRaw.(int))
	} else if req.Operation == logical.CreateOperation {
		role.SecretIDTTL = time.Second * time.Duration(data.Get("secret_id_ttl").(int))
	}

	// handle upgrade cases
	{
		if err := tokenutil.UpgradeValue(data, "policies", "token_policies", &role.Policies, &role.TokenPolicies); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "period", "token_period", &role.Period, &role.TokenPeriod); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	if role.Period > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse(fmt.Sprintf("period of %q is greater than the backend's maximum lease TTL of %q", role.Period.String(), b.System().MaxLeaseTTL().String())), nil
	}

	if role.TokenMaxTTL > b.System().MaxLeaseTTL() {
		resp = &logical.Response{}
		resp.AddWarning("token_max_ttl is greater than the backend mount's maximum TTL value; issued tokens' max TTL value will be truncated")
	}

	// Store the entry.
	return resp, b.setRoleEntry(ctx, req.Storage, role.name, role, previousRoleID)
}

// pathRoleRead grabs a read lock and reads the options set on the role from the storage
func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	lockRelease := lock.RUnlock

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		lockRelease()
		return nil, err
	}

	if role == nil {
		lockRelease()
		return nil, nil
	}

	respData := map[string]interface{}{
		"bind_secret_id":        role.BindSecretID,
		"secret_id_bound_cidrs": role.SecretIDBoundCIDRs,
		"secret_id_num_uses":    role.SecretIDNumUses,
		"secret_id_ttl":         role.SecretIDTTL / time.Second,
		"local_secret_ids":      false,
	}
	role.PopulateTokenData(respData)

	if role.SecretIDPrefix == secretIDLocalPrefix {
		respData["local_secret_ids"] = true
	}

	// Backwards compat data
	if role.Period != 0 {
		respData["period"] = role.Period / time.Second
	}
	if len(role.Policies) > 0 {
		respData["policies"] = role.Policies
	}

	resp := &logical.Response{
		Data: respData,
	}

	if err := validateRoleConstraints(role); err != nil {
		resp.AddWarning("Role does not have any constraints set on it. Updates to this role will require a constraint to be set")
	}

	// For sanity, verify that the index still exists. If the index is missing,
	// add one and return a warning so it can be reported.
	roleIDIndex, err := b.roleIDEntry(ctx, req.Storage, role.RoleID)
	if err != nil {
		lockRelease()
		return nil, err
	}

	if roleIDIndex == nil {
		// Switch to a write lock
		lock.RUnlock()
		lock.Lock()
		lockRelease = lock.Unlock

		// Check again if the index is missing
		roleIDIndex, err = b.roleIDEntry(ctx, req.Storage, role.RoleID)
		if err != nil {
			lockRelease()
			return nil, err
		}

		if roleIDIndex == nil {
			// Create a new index
			err = b.setRoleIDEntry(ctx, req.Storage, role.RoleID, &roleIDStorageEntry{
				Name: role.name,
			})
			if err != nil {
				lockRelease()
				return nil, errwrap.Wrapf(fmt.Sprintf("failed to create secondary index for role_id %q: {{err}}", role.RoleID), err)
			}
			resp.AddWarning("Role identifier was missing an index back to role name. A new index has been added. Please report this observation.")
		}
	}

	lockRelease()

	return resp, nil
}

// pathRoleDelete removes the role from the storage
func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Just before the role is deleted, remove all the SecretIDs issued as part of the role.
	if err = b.flushRoleSecrets(ctx, req.Storage, role.name, role.HMACKey, role.SecretIDPrefix); err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to invalidate the secrets belonging to role %q: {{err}}", role.name), err)
	}

	// Delete the reverse mapping from RoleID to the role
	if err = b.roleIDEntryDelete(ctx, req.Storage, role.RoleID); err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to delete the mapping from RoleID to role %q: {{err}}", role.name), err)
	}

	// After deleting the SecretIDs and the RoleID, delete the role itself
	if err = req.Storage.Delete(ctx, "role/"+strings.ToLower(role.name)); err != nil {
		return nil, err
	}

	return nil, nil
}

// Returns the properties of the SecretID
func (b *backend) pathRoleSecretIDLookupUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	secretID := data.Get("secret_id").(string)
	if secretID == "" {
		return logical.ErrorResponse("missing secret_id"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	// Fetch the role
	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role %q does not exist", roleName)
	}

	// Create the HMAC of the secret ID using the per-role HMAC key
	secretIDHMAC, err := createHMAC(role.HMACKey, secretID)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create HMAC of secret_id: {{err}}", err)
	}

	// Create the HMAC of the roleName using the per-role HMAC key
	roleNameHMAC, err := createHMAC(role.HMACKey, role.name)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create HMAC of role_name: {{err}}", err)
	}

	// Create the index at which the secret_id would've been stored
	entryIndex := fmt.Sprintf("%s%s/%s", role.SecretIDPrefix, roleNameHMAC, secretIDHMAC)

	secretLock := b.secretIDLock(secretIDHMAC)
	secretLock.Lock()
	defer secretLock.Unlock()

	secretIDEntry, err := b.nonLockedSecretIDStorageEntry(ctx, req.Storage, role.SecretIDPrefix, roleNameHMAC, secretIDHMAC)
	if err != nil {
		return nil, err
	}
	if secretIDEntry == nil {
		return nil, nil
	}

	// If a secret ID entry does not have a corresponding accessor
	// entry, revoke the secret ID immediately
	accessorEntry, err := b.secretIDAccessorEntry(ctx, req.Storage, secretIDEntry.SecretIDAccessor, role.SecretIDPrefix)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read secret ID accessor entry: {{err}}", err)
	}
	if accessorEntry == nil {
		if err := req.Storage.Delete(ctx, entryIndex); err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("error deleting secret ID %q from storage: {{err}}", secretIDHMAC), err)
		}
		return logical.ErrorResponse("invalid secret id"), nil
	}

	return &logical.Response{
		Data: secretIDEntry.ToResponseData(),
	}, nil
}

func (entry *secretIDStorageEntry) ToResponseData() map[string]interface{} {
	ret := map[string]interface{}{
		"secret_id_accessor": entry.SecretIDAccessor,
		"secret_id_num_uses": entry.SecretIDNumUses,
		"secret_id_ttl":      entry.SecretIDTTL / time.Second,
		"creation_time":      entry.CreationTime,
		"expiration_time":    entry.ExpirationTime,
		"last_updated_time":  entry.LastUpdatedTime,
		"metadata":           entry.Metadata,
		"cidr_list":          entry.CIDRList,
		"token_bound_cidrs":  entry.TokenBoundCIDRs,
	}
	if len(entry.TokenBoundCIDRs) == 0 {
		ret["token_bound_cidrs"] = []string{}
	}
	return ret
}

func (b *backend) pathRoleSecretIDDestroyUpdateDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	secretID := data.Get("secret_id").(string)
	if secretID == "" {
		return logical.ErrorResponse("missing secret_id"), nil
	}

	roleLock := b.roleLock(roleName)
	roleLock.RLock()
	defer roleLock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role %q does not exist", roleName)
	}

	secretIDHMAC, err := createHMAC(role.HMACKey, secretID)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create HMAC of secret_id: {{err}}", err)
	}

	roleNameHMAC, err := createHMAC(role.HMACKey, role.name)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create HMAC of role_name: {{err}}", err)
	}

	entryIndex := fmt.Sprintf("%s%s/%s", role.SecretIDPrefix, roleNameHMAC, secretIDHMAC)

	lock := b.secretIDLock(secretIDHMAC)
	lock.Lock()
	defer lock.Unlock()

	entry, err := b.nonLockedSecretIDStorageEntry(ctx, req.Storage, role.SecretIDPrefix, roleNameHMAC, secretIDHMAC)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	// Delete the accessor of the SecretID first
	if err := b.deleteSecretIDAccessorEntry(ctx, req.Storage, entry.SecretIDAccessor, role.SecretIDPrefix); err != nil {
		return nil, err
	}

	// Delete the storage entry that corresponds to the SecretID
	if err := req.Storage.Delete(ctx, entryIndex); err != nil {
		return nil, errwrap.Wrapf("failed to delete secret_id: {{err}}", err)
	}

	return nil, nil
}

// pathRoleSecretIDAccessorLookupUpdate returns the properties of the SecretID
// given its accessor
func (b *backend) pathRoleSecretIDAccessorLookupUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	secretIDAccessor := data.Get("secret_id_accessor").(string)
	if secretIDAccessor == "" {
		return logical.ErrorResponse("missing secret_id_accessor"), nil
	}

	// SecretID is indexed based on HMACed roleName and HMACed SecretID.
	// Get the role details to fetch the RoleID and accessor to get
	// the HMACed SecretID.

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role %q does not exist", roleName)
	}

	accessorEntry, err := b.secretIDAccessorEntry(ctx, req.Storage, secretIDAccessor, role.SecretIDPrefix)
	if err != nil {
		return nil, err
	}
	if accessorEntry == nil {
		return nil, fmt.Errorf("failed to find accessor entry for secret_id_accessor: %q", secretIDAccessor)
	}

	roleNameHMAC, err := createHMAC(role.HMACKey, role.name)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create HMAC of role_name: {{err}}", err)
	}

	secretLock := b.secretIDLock(accessorEntry.SecretIDHMAC)
	secretLock.RLock()
	defer secretLock.RUnlock()

	secretIDEntry, err := b.nonLockedSecretIDStorageEntry(ctx, req.Storage, role.SecretIDPrefix, roleNameHMAC, accessorEntry.SecretIDHMAC)
	if err != nil {
		return nil, err
	}
	if secretIDEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: secretIDEntry.ToResponseData(),
	}, nil
}

func (b *backend) pathRoleSecretIDAccessorDestroyUpdateDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	secretIDAccessor := data.Get("secret_id_accessor").(string)
	if secretIDAccessor == "" {
		return logical.ErrorResponse("missing secret_id_accessor"), nil
	}

	// SecretID is indexed based on HMACed roleName and HMACed SecretID.
	// Get the role details to fetch the RoleID and accessor to get
	// the HMACed SecretID.

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role %q does not exist", roleName)
	}

	accessorEntry, err := b.secretIDAccessorEntry(ctx, req.Storage, secretIDAccessor, role.SecretIDPrefix)
	if err != nil {
		return nil, err
	}
	if accessorEntry == nil {
		return nil, fmt.Errorf("failed to find accessor entry for secret_id_accessor: %q", secretIDAccessor)
	}

	roleNameHMAC, err := createHMAC(role.HMACKey, role.name)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create HMAC of role_name: {{err}}", err)
	}

	entryIndex := fmt.Sprintf("%s%s/%s", role.SecretIDPrefix, roleNameHMAC, accessorEntry.SecretIDHMAC)

	lock := b.secretIDLock(accessorEntry.SecretIDHMAC)
	lock.Lock()
	defer lock.Unlock()

	// Delete the accessor of the SecretID first
	if err := b.deleteSecretIDAccessorEntry(ctx, req.Storage, secretIDAccessor, role.SecretIDPrefix); err != nil {
		return nil, err
	}

	// Delete the storage entry that corresponds to the SecretID
	if err := req.Storage.Delete(ctx, entryIndex); err != nil {
		return nil, errwrap.Wrapf("failed to delete secret_id: {{err}}", err)
	}

	return nil, nil
}

func (b *backend) pathRoleBoundCIDRUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	delete(data.Raw, "token_bound_cidrs")
	delete(data.Raw, "secret_id_bound_cidrs")
	return b.pathRoleBoundCIDRUpdateCommon(ctx, req, data)
}

func (b *backend) pathRoleSecretIDBoundCIDRUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	delete(data.Raw, "bound_cidr_list")
	delete(data.Raw, "token_bound_cidrs")
	return b.pathRoleBoundCIDRUpdateCommon(ctx, req, data)
}

func (b *backend) pathRoleTokenBoundCIDRUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	delete(data.Raw, "bound_cidr_list")
	delete(data.Raw, "secret_id_bound_cidrs")
	return b.pathRoleBoundCIDRUpdateCommon(ctx, req, data)
}

func (b *backend) pathRoleBoundCIDRUpdateCommon(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	// Re-read the role after grabbing the lock
	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	if cidrsIfc, ok := data.GetFirst("secret_id_bound_cidrs", "bound_cidr_list"); ok {
		cidrs := cidrsIfc.([]string)
		if len(cidrs) == 0 {
			return logical.ErrorResponse("missing bound_cidr_list"), nil
		}
		valid, err := cidrutil.ValidateCIDRListSlice(cidrs)
		if err != nil {
			return logical.ErrorResponse(errwrap.Wrapf("failed to validate CIDR blocks: {{err}}", err).Error()), nil
		}
		if !valid {
			return logical.ErrorResponse("failed to validate CIDR blocks"), nil
		}
		role.SecretIDBoundCIDRs = cidrs

	} else if cidrsIfc, ok := data.GetOk("token_bound_cidrs"); ok {
		cidrs, err := parseutil.ParseAddrs(cidrsIfc.([]string))
		if err != nil {
			return logical.ErrorResponse(errwrap.Wrapf("failed to parse token_bound_cidrs: {{err}}", err).Error()), nil
		}
		role.TokenBoundCIDRs = cidrs
	}

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRoleSecretIDBoundCIDRRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.pathRoleFieldRead(ctx, req, data, "secret_id_bound_cidrs")
}

func (b *backend) pathRoleTokenBoundCIDRRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.pathRoleFieldRead(ctx, req, data, "token_bound_cidrs")
}

func (b *backend) pathRoleBoundCIDRListRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.pathRoleFieldRead(ctx, req, data, "bound_cidr_list")
}

func (b *backend) pathRoleFieldRead(ctx context.Context, req *logical.Request, data *framework.FieldData, fieldName string) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	} else {
		switch fieldName {
		case "secret_id_bound_cidrs":
			return &logical.Response{
				Data: map[string]interface{}{
					"secret_id_bound_cidrs": role.SecretIDBoundCIDRs,
				},
			}, nil
		case "token_bound_cidrs":
			return &logical.Response{
				Data: map[string]interface{}{
					"token_bound_cidrs": role.TokenBoundCIDRs,
				},
			}, nil
		case "bound_cidr_list":
			resp := &logical.Response{
				Data: map[string]interface{}{
					"bound_cidr_list": role.BoundCIDRList,
				},
			}
			resp.AddWarning(`The "bound_cidr_list" parameter is deprecated and will be removed. Please use "secret_id_bound_cidrs" instead.`)
			return resp, nil
		default:
			// shouldn't occur IRL
			return nil, errors.New("unrecognized field provided: " + fieldName)
		}
	}
}

func (b *backend) pathRoleBoundCIDRDelete(ctx context.Context, req *logical.Request, data *framework.FieldData, fieldName string) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	switch fieldName {
	case "bound_cidr_list":
		role.BoundCIDRList = nil
	case "secret_id_bound_cidrs":
		role.SecretIDBoundCIDRs = nil
	case "token_bound_cidrs":
		role.TokenBoundCIDRs = nil
	}

	return nil, b.setRoleEntry(ctx, req.Storage, roleName, role, "")
}

func (b *backend) pathRoleBoundCIDRListDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.pathRoleBoundCIDRDelete(ctx, req, data, "bound_cidr_list")
}

func (b *backend) pathRoleSecretIDBoundCIDRDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.pathRoleBoundCIDRDelete(ctx, req, data, "secret_id_bound_cidrs")
}

func (b *backend) pathRoleTokenBoundCIDRDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.pathRoleBoundCIDRDelete(ctx, req, data, "token_bound_cidrs")
}

func (b *backend) pathRoleBindSecretIDUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	if bindSecretIDRaw, ok := data.GetOk("bind_secret_id"); ok {
		role.BindSecretID = bindSecretIDRaw.(bool)
		return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
	} else {
		return logical.ErrorResponse("missing bind_secret_id"), nil
	}
}

func (b *backend) pathRoleBindSecretIDRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"bind_secret_id": role.BindSecretID,
		},
	}, nil
}

func (b *backend) pathRoleBindSecretIDDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Deleting a field implies setting the value to it's default value.
	role.BindSecretID = data.GetDefaultOrZero("bind_secret_id").(bool)

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRoleLocalSecretIDsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	localSecretIDs := false
	if role.SecretIDPrefix == secretIDLocalPrefix {
		localSecretIDs = true
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"local_secret_ids": localSecretIDs,
		},
	}, nil
}

func (b *backend) pathRolePoliciesUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	policiesRaw, ok := data.GetOk("token_policies")
	if !ok {
		policiesRaw, ok = data.GetOk("policies")
		if ok {
			role.Policies = policyutil.ParsePolicies(policiesRaw)
			role.TokenPolicies = role.Policies
		} else {
			return logical.ErrorResponse("missing token_policies"), nil
		}
	} else {
		role.TokenPolicies = policyutil.ParsePolicies(policiesRaw)
		_, ok = data.GetOk("policies")
		if ok {
			role.Policies = role.TokenPolicies
		} else {
			role.Policies = nil
		}
	}

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRolePoliciesRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	p := role.TokenPolicies
	if p == nil {
		p = []string{}
	}
	d := map[string]interface{}{
		"token_policies": p,
	}

	if len(role.Policies) > 0 {
		d["policies"] = role.Policies
	}

	return &logical.Response{
		Data: d,
	}, nil
}

func (b *backend) pathRolePoliciesDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.TokenPolicies = nil
	role.Policies = nil

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRoleSecretIDNumUsesUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	if numUsesRaw, ok := data.GetOk("secret_id_num_uses"); ok {
		role.SecretIDNumUses = numUsesRaw.(int)
		if role.SecretIDNumUses < 0 {
			return logical.ErrorResponse("secret_id_num_uses cannot be negative"), nil
		}
		return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
	} else {
		return logical.ErrorResponse("missing secret_id_num_uses"), nil
	}
}

func (b *backend) pathRoleRoleIDUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	previousRoleID := role.RoleID
	role.RoleID = data.Get("role_id").(string)
	if role.RoleID == "" {
		return logical.ErrorResponse("missing role_id"), nil
	}

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, previousRoleID)
}

func (b *backend) pathRoleRoleIDRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"role_id": role.RoleID,
		},
	}, nil
}

func (b *backend) pathRoleSecretIDNumUsesRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"secret_id_num_uses": role.SecretIDNumUses,
		},
	}, nil
}

func (b *backend) pathRoleSecretIDNumUsesDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.SecretIDNumUses = data.GetDefaultOrZero("secret_id_num_uses").(int)

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRoleSecretIDTTLUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	if secretIDTTLRaw, ok := data.GetOk("secret_id_ttl"); ok {
		role.SecretIDTTL = time.Second * time.Duration(secretIDTTLRaw.(int))
		return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
	} else {
		return logical.ErrorResponse("missing secret_id_ttl"), nil
	}
}

func (b *backend) pathRoleSecretIDTTLRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"secret_id_ttl": role.SecretIDTTL / time.Second,
		},
	}, nil
}

func (b *backend) pathRoleSecretIDTTLDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.SecretIDTTL = time.Second * time.Duration(data.GetDefaultOrZero("secret_id_ttl").(int))

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRolePeriodUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	periodRaw, ok := data.GetOk("token_period")
	if !ok {
		periodRaw, ok = data.GetOk("period")
		if ok {
			role.Period = time.Second * time.Duration(periodRaw.(int))
			role.TokenPeriod = role.Period
		} else {
			return logical.ErrorResponse("missing period"), nil
		}
	} else {
		role.TokenPeriod = time.Second * time.Duration(periodRaw.(int))
		_, ok = data.GetOk("period")
		if ok {
			role.Period = role.TokenPeriod
		} else {
			role.Period = 0
		}
	}

	if role.TokenPeriod > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse(fmt.Sprintf("period of %q is greater than the backend's maximum lease TTL of %q", role.Period.String(), b.System().MaxLeaseTTL().String())), nil
	}
	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRolePeriodRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	d := map[string]interface{}{
		"token_period": role.TokenPeriod / time.Second,
	}

	if role.Period > 0 {
		d["period"] = role.Period / time.Second
	}

	return &logical.Response{
		Data: d,
	}, nil
}

func (b *backend) pathRolePeriodDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.TokenPeriod = 0
	role.Period = 0

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRoleTokenNumUsesUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	if tokenNumUsesRaw, ok := data.GetOk("token_num_uses"); ok {
		role.TokenNumUses = tokenNumUsesRaw.(int)
		return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
	} else {
		return logical.ErrorResponse("missing token_num_uses"), nil
	}
}

func (b *backend) pathRoleTokenNumUsesRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token_num_uses": role.TokenNumUses,
		},
	}, nil
}

func (b *backend) pathRoleTokenNumUsesDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.TokenNumUses = data.GetDefaultOrZero("token_num_uses").(int)

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRoleTokenTTLUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	if tokenTTLRaw, ok := data.GetOk("token_ttl"); ok {
		role.TokenTTL = time.Second * time.Duration(tokenTTLRaw.(int))
		if role.TokenMaxTTL > time.Duration(0) && role.TokenTTL > role.TokenMaxTTL {
			return logical.ErrorResponse("token_ttl should not be greater than token_max_ttl"), nil
		}
		return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
	} else {
		return logical.ErrorResponse("missing token_ttl"), nil
	}
}

func (b *backend) pathRoleTokenTTLRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token_ttl": role.TokenTTL / time.Second,
		},
	}, nil
}

func (b *backend) pathRoleTokenTTLDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.TokenTTL = time.Second * time.Duration(data.GetDefaultOrZero("token_ttl").(int))

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRoleTokenMaxTTLUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, logical.ErrUnsupportedPath
	}

	if tokenMaxTTLRaw, ok := data.GetOk("token_max_ttl"); ok {
		role.TokenMaxTTL = time.Second * time.Duration(tokenMaxTTLRaw.(int))
		if role.TokenMaxTTL > time.Duration(0) && role.TokenTTL > role.TokenMaxTTL {
			return logical.ErrorResponse("token_max_ttl should be greater than or equal to token_ttl"), nil
		}
		return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
	} else {
		return logical.ErrorResponse("missing token_max_ttl"), nil
	}
}

func (b *backend) pathRoleTokenMaxTTLRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token_max_ttl": role.TokenMaxTTL / time.Second,
		},
	}, nil
}

func (b *backend) pathRoleTokenMaxTTLDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	lock := b.roleLock(roleName)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.TokenMaxTTL = time.Second * time.Duration(data.GetDefaultOrZero("token_max_ttl").(int))

	return nil, b.setRoleEntry(ctx, req.Storage, role.name, role, "")
}

func (b *backend) pathRoleSecretIDUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	secretID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, errwrap.Wrapf("failed to generate secret_id: {{err}}", err)
	}
	return b.handleRoleSecretIDCommon(ctx, req, data, secretID)
}

func (b *backend) pathRoleCustomSecretIDUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.handleRoleSecretIDCommon(ctx, req, data, data.Get("secret_id").(string))
}

func (b *backend) handleRoleSecretIDCommon(ctx context.Context, req *logical.Request, data *framework.FieldData, secretID string) (*logical.Response, error) {
	roleName := data.Get("role_name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role_name"), nil
	}

	if secretID == "" {
		return logical.ErrorResponse("missing secret_id"), nil
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role %q does not exist", roleName)), logical.ErrUnsupportedPath
	}

	if !role.BindSecretID {
		return logical.ErrorResponse("bind_secret_id is not set on the role"), nil
	}

	secretIDCIDRs := data.Get("cidr_list").([]string)

	// Validate the list of CIDR blocks
	if len(secretIDCIDRs) != 0 {
		valid, err := cidrutil.ValidateCIDRListSlice(secretIDCIDRs)
		if err != nil {
			return nil, errwrap.Wrapf("failed to validate CIDR blocks: {{err}}", err)
		}
		if !valid {
			return logical.ErrorResponse("failed to validate CIDR blocks"), nil
		}
	}
	// Ensure that the CIDRs on the secret ID are a subset of that of role's
	if err := verifyCIDRRoleSecretIDSubset(secretIDCIDRs, role.SecretIDBoundCIDRs); err != nil {
		return nil, err
	}

	secretIDTokenCIDRs := data.Get("token_bound_cidrs").([]string)
	if len(secretIDTokenCIDRs) != 0 {
		valid, err := cidrutil.ValidateCIDRListSlice(secretIDTokenCIDRs)
		if err != nil {
			return nil, errwrap.Wrapf("failed to validate token CIDR blocks: {{err}}", err)
		}
		if !valid {
			return logical.ErrorResponse("failed to validate token CIDR blocks"), nil
		}
	}
	// Ensure that the token CIDRs on the secret ID are a subset of that of role's
	var roleCIDRs []string
	for _, v := range role.TokenBoundCIDRs {
		roleCIDRs = append(roleCIDRs, v.String())
	}
	if err := verifyCIDRRoleSecretIDSubset(secretIDTokenCIDRs, roleCIDRs); err != nil {
		return nil, err
	}

	secretIDStorage := &secretIDStorageEntry{
		SecretIDNumUses: role.SecretIDNumUses,
		SecretIDTTL:     role.SecretIDTTL,
		Metadata:        make(map[string]string),
		CIDRList:        secretIDCIDRs,
		TokenBoundCIDRs: secretIDTokenCIDRs,
	}

	if err = strutil.ParseArbitraryKeyValues(data.Get("metadata").(string), secretIDStorage.Metadata, ","); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse metadata: %v", err)), nil
	}

	if secretIDStorage, err = b.registerSecretIDEntry(ctx, req.Storage, role.name, secretID, role.HMACKey, role.SecretIDPrefix, secretIDStorage); err != nil {
		return nil, errwrap.Wrapf("failed to store secret_id: {{err}}", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"secret_id":          secretID,
			"secret_id_accessor": secretIDStorage.SecretIDAccessor,
		},
	}, nil
}

func (b *backend) roleIDLock(roleID string) *locksutil.LockEntry {
	return locksutil.LockForKey(b.roleIDLocks, roleID)
}

func (b *backend) roleLock(roleName string) *locksutil.LockEntry {
	return locksutil.LockForKey(b.roleLocks, strings.ToLower(roleName))
}

// setRoleIDEntry creates a storage entry that maps RoleID to Role
func (b *backend) setRoleIDEntry(ctx context.Context, s logical.Storage, roleID string, roleIDEntry *roleIDStorageEntry) error {
	lock := b.roleIDLock(roleID)
	lock.Lock()
	defer lock.Unlock()

	salt, err := b.Salt(ctx)
	if err != nil {
		return err
	}
	entryIndex := "role_id/" + salt.SaltID(roleID)

	entry, err := logical.StorageEntryJSON(entryIndex, roleIDEntry)
	if err != nil {
		return err
	}
	if err = s.Put(ctx, entry); err != nil {
		return err
	}
	return nil
}

// roleIDEntry is used to read the storage entry that maps RoleID to Role
func (b *backend) roleIDEntry(ctx context.Context, s logical.Storage, roleID string) (*roleIDStorageEntry, error) {
	if roleID == "" {
		return nil, fmt.Errorf("missing role id")
	}

	lock := b.roleIDLock(roleID)
	lock.RLock()
	defer lock.RUnlock()

	var result roleIDStorageEntry

	salt, err := b.Salt(ctx)
	if err != nil {
		return nil, err
	}
	entryIndex := "role_id/" + salt.SaltID(roleID)

	if entry, err := s.Get(ctx, entryIndex); err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	} else if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// roleIDEntryDelete is used to remove the secondary index that maps the
// RoleID to the Role itself.
func (b *backend) roleIDEntryDelete(ctx context.Context, s logical.Storage, roleID string) error {
	if roleID == "" {
		return fmt.Errorf("missing role id")
	}

	lock := b.roleIDLock(roleID)
	lock.Lock()
	defer lock.Unlock()

	salt, err := b.Salt(ctx)
	if err != nil {
		return err
	}
	entryIndex := "role_id/" + salt.SaltID(roleID)

	return s.Delete(ctx, entryIndex)
}

var roleHelp = map[string][2]string{
	"role-list": {
		"Lists all the roles registered with the backend.",
		"The list will contain the names of the roles.",
	},
	"role": {
		"Register an role with the backend.",
		`A role can represent a service, a machine or anything that can be IDed.
The set of policies on the role defines access to the role, meaning, any
Vault token with a policy set that is a superset of the policies on the
role registered here will have access to the role. If a SecretID is desired
to be generated against only this specific role, it can be done via
'role/<role_name>/secret-id' and 'role/<role_name>/custom-secret-id' endpoints.
The properties of the SecretID created against the role and the properties
of the token issued with the SecretID generated against the role, can be
configured using the parameters of this endpoint.`,
	},
	"role-bind-secret-id": {
		"Impose secret_id to be presented during login using this role.",
		`By setting this to 'true', during login the parameter 'secret_id' becomes a mandatory argument.
The value of 'secret_id' can be retrieved using 'role/<role_name>/secret-id' endpoint.`,
	},
	"role-bound-cidr-list": {
		`Deprecated: Comma separated list of CIDR blocks, if set, specifies blocks of IP
addresses which can perform the login operation`,
		`During login, the IP address of the client will be checked to see if it
belongs to the CIDR blocks specified. If CIDR blocks were set and if the
IP is not encompassed by it, login fails`,
	},
	"secret-id-bound-cidrs": {
		`Comma separated list of CIDR blocks, if set, specifies blocks of IP
addresses which can perform the login operation`,
		`During login, the IP address of the client will be checked to see if it
belongs to the CIDR blocks specified. If CIDR blocks were set and if the
IP is not encompassed by it, login fails`,
	},
	"token-bound-cidrs": {
		`Comma separated string or list of CIDR blocks. If set, specifies the blocks of
IP addresses which can use the returned token.`,
		`During use of the returned token, the IP address of the client will be checked to see if it
belongs to the CIDR blocks specified. If CIDR blocks were set and if the
IP is not encompassed by it, token use fails`,
	},
	"role-policies": {
		"Policies of the role.",
		`A comma-delimited set of Vault policies that defines access to the role.
All the Vault tokens with policies that encompass the policy set
defined on the role, can access the role.`,
	},
	"role-secret-id-num-uses": {
		"Use limit of the SecretID generated against the role.",
		`If the SecretIDs are generated/assigned against the role using the
'role/<role_name>/secret-id' or 'role/<role_name>/custom-secret-id' endpoints,
then the number of times that SecretID can access the role is defined by
this option.`,
	},
	"role-secret-id-ttl": {
		`Duration in seconds, representing the lifetime of the SecretIDs
that are generated against the role using 'role/<role_name>/secret-id' or
'role/<role_name>/custom-secret-id' endpoints.`,
		``,
	},
	"role-secret-id-lookup": {
		"Read the properties of an issued secret_id",
		`This endpoint is used to read the properties of a secret_id associated to a
role.`},
	"role-secret-id-destroy": {
		"Invalidate an issued secret_id",
		`This endpoint is used to delete the properties of a secret_id associated to a
role.`},
	"role-secret-id-accessor-lookup": {
		"Read an issued secret_id, using its accessor",
		`This is particularly useful to lookup the non-expiring 'secret_id's.
The list operation on the 'role/<role_name>/secret-id' endpoint will return
the 'secret_id_accessor's. This endpoint can be used to read the properties
of the secret. If the 'secret_id_num_uses' field in the response is 0, it
represents a non-expiring 'secret_id'.`,
	},
	"role-secret-id-accessor-destroy": {
		"Delete an issued secret_id, using its accessor",
		`This is particularly useful to clean-up the non-expiring 'secret_id's.
The list operation on the 'role/<role_name>/secret-id' endpoint will return
the 'secret_id_accessor's. This endpoint can be used to read the properties
of the secret. If the 'secret_id_num_uses' field in the response is 0, it
represents a non-expiring 'secret_id'.`,
	},
	"role-token-num-uses": {
		"Number of times issued tokens can be used",
		`By default, this will be set to zero, indicating that the issued
tokens can be used any number of times.`,
	},
	"role-token-ttl": {
		`Duration in seconds, the lifetime of the token issued by using the SecretID that
is generated against this role, before which the token needs to be renewed.`,
		`If SecretIDs are generated against the role, using 'role/<role_name>/secret-id' or the
'role/<role_name>/custom-secret-id' endpoints, and if those SecretIDs are used
to perform the login operation, then the value of 'token-ttl' defines the
lifetime of the token issued, before which the token needs to be renewed.`,
	},
	"role-token-max-ttl": {
		`Duration in seconds, the maximum lifetime of the tokens issued by using
the SecretIDs that were generated against this role, after which the
tokens are not allowed to be renewed.`,
		`If SecretIDs are generated against the role using 'role/<role_name>/secret-id'
or the 'role/<role_name>/custom-secret-id' endpoints, and if those SecretIDs
are used to perform the login operation, then the value of 'token-max-ttl'
defines the maximum lifetime of the tokens issued, after which the tokens
cannot be renewed. A reauthentication is required after this duration.
This value will be capped by the backend mount's maximum TTL value.`,
	},
	"role-id": {
		"Returns the 'role_id' of the role.",
		`If login is performed from an role, then its 'role_id' should be presented
as a credential during the login. This 'role_id' can be retrieved using
this endpoint.`,
	},
	"role-secret-id": {
		"Generate a SecretID against this role.",
		`The SecretID generated using this endpoint will be scoped to access
just this role and none else. The properties of this SecretID will be
based on the options set on the role. It will expire after a period
defined by the 'secret_id_ttl' option on the role and/or the backend
mount's maximum TTL value.`,
	},
	"role-custom-secret-id": {
		"Assign a SecretID of choice against the role.",
		`This option is not recommended unless there is a specific need
to do so. This will assign a client supplied SecretID to be used to access
the role. This SecretID will behave similarly to the SecretIDs generated by
the backend. The properties of this SecretID will be based on the options
set on the role. It will expire after a period defined by the 'secret_id_ttl'
option on the role and/or the backend mount's maximum TTL value.`,
	},
	"role-period": {
		"Updates the value of 'period' on the role",
		`If set,  indicates that the token generated using this role
should never expire. The token should be renewed within the
duration specified by this value. The renewal duration will
be fixed. If the Period in the role is modified, the token
will pick up the new value during its next renewal.`,
	},
	"role-local-secret-ids": {
		"Enables cluster local secret IDs",
		`If set, the secret IDs generated using this role will be cluster local.
This can only be set during role creation and once set, it can't be
reset later.`,
	},
}
