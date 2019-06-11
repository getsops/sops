/*
Package pgp contains an implementation of the go.mozilla.org/sops.MasterKey interface that encrypts and decrypts the
data key by first trying with the golang.org/x/crypto/openpgp package and if that fails, by calling the "gpg" binary.
*/
package pgp //import "go.mozilla.org/sops/pgp"

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"os/exec"

	"github.com/howeyc/gopass"
	"github.com/sirupsen/logrus"
	gpgagent "go.mozilla.org/gopgagent"
	"go.mozilla.org/sops/logging"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("PGP")
}

// MasterKey is a PGP key used to securely store sops' data key by encrypting it and decrypting it
type MasterKey struct {
	Fingerprint  string
	EncryptedKey string
	CreationDate time.Time
}

// EncryptedDataKey returns the encrypted data key this master key holds
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key
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
	fingerprint := key.Fingerprint
	if offset := len(fingerprint) - 16; offset > 0 {
		fingerprint = fingerprint[offset:]
	}
	args := []string{
		"--no-default-recipient",
		"--yes",
		"--encrypt",
		"-a",
		"-r",
		key.Fingerprint,
		"--trusted-key",
		fingerprint,
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

func getKeyFromKeyServer(keyserver string, fingerprint string) (openpgp.Entity, error) {
	url := fmt.Sprintf("https://%s/pks/lookup?op=get&options=mr&search=0x%s", keyserver, fingerprint)
	resp, err := http.Get(url)
	if err != nil {
		return openpgp.Entity{}, fmt.Errorf("error getting key from keyserver: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return openpgp.Entity{}, fmt.Errorf("keyserver returned non-200 status code %s", resp.Status)
	}
	ents, err := openpgp.ReadArmoredKeyRing(resp.Body)
	if err != nil {
		return openpgp.Entity{}, fmt.Errorf("could not read entities: %s", err)
	}
	return *ents[0], nil
}

func gpgKeyServer() string {
	keyServer := "gpg.mozilla.org"
	if envKeyServer := os.Getenv("SOPS_GPG_KEYSERVER"); envKeyServer != "" {
		keyServer = envKeyServer
	}
	return keyServer
}

func (key *MasterKey) getPubKey() (openpgp.Entity, error) {
	ring, err := key.pubRing()
	if err == nil {
		fingerprints := key.fingerprintMap(ring)
		entity, ok := fingerprints[key.Fingerprint]
		if ok {
			return entity, nil
		}
	}
	keyServer := gpgKeyServer()
	entity, err := getKeyFromKeyServer(keyServer, key.Fingerprint)
	if err != nil {
		return openpgp.Entity{},
			fmt.Errorf("key with fingerprint %s is not available "+
				"in keyring and could not be retrieved from keyserver", key.Fingerprint)
	}
	return entity, nil
}

func (key *MasterKey) encryptWithCryptoOpenPGP(dataKey []byte) error {
	entity, err := key.getPubKey()
	if err != nil {
		return err
	}
	encbuf := new(bytes.Buffer)
	armorbuf, err := armor.Encode(encbuf, "PGP MESSAGE", nil)
	if err != nil {
		return err
	}
	plaintextbuf, err := openpgp.Encrypt(armorbuf, []*openpgp.Entity{&entity}, nil, &openpgp.FileHints{IsBinary: true}, nil)
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
	openpgpErr := key.encryptWithCryptoOpenPGP(dataKey)
	if openpgpErr == nil {
		log.WithField("fingerprint", key.Fingerprint).Info("Encryption succeeded")
		return nil
	}
	binaryErr := key.encryptWithGPGBinary(dataKey)
	if binaryErr == nil {
		log.WithField("fingerprint", key.Fingerprint).Info("Encryption succeeded")
		return nil
	}
	log.WithField("fingerprint", key.Fingerprint).Info("Encryption failed")
	return fmt.Errorf(
		`could not encrypt data key with PGP key: golang.org/x/crypto/openpgp error: %v; GPG binary error: %v`,
		openpgpErr, binaryErr)
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
		return b, nil
	}
	return nil, fmt.Errorf("The key could not be decrypted with any of the PGP entries")
}

// Decrypt uses PGP to obtain the data key from the EncryptedKey store in the MasterKey and returns it
func (key *MasterKey) Decrypt() ([]byte, error) {
	dataKey, openpgpErr := key.decryptWithCryptoOpenpgp()
	if openpgpErr == nil {
		log.WithField("fingerprint", key.Fingerprint).Info("Decryption succeeded")
		return dataKey, nil
	}
	dataKey, binaryErr := key.decryptWithGPGBinary()
	if binaryErr == nil {
		log.WithField("fingerprint", key.Fingerprint).Info("Decryption succeeded")
		return dataKey, nil
	}
	log.WithField("fingerprint", key.Fingerprint).Info("Decryption failed")
	return nil, fmt.Errorf(
		`could not decrypt data key with PGP key: golang.org/x/crypto/openpgp error: %v; GPG binary error: %v`,
		openpgpErr, binaryErr)
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
			return path.Join(os.Getenv("HOME"), "/.gnupg")
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
		log.Infof("gpg-agent not found, continuing with manual passphrase " +
			"input...")
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
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["fp"] = key.Fingerprint
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}
