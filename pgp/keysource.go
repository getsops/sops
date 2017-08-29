package pgp //import "go.mozilla.org/sops/pgp"

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"log"

	"os/exec"

	"github.com/howeyc/gopass"
	gpgagent "go.mozilla.org/gopgagent"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// MasterKey is a PGP key used to securely store sops' data key by encrypting it and decrypting it
type MasterKey struct {
	Fingerprint  string
	EncryptedKey string
	CreationDate time.Time
}

func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

func gpgBinary() string {
	binary := "gpg"
	if envBinary := os.Getenv("SOPS_GPG_EXEC"); envBinary != "" {
		binary = envBinary
	}
	return binary
}

func (key *MasterKey) encryptWithGPGBinary(dataKey []byte) error {
	args := []string{
		"--no-default-recipient",
		"--yes",
		"--encrypt",
		"-a",
		"-r",
		key.Fingerprint,
		"--trusted-key",
		key.Fingerprint[len(key.Fingerprint)-16:],
		"--no-encrypt-to",
	}
	cmd := exec.Command(gpgBinary(), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(dataKey)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	key.EncryptedKey = stdout.String()
	return nil
}

func (key *MasterKey) encryptWithCryptoOpenPGP(dataKey []byte) error {
	ring, err := key.pubRing()
	if err != nil {
		return err
	}
	fingerprints := key.fingerprintMap(ring)
	entity, ok := fingerprints[key.Fingerprint]
	if !ok {
		return fmt.Errorf("Key with fingerprint %s is not available in keyring.", key.Fingerprint)
	}
	encbuf := new(bytes.Buffer)
	armorbuf, err := armor.Encode(encbuf, "PGP MESSAGE", nil)
	if err != nil {
		return err
	}
	plaintextbuf, err := openpgp.Encrypt(armorbuf, []*openpgp.Entity{&entity}, nil, nil, nil)
	if err != nil {
		return err
	}
	_, err = plaintextbuf.Write(dataKey)
	if err != nil {
		return err
	}
	err = plaintextbuf.Close()
	if err != nil {
		return err
	}
	err = armorbuf.Close()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(encbuf)
	if err != nil {
		return err
	}
	key.EncryptedKey = string(bytes)
	return nil
}

// Encrypt encrypts the data key with the PGP key with the same fingerprint as the MasterKey. It looks for PGP public keys in $PGPHOME/pubring.gpg.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	log.Printf("Attempting encryption of GPG MasterKey with fingerprint %s", key.Fingerprint)
	openpgpErr := key.encryptWithCryptoOpenPGP(dataKey)
	if openpgpErr == nil {
		log.Printf("Encryption of GPG MasterKey with fingerprint %s succeeded", key.Fingerprint)
		return nil
	}
	log.Print("Encryption with golang's openpgp package failed, falling back to the GPG binary")
	binaryErr := key.encryptWithGPGBinary(dataKey)
	if binaryErr == nil {
		log.Printf("Encryption of GPG MasterKey with fingerprint %s succeeded", key.Fingerprint)
		return nil
	}
	log.Printf("Encryption of GPG MasterKey with fingerprint %s failed", key.Fingerprint)
	return fmt.Errorf(`could not encrypt data key with PGP key.
	\tgolang.org/x/crypto/openpgp error: %s
	\tGPG binary error: %s`, openpgpErr, binaryErr)
}

// EncryptIfNeeded encrypts the data key with PGP only if it's needed, that is, if it hasn't been encrypted already
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

func (key *MasterKey) decryptWithGPGBinary() ([]byte, error) {
	args := []string{
		"--use-agent",
		"-d",
	}
	cmd := exec.Command(gpgBinary(), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(key.EncryptedKey)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}

func (key *MasterKey) decryptWithCryptoOpenpgp() ([]byte, error) {
	ring, err := key.secRing()
	if err != nil {
		return nil, fmt.Errorf("Could not load secring: %s", err)
	}
	block, err := armor.Decode(strings.NewReader(key.EncryptedKey))
	if err != nil {
		return nil, fmt.Errorf("Armor decoding failed: %s", err)
	}
	md, err := openpgp.ReadMessage(block.Body, ring, key.passphrasePrompt, nil)
	if err != nil {
		return nil, fmt.Errorf("Reading PGP message failed: %s", err)
	}
	if b, err := ioutil.ReadAll(md.UnverifiedBody); err == nil {
		log.Printf("Decryption of GPG MasterKey with fingerprint %s successful", key.Fingerprint)
		return b, nil
	}
	return nil, fmt.Errorf("The key could not be decrypted with any of the GPG entries")
}

// Decrypt uses PGP to obtain the data key from the EncryptedKey store in the MasterKey and returns it
func (key *MasterKey) Decrypt() ([]byte, error) {
	log.Printf("Attempting decryption of GPG MasterKey with fingerprint %s", key.Fingerprint)
	dataKey, openpgpErr := key.decryptWithCryptoOpenpgp()
	if openpgpErr == nil {
		log.Printf("Decryption of GPG MasterKey with fingerprint %s succeeded", key.Fingerprint)
		return dataKey, nil
	}
	log.Print("Decryption with golang's openpgp package failed, falling back to the GPG binary")
	dataKey, binaryErr := key.decryptWithGPGBinary()
	if binaryErr == nil {
		log.Printf("Decryption of GPG MasterKey with fingerprint %s succeeded", key.Fingerprint)
		return dataKey, nil
	}
	log.Printf("Decryption of GPG MasterKey with fingerprint %s failed", key.Fingerprint)
	return nil, fmt.Errorf(`could not encrypt data key with PGP key.
	\tgolang.org/x/crypto/openpgp error: %s
	\tGPG binary error: %s`, openpgpErr, binaryErr)
}

// NeedsRotation returns whether the data key needs to be rotated or not
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate).Hours() > 24*30*6
}

// ToString returns the string representation of the key, i.e. its fingerprint
func (key *MasterKey) ToString() string {
	return key.Fingerprint
}

func (key *MasterKey) gpgHome() string {
	dir := os.Getenv("GNUPGHOME")
	if dir == "" {
		usr, err := user.Current()
		if err != nil {
			return "~/.gnupg"
		}
		return path.Join(usr.HomeDir, ".gnupg")
	}
	return dir
}

// NewMasterKeyFromFingerprint takes a PGP fingerprint and returns a new MasterKey with that fingerprint
func NewMasterKeyFromFingerprint(fingerprint string) *MasterKey {
	return &MasterKey{
		Fingerprint:  strings.Replace(fingerprint, " ", "", -1),
		CreationDate: time.Now().UTC(),
	}
}

// MasterKeysFromFingerprintString takes a comma separated list of PGP fingerprints and returns a slice of new MasterKeys with those fingerprints
func MasterKeysFromFingerprintString(fingerprint string) []*MasterKey {
	var keys []*MasterKey
	if fingerprint == "" {
		return keys
	}
	for _, s := range strings.Split(fingerprint, ",") {
		keys = append(keys, NewMasterKeyFromFingerprint(s))
	}
	return keys
}

func (key *MasterKey) loadRing(path string) (openpgp.EntityList, error) {
	f, err := os.Open(path)
	if err != nil {
		return openpgp.EntityList{}, err
	}
	defer f.Close()
	keyring, err := openpgp.ReadKeyRing(f)
	if err != nil {
		return keyring, err
	}
	return keyring, nil
}

func (key *MasterKey) secRing() (openpgp.EntityList, error) {
	return key.loadRing(key.gpgHome() + "/secring.gpg")
}

func (key *MasterKey) pubRing() (openpgp.EntityList, error) {
	return key.loadRing(key.gpgHome() + "/pubring.gpg")
}

func (key *MasterKey) fingerprintMap(ring openpgp.EntityList) map[string]openpgp.Entity {
	fps := make(map[string]openpgp.Entity)
	for _, entity := range ring {
		fp := strings.ToUpper(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
		if entity != nil {
			fps[fp] = *entity
		}
	}
	return fps
}

func (key *MasterKey) passphrasePrompt(keys []openpgp.Key, symmetric bool) ([]byte, error) {
	conn, err := gpgagent.NewConn()
	if err == gpgagent.ErrNoAgent {
		log.Printf("gpg-agent not found, continuing with manual passphrase input...")
		log.Print("Enter PGP key passphrase: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return nil, err
		}
		for _, k := range keys {
			k.PrivateKey.Decrypt(pass)
		}
		return pass, err
	}
	if err != nil {
		return nil, fmt.Errorf("Could not establish connection with gpg-agent: %s", err)
	}
	defer conn.Close()
	for _, k := range keys {
		req := gpgagent.PassphraseRequest{
			CacheKey: k.PublicKey.KeyIdShortString(),
			Prompt:   "Passphrase",
			Desc:     fmt.Sprintf("Unlock key %s to decrypt sops's key", k.PublicKey.KeyIdShortString()),
		}
		pass, err := conn.GetPassphrase(&req)
		if err != nil {
			return nil, fmt.Errorf("gpg-agent passphrase request errored: %s", err)
		}
		k.PrivateKey.Decrypt([]byte(pass))
		return []byte(pass), nil
	}
	return nil, fmt.Errorf("No key to unlock")
}

// ToMap converts the MasterKey into a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["fp"] = key.Fingerprint
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}
