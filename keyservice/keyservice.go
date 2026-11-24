/*
Package keyservice implements a gRPC API that can be used by SOPS to encrypt and decrypt the data key using remote
master keys.
*/
package keyservice

import (
	"encoding/json"
	"fmt"

	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/azkv"
	"github.com/getsops/sops/v3/gcpkms"
	"github.com/getsops/sops/v3/hckms"
	"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/pgp"
	sopsplugin "github.com/getsops/sops/v3/plugin"
)

// KeyFromMasterKey converts a SOPS internal MasterKey to an RPC Key that can be serialized with Protocol Buffers.
// It panics if the key cannot be converted.
func KeyFromMasterKey(mk keys.MasterKey) Key {
	key, err := KeyFromMasterKeyWithError(mk)
	if err != nil {
		panic(fmt.Sprintf("Tried to convert MasterKey type %T to keyservice.Key: %v", mk, err))
	}
	return Key{KeyType: key.KeyType}
}

// KeyFromMasterKeyWithError converts a SOPS internal MasterKey to an RPC Key,
// returning conversion errors instead of panicking.
func KeyFromMasterKeyWithError(mk keys.MasterKey) (Key, error) {
	switch mk := mk.(type) {
	case *pgp.MasterKey:
		return Key{
			KeyType: &Key_PgpKey{
				PgpKey: &PgpKey{
					Fingerprint: mk.Fingerprint,
				},
			},
		}, nil
	case *gcpkms.MasterKey:
		return Key{
			KeyType: &Key_GcpKmsKey{
				GcpKmsKey: &GcpKmsKey{
					ResourceId: mk.ResourceID,
				},
			},
		}, nil
	case *hcvault.MasterKey:
		return Key{
			KeyType: &Key_VaultKey{
				VaultKey: &VaultKey{
					VaultAddress: mk.VaultAddress,
					EnginePath:   mk.EnginePath,
					KeyName:      mk.KeyName,
				},
			},
		}, nil
	case *kms.MasterKey:
		ctx := make(map[string]string)
		for k, v := range mk.EncryptionContext {
			ctx[k] = *v
		}
		return Key{
			KeyType: &Key_KmsKey{
				KmsKey: &KmsKey{
					Arn:        mk.Arn,
					Role:       mk.Role,
					Context:    ctx,
					AwsProfile: mk.AwsProfile,
				},
			},
		}, nil
	case *azkv.MasterKey:
		return Key{
			KeyType: &Key_AzureKeyvaultKey{
				AzureKeyvaultKey: &AzureKeyVaultKey{
					VaultUrl: mk.VaultURL,
					Name:     mk.Name,
					Version:  mk.Version,
				},
			},
		}, nil
	case *age.MasterKey:
		return Key{
			KeyType: &Key_AgeKey{
				AgeKey: &AgeKey{
					Recipient: mk.Recipient,
				},
			},
		}, nil
	case *hckms.MasterKey:
		return Key{
			KeyType: &Key_HckmsKey{
				HckmsKey: &HckmsKey{
					KeyId: mk.KeyID,
				},
			},
		}, nil
	case *sopsplugin.MasterKey:
		config, err := marshalPluginObject(mk.Config)
		if err != nil {
			return Key{}, fmt.Errorf("marshal plugin %q configuration: %w", mk.Name, err)
		}
		metadata, err := marshalPluginObject(mk.ProviderMetadata)
		if err != nil {
			return Key{}, fmt.Errorf("marshal plugin %q metadata: %w", mk.Name, err)
		}
		return Key{
			KeyType: &Key_PluginKey{
				PluginKey: &PluginKey{
					Name:         mk.Name,
					ConfigJson:   config,
					MetadataJson: metadata,
				},
			},
		}, nil
	default:
		return Key{}, fmt.Errorf("cannot convert unknown MasterKey type %T to keyservice.Key", mk)
	}
}

func marshalPluginObject(value map[string]interface{}) ([]byte, error) {
	if value == nil {
		value = map[string]interface{}{}
	}
	return json.Marshal(value)
}
