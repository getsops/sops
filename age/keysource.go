package age

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
	"filippo.io/age/plugin"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
	"golang.org/x/term"
)

const (
	// SopsAgeKeyEnv can be set as an environment variable with a string list
	// of age keys as value.
	SopsAgeKeyEnv = "SOPS_AGE_KEY"
	// SopsAgeKeyFileEnv can be set as an environment variable pointing to an
	// age keys file.
	SopsAgeKeyFileEnv = "SOPS_AGE_KEY_FILE"
	// SopsAgeKeyUserConfigPath is the default age keys file path in
	// getUserConfigDir().
	SopsAgeKeyUserConfigPath = "sops/age/keys.txt"
	// On macOS, os.UserConfigDir() ignores XDG_CONFIG_HOME. So we handle that manually.
	xdgConfigHome = "XDG_CONFIG_HOME"
	// KeyTypeIdentifier is the string used to identify an age MasterKey.
	KeyTypeIdentifier = "age"
)

// log is the global logger for any age MasterKey.
var log *logrus.Logger

func init() {
	log = logging.NewLogger("AGE")
}

// MasterKey is an age key used to Encrypt and Decrypt SOPS' data key.
type MasterKey struct {
	// Identity used to contain a Bench32-encoded private key.
	// Deprecated: private keys are no longer publicly exposed.
	// Instead, they are either injected by a (local) key service server
	// using ParsedIdentities.ApplyToMasterKey, or loaded from the runtime
	// environment (variables) as defined by the `SopsAgeKey*` constants.
	Identity string
	// Recipient contains the Bench32-encoded age public key used to Encrypt.
	Recipient string
	// EncryptedKey contains the SOPS data key encrypted with age.
	EncryptedKey string

	// parsedIdentities contains a slice of parsed age identities.
	// It is used to lazy-load the Identities at-most once.
	// It can also be injected by a (local) keyservice.KeyServiceServer using
	// ParsedIdentities.ApplyToMasterKey().
	parsedIdentities []age.Identity
	// parsedRecipient contains a parsed age public key.
	// It is used to lazy-load the Recipient at-most once.
	parsedRecipient age.Recipient
}

// MasterKeysFromRecipients takes a comma-separated list of Bech32-encoded
// public keys, parses them, and returns a slice of new MasterKeys.
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

// MasterKeyFromRecipient takes a Bech32-encoded age public key, parses it, and
// returns a new MasterKey.
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

// ParsedIdentities contains a set of parsed age identities.
// It allows for creating a (local) keyservice.KeyServiceServer which parses
// identities only once, to then inject them using ApplyToMasterKey() for all
// requests.
type ParsedIdentities []age.Identity

// Import attempts to parse the given identities, to then add them to itself.
// It returns any parsing error.
// A single identity argument is allowed to be a multiline string containing
// multiple identities. Empty lines and lines starting with "#" are ignored.
// It is not thread safe, and parallel importing would better be done by
// parsing (using age.ParseIdentities) and appending to the slice yourself, in
// combination with e.g. a sync.Mutex.
func (i *ParsedIdentities) Import(identity ...string) error {
	identities, err := parseIdentities(identity...)
	if err != nil {
		return fmt.Errorf("failed to parse and add to age identities: %w", err)
	}
	*i = append(*i, identities...)
	return nil
}

// ApplyToMasterKey configures the ParsedIdentities on the provided key.
func (i ParsedIdentities) ApplyToMasterKey(key *MasterKey) {
	key.parsedIdentities = i
}

// Encrypt takes a SOPS data key, encrypts it with the Recipient, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	if key.parsedRecipient == nil {
		parsedRecipient, err := parseRecipient(key.Recipient)
		if err != nil {
			log.WithField("recipient", key.parsedRecipient).Info("Encryption failed")
			return err
		}
		key.parsedRecipient = parsedRecipient
	}

	var buffer bytes.Buffer
	aw := armor.NewWriter(&buffer)
	w, err := age.Encrypt(aw, key.parsedRecipient)
	if err != nil {
		log.WithField("recipient", key.parsedRecipient).Info("Encryption failed")
		return fmt.Errorf("failed to create writer for encrypting sops data key with age: %w", err)
	}
	if _, err := w.Write(dataKey); err != nil {
		log.WithField("recipient", key.parsedRecipient).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with age: %w", err)
	}
	if err := w.Close(); err != nil {
		log.WithField("recipient", key.parsedRecipient).Info("Encryption failed")
		return fmt.Errorf("failed to close writer for encrypting sops data key with age: %w", err)
	}
	if err := aw.Close(); err != nil {
		log.WithField("recipient", key.parsedRecipient).Info("Encryption failed")
		return fmt.Errorf("failed to close armored writer: %w", err)
	}

	key.SetEncryptedDataKey(buffer.Bytes())
	log.WithField("recipient", key.parsedRecipient).Info("Encryption succeeded")
	return nil
}

// EncryptIfNeeded encrypts the provided SOPS data key, if it has not been
// encrypted yet.
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// EncryptedDataKey returns the encrypted SOPS data key this master key holds.
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted SOPS data key for this master key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Decrypt decrypts the EncryptedKey with the parsed or loaded identities, and
// returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	if len(key.parsedIdentities) == 0 {
		ids, err := key.loadIdentities()
		if err != nil {
			log.Info("Decryption failed")
			return nil, fmt.Errorf("failed to load age identities: %w", err)
		}
		ids.ApplyToMasterKey(key)
	}

	src := bytes.NewReader([]byte(key.EncryptedKey))
	ar := armor.NewReader(src)
	r, err := age.Decrypt(ar, key.parsedIdentities...)
	if err != nil {
		log.Info("Decryption failed")
		return nil, fmt.Errorf("failed to create reader for decrypting sops data key with age: %w", err)
	}

	var b bytes.Buffer
	if _, err := io.Copy(&b, r); err != nil {
		log.Info("Decryption failed")
		return nil, fmt.Errorf("failed to copy age decrypted data into bytes.Buffer: %w", err)
	}

	log.Info("Decryption succeeded")
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
	out := make(map[string]interface{})
	out["recipient"] = key.Recipient
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

func getUserConfigDir() (string, error) {
	if runtime.GOOS == "darwin" {
		if userConfigDir, ok := os.LookupEnv(xdgConfigHome); ok && userConfigDir != "" {
			return userConfigDir, nil
		}
	}
	return os.UserConfigDir()
}

// loadIdentities attempts to load the age identities based on runtime
// environment configurations (e.g. SopsAgeKeyEnv, SopsAgeKeyFileEnv,
// SopsAgeKeyUserConfigPath). It will load all found references, and expects
// at least one configuration to be present.
func (key *MasterKey) loadIdentities() (ParsedIdentities, error) {
	readers := make(map[string]io.Reader, 0)

	if ageKey, ok := os.LookupEnv(SopsAgeKeyEnv); ok {
		readers[SopsAgeKeyEnv] = strings.NewReader(ageKey)
	}

	if ageKeyFile, ok := os.LookupEnv(SopsAgeKeyFileEnv); ok {
		f, err := os.Open(ageKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s file: %w", SopsAgeKeyFileEnv, err)
		}
		defer f.Close()
		readers[SopsAgeKeyFileEnv] = f
	}

	userConfigDir, err := getUserConfigDir()
	if err != nil && len(readers) == 0 {
		return nil, fmt.Errorf("user config directory could not be determined: %w", err)
	}
	if userConfigDir != "" {
		ageKeyFilePath := filepath.Join(userConfigDir, filepath.FromSlash(SopsAgeKeyUserConfigPath))
		f, err := os.Open(ageKeyFilePath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		if errors.Is(err, os.ErrNotExist) && len(readers) == 0 {
			// If we have no other readers, presence of the file is required.
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		if err == nil {
			defer f.Close()
			readers[ageKeyFilePath] = f
		}
	}

	var identities ParsedIdentities
	for n, r := range readers {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, r)
		if err != nil {
			return nil, fmt.Errorf("failed to read '%s' age identities: %w", n, err)
		}
		ids, err := parseIdentities(buf.String())
		if err != nil {
			return nil, fmt.Errorf("failed to parse '%s' age identities: %w", n, err)
		}
		identities = append(identities, ids...)
	}
	return identities, nil
}

// clearLine clears the current line on the terminal, or opens a new line if
// terminal escape codes don't work.
func clearLine(out io.Writer) {
	const (
		CUI = "\033["   // Control Sequence Introducer
		CPL = CUI + "F" // Cursor Previous Line
		EL  = CUI + "K" // Erase in Line
	)

	// First, open a new line, which is guaranteed to work everywhere. Then, try
	// to erase the line above with escape codes.
	//
	// (We use CRLF instead of LF to work around an apparent bug in WSL2's
	// handling of CONOUT$. Only when running a Windows binary from WSL2, the
	// cursor would not go back to the start of the line with a simple LF.
	// Honestly, it's impressive CONIN$ and CONOUT$ work at all inside WSL2.)
	fmt.Fprintf(out, "\r\n"+CPL+EL)
}

func withTerminal(f func(in, out *os.File) error) error {
	if runtime.GOOS == "windows" {
		in, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		defer out.Close()
		return f(in, out)
	} else if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		defer tty.Close()
		return f(tty, tty)
	} else if term.IsTerminal(int(os.Stdin.Fd())) {
		return f(os.Stdin, os.Stdin)
	} else {
		return fmt.Errorf("standard input is not a terminal, and /dev/tty is not available: %v", err)
	}
}

// readSecret reads a value from the terminal with no echo. The prompt is ephemeral.
func readSecret(prompt string) (s []byte, err error) {
	err = withTerminal(func(in, out *os.File) error {
		fmt.Fprintf(out, "%s ", prompt)
		defer clearLine(out)
		s, err = term.ReadPassword(int(in.Fd()))
		return err
	})
	return
}

// readCharacter reads a single character from the terminal with no echo. The
// prompt is ephemeral.
func readCharacter(prompt string) (c byte, err error) {
	err = withTerminal(func(in, out *os.File) error {
		fmt.Fprintf(out, "%s ", prompt)
		defer clearLine(out)

		oldState, err := term.MakeRaw(int(in.Fd()))
		if err != nil {
			return err
		}
		defer term.Restore(int(in.Fd()), oldState)

		b := make([]byte, 1)
		if _, err := in.Read(b); err != nil {
			return err
		}

		c = b[0]
		return nil
	})
	return
}

var pluginTerminalUI = &plugin.ClientUI{
	DisplayMessage: func(name, message string) error {
		log.Infof("%s plugin: %s", name, message)
		return nil
	},
	RequestValue: func(name, message string, _ bool) (s string, err error) {
		defer func() {
			if err != nil {
				log.Warnf("could not read value for age-plugin-%s: %v", name, err)
			}
		}()
		secret, err := readSecret(message)
		if err != nil {
			return "", err
		}
		return string(secret), nil
	},
	Confirm: func(name, message, yes, no string) (choseYes bool, err error) {
		defer func() {
			if err != nil {
				log.Warnf("could not read value for age-plugin-%s: %v", name, err)
			}
		}()
		if no == "" {
			message += fmt.Sprintf(" (press enter for %q)", yes)
			_, err := readSecret(message)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		message += fmt.Sprintf(" (press [1] for %q or [2] for %q)", yes, no)
		for {
			selection, err := readCharacter(message)
			if err != nil {
				return false, err
			}
			switch selection {
			case '1':
				return true, nil
			case '2':
				return false, nil
			case '\x03': // CTRL-C
				return false, errors.New("user cancelled prompt")
			default:
				log.Warnf("reading value for age-plugin-%s: invalid selection %q", name, selection)
			}
		}
	},
	WaitTimer: func(name string) {
		log.Infof("waiting on %s plugin...", name)
	},
}

// parseRecipient attempts to parse a string containing an encoded age public
// key.
func parseRecipient(recipient string) (age.Recipient, error) {
	switch {
	case strings.HasPrefix(recipient, "age1") && strings.Count(recipient, "1") > 1:
		return plugin.NewRecipient(recipient, pluginTerminalUI)
	case strings.HasPrefix(recipient, "age1"):
		return age.ParseX25519Recipient(recipient)
	}

	return nil, fmt.Errorf("unknown recipient type: %q", recipient)
}

// parseIdentities attempts to parse the string set of encoded age identities.
// A single identity argument is allowed to be a multiline string containing
// multiple identities. Empty lines and lines starting with "#" are ignored.
func parseIdentities(identity ...string) (ParsedIdentities, error) {
	var identities []age.Identity
	for _, i := range identity {
		parsed, err := _parseIdentities(strings.NewReader(i))
		if err != nil {
			return nil, err
		}
		identities = append(identities, parsed...)
	}
	return identities, nil
}

func parseIdentity(s string) (age.Identity, error) {
	switch {
	case strings.HasPrefix(s, "AGE-PLUGIN-"):
		return plugin.NewIdentity(s, pluginTerminalUI)
	case strings.HasPrefix(s, "AGE-SECRET-KEY-1"):
		return age.ParseX25519Identity(s)
	default:
		return nil, fmt.Errorf("unknown identity type")
	}
}

// parseIdentities is like age.ParseIdentities, but supports plugin identities.
func _parseIdentities(f io.Reader) (ParsedIdentities, error) {
	const privateKeySizeLimit = 1 << 24 // 16 MiB
	var ids []age.Identity
	scanner := bufio.NewScanner(io.LimitReader(f, privateKeySizeLimit))
	var n int
	for scanner.Scan() {
		n++
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		i, err := parseIdentity(line)
		if err != nil {
			return nil, fmt.Errorf("error at line %d: %v", n, err)
		}
		ids = append(ids, i)

	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read secret keys file: %v", err)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no secret keys found")
	}
	return ids, nil
}
