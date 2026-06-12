package hckms

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
	return "hckms"
}

func (p *Provider) MarshalKey(key keys.MasterKey) (map[string]any, error) {
	k, ok := key.(*MasterKey)
	if !ok {
		return nil, nil
	}
	return map[string]any{
		"key_id":     k.KeyID,
		"enc":        k.EncryptedKey,
		"created_at": k.CreationDate.Format(time.RFC3339),
	}, nil
}

func (p *Provider) UnmarshalKey(data map[string]any) (keys.MasterKey, error) {
	var keyID, enc, createdAtStr string
	if v, ok := data["key_id"].(string); ok {
		keyID = v
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
	key, err := NewMasterKey(keyID)
	if err != nil {
		return nil, err
	}
	key.EncryptedKey = enc
	key.CreationDate = createdAt
	return key, nil
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
				var keyID string
				if v, ok := m["key_id"].(string); ok {
					keyID = v
				}
				keysList, err := NewMasterKeyFromKeyIDString(keyID)
				if err != nil {
					return nil, err
				}
				for _, k := range keysList {
					res = append(res, k)
				}
			}
			return res, nil
		}
	}

	ids, err := keys.ParseStringSlice(config, "hckms")
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	hckmsMasterKeys, err := NewMasterKeyFromKeyIDString(strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	var res []keys.MasterKey
	for _, k := range hckmsMasterKeys {
		res = append(res, k)
	}
	return res, nil
}

func (p *Provider) CLIConfig() []keys.ProviderFlag {
	return []keys.ProviderFlag{
		{
			Name:            "hckms",
			Usage:           "comma separated list of HuaweiCloud KMS key IDs (format: region:key-uuid)",
			EnvVar:          "SOPS_HUAWEICLOUD_KMS_IDS",
			IsKeyIdentifier: true,
		},
	}
}

func (p *Provider) MasterKeysFromCLI(c keys.FlagGetter, prefix string) ([]keys.MasterKey, error) {
	var masterKeys []keys.MasterKey
	flagName := prefix + "hckms"
	
	if prefix == "" {
		slices := c.StringSlice(flagName)
		if len(slices) > 0 {
			ids := strings.Join(slices, ",")
			hckmsKeys, err := NewMasterKeyFromKeyIDString(ids)
			if err != nil {
				return nil, err
			}
			for _, k := range hckmsKeys {
				masterKeys = append(masterKeys, k)
			}
			return masterKeys, nil
		}
	}

	ids := c.String(flagName)
	if ids == "" {
		return masterKeys, nil
	}
	hckmsKeys, err := NewMasterKeyFromKeyIDString(ids)
	if err != nil {
		return nil, err
	}
	for _, k := range hckmsKeys {
		masterKeys = append(masterKeys, k)
	}
	return masterKeys, nil
}

