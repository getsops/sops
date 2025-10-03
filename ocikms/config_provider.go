package ocikms

import (
	"os"

	ocep "github.com/ontariosystems/oci-cli-env-provider"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
)

// newIPProvider is a variable to allow tests to stub the Instance Principal provider factory
var newIPProvider = auth.InstancePrincipalConfigurationProvider

// configurationProvider composes multiple OCI configuration providers to make
// authentication work seamlessly across environments.
// Order of precedence:
// 1) OCI_CLI_* environment variables (via ontariosystems/oci-cli-env-provider)
// 2) OCI_* environment variables (native SDK env provider)
// 3) Config file providers (OCI_CLI_CONFIG_FILE/PROFILE if set)
// 4) Instance Principals (when running on OCI compute) - only if env vars don't work
// 5) Default config provider (~/.oci/config, TF_VAR_*), as a last resort
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

	// EARLY EXIT: If we have working credentials from env vars or config files, use them
	// and skip Instance Principal (which can be slow when not on OCI compute).
	if len(providers) > 0 {
		if p, err := common.ComposingConfigurationProvider(providers); err == nil {
			// Test if the provider actually has valid credentials by checking TenancyOCID
			if _, err := p.TenancyOCID(); err == nil {
				// Valid credentials found, return early without trying Instance Principal
				return p, nil
			}
		}
	}

	// 4) Instance principals for compute instances (only if env vars/config didn't work)
	if ip, err := newIPProvider(); err == nil {
		providers = append(providers, ip)
	}

	// 5) Always keep a last-resort default config provider at the end
	providers = append(providers, common.DefaultConfigProvider())

	return common.ComposingConfigurationProvider(providers)
}
