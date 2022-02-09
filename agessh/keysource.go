package agessh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"filippo.io/age/armor"
	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops/v3/logging"
)

var log *logrus.Logger

const (
	fileEnv = "SOPS_AGE_SSH_PRIVATE_KEY"
)

func init() {
	log = logging.NewLogger("AGE-SSH")
}

// MasterKey is an age key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	Identity     string // ssh private key
	PublicKey    string // ssh public key
	EncryptedKey string // a sops data key encrypted with age

	parsedRecipient age.Recipient // a parsed age recipient from ssh public key
}

// Encrypt takes a sops data key, encrypts it with age using the ssh public key and stores the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(datakey []byte) error {
	buffer := &bytes.Buffer{}

	if key.parsedRecipient == nil {
		parsedRecipient, err := parseRecipient(key.PublicKey)

		if err != nil {
			log.WithField("recipient", key.parsedRecipient).Error("Encryption failed")
			return err
		}

		key.parsedRecipient = parsedRecipient
	}

	aw := armor.NewWriter(buffer)
	w, err := age.Encrypt(aw, key.parsedRecipient)
	if err != nil {
		return fmt.Errorf("failed to open file for encrypting sops data key with age ssh: %w", err)
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
	ageKeyFilePath, ok := os.LookupEnv(fileEnv)

	if !ok {
		userConfigDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("user config directory could not be determined: %w", err)
		}

		// Check if there is a Ed25519 Key
		ageKeyFilePath = filepath.Join(userConfigDir, ".ssh", "id_ed25519")
		_, err = os.Stat(ageKeyFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				// Fallback to RSA key
				ageKeyFilePath = filepath.Join(userConfigDir, ".ssh", "id_rsa")
			} else {
				return nil, fmt.Errorf("failed to find ssh private key")
			}
		}
	}

	ageKeyFile, err := os.Open(ageKeyFilePath)

	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer ageKeyFile.Close()

	keyBytes, err := io.ReadAll(ageKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read ssh private key")
	}

	identity, err := parseIdentity(keyBytes)
	if err != nil {
		return nil, err
	}

	src := bytes.NewReader([]byte(key.EncryptedKey))
	ar := armor.NewReader(src)
	r, err := age.Decrypt(ar, identity)

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
	return key.PublicKey
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key *MasterKey) ToMap() map[string]interface{} {
	return map[string]interface{}{"recipient": key.PublicKey, "enc": key.EncryptedKey}
}

// MasterKeysFromPublicKeyFiles takes a comma-separated list of files with pem encoded public keys and returns a
// slice of new MasterKeys.
func MasterKeysFromPublicKeyFiles(commaSeparatedRecipients string) ([]*MasterKey, error) {
	if commaSeparatedRecipients == "" {
		// otherwise Split returns [""] and MasterKeyFromFile is unhappy
		return make([]*MasterKey, 0), nil
	}
	recipients := strings.Split(commaSeparatedRecipients, ",")

	var keys []*MasterKey

	for _, recipient := range recipients {
		key, err := MasterKeyFromFile(recipient)

		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// MasterKeyFromFile takes ssh public key and returns a new MasterKey.
func MasterKeyFromFile(filename string) (*MasterKey, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error while opening public key file: %w", err)
	}
	defer f.Close()

	key, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %w", err)
	}

	recipient := string(key)

	parsedRecipient, err := parseRecipient(recipient)

	if err != nil {
		return nil, err
	}

	return &MasterKey{
		PublicKey:       recipient,
		parsedRecipient: parsedRecipient,
	}, nil
}

// parseRecipient attempts to parse a string containing an PEM encoded ssh public key
func parseRecipient(recipient string) (age.Recipient, error) {
	parsedRecipient, err := agessh.ParseRecipient(recipient)

	if err != nil {
		return nil, fmt.Errorf("failed to parse input as PEM encoded ssh public key: %w", err)
	}

	return parsedRecipient, nil
}

// parseIdentity attempts to parse a PEM encoded ssh private key
func parseIdentity(identity []byte) (age.Identity, error) {
	parsedIdentity, err := agessh.ParseIdentity(identity)
	if err != nil {
		return nil, fmt.Errorf("failed to parse identity from PEM encoded ssh private key: %w", err)
	}

	return parsedIdentity, nil
}
