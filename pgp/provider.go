package pgp

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
	return "pgp"
}

func (p *Provider) MarshalKey(key keys.MasterKey) (map[string]any, error) {
	k, ok := key.(*MasterKey)
	if !ok {
		return nil, nil
	}
	return map[string]any{
		"fp":         k.Fingerprint,
		"enc":        k.EncryptedKey,
		"created_at": k.CreationDate.Format(time.RFC3339),
	}, nil
}

func (p *Provider) UnmarshalKey(data map[string]any) (keys.MasterKey, error) {
	var fingerprint, enc, createdAtStr string
	if v, ok := data["fp"].(string); ok {
		fingerprint = v
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
		Fingerprint:  fingerprint,
		EncryptedKey: enc,
		CreationDate: createdAt,
	}, nil
}

func (p *Provider) KeysFromConfig(config any, opts keys.CreationOptions) ([]keys.MasterKey, error) {
	fps, err := keys.ParseStringSlice(config, "pgp")
	if err != nil {
		return nil, err
	}
	var res []keys.MasterKey
	for _, k := range MasterKeysFromFingerprintString(strings.Join(fps, ",")) {
		res = append(res, k)
	}
	return res, nil
}
