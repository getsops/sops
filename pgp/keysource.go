/*
Package pgp contains an implementation of the github.com/getsops/sops/v3.MasterKey
interface that encrypts and decrypts the data key by first trying with the
github.com/ProtonMail/go-crypto/openpgp package and if that fails, by calling
the "gpg" binary.
*/
package pgp // import "github.com/getsops/sops/v3/pgp"

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	gpgagent "github.com/getsops/gopgagent"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify a PGP MasterKey.
	KeyTypeIdentifier = "pgp"
	// SopsGpgExecEnv can be set as an environment variable to overwrite the
	// GnuPG binary used.
	SopsGpgExecEnv = "SOPS_GPG_EXEC"
)

var (
	// pgpTTL is the duration after which a MasterKey requires rotation.
	pgpTTL = time.Hour * 24 * 30 * 6
	// defaultPubRing is the relative path to the pubring in the GnuPG
	// home.
	// NB: This format is no longer in use since GnuPG >=2.1, which switched
	// to .kbx for new installations, and merged secring.gpg into pubring.gpg.
	defaultPubRing = "pubring.gpg"
	// defaultSecRing is the relative path to the secring in the GnuPG
	// home.
	// NB: GnuPG >= 2.1 merged this together with pubring.gpg, see
	// defaultPubRing.
	defaultSecRing = "secring.gpg"
)

// log is the global logger for any PGP MasterKey.
// TODO(hidde): this is not-so-nice for any implementation other than the CLI,
// as it becomes difficult to sugar the logger with data for e.g. individual
// processes.
var log *logrus.Logger

func init() {
	log = logging.NewLogger("PGP")
}

// MasterKey is a PGP key used to securely store SOPS' data key by
// encrypting it and decrypting it.
type MasterKey struct {
	// Fingerprint contains the fingerprint of the PGP key used to Encrypt
	// or Decrypt the data key with.
	Fingerprint string
	// EncryptedKey contains the SOPS data key encrypted with PGP.
	EncryptedKey string
	// CreationDate of the MasterKey, used to determine if the EncryptedKey
	// needs rotation.
	CreationDate time.Time

	// gnuPGHomeDir contains the absolute path to a GnuPG home directory.
	// It can be injected by a (local) keyservice.KeyServiceServer using
	// GnuPGHome.ApplyToMasterKey().
	gnuPGHomeDir string
	// disableOpenPGP instructs the MasterKey to skip attempting to open any
	// pubRing or secRing using OpenPGP.
	disableOpenPGP bool
	// pubRing contains the absolute path to a public keyring used by OpenPGP.
	// When empty, defaultPubRing relative to GnuPG home is assumed.
	pubRing string
	// secRing contains the absolute path to a sec keyring used by OpenPGP.
	// When empty, defaultSecRing relative to GnuPG home is assumed.
	secRing string
}

// NewMasterKeyFromFingerprint takes a PGP fingerprint and returns a new
// MasterKey with that fingerprint.
func NewMasterKeyFromFingerprint(fingerprint string) *MasterKey {
	return &MasterKey{
		Fingerprint:  strings.Replace(fingerprint, " ", "", -1),
		CreationDate: time.Now().UTC(),
	}
}

// MasterKeysFromFingerprintString takes a comma separated list of PGP
// fingerprints and returns a slice of new MasterKeys with those fingerprints.
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

// GnuPGHome is the absolute path to a GnuPG home directory.
// A new keyring can be constructed by combining the use of NewGnuPGHome() and
// Import() or ImportFile().
type GnuPGHome string

// NewGnuPGHome initializes a new GnuPGHome in a temporary directory.
// The caller is expected to handle the garbage collection of the created
// directory.
func NewGnuPGHome() (GnuPGHome, error) {
	tmpDir, err := os.MkdirTemp("", "sops-gnupghome-")
	if err != nil {
		return "", fmt.Errorf("failed to create new GnuPG home: %w", err)
	}
	return GnuPGHome(tmpDir), nil
}

// Import attempts to import the armored key bytes into the GnuPGHome keyring.
// It returns an error if the GnuPGHome does not pass Validate, or if the
// import failed.
func (d GnuPGHome) Import(armoredKey []byte) error {
	if err := d.Validate(); err != nil {
		return fmt.Errorf("cannot import armored key data into GnuPG keyring: %w", err)
	}

	args := []string{"--batch", "--import"}
	_, stderr, err := gpgExec(d.String(), args, bytes.NewReader(armoredKey))
	if err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		errStr := err.Error()
		var sb strings.Builder
		sb.WriteString("failed to import armored key data into GnuPG keyring")
		if len(stderrStr) > 0 {
			if len(errStr) > 0 {
				fmt.Fprintf(&sb, " (%s)", errStr)
			}
			fmt.Fprintf(&sb, ": %s", stderrStr)
		} else if len(errStr) > 0 {
			fmt.Fprintf(&sb, ": %s", errStr)
		}
		return errors.New(sb.String())
	}
	return nil
}

// ImportFile attempts to import the armored key file into the GnuPGHome
// keyring.
// It returns an error if the GnuPGHome does not pass Validate, or if the
// import failed.
func (d GnuPGHome) ImportFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read armored key data from file: %w", err)
	}
	return d.Import(b)
}

// Cleanup deletes the GnuPGHome if it passes Validate.
// It returns an error if the GnuPGHome does not pass Validate, or if the
// removal failed.
func (d GnuPGHome) Cleanup() error {
	if err := d.Validate(); err != nil {
		return err
	}
	return os.RemoveAll(d.String())
}

// Validate ensures the GnuPGHome is a valid GnuPG home directory path.
// When validation fails, it returns a descriptive reason as error.
func (d GnuPGHome) Validate() error {
	if d == "" {
		return fmt.Errorf("empty GNUPGHOME path")
	}
	if !filepath.IsAbs(d.String()) {
		return fmt.Errorf("GNUPGHOME must be an absolute path")
	}
	fi, err := os.Lstat(d.String())
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("GNUPGHOME does not exist")
		}
		return fmt.Errorf("cannot stat GNUPGHOME: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("GNUGPHOME is not a directory")
	}
	if perm := fi.Mode().Perm(); perm != 0o700 {
		return fmt.Errorf("GNUPGHOME has invalid permissions: got %#o wanted %#o", perm, 0o700)
	}
	return nil
}

// String returns the GnuPGHome as a string. It does not Validate.
func (d GnuPGHome) String() string {
	return string(d)
}

// ApplyToMasterKey configures the GnuPGHome on the provided key if it passes
// Validate.
func (d GnuPGHome) ApplyToMasterKey(key *MasterKey) {
	if err := d.Validate(); err == nil {
		key.gnuPGHomeDir = d.String()
	}
}

// DisableOpenPGP disables encrypt and decrypt operations using OpenPGP.
type DisableOpenPGP struct{}

// ApplyToMasterKey configures the provided key to not use OpenPGP.
func (d DisableOpenPGP) ApplyToMasterKey(key *MasterKey) {
	key.disableOpenPGP = true
}

// PubRing can be used to configure the absolute path to a public keyring
// used by OpenPGP.
type PubRing string

// ApplyToMasterKey configures the provided key to not use the GnuPG agent.
func (r PubRing) ApplyToMasterKey(key *MasterKey) {
	key.pubRing = string(r)
}

// SecRing can be used to configure the absolute path to a sec keyring
// used by OpenPGP.
type SecRing string

// ApplyToMasterKey configures the provided key to not use the GnuPG agent.
func (r SecRing) ApplyToMasterKey(key *MasterKey) {
	key.secRing = string(r)
}

// errSet is a collection of captured errors.
type errSet []error

// Error joins the errors into a "; " separated string.
func (e errSet) Error() string {
	str := make([]string, len(e))
	for i, err := range e {
		str[i] = err.Error()
	}
	return strings.Join(str, "; ")
}

// Encrypt encrypts the data key with the PGP key with the same
// fingerprint as the MasterKey.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	var errs errSet

	if !key.disableOpenPGP {
		openpgpErr := key.encryptWithOpenPGP(dataKey)
		if openpgpErr == nil {
			log.WithField("fingerprint", key.Fingerprint).Info("Encryption succeeded")
			return nil
		}
		errs = append(errs, fmt.Errorf("github.com/ProtonMail/go-crypto/openpgp error: %w", openpgpErr))
	}

	binaryErr := key.encryptWithGnuPG(dataKey)
	if binaryErr == nil {
		log.WithField("fingerprint", key.Fingerprint).Info("Encryption succeeded")
		return nil
	}
	errs = append(errs, fmt.Errorf("GnuPG binary error: %w", binaryErr))

	log.WithField("fingerprint", key.Fingerprint).Info("Encryption failed")
	return fmt.Errorf("could not encrypt data key with PGP key: %w", errs)
}

// encryptWithOpenPGP attempts to encrypt the data key using OpenPGP with the
// PGP key that belongs to Fingerprint. It sets EncryptedDataKey, or returns
// an error.
func (key *MasterKey) encryptWithOpenPGP(dataKey []byte) error {
	entity, err := key.retrievePubKey()
	if err != nil {
		return err
	}

	encBuf := new(bytes.Buffer)
	armorBuf, err := armor.Encode(encBuf, "PGP MESSAGE", nil)
	if err != nil {
		return err
	}
	plainBuf, err := openpgp.Encrypt(armorBuf, []*openpgp.Entity{&entity}, nil, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return err
	}
	_, err = plainBuf.Write(dataKey)
	if err != nil {
		return err
	}
	err = plainBuf.Close()
	if err != nil {
		return err
	}
	err = armorBuf.Close()
	if err != nil {
		return err
	}

	b, err := io.ReadAll(encBuf)
	if err != nil {
		return err
	}

	key.SetEncryptedDataKey(b)
	return nil
}

// encryptWithOpenPGP attempts to encrypt the data key using GnuPG with the
// PGP key that belongs to Fingerprint. It sets EncryptedDataKey, or returns
// an error.
func (key *MasterKey) encryptWithGnuPG(dataKey []byte) error {
	fingerprint := shortenFingerprint(key.Fingerprint)

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
	stdout, stderr, err := gpgExec(key.gnuPGHomeDir, args, bytes.NewReader(dataKey))
	if err != nil {
		return fmt.Errorf("failed to encrypt sops data key with pgp: %s", strings.TrimSpace(stderr.String()))
	}

	key.SetEncryptedDataKey(bytes.TrimSpace(stdout.Bytes()))
	return nil
}

// EncryptIfNeeded encrypts the data key with PGP only if it's needed,
// that is, if it hasn't been encrypted already.
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
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

// Decrypt first attempts to obtain the data key from the EncryptedKey
// stored in the MasterKey using OpenPGP, before falling back to GnuPG.
// When both attempts fail, an error is returned.
func (key *MasterKey) Decrypt() ([]byte, error) {
	var errs errSet

	if !key.disableOpenPGP {
		dataKey, openpgpErr := key.decryptWithOpenPGP()
		if openpgpErr == nil {
			log.WithField("fingerprint", key.Fingerprint).Info("Decryption succeeded")
			return dataKey, nil
		}
		errs = append(errs, fmt.Errorf("github.com/ProtonMail/go-crypto/openpgp error: %w", openpgpErr))
	}

	dataKey, binaryErr := key.decryptWithGnuPG()
	if binaryErr == nil {
		log.WithField("fingerprint", key.Fingerprint).Info("Decryption succeeded")
		return dataKey, nil
	}
	errs = append(errs, fmt.Errorf("GnuPG binary error: %w", binaryErr))

	log.WithField("fingerprint", key.Fingerprint).Info("Decryption failed")
	return nil, fmt.Errorf("could not decrypt data key with PGP key: %w", errs)
}

// decryptWithOpenPGP attempts to obtain the data key from the EncryptedKey
// using OpenPGP and returns the result.
//
// Note: the current development of OpenPGP vs GnuPG has moved in separate
// directions. This means that e.g. GnuPG >=2.1 works with a .kbx format which
// can not be read by OpenPGP. Given the further assumptions around the
// placement of the files, and the generic fallback Decrypt uses, this raises
// the question of how widely utilized this method still is.
func (key *MasterKey) decryptWithOpenPGP() ([]byte, error) {
	ring, err := key.getSecRing()
	if err != nil {
		return nil, fmt.Errorf("could not load secring: %s", err)
	}
	block, err := armor.Decode(strings.NewReader(key.EncryptedKey))
	if err != nil {
		return nil, fmt.Errorf("armor decoding failed: %s", err)
	}
	md, err := openpgp.ReadMessage(block.Body, ring, key.passphrasePrompt(), nil)
	if err != nil {
		return nil, fmt.Errorf("reading PGP message failed: %s", err)
	}
	if b, err := io.ReadAll(md.UnverifiedBody); err == nil {
		return b, nil
	}
	return nil, fmt.Errorf("the key could not be decrypted with any of the PGP entries")
}

// decryptWithGnuPG attempts to obtain the data key from the EncryptedKey using
// GnuPG and returns the result. If DisableAgent is configured on the MasterKey,
// the GnuPG agent is not enabled. When the decryption command fails, it returns
// the error from stdout.
func (key *MasterKey) decryptWithGnuPG() ([]byte, error) {
	args := []string{
		"-d",
	}
	stdout, stderr, err := gpgExec(key.gnuPGHomeDir, args, strings.NewReader(key.EncryptedKey))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt sops data key with pgp: %s",
			strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

// NeedsRotation returns whether the data key needs to be rotated
// or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (pgpTTL)
}

// ToString returns the string representation of the key, i.e. its
// fingerprint.
func (key *MasterKey) ToString() string {
	return key.Fingerprint
}

// ToMap converts the MasterKey into a map for serialization purposes.
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["fp"] = key.Fingerprint
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// retrievePubKey attempts to retrieve the public key from the public keyring
// by Fingerprint.
func (key *MasterKey) retrievePubKey() (openpgp.Entity, error) {
	ring, err := key.getPubRing()
	if err == nil {
		fingerprints := fingerprintIndex(ring)
		entity, ok := fingerprints[key.Fingerprint]
		if ok {
			return entity, nil
		}
	}
	return openpgp.Entity{},
		fmt.Errorf("key with fingerprint '%s' is not available "+
			"in keyring", key.Fingerprint)
}

// getPubRing loads the public keyring from the configured path, or falls back
// to defaultPubRing relative to the GnuPG home. It returns an openpgp.EntityList
// read from the keyring, or an error.
func (key *MasterKey) getPubRing() (openpgp.EntityList, error) {
	path := key.pubRing
	if path == "" {
		path = filepath.Join(gnuPGHome(key.gnuPGHomeDir), defaultPubRing)
	}
	return loadRing(path)
}

// getSecRing loads the sec keyring from the configured path, or falls back
// to defaultSecRing relative to the GnuPG home. It returns an openpgp.EntityList
// read from the keyring, or an error.
func (key *MasterKey) getSecRing() (openpgp.EntityList, error) {
	path := key.secRing
	if path == "" {
		path = filepath.Join(gnuPGHome(key.gnuPGHomeDir), defaultSecRing)
	}
	if _, err := os.Lstat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return key.getPubRing()
	}
	return loadRing(path)
}

// passphrasePrompt prompts the user for a passphrase when this is required for
// encryption or decryption.
func (key *MasterKey) passphrasePrompt() func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
	callCounter := 0
	maxCalls := 3
	return func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if callCounter >= maxCalls {
			return nil, fmt.Errorf("function passphrasePrompt called too many times")
		}
		callCounter++

		conn, err := gpgagent.NewConn()
		if err == gpgagent.ErrNoAgent {
			log.Infof("gpg-agent not found, continuing with manual passphrase " +
				"input...")
			fmt.Print("Enter PGP key passphrase: ")
			pass, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return nil, err
			}
			for _, k := range keys {
				k.PrivateKey.Decrypt(pass)
			}
			return pass, err
		}
		if err != nil {
			return nil, fmt.Errorf("could not establish connection with gpg-agent: %s", err)
		}
		defer func(conn *gpgagent.Conn) {
			if err := conn.Close(); err != nil {
				log.Errorf("failed to close connection with gpg-agent: %s", err)
			}
		}(conn)

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

		return nil, fmt.Errorf("no key to unlock")
	}
}

// loadRing attempts to load the keyring from the provided path.
// Unsupported keys are ignored as long as at least a single valid key is
// found.
func loadRing(path string) (openpgp.EntityList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	keyring, err := openpgp.ReadKeyRing(f)
	if err != nil {
		return nil, err
	}
	return keyring, nil
}

// fingerprintIndex indexes the openpgp.Entity objects from the given ring
// by their fingerprint, and returns the result.
func fingerprintIndex(ring openpgp.EntityList) map[string]openpgp.Entity {
	fps := make(map[string]openpgp.Entity)
	for _, entity := range ring {
		if entity != nil {
			fp := strings.ToUpper(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
			fps[fp] = *entity
		}
	}
	return fps
}

// gpgExec runs the provided args with the gpgBinary, while restricting it to
// homeDir when provided. Stdout and stderr can be read from the returned
// buffers. When the command fails, an error is returned.
func gpgExec(homeDir string, args []string, stdin io.Reader) (stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	if homeDir != "" {
		args = append([]string{"--homedir", homeDir}, args...)
	}

	cmd := exec.Command(gpgBinary(), args...)
	cmd.Stdin = stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return
}

// gpgBinary returns the GnuPG binary which must be used.
// It allows for runtime modifications by setting the replacement binary
// via the environment variable SopsGpgExecEnv.
func gpgBinary() string {
	if envBinary := os.Getenv(SopsGpgExecEnv); envBinary != "" {
		return envBinary
	}
	return "gpg"
}

// gnuPGHome determines the GnuPG home directory, and returns its path.
// In order of preference:
//  1. customPath
//  2. $GNUPGHOME
//  3. user.Current().HomeDir/.gnupg
//  4. $HOME/.gnupg
func gnuPGHome(customPath string) string {
	if customPath != "" {
		return customPath
	}

	dir := os.Getenv("GNUPGHOME")
	if dir == "" {
		usr, err := user.Current()
		if err != nil {
			return filepath.Join(os.Getenv("HOME"), ".gnupg")
		}
		return filepath.Join(usr.HomeDir, ".gnupg")
	}
	return dir
}

// shortenFingerprint returns the short ID of the given fingerprint.
// This is mostly used for compatibility reasons, as older versions of GnuPG
// do not always like long IDs.
func shortenFingerprint(fingerprint string) string {
	if offset := len(fingerprint) - 16; offset > 0 {
		fingerprint = fingerprint[offset:]
	}
	return fingerprint
}
