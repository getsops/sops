/*
Package keyservice implements a gRPC API that can be used by SOPS to encrypt and decrypt the data key using remote
master keys.
*/
package keyservice

import (
	"fmt"

	"go.mozilla.org/sops/v3/azkv"
	"go.mozilla.org/sops/v3/gcpkms"
	"go.mozilla.org/sops/v3/hcvault"
	"go.mozilla.org/sops/v3/keys"
	"go.mozilla.org/sops/v3/kms"
	"go.mozilla.org/sops/v3/pgp"
)

// KeyFromMasterKey converts a SOPS internal MasterKey to an RPC Key that can be serialized with Protocol Buffers
func KeyFromMasterKey(mk keys.MasterKey) Key {
	switch mk := mk.(type) {
	case *pgp.MasterKey:
		return Key{
			KeyType: &Key_PgpKey{
				PgpKey: &PgpKey{
					Fingerprint: mk.Fingerprint,
				},
			},
		}
	case *gcpkms.MasterKey:
		return Key{
			KeyType: &Key_GcpKmsKey{
				GcpKmsKey: &GcpKmsKey{
					ResourceId: mk.ResourceID,
				},
			},
		}
	case *hcvault.MasterKey:
		return Key{
			KeyType: &Key_VaultKey{
				VaultKey: &VaultKey{
					VaultAddress: mk.VaultAddress,
					EnginePath:    mk.EnginePath,
					KeyName:        mk.KeyName,
				},
			},
		}
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
		}
	case *azkv.MasterKey:
		return Key{
			KeyType: &Key_AzureKeyvaultKey{
				AzureKeyvaultKey: &AzureKeyVaultKey{
					VaultUrl: mk.VaultURL,
					Name:     mk.Name,
					Version:  mk.Version,
				},
			},
		}
	default:
		panic(fmt.Sprintf("Tried to convert unknown MasterKey type %T to keyservice.Key", mk))
	}
}
