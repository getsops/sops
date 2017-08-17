package keyservice

import (
	"fmt"

	"go.mozilla.org/sops/keys"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
)

func KeyFromMasterKey(mk keys.MasterKey) Key {
	switch mk := mk.(type) {
	case *pgp.MasterKey:
		return Key{
			KeyType: &Key_GpgKey{
				GpgKey: &GpgKey{
					Fingerprint: mk.Fingerprint,
				},
			},
		}
	case *kms.MasterKey:
		var ctx map[string]string
		for k, v := range mk.EncryptionContext {
			ctx[k] = *v
		}
		return Key{
			KeyType: &Key_KmsKey{
				KmsKey: &KmsKey{
					Arn:     mk.Arn,
					Role:    mk.Role,
					Context: ctx,
				},
			},
		}
	default:
		panic(fmt.Sprintf("Tried to convert unknown MasterKey type %T to keyservice.Key", mk))
	}
}
