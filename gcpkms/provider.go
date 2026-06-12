package gcpkms

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
	return "gcp_kms"
}

func (p *Provider) MarshalKey(key keys.MasterKey) (map[string]any, error) {
	k, ok := key.(*MasterKey)
	if !ok {
		return nil, nil
	}
	return map[string]any{
		"resource_id": k.ResourceID,
		"enc":         k.EncryptedKey,
		"created_at":  k.CreationDate.Format(time.RFC3339),
	}, nil
}

func (p *Provider) UnmarshalKey(data map[string]any) (keys.MasterKey, error) {
	var resourceID, enc, createdAtStr string
	if v, ok := data["resource_id"].(string); ok {
		resourceID = v
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
		ResourceID:   resourceID,
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
				var resourceID string
				if v, ok := m["resource_id"].(string); ok {
					resourceID = v
				}
				for _, k := range MasterKeysFromResourceIDString(resourceID) {
					res = append(res, k)
				}
			}
			return res, nil
		}
	}

	ids, err := keys.ParseStringSlice(config, "gcp_kms")
	if err != nil {
		return nil, err
	}
	var res []keys.MasterKey
	for _, k := range MasterKeysFromResourceIDString(strings.Join(ids, ",")) {
		res = append(res, k)
	}
	return res, nil
}

func (p *Provider) CLIConfig() []keys.ProviderFlag {
	return []keys.ProviderFlag{
		{
			Name:            "gcp-kms",
			Usage:           "comma separated list of GCP KMS resource IDs",
			EnvVar:          "SOPS_GCP_KMS_IDS",
			IsKeyIdentifier: true,
		},
	}
}

func (p *Provider) MasterKeysFromCLI(c keys.FlagGetter, prefix string) ([]keys.MasterKey, error) {
	var masterKeys []keys.MasterKey
	flagName := prefix + "gcp-kms"
	
	if prefix == "" {
		slices := c.StringSlice(flagName)
		if len(slices) > 0 {
			ids := strings.Join(slices, ",")
			for _, k := range MasterKeysFromResourceIDString(ids) {
				masterKeys = append(masterKeys, k)
			}
			return masterKeys, nil
		}
	}

	ids := c.String(flagName)
	if ids == "" {
		return masterKeys, nil
	}
	for _, k := range MasterKeysFromResourceIDString(ids) {
		masterKeys = append(masterKeys, k)
	}
	return masterKeys, nil
}

