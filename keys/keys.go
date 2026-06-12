package keys

import (
	"fmt"
	"strings"
)

// MasterKey provides a way of securing the key used to encrypt the Tree by encrypting and decrypting said key.
type MasterKey interface {
	Encrypt(dataKey []byte) error
	EncryptIfNeeded(dataKey []byte) error
	EncryptedDataKey() []byte
	SetEncryptedDataKey([]byte)
	Decrypt() ([]byte, error)
	NeedsRotation() bool
	ToString() string
	ToMap() map[string]interface{}
	TypeToIdentifier() string
}

type CreationOptions struct {
	KmsEncryptionContext map[string]*string
	GlobalConfig         map[string]any
}

// KeyProvider is responsible for marshaling and unmarshaling MasterKeys
// to and from a generic map representation for a specific backend.
type KeyProvider interface {
	Type() string
	MarshalKey(key MasterKey) (map[string]any, error)
	UnmarshalKey(data map[string]any) (MasterKey, error)
	KeysFromConfig(config any, opts CreationOptions) ([]MasterKey, error)
}

// BugFixer is an optional interface that a KeyProvider can implement
// if it needs to apply legacy bug fixes directly to the SOPS tree.
type BugFixer interface {
	DetectTreeBugs(version string, keyGroups [][]MasterKey) bool
	BugExplanation() string
	RecoverDataKey(keyGroups [][]MasterKey, decryptFn func([][]MasterKey) ([]byte, error)) []byte
}

var KeyProviders = make(map[string]KeyProvider)

func RegisterProvider(provider KeyProvider) {
	KeyProviders[provider.Type()] = provider
}

func GetProvider(name string) KeyProvider {
	return KeyProviders[name]
}

func ParseStringSlice(field interface{}, fieldName string) ([]string, error) {
	if field == nil {
		return []string{}, nil
	}
	switch v := field.(type) {
	case string:
		if v == "" {
			return []string{}, nil
		}
		keys := strings.Split(v, ",")
		var result []string
		for _, key := range keys {
			if trimmed := strings.TrimSpace(key); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result, nil
	case []interface{}:
		var result []string
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			} else {
				return nil, fmt.Errorf(
					"invalid %s key configuration: expected string in list, got %T",
					fieldName,
					item,
				)
			}
		}
		return result, nil
	case []string:
		return v, nil
	default:
		return nil, fmt.Errorf(
			"invalid %s key configuration: expected string, []string, or nil, got %T",
			fieldName,
			field,
		)
	}
}

// ParseStringMap parses a comma separated key:value string into a map
func ParseStringMap(ctx string) map[string]*string {
	contextMap := make(map[string]*string)
	if ctx == "" {
		return contextMap
	}

	contexts := strings.Split(ctx, ",")

	for _, context := range contexts {
		// Only splitting on the first colon so that the values can contain colons
		kv := strings.SplitN(context, ":", 2)
		if len(kv) == 2 {
			contextMap[kv[0]] = &kv[1]
		}
	}
	return contextMap
}

// ProviderFlag defines a CLI flag provided by a KeyProvider.
type ProviderFlag struct {
	Name            string
	Usage           string
	EnvVar          string
	IsKeyIdentifier bool // If true, this flag identifies keys (used for rotation add/rm and slice flags)
}

// FlagGetter defines an interface for retrieving flag values, decoupling providers from the specific CLI framework.
type FlagGetter interface {
	String(name string) string
	StringSlice(name string) []string
}

// CLIProvider is an optional interface that a KeyProvider can implement
// to dynamically register CLI flags and parse keys from them.
type CLIProvider interface {
	CLIConfig() []ProviderFlag
	MasterKeysFromCLI(c FlagGetter, prefix string) ([]MasterKey, error)
}
