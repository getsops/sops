/*
Package keyservice implements a gRPC API that can be used by SOPS to encrypt and decrypt the data key using remote
master keys.
*/
package keyservice

import (
	"errors"
	"fmt"

	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/azkv"
	"github.com/getsops/sops/v3/gcpkms"
	"github.com/getsops/sops/v3/hckms"
	"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/pgp"
)

// ErrUnsupportedMasterKeyType indicates that the provided master key type cannot be converted
// into a keyservice RPC key.
var ErrUnsupportedMasterKeyType = errors.New("unsupported master key type")

// KeyFromMasterKey converts a SOPS internal MasterKey to an RPC Key that can be serialized with Protocol Buffers
func KeyFromMasterKey(mk keys.MasterKey) Key {
	k, _ := KeyFromMasterKeyOrError(mk)
	return k
}

// KeyFromMasterKeyOrError converts a SOPS internal MasterKey to an RPC Key and
// returns an error for unsupported key types.
func KeyFromMasterKeyOrError(mk keys.MasterKey) (Key, error) {
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
	default:
		return Key{}, fmt.Errorf("%w: %T", ErrUnsupportedMasterKeyType, mk)
	}
}
