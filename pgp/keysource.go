package pgp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/howeyc/gopass"
	"go.mozilla.org/sops/gpgagent"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

type GPGMasterKey struct {
	Fingerprint  string
	EncryptedKey string
	CreationDate time.Time
}

func (key *GPGMasterKey) Encrypt(dataKey string) error {
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
	_, err = plaintextbuf.Write([]byte(dataKey))
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

func (key *GPGMasterKey) EncryptIfNeeded(dataKey string) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

func (key *GPGMasterKey) Decrypt() (string, error) {
	ring, err := key.secRing()
	if err != nil {
		return "", fmt.Errorf("Could not load secring: %s", err)
	}
	block, err := armor.Decode(strings.NewReader(key.EncryptedKey))
	if err != nil {
		return "", fmt.Errorf("Armor decoding failed: %s", err)
	}
	md, err := openpgp.ReadMessage(block.Body, ring, key.passphrasePrompt, nil)
	if err != nil {
		return "", fmt.Errorf("Reading PGP message failed: %s", err)
	}
	if b, err := ioutil.ReadAll(md.UnverifiedBody); err == nil {
		return string(b), nil
	}
	return "", fmt.Errorf("The key could not be decrypted with any of the GPG entries")
}

func (key *GPGMasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate).Hours() > 24*30*6
}

func (key *GPGMasterKey) ToString() string {
	return key.Fingerprint
}

func (key *GPGMasterKey) gpgHome() string {
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

func NewGPGMasterKeyFromFingerprint(fingerprint string) GPGMasterKey {
	return GPGMasterKey{
		Fingerprint: strings.Replace(fingerprint, " ", "", -1),
	}
}

func GPGMasterKeysFromFingerprintString(fingerprint string) []GPGMasterKey {
	var keys []GPGMasterKey
	for _, s := range strings.Split(fingerprint, ",") {
		keys = append(keys, NewGPGMasterKeyFromFingerprint(s))
	}
	return keys
}

func (key *GPGMasterKey) loadRing(path string) (openpgp.EntityList, error) {
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

func (key *GPGMasterKey) secRing() (openpgp.EntityList, error) {
	return key.loadRing(key.gpgHome() + "/secring.gpg")
}

func (key *GPGMasterKey) pubRing() (openpgp.EntityList, error) {
	return key.loadRing(key.gpgHome() + "/pubring.gpg")
}

func (key *GPGMasterKey) fingerprintMap(ring openpgp.EntityList) map[string]openpgp.Entity {
	fps := make(map[string]openpgp.Entity)
	for _, entity := range ring {
		fp := strings.ToUpper(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
		if entity != nil {
			fps[fp] = *entity
		}
	}
	return fps
}

func (key *GPGMasterKey) passphrasePrompt(keys []openpgp.Key, symmetric bool) ([]byte, error) {
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
