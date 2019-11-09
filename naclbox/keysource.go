package naclbox //import "go.mozilla.org/sops/naclbox"

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops/logging"
	"golang.org/x/crypto/nacl/box"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("NACLBOX")
}

// MasterKey is a NACL key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	PublicKey       string    // a base64 encoded public key
	pub, priv       *[32]byte // decoded public and private nacl box keys
	EncryptedKey    string
	EphemeralPubKey string // base64 encoded public key used for authentication
	Nonce           string // base64 encoded unique nonce
	CreationDate    time.Time
}

// KeyFile is the file representation of a public and private NACL BOX keypair
type KeyFile struct {
	PublicKey, PrivateKey string
}

// EncryptedDataKey returns the encrypted data key this master key holds
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Encrypt takes a sops data key, encrypts it with NACL BOX
// and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {
	// NACL BOX uses two keys pairs to seal secret data: sender and recipient.
	// When used in Sops, we don't care much about the sender keypair, but it is
	// required for the protocol to work, so we issue an ephemeral keypair instead
	// and store its pubkey alongside the encrypted data key, so it can later
	// be decrypted.
	ephemeralPubKey, ephemeralPrivKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	key.EphemeralPubKey = base64.StdEncoding.EncodeToString(ephemeralPubKey[:])

	// We must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	_, err = io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return err
	}
	key.Nonce = base64.StdEncoding.EncodeToString(nonce[:])

	var pubKey *[32]byte
	p, err := base64.StdEncoding.DecodeString(key.PublicKey)
	if err != nil {
		return err
	}
	pubKey = new([32]byte)
	copy(pubKey[:], p)
	encrypted := box.Seal(nonce[:], dataKey, &nonce, pubKey, ephemeralPrivKey)
	key.EncryptedKey = base64.StdEncoding.EncodeToString(encrypted)
	log.WithField("PublicKey", key.PublicKey).Info("Encryption succeeded")
	return nil
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with NACL BOX private key and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	// All we know at this stage is the public key used to encrypt the data key.
	// We need to find the corresponding private key to perform the decryption.
	// As a convention, we require that the private key is located under
	// $HOME/.sops/naclbox/<sha256(pubkey)>.key
	// where <sha256(pubkey)> is the sha256 hash of the raw public key
	h := sha256.Sum256(key.pub[:])
	path := fmt.Sprintf("%s/.sops/naclbox/%x.key", os.Getenv("HOME"), h[:])
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.WithField("PublicKey", key.PublicKey).Errorf("no private key found at %s", path)
		return nil, fmt.Errorf("no private key found at %s", path)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithField("PublicKey", key.PublicKey).Errorf("failed to read key file from %s: %v", path, err)
		return nil, fmt.Errorf("failed to read key file from %s: %v", path, err)
	}
	var keyFile KeyFile
	err = json.Unmarshal(data, &keyFile)
	if err != nil {
		log.WithField("PublicKey", key.PublicKey).Errorf("failed to parse json key file from %s: %v", path, err)
		return nil, fmt.Errorf("failed to parse json key file from %s: %v", path, err)
	}
	priv, err := base64.StdEncoding.DecodeString(keyFile.PrivateKey)
	if err != nil {
		log.WithField("PublicKey", key.PublicKey).Errorf("failed to decode base64 private key from %s: %v", path, err)
		return nil, fmt.Errorf("failed to decode base64 private key from %s: %v", path, err)
	}
	copy(priv[:32], key.priv[:])

	// Get the nonce value stored in the master key metadata. This was set
	// at encryption and must be reused to decrypt the data key.
	if key.Nonce == "" {
		log.WithField("PublicKey", key.PublicKey).Error("missing nonce value, required for decryption")
		return nil, fmt.Errorf("missing nonce value, required for decryption")
	}
	decodedNonce, err := base64.StdEncoding.DecodeString(key.Nonce)
	if err != nil {
		log.WithField("PublicKey", key.PublicKey).Errorf("failed to decode base64 nonce: %v", err)
		return nil, fmt.Errorf("failed to decode base64 nonce: %v", err)
	}
	var nonce [24]byte
	copy(nonce[:24], decodedNonce[:])

	// Get the ephemeral public key value stored in the master key metadata. This was set
	// at encryption and must be reused to decrypt the data key.
	if key.EphemeralPubKey == "" {
		log.WithField("PublicKey", key.PublicKey).Error("missing ephemeral public key value, required for decryption")
		return nil, fmt.Errorf("missing ephemeral public key value, required for decryption")
	}
	ePubKey, err := base64.StdEncoding.DecodeString(key.EphemeralPubKey)
	if err != nil {
		log.WithField("PublicKey", key.PublicKey).Errorf("failed to decode base64 ephemeral public key: %v", err)
		return nil, fmt.Errorf("failed to decode base64 ephemeral public key: %v", err)
	}
	var ephemeralPubKey *[32]byte
	copy(ephemeralPubKey[:32], ePubKey)

	encrypted, err := base64.StdEncoding.DecodeString(key.EncryptedKey)
	if err != nil {
		log.WithField("PublicKey", key.PublicKey).Errorf("failed to decode base64 encrypted data key: %v", err)
		return nil, fmt.Errorf("failed to decode base64 encrypted data key: %v", err)
	}

	dataKey, ok := box.Open(nil, encrypted[:], &nonce, ephemeralPubKey, key.priv)
	if !ok {
		log.WithField("PublicKey", key.PublicKey).Errorf("decryption failed")
		return nil, fmt.Errorf("decryption failed")
	}
	return dataKey, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return key.PublicKey
}

// MasterKeysFromPublicKeysString takes a list of comma-separated public keys and returns a slice of master keys
func MasterKeysFromPublicKeysString(pubkeys string) []*MasterKey {
	var keys []*MasterKey
	if pubkeys == "" {
		return keys
	}
	for _, s := range strings.Split(pubkeys, ",") {
		keys = append(keys, NewMasterKeyFromPublicKey(s))
	}
	return keys
}

// NewMasterKeyFromPublicKey takes a NACL BOX Public Key and returns a new MasterKey
func NewMasterKeyFromPublicKey(publickey string) *MasterKey {
	if publickey == "" {
		panic("cannot set master key from empty public key")
	}
	mk := &MasterKey{
		PublicKey:    strings.Replace(publickey, " ", "", -1),
		CreationDate: time.Now().UTC(),
	}
	pub, err := base64.StdEncoding.DecodeString(publickey)
	if err != nil {
		panic(err)
	}
	mk.pub = new([32]byte)
	copy(mk.pub[:], pub)
	return mk
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["publickey"] = key.PublicKey
	out["enc"] = key.EncryptedKey
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["nonce"] = key.Nonce
	out["ephemeralpubkey"] = key.EphemeralPubKey
	return out
}
