package age

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
)

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

		ageKeyFile = filepath.Join(userConfigDir, "sops", "age", "keys.txt")
	}

	identities, err := parseIdentitiesFile(ageKeyFile)

	if err != nil {
		return nil, err
	}

	var buffer *bytes.Buffer

	for _, identity := range identities {
		buffer = &bytes.Buffer{}
		reader := bytes.NewReader([]byte(key.EncryptedKey))

		r, err := age.Decrypt(reader, identity)

		if err != nil {
			continue
		}

		if _, err := io.Copy(buffer, r); err != nil {
			continue
		}

		return buffer.Bytes(), nil
	}

	return nil, fmt.Errorf("no age identity found in %q that could decrypt the data", ageKeyFile)
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

// parseIdentitiesFile parses a file containing age private keys. Derived from
// https://github.com/FiloSottile/age/blob/189041b668629795593766bcb8d3f70ee248b842/cmd/age/parse.go
// but should be replaced with a library function if a future version of the age library exposes
// this functionality.
func parseIdentitiesFile(name string) ([]age.Identity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(io.LimitReader(f, privateKeySizeLimit))
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %v", name, err)
	}
	if len(contents) == privateKeySizeLimit {
		return nil, fmt.Errorf("failed to read %q: file too long", name)
	}

	var ids []age.Identity
	var ageParsingError error
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if strings.HasPrefix(line, "-----BEGIN") {
			ageParsingError = fmt.Errorf("sops does not yet support SSH keys via age. SSH key found in file at %q", name)
			continue
		}
		if ageParsingError != nil {
			continue
		}
		i, err := age.ParseX25519Identity(line)
		if err != nil {
			ageParsingError = fmt.Errorf("malformed secret keys file %q: %v", name, err)
			continue
		}
		ids = append(ids, i)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read %q: %v", name, err)
	}
	if ageParsingError != nil {
		return nil, ageParsingError
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("no secret keys found in %q", name)
	}
	return ids, nil
}
