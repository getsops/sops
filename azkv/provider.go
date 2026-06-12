package azkv

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
	return "azure_kv"
}

func (p *Provider) MarshalKey(key keys.MasterKey) (map[string]any, error) {
	k, ok := key.(*MasterKey)
	if !ok {
		return nil, nil
	}
	return map[string]any{
		"vault_url":  k.VaultURL,
		"name":       k.Name,
		"version":    k.Version,
		"enc":        k.EncryptedKey,
		"created_at": k.CreationDate.Format(time.RFC3339),
	}, nil
}

func (p *Provider) UnmarshalKey(data map[string]any) (keys.MasterKey, error) {
	var vaultURL, name, version, enc, createdAtStr string
	if v, ok := data["vault_url"].(string); ok {
		vaultURL = v
	}
	if v, ok := data["name"].(string); ok {
		name = v
	}
	if v, ok := data["version"].(string); ok {
		version = v
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
		VaultURL:     vaultURL,
		Name:         name,
		Version:      version,
		EncryptedKey: enc,
		CreationDate: createdAt,
	}, nil
}

func (p *Provider) KeysFromConfig(config any, opts keys.CreationOptions) ([]keys.MasterKey, error) {
	if maps, ok := config.([]interface{}); ok {
		var isMap bool
		if len(maps) > 0 {
			_, isMap = maps[0].(map[string]interface{})
		}
		if isMap {
			var res []keys.MasterKey
			for _, item := range maps {
				m := item.(map[string]interface{})
				var vaultURL, key, version string
				if v, ok := m["vaultUrl"].(string); ok {
					vaultURL = v
				}
				if v, ok := m["vault_url"].(string); ok {
					vaultURL = v
				}
				if v, ok := m["key"].(string); ok {
					key = v
				}
				if v, ok := m["version"].(string); ok {
					version = v
				}

				azureKey, err := NewMasterKeyWithOptionalVersion(vaultURL, key, version)
				if err != nil {
					return nil, err
				}
				res = append(res, azureKey)
			}
			return res, nil
		}
	}

	urls, err := keys.ParseStringSlice(config, "azure_kv")
	if err != nil {
		return nil, err
	}
	if len(urls) == 0 {
		return nil, nil
	}

	azureKeys, err := MasterKeysFromURLs(strings.Join(urls, ","))
	if err != nil {
		return nil, err
	}

	var res []keys.MasterKey
	for _, k := range azureKeys {
		res = append(res, k)
	}
	return res, nil
}

func (p *Provider) CLIConfig() []keys.ProviderFlag {
	return []keys.ProviderFlag{
		{
			Name:            "azure-kv",
			Usage:           "comma separated list of Azure Key Vault URLs",
			EnvVar:          "SOPS_AZURE_KEYVAULT_URL",
			IsKeyIdentifier: true,
		},
	}
}

func (p *Provider) MasterKeysFromCLI(c keys.FlagGetter, prefix string) ([]keys.MasterKey, error) {
	var masterKeys []keys.MasterKey
	flagName := prefix + "azure-kv"

	if prefix == "" {
		slices := c.StringSlice(flagName)
		if len(slices) > 0 {
			urls := strings.Join(slices, ",")
			azureKeys, err := MasterKeysFromURLs(urls)
			if err != nil {
				return nil, err
			}
			for _, k := range azureKeys {
				masterKeys = append(masterKeys, k)
			}
			return masterKeys, nil
		}
	}

	urls := c.String(flagName)
	if urls == "" {
		return masterKeys, nil
	}
	azureKeys, err := MasterKeysFromURLs(urls)
	if err != nil {
		return nil, err
	}
	for _, k := range azureKeys {
		masterKeys = append(masterKeys, k)
	}
	return masterKeys, nil
}
