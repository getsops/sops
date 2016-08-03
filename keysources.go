package sops

type KeySource interface {
	DecryptKey(encryptedKey string) string
	EncryptKey(key string) string
}

type KMSKeySource struct {
	Role string
	Arn  string
}
type GPGKeySource struct{}

func (kms KMSKeySource) DecryptKey(encryptedKey string) string {
	return encryptedKey
}

func (kms KMSKeySource) EncryptKey(key string) string {
	return key
}

func (gpg GPGKeySource) DecryptKey(encryptedKey string) string {
	return encryptedKey
}

func (gpg GPGKeySource) EncryptKey(key string) string {
	return key
}
