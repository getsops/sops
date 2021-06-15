package age

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops/v3/logging"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("AGE")
}

const privateKeySizeLimit = 1 << 24 // 16 MiB

// MasterKey is an age key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	Identity     string // a Bech32-encoded private key
	Recipient    string // a Bech32-encoded public key
	EncryptedKey string // a sops data key encrypted with age

	parsedIdentity  *age.X25519Identity  // a parsed age private key
	parsedRecipient *age.X25519Recipient // a parsed age public key
}

// Encrypt takes a sops data key, encrypts it with age and stores the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(datakey []byte) error {
	buffer := &bytes.Buffer{}

	if key.parsedRecipient == nil {
		parsedRecipient, err := parseRecipient(key.Recipient)

		if err != nil {
			log.WithField("recipient", key.parsedRecipient).Error("Encryption failed")
			return err
		}

		key.parsedRecipient = parsedRecipient
	}

	aw := armor.NewWriter(buffer)
	w, err := age.Encrypt(aw, key.parsedRecipient)
	if err != nil {
		return fmt.Errorf("failed to open file for encrypting sops data key with age: %w", err)
	}

	if _, err := w.Write(datakey); err != nil {
		log.WithField("recipient", key.parsedRecipient).Error("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with age: %w", err)
	}

	if err := w.Close(); err != nil {
		log.WithField("recipient", key.parsedRecipient).Error("Encryption failed")
		return fmt.Errorf("failed to close file for encrypting sops data key with age: %w", err)
	}

	if err := aw.Close(); err != nil {
		log.WithField("recipient", key.parsedRecipient).Error("Encryption failed")
		return fmt.Errorf("failed to close armored writer: %w", err)
	}

	key.EncryptedKey = buffer.String()

	log.WithField("recipient", key.parsedRecipient).Info("Encryption succeeded")

	return nil
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet.
func (key *MasterKey) EncryptIfNeeded(datakey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(datakey)
	}

	return nil
}

// EncryptedDataKey returns the encrypted data key this master key holds.
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Decrypt decrypts the EncryptedKey field with the age identity and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	ageKeyFilePath, ok := os.LookupEnv("SOPS_AGE_KEY_FILE")

	if !ok {
		userConfigDir, err := os.UserConfigDir()

		if err != nil {
			return nil, fmt.Errorf("user config directory could not be determined: %w", err)
		}

		ageKeyFilePath = filepath.Join(userConfigDir, "sops", "age", "keys.txt")
	}

	ageKeyFile, err := os.Open(ageKeyFilePath)

	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer ageKeyFile.Close()

	identities, err := age.ParseIdentities(ageKeyFile)

	if err != nil {
		return nil, err
	}

	src := bytes.NewReader([]byte(key.EncryptedKey))
	ar := armor.NewReader(src)
	r, err := age.Decrypt(ar, identities...)

	if err != nil {
		return nil, fmt.Errorf("no age identity found in %q that could decrypt the data", ageKeyFilePath)
	}

	var b bytes.Buffer
	if _, err := io.Copy(&b, r); err != nil {
		return nil, fmt.Errorf("failed to copy decrypted data into bytes.Buffer: %w", err)
	}

	return b.Bytes(), nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return false
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return key.Recipient
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key *MasterKey) ToMap() map[string]interface{} {
	return map[string]interface{}{"recipient": key.Recipient, "enc": key.EncryptedKey}
}

// MasterKeysFromRecipients takes a comma-separated list of Bech32-encoded public keys and returns a
// slice of new MasterKeys.
func MasterKeysFromRecipients(commaSeparatedRecipients string) ([]*MasterKey, error) {
	if commaSeparatedRecipients == "" {
		// otherwise Split returns [""] and MasterKeyFromRecipient is unhappy
		return make([]*MasterKey, 0), nil
	}
	recipients := strings.Split(commaSeparatedRecipients, ",")

	var keys []*MasterKey

	for _, recipient := range recipients {
		key, err := MasterKeyFromRecipient(recipient)

		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// MasterKeyFromRecipient takes a Bech32-encoded public key and returns a new MasterKey.
func MasterKeyFromRecipient(recipient string) (*MasterKey, error) {
	recipient = strings.TrimSpace(recipient)
	parsedRecipient, err := parseRecipient(recipient)

	if err != nil {
		return nil, err
	}

	return &MasterKey{
		Recipient:       recipient,
		parsedRecipient: parsedRecipient,
	}, nil
}

// parseRecipient attempts to parse a string containing an encoded age public key
func parseRecipient(recipient string) (*age.X25519Recipient, error) {
	parsedRecipient, err := age.ParseX25519Recipient(recipient)

	if err != nil {
		return nil, fmt.Errorf("failed to parse input as Bech32-encoded age public key: %w", err)
	}

	return parsedRecipient, nil
}
