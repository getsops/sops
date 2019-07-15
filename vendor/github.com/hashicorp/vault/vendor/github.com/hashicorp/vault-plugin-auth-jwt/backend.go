package jwtauth

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"
)

const (
	configPath string = "config"
	rolePrefix string = "role/"
)

// Factory is used by framework
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

type jwtAuthBackend struct {
	*framework.Backend

	l            sync.RWMutex
	provider     *oidc.Provider
	keySet       oidc.KeySet
	cachedConfig *jwtConfig
	oidcStates   *cache.Cache

	providerCtx       context.Context
	providerCtxCancel context.CancelFunc
}

func backend() *jwtAuthBackend {
	b := new(jwtAuthBackend)
	b.providerCtx, b.providerCtxCancel = context.WithCancel(context.Background())
	b.oidcStates = cache.New(oidcStateTimeout, 1*time.Minute)

	b.Backend = &framework.Backend{
		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
		Invalidate:  b.invalidate,
		Help:        backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
				"oidc/auth_url",
				"oidc/callback",

				// Uncomment to mount simple UI handler for local development
				// "ui",
			},
			SealWrapStorage: []string{
				"config",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathLogin(b),
				pathRoleList(b),
				pathRole(b),
				pathConfig(b),

				// Uncomment to mount simple UI handler for local development
				// pathUI(b),
			},
			pathOIDC(b),
		),
		Clean: b.cleanup,
	}

	return b
}

func (b *jwtAuthBackend) cleanup(_ context.Context) {
	b.l.Lock()
	if b.providerCtxCancel != nil {
		b.providerCtxCancel()
	}
	b.l.Unlock()
}

func (b *jwtAuthBackend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config":
		b.reset()
	}
}

func (b *jwtAuthBackend) reset() {
	b.l.Lock()
	b.provider = nil
	b.cachedConfig = nil
	b.l.Unlock()
}

func (b *jwtAuthBackend) getProvider(config *jwtConfig) (*oidc.Provider, error) {
	b.l.RLock()
	unlockFunc := b.l.RUnlock
	defer func() { unlockFunc() }()

	if b.provider != nil {
		return b.provider, nil
	}

	b.l.RUnlock()
	b.l.Lock()
	unlockFunc = b.l.Unlock

	if b.provider != nil {
		return b.provider, nil
	}

	provider, err := b.createProvider(config)
	if err != nil {
		return nil, err
	}

	b.provider = provider
	return provider, nil
}

// getKeySet returns a new JWKS KeySet based on the provided config.
func (b *jwtAuthBackend) getKeySet(config *jwtConfig) (oidc.KeySet, error) {
	b.l.Lock()
	defer b.l.Unlock()

	if b.keySet != nil {
		return b.keySet, nil
	}

	if config.JWKSURL == "" {
		return nil, errors.New("keyset error: jwks_url not configured")
	}

	ctx, err := b.createCAContext(b.providerCtx, config.JWKSCAPEM)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing jwks_ca_pem: {{err}}", err)
	}

	b.keySet = oidc.NewRemoteKeySet(ctx, config.JWKSURL)

	return b.keySet, nil
}

const (
	backendHelp = `
The JWT backend plugin allows authentication using JWTs (including OIDC).
`
)
