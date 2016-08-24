package pgp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/howeyc/gopass"
	gpgagent "go.mozilla.org/gopgagent"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

// MasterKey is a PGP key used to securely store sops' data key by encrypting it and decrypting it
type MasterKey struct {
	Fingerprint  string
	EncryptedKey string
	CreationDate time.Time
}

// Encrypt encrypts the data key with the PGP key with the same fingerprint as the MasterKey. It looks for PGP public keys in $PGPHOME/pubring.gpg.
func (key *MasterKey) Encrypt(dataKey []byte) error {
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

// EncryptIfNeeded encrypts the data key with PGP only if it's needed, that is, if it hasn't been encrypted already
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt uses PGP to obtain the data key from the EncryptedKey store in the MasterKey and returns it
func (key *MasterKey) Decrypt() ([]byte, error) {
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
		return b, nil
	}
	return nil, fmt.Errorf("The key could not be decrypted with any of the GPG entries")
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
func NewMasterKeyFromFingerprint(fingerprint string) MasterKey {
	return MasterKey{
		Fingerprint:  strings.Replace(fingerprint, " ", "", -1),
		CreationDate: time.Now().UTC(),
	}
}

// MasterKeysFromFingerprintString takes a comma separated list of PGP fingerprints and returns a slice of new MasterKeys with those fingerprints
func MasterKeysFromFingerprintString(fingerprint string) []MasterKey {
	var keys []MasterKey
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
		fmt.Println("gpg-agent not found, continuing with manual passphrase input...")
		fmt.Print("Enter PGP key passphrase: ")
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
func (key MasterKey) ToMap() map[string]string {
	out := make(map[string]string)
	out["fp"] = key.Fingerprint
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}
