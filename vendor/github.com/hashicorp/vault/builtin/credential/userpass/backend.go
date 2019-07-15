package userpass

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
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathUsers(&b),
			pathUsersList(&b),
			pathUserPolicies(&b),
			pathUserPassword(&b),
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
The "userpass" credential provider allows authentication using
a combination of a username and password. No additional factors
are supported.

The username/password combination is configured using the "users/"
endpoints by a user with root access. Authentication is then done
by supplying the two fields for "login".
`
