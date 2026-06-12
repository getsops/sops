package hcvault

import (
	"strings"
	"time"

	"github.com/getsops/sops/v3/keys"
)

func init() {
	keys.RegisterProvider(&Provider{})
}

type Provider struct{}

func (p *Provider) Type() string {
	return "hc_vault"
}

func (p *Provider) MarshalKey(key keys.MasterKey) (map[string]any, error) {
	k, ok := key.(*MasterKey)
	if !ok {
		return nil, nil
	}
	return map[string]any{
		"vault_address": k.VaultAddress,
		"engine_path":   k.EnginePath,
		"key_name":      k.KeyName,
		"enc":           k.EncryptedKey,
		"created_at":    k.CreationDate.Format(time.RFC3339),
	}, nil
}

func (p *Provider) UnmarshalKey(data map[string]any) (keys.MasterKey, error) {
	var vaultAddress, enginePath, keyName, enc, createdAtStr string
	if v, ok := data["vault_address"].(string); ok {
		vaultAddress = v
	}
	if v, ok := data["engine_path"].(string); ok {
		enginePath = v
	}
	if v, ok := data["key_name"].(string); ok {
		keyName = v
	}
	if v, ok := data["enc"].(string); ok {
		enc = v
	}
	if v, ok := data["created_at"].(string); ok {
		createdAtStr = v
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}
	return &MasterKey{
		VaultAddress: vaultAddress,
		EnginePath:   enginePath,
		KeyName:      keyName,
		EncryptedKey: enc,
		CreationDate: createdAt,
	}, nil
}

func (p *Provider) KeysFromConfig(config any, opts keys.CreationOptions) ([]keys.MasterKey, error) {
	var uris []string
	if config != nil {
		u, err := keys.ParseStringSlice(config, "hc_vault")
		if err != nil {
			return nil, err
		}
		uris = u
	}

	if len(uris) == 0 {
		return nil, nil
	}

	vaultKeys, err := NewMasterKeysFromURIs(strings.Join(uris, ","))
	if err != nil {
		return nil, err
	}
	var res []keys.MasterKey
	for _, k := range vaultKeys {
		res = append(res, k)
	}
	return res, nil
}

func (p *Provider) CLIConfig() []keys.ProviderFlag {
	return []keys.ProviderFlag{
		{
			Name:            "hc-vault-transit",
			Usage:           "comma separated list of Vault's URI keys (e.g. 'https://vault.example.org:8200/v1/transit/keys/dev')",
			EnvVar:          "SOPS_VAULT_URIS",
			IsKeyIdentifier: true,
		},
	}
}

func (p *Provider) MasterKeysFromCLI(c keys.FlagGetter, prefix string) ([]keys.MasterKey, error) {
	var masterKeys []keys.MasterKey
	flagName := prefix + "hc-vault-transit"
	
	if prefix == "" {
		slices := c.StringSlice(flagName)
		if len(slices) > 0 {
			uris := strings.Join(slices, ",")
			hcVaultKeys, err := NewMasterKeysFromURIs(uris)
			if err != nil {
				return nil, err
			}
			for _, k := range hcVaultKeys {
				masterKeys = append(masterKeys, k)
			}
			return masterKeys, nil
		}
	}

	uris := c.String(flagName)
	if uris == "" {
		return masterKeys, nil
	}
	hcVaultKeys, err := NewMasterKeysFromURIs(uris)
	if err != nil {
		return nil, err
	}
	for _, k := range hcVaultKeys {
		masterKeys = append(masterKeys, k)
	}
	return masterKeys, nil
}

