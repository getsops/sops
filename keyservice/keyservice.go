/*
Package keyservice implements a gRPC API that can be used by SOPS to encrypt and decrypt the data key using remote
master keys.
*/
package keyservice

import (
	"fmt"

	"go.mozilla.org/sops/azkv"
	"go.mozilla.org/sops/gcpkms"
	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
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
