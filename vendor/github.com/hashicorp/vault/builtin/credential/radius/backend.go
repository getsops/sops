package radius

import (
	"context"

	"github.com/hashicorp/vault/helper/mfa"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: mfa.MFARootPaths(),

			Unauthenticated: []string{
				"login",
				"login/*",
			},

			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathUsers(&b),
			pathUsersList(&b),
		},
			mfa.MFAPaths(b.Backend, pathLogin(&b))...,
		),

		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
	}

	return &b
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The "radius" credential provider allows authentication against
a RADIUS server, checking username and associating users
to set of policies.

Configuration of the server is done through the "config" and "users"
endpoints by a user with appropriate access mandated by policy.
Authentication is then done by supplying the two fields for "login".

The backend optionally allows to grant a set of policies to any 
user that successfully authenticates against the RADIUS server, 
without them being explicitly mapped in vault.
`
