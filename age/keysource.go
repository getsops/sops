package age

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
)

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
			return err
		}

		key.parsedRecipient = parsedRecipient
	}

	w, err := age.Encrypt(buffer, key.parsedRecipient)

	if err != nil {
		return fmt.Errorf("failed to open file for encrypting sops data key with age: %v", err)
	}

	if _, err := w.Write(datakey); err != nil {
		return fmt.Errorf("failed to encrypt sops data key with age: %v", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close file for encrypting sops data key with age: %v", err)
	}

	key.EncryptedKey = buffer.String()

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
	ageKeyFile, ok := os.LookupEnv("SOPS_AGE_KEY_FILE")

	if !ok {
		userConfigDir, err := os.UserConfigDir()

		if err != nil {
			return nil, fmt.Errorf("user config directory could not be determined: %v", err)
		}

		ageKeyFile = filepath.Join(userConfigDir, ".sops", "age", "keys.txt")
	}

	_, err := os.Stat(ageKeyFile)

	if err != nil {
		return nil, err
	}

	file, err := os.Open(ageKeyFile)

	if err != nil {
		return nil, fmt.Errorf("no key file found at %s: %s", ageKeyFile, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var privateKey string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "AGE-SECRET-KEY") {
			privateKey = line
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning lines in age key file: %v", err)
	}

	if privateKey == "" {
		return nil, fmt.Errorf("no age private key found in file at: %v", ageKeyFile)
	}

	parsedIdentity, err := age.ParseX25519Identity(string(privateKey))

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key as age X25519Identity at %s: %v", ageKeyFile, err)
	}

	buffer := &bytes.Buffer{}
	reader := bytes.NewReader([]byte(key.EncryptedKey))

	r, err := age.Decrypt(reader, parsedIdentity)

	if err != nil {
		return nil, fmt.Errorf("failed to open encrypted data key: %v", err)
	}

	if _, err := io.Copy(buffer, r); err != nil {
		return nil, fmt.Errorf("failed to read encrypted data key: %v", err)
	}

	return buffer.Bytes(), nil
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
	out := make(map[string]interface{})
	out["recipient"] = key.Recipient
	out["enc"] = key.EncryptedKey
	return out
}

// MasterKeysFromRecipients takes a comma-separated list of Bech32-encoded public keys and returns a
// slice of new MasterKeys.
func MasterKeysFromRecipients(commaSeparatedRecipients string) ([]*MasterKey, error) {
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
		return nil, fmt.Errorf("failed to parse input as Bech32-encoded age public key: %v", err)
	}

	return parsedRecipient, nil
}
