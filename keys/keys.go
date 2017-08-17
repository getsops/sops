package keys

// MasterKey provides a way of securing the key used to encrypt the Tree by encrypting and decrypting said key.
type MasterKey interface {
	Encrypt(dataKey []byte) error
	EncryptIfNeeded(dataKey []byte) error
	EncryptedDataKey() []byte
	SetEncryptedDataKey([]byte)
	Decrypt() ([]byte, error)
	NeedsRotation() bool
	ToString() string
	ToMap() map[string]interface{}
}
