package kms

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
	return "kms"
}

func (p *Provider) MarshalKey(key keys.MasterKey) (map[string]any, error) {
	k, ok := key.(*MasterKey)
	if !ok {
		return nil, nil
	}
	return map[string]any{
		"arn":         k.Arn,
		"role":        k.Role,
		"context":     k.EncryptionContext,
		"created_at":  k.CreationDate.Format(time.RFC3339),
		"enc":         k.EncryptedKey,
		"aws_profile": k.AwsProfile,
	}, nil
}

func (p *Provider) UnmarshalKey(data map[string]any) (keys.MasterKey, error) {
	var arn, role, enc, awsProfile, createdAtStr string
	if v, ok := data["arn"].(string); ok {
		arn = v
	}
	if v, ok := data["role"].(string); ok {
		role = v
	}
	if v, ok := data["enc"].(string); ok {
		enc = v
	}
	if v, ok := data["aws_profile"].(string); ok {
		awsProfile = v
	}
	if v, ok := data["created_at"].(string); ok {
		createdAtStr = v
	}

	var context map[string]*string
	if v, ok := data["context"].(map[string]any); ok {
		context = make(map[string]*string)
		for key, val := range v {
			if s, ok := val.(string); ok {
				context[key] = &s
			}
		}
	} else if v, ok := data["context"].(map[string]*string); ok {
		context = v
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}
	return &MasterKey{
		Arn:               arn,
		Role:              role,
		EncryptionContext: context,
		CreationDate:      createdAt,
		EncryptedKey:      enc,
		AwsProfile:        awsProfile,
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
				var arn, role, awsProfile string
				if v, ok := m["arn"].(string); ok {
					arn = v
				}
				if v, ok := m["role"].(string); ok {
					role = v
				}
				if v, ok := m["aws_profile"].(string); ok {
					awsProfile = v
				}
				var context map[string]*string
				if v, ok := m["context"].(map[string]interface{}); ok {
					context = make(map[string]*string)
					for k, val := range v {
						if s, ok := val.(string); ok {
							context[k] = &s
						}
					}
				}
				res = append(res, NewMasterKeyWithProfile(arn, role, context, awsProfile))
			}
			return res, nil
		}
	}

	arns, err := keys.ParseStringSlice(config, "kms")
	if err != nil {
		return nil, err
	}
	if len(arns) == 0 {
		return nil, nil
	}

	var awsProfile string
	if opts.GlobalConfig != nil {
		if v, ok := opts.GlobalConfig["aws_profile"].(string); ok {
			awsProfile = v
		}
	}

	var res []keys.MasterKey
	for _, k := range MasterKeysFromArnString(strings.Join(arns, ","), opts.KmsEncryptionContext, awsProfile) {
		res = append(res, k)
	}
	return res, nil
}
