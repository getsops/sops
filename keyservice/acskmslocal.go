package keyservice

// AcsKmsKey holds the ARN of an Alibaba Cloud KMS key.
// This type is not protobuf-generated; it is handled only in the local keyservice.
type AcsKmsKey struct {
	Arn          string
	EncryptedKey string
}

// Key_AcsKmsKey is the oneof wrapper for AcsKmsKey inside a Key.
type Key_AcsKmsKey struct {
	AcsKmsKey *AcsKmsKey
}

func (*Key_AcsKmsKey) isKey_KeyType() {}

// GetAcsKmsKey returns the AcsKmsKey if the Key holds one.
func (k *Key) GetAcsKmsKey() *AcsKmsKey {
	if x, ok := k.GetKeyType().(*Key_AcsKmsKey); ok {
		return x.AcsKmsKey
	}
	return nil
}
