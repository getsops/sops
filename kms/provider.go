package kms

import (
	"strings"
	"time"

	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/version"
)

func init() {
	keys.RegisterProvider(&Provider{})
}

type Provider struct{}

func (p *Provider) Type() string {
	return "kms"
}

func (p *Provider) DetectTreeBugs(sopsVersion string, keyGroups [][]keys.MasterKey) bool {
	versionCheck, err := version.AIsNewerThanB("3.3.0", sopsVersion)
	if err != nil || !versionCheck {
		return false
	}
	_, _, key := p.getBuggyKMSKey(keyGroups)
	return key != nil
}

func (p *Provider) BugExplanation() string {
	return "Up until version 3.3.0 of sops there was a bug surrounding the " +
		"use of encryption context with AWS KMS." +
		"\nYou can read the full description of the issue here:" +
		"\nhttps://github.com/mozilla/sops/pull/435"
}

func (p *Provider) RecoverDataKey(keyGroups [][]keys.MasterKey, decryptFn func([][]keys.MasterKey) ([]byte, error)) []byte {
	kgndx, kndx, originalKey := p.getBuggyKMSKey(keyGroups)
	if originalKey == nil {
		return nil
	}

	keyToEdit := *originalKey

	encCtxVals := map[string]interface{}{}
	for _, v := range keyToEdit.EncryptionContext {
		encCtxVals[*v] = ""
	}

	var encCtxVariations []map[string]*string
	for ctxVal := range encCtxVals {
		encCtxVariation := map[string]*string{}
		for key := range keyToEdit.EncryptionContext {
			val := ctxVal
			encCtxVariation[key] = &val
		}
		encCtxVariations = append(encCtxVariations, encCtxVariation)
	}

	for _, encCtxVar := range encCtxVariations {
		keyToEdit.EncryptionContext = encCtxVar
		keyGroups[kgndx][kndx] = &keyToEdit
		dataKey, err := decryptFn(keyGroups)
		if err == nil {
			keyGroups[kgndx][kndx] = originalKey
			return dataKey
		}
	}
	return nil
}

func (p *Provider) getBuggyKMSKey(keyGroups [][]keys.MasterKey) (int, int, *MasterKey) {
	for i, kg := range keyGroups {
		for n, k := range kg {
			kmsKey, ok := k.(*MasterKey)
			if ok && len(kmsKey.EncryptionContext) >= 2 {
				duplicateValues := map[string]int{}
				for _, v := range kmsKey.EncryptionContext {
					duplicateValues[*v] = duplicateValues[*v] + 1
				}
				if len(duplicateValues) > 1 {
					return i, n, kmsKey
				}
			}
		}
	}
	return 0, 0, nil
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
	var kmsCtx map[string]*string
	if opts.GlobalConfig != nil {
		if v, ok := opts.GlobalConfig["aws_profile"].(string); ok {
			awsProfile = v
		}
		if v, ok := opts.GlobalConfig["encryption-context"].(string); ok {
			kmsCtx = keys.ParseStringMap(v)
		}
	}
	if kmsCtx == nil {
		kmsCtx = opts.KmsEncryptionContext
	}

	var res []keys.MasterKey
	for _, k := range MasterKeysFromArnString(strings.Join(arns, ","), kmsCtx, awsProfile) {
		res = append(res, k)
	}
	return res, nil
}

func (p *Provider) CLIConfig() []keys.ProviderFlag {
	return []keys.ProviderFlag{
		{
			Name:            "kms, k",
			Usage:           "comma separated list of KMS ARNs",
			EnvVar:          "SOPS_KMS_ARN",
			IsKeyIdentifier: true,
		},
		{
			Name:  "aws-profile",
			Usage: "The AWS profile to use for requests to AWS",
		},
		{
			Name:  "encryption-context",
			Usage: "comma separated list of KMS encryption context key:value pairs",
		},
	}
}

func (p *Provider) MasterKeysFromCLI(c keys.FlagGetter, prefix string) ([]keys.MasterKey, error) {
	var keys []keys.MasterKey
	
	// 'kms' can be requested as 'kms', 'add-kms', 'rm-kms'
	flagName := prefix + "kms"
	if prefix == "" { // for slice or backward compatibility, 'kms' is both the slice and the string global
		// Wait, 'groups add' uses slice.
		slices := c.StringSlice(flagName)
		if len(slices) > 0 {
			arns := strings.Join(slices, ",")
			for _, k := range MasterKeysFromArnString(arns, ParseKMSContext(c.String("encryption-context")), c.String("aws-profile")) {
				keys = append(keys, k)
			}
			return keys, nil
		}
	}
	
	arns := c.String(flagName)
	if arns == "" {
		return keys, nil
	}
	
	for _, k := range MasterKeysFromArnString(arns, ParseKMSContext(c.String("encryption-context")), c.String("aws-profile")) {
		keys = append(keys, k)
	}
	return keys, nil
}
