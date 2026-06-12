package age

import (
	"strings"

	"github.com/getsops/sops/v3/keys"
)

func init() {
	keys.RegisterProvider(&Provider{})
}

type Provider struct{}

func (p *Provider) Type() string {
	return "age"
}

func (p *Provider) MarshalKey(key keys.MasterKey) (map[string]any, error) {
	k, ok := key.(*MasterKey)
	if !ok {
		return nil, nil
	}
	return map[string]any{
		"recipient": k.Recipient,
		"enc":       k.EncryptedKey,
	}, nil
}

func (p *Provider) UnmarshalKey(data map[string]any) (keys.MasterKey, error) {
	var recipient, enc string
	if v, ok := data["recipient"].(string); ok {
		recipient = v
	}
	if v, ok := data["enc"].(string); ok {
		enc = v
	}

	return &MasterKey{
		Recipient:    recipient,
		EncryptedKey: enc,
	}, nil
}

func (p *Provider) KeysFromConfig(config any, opts keys.CreationOptions) ([]keys.MasterKey, error) {
	recipients, err := keys.ParseStringSlice(config, "age")
	if err != nil {
		return nil, err
	}
	if len(recipients) == 0 {
		return nil, nil
	}

	ageKeys, err := MasterKeysFromRecipients(strings.Join(recipients, ","))
	if err != nil {
		return nil, err
	}

	var res []keys.MasterKey
	for _, ak := range ageKeys {
		res = append(res, ak)
	}
	return res, nil
}

func (p *Provider) CLIConfig() []keys.ProviderFlag {
	return []keys.ProviderFlag{
		{
			Name:            "age, a",
			Usage:           "comma separated list of age recipients",
			EnvVar:          "SOPS_AGE_RECIPIENTS",
			IsKeyIdentifier: true,
		},
	}
}

func (p *Provider) MasterKeysFromCLI(c keys.FlagGetter, prefix string) ([]keys.MasterKey, error) {
	var masterKeys []keys.MasterKey
	flagName := prefix + "age"

	if prefix == "" {
		slices := c.StringSlice(flagName)
		if len(slices) > 0 {
			recs := strings.Join(slices, ",")
			ageKeys, err := MasterKeysFromRecipients(recs)
			if err != nil {
				return nil, err
			}
			for _, k := range ageKeys {
				masterKeys = append(masterKeys, k)
			}
			return masterKeys, nil
		}
	}

	recs := c.String(flagName)
	if recs == "" {
		return masterKeys, nil
	}
	ageKeys, err := MasterKeysFromRecipients(recs)
	if err != nil {
		return nil, err
	}
	for _, k := range ageKeys {
		masterKeys = append(masterKeys, k)
	}
	return masterKeys, nil
}
