package ocikms

import (
	"crypto/rsa"
	"os"
	"sync"

	ocep "github.com/ontariosystems/oci-cli-env-provider"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
)

// newIPProvider is a variable to allow tests to stub the Instance Principal provider factory
var newIPProvider = auth.InstancePrincipalConfigurationProvider

// configurationProvider composes multiple OCI configuration providers to make
// authentication work seamlessly across environments.
// Order of precedence (composing provider will try each in order until one works):
// 1) OCI_CLI_* environment variables (via ontariosystems/oci-cli-env-provider)
// 2) OCI_* environment variables (native SDK env provider)
// 3) Config file providers (OCI_CLI_CONFIG_FILE/PROFILE if set)
// 4) Default config provider (~/.oci/config, TF_VAR_*)
// 5) Instance Principals (when running on OCI compute) - lazily evaluated as last resort
func configurationProvider() (common.ConfigurationProvider, error) {
	var providers []common.ConfigurationProvider

	// 1) Prefer the CLI-compatible envs used widely in CI/containers (envs only; no implicit fallbacks)
	providers = append(providers, ocep.OciCliEnvironmentConfigurationProvider())

	// 2) Native SDK envs (OCI_tenancy_ocid, OCI_user_ocid, OCI_fingerprint, OCI_private_key_path, OCI_region)
	providers = append(providers, common.ConfigurationProviderEnvironmentVariables("OCI", ""))

	// 3) File-based fallbacks
	if cfg := os.Getenv(OCICLIConfigFile); cfg != "" {
		if prof := os.Getenv(OCICLIProfile); prof != "" {
			if p, err := common.ConfigurationProviderFromFileWithProfile(cfg, prof, ""); err == nil {
				providers = append(providers, p)
			}
		} else {
			if p, err := common.ConfigurationProviderFromFile(cfg, ""); err == nil {
				providers = append(providers, p)
			}
		}
	}

	// 4) Default config provider (~/.oci/config, TF_VAR_*)
	providers = append(providers, common.DefaultConfigProvider())

	// 5) Instance principals for compute instances (lazy, only called if nothing else works)
	providers = append(providers, &lazyConfigurationProvider{factory: newIPProvider})

	return common.ComposingConfigurationProvider(providers)
}

// lazyConfigurationProvider wraps a ConfigurationProvider factory function and defers its
// creation until the first method call. This is useful for expensive providers
// like Instance Principal that may timeout or fail in non-OCI environments.
type lazyConfigurationProvider struct {
	factory  func() (common.ConfigurationProvider, error)
	provider common.ConfigurationProvider
	once     sync.Once
	err      error
}

var _ common.ConfigurationProvider = (*lazyConfigurationProvider)(nil)

func (l *lazyConfigurationProvider) init() {
	l.provider, l.err = l.factory()
}

func (l *lazyConfigurationProvider) TenancyOCID() (string, error) {
	l.once.Do(l.init)
	if l.err != nil {
		return "", l.err
	}
	return l.provider.TenancyOCID()
}

func (l *lazyConfigurationProvider) UserOCID() (string, error) {
	l.once.Do(l.init)
	if l.err != nil {
		return "", l.err
	}
	return l.provider.UserOCID()
}

func (l *lazyConfigurationProvider) KeyFingerprint() (string, error) {
	l.once.Do(l.init)
	if l.err != nil {
		return "", l.err
	}
	return l.provider.KeyFingerprint()
}

func (l *lazyConfigurationProvider) Region() (string, error) {
	l.once.Do(l.init)
	if l.err != nil {
		return "", l.err
	}
	return l.provider.Region()
}

func (l *lazyConfigurationProvider) KeyID() (string, error) {
	l.once.Do(l.init)
	if l.err != nil {
		return "", l.err
	}
	return l.provider.KeyID()
}

func (l *lazyConfigurationProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	l.once.Do(l.init)
	if l.err != nil {
		return nil, l.err
	}
	return l.provider.PrivateRSAKey()
}

func (l *lazyConfigurationProvider) AuthType() (common.AuthConfig, error) {
	l.once.Do(l.init)
	if l.err != nil {
		return common.AuthConfig{}, l.err
	}
	return l.provider.AuthType()
}
