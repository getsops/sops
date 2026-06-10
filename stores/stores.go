/*
Package stores acts as a layer between the internal representation of encrypted files and the encrypted files
themselves.

Subpackages implement serialization and deserialization to multiple formats.

This package defines the structure SOPS files should have and conversions to and from the internal representation. Part
of the purpose of this package is to make it easy to change the SOPS file format while remaining backwards-compatible.
*/
package stores

import (
	"slices"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/keys"
)

const (
	// SopsMetadataKey is the key used to store SOPS metadata at in SOPS encrypted files.
	SopsMetadataKey = "sops"
)

// Metadata is stored in SOPS encrypted files, and it contains the information necessary to decrypt the file.
// This struct is just used for serialization, and SOPS uses another struct internally, sops.Metadata. It exists
// in order to allow the binary format to stay backwards compatible over time, but at the same time allow the internal
// representation SOPS uses to change over time.
type metadata struct {
	ShamirThreshold           int                           `mapstructure:"shamir_threshold,omitempty"`
	KeyGroups                 []map[string][]map[string]any `mapstructure:"key_groups,omitempty,deep"`
	LastModified              string                        `mapstructure:"lastmodified"`
	MessageAuthenticationCode string                        `mapstructure:"mac"`
	UnencryptedSuffix         string                        `mapstructure:"unencrypted_suffix,omitempty"`
	EncryptedSuffix           string                        `mapstructure:"encrypted_suffix,omitempty"`
	UnencryptedRegex          string                        `mapstructure:"unencrypted_regex,omitempty"`
	EncryptedRegex            string                        `mapstructure:"encrypted_regex,omitempty"`
	UnencryptedCommentRegex   string                        `mapstructure:"unencrypted_comment_regex,omitempty"`
	EncryptedCommentRegex     string                        `mapstructure:"encrypted_comment_regex,omitempty"`
	MACOnlyEncrypted          bool                          `mapstructure:"mac_only_encrypted,omitempty"`
	Version                   string                        `mapstructure:"version"`

	// Legacy top level keys
	PGPKeys           []map[string]any `mapstructure:"pgp,omitempty,deep"`
	PluginKeys        []map[string]any `mapstructure:"plugins,omitempty,deep"`
	KMSKeys           []map[string]any `mapstructure:"kms,omitempty,deep"`
	GCPKMSKeys        []map[string]any `mapstructure:"gcp_kms,omitempty,deep"`
	HCKmsKeys         []map[string]any `mapstructure:"hckms,omitempty,deep"`
	AzureKeyVaultKeys []map[string]any `mapstructure:"azure_kv,omitempty,deep"`
	VaultKeys         []map[string]any `mapstructure:"hc_vault,omitempty,deep"`
	AgeKeys           []map[string]any `mapstructure:"age,omitempty,deep"`
}

// metadataFromInternal converts an internal SOPS metadata representation to a
// representation appropriate for storage.
func metadataFromInternal(sopsMetadata sops.Metadata) metadata {
	var m metadata
	m.LastModified = sopsMetadata.LastModified.Format(time.RFC3339)
	m.UnencryptedSuffix = sopsMetadata.UnencryptedSuffix
	m.EncryptedSuffix = sopsMetadata.EncryptedSuffix
	m.UnencryptedRegex = sopsMetadata.UnencryptedRegex
	m.EncryptedRegex = sopsMetadata.EncryptedRegex
	m.UnencryptedCommentRegex = sopsMetadata.UnencryptedCommentRegex
	m.EncryptedCommentRegex = sopsMetadata.EncryptedCommentRegex
	m.MessageAuthenticationCode = sopsMetadata.MessageAuthenticationCode
	m.MACOnlyEncrypted = sopsMetadata.MACOnlyEncrypted
	m.Version = sopsMetadata.Version
	m.ShamirThreshold = sopsMetadata.ShamirThreshold

	if len(sopsMetadata.KeyGroups) == 1 {
		group := sopsMetadata.KeyGroups[0]
		keys := keysFromGroup(group)
		m.PGPKeys = keys["pgp"]
		m.KMSKeys = keys["kms"]
		m.GCPKMSKeys = keys["gcp_kms"]
		m.HCKmsKeys = keys["hckms"]
		m.VaultKeys = keys["hc_vault"]
		m.AzureKeyVaultKeys = keys["azure_kv"]
		m.AgeKeys = keys["age"]
		m.PluginKeys = keys["plugins"]
	} else {
		for _, group := range sopsMetadata.KeyGroups {
			m.KeyGroups = append(m.KeyGroups, keysFromGroup(group))
		}
	}
	return m
}

func keysFromGroup(group sops.KeyGroup) map[string][]map[string]any {
	result := make(map[string][]map[string]any)
	for _, key := range group {
		for name, provider := range keys.KeyProviders {
			data, err := provider.MarshalKey(key)
			if err == nil && data != nil {
				result[name] = append(result[name], data)
				break
			}
		}
	}
	return result
}

// ToInternal converts a storage-appropriate Metadata struct to a SOPS internal representation
func (m *metadata) ToInternal() (sops.Metadata, error) {
	lastModified, err := time.Parse(time.RFC3339, m.LastModified)
	if err != nil {
		return sops.Metadata{}, err
	}
	groups, err := m.internalKeygroups()
	if err != nil {
		return sops.Metadata{}, err
	}

	cryptRuleCount := 0
	if m.UnencryptedSuffix != "" {
		cryptRuleCount++
	}
	if m.EncryptedSuffix != "" {
		cryptRuleCount++
	}
	if m.UnencryptedRegex != "" {
		cryptRuleCount++
	}
	if m.EncryptedRegex != "" {
		cryptRuleCount++
	}
	if m.UnencryptedCommentRegex != "" {
		cryptRuleCount++
	}
	if m.EncryptedCommentRegex != "" {
		cryptRuleCount++
	}

	if cryptRuleCount > 1 {
		return sops.Metadata{}, fmt.Errorf(
			"Cannot use more than one of encrypted_suffix, unencrypted_suffix, encrypted_regex, unencrypted_regex, encrypted_comment_regex, or unencrypted_comment_regex in the same file",
		)
	}

	if cryptRuleCount == 0 {
		m.UnencryptedSuffix = sops.DefaultUnencryptedSuffix
	}
	return sops.Metadata{
		KeyGroups:                 groups,
		ShamirThreshold:           m.ShamirThreshold,
		Version:                   m.Version,
		MessageAuthenticationCode: m.MessageAuthenticationCode,
		UnencryptedSuffix:         m.UnencryptedSuffix,
		EncryptedSuffix:           m.EncryptedSuffix,
		UnencryptedRegex:          m.UnencryptedRegex,
		EncryptedRegex:            m.EncryptedRegex,
		UnencryptedCommentRegex:   m.UnencryptedCommentRegex,
		EncryptedCommentRegex:     m.EncryptedCommentRegex,
		MACOnlyEncrypted:          m.MACOnlyEncrypted,
		LastModified:              lastModified,
	}, nil
}

func internalGroupFrom(groupMap map[string][]map[string]any) (sops.KeyGroup, error) {
	var internalGroup sops.KeyGroup
	// to match the old behavior of the code, we use this order
	order := []string{"kms", "gcp_kms", "hckms", "azure_kv", "hc_vault", "pgp", "age", "plugins"}
	
	for _, name := range order {
		dataList, ok := groupMap[name]
		if !ok {
			continue
		}
		provider := keys.GetProvider(name)
		if provider == nil {
			continue
		}
		for _, data := range dataList {
			key, err := provider.UnmarshalKey(data)
			if err != nil {
				return nil, err
			}
			internalGroup = append(internalGroup, key)
		}
	}
	
	for name, dataList := range groupMap {
		found := slices.Contains(order, name)
		if found {
			continue
		}
		
		provider := keys.GetProvider(name)
		if provider == nil {
			continue // skip unknown providers
		}
		for _, data := range dataList {
			key, err := provider.UnmarshalKey(data)
			if err != nil {
				return nil, err
			}
			internalGroup = append(internalGroup, key)
		}
	}
	
	return internalGroup, nil
}

func (m *metadata) internalKeygroups() ([]sops.KeyGroup, error) {
	var internalGroups []sops.KeyGroup
	hasTopLevelKeys := len(m.PGPKeys) > 0 || len(m.KMSKeys) > 0 || len(m.GCPKMSKeys) > 0 ||
		len(m.HCKmsKeys) > 0 ||
		len(m.AzureKeyVaultKeys) > 0 ||
		len(m.VaultKeys) > 0 ||
		len(m.AgeKeys) > 0 ||
		len(m.PluginKeys) > 0

	if hasTopLevelKeys {
		topLevelGroup := map[string][]map[string]any{
			"pgp":      m.PGPKeys,
			"kms":      m.KMSKeys,
			"gcp_kms":  m.GCPKMSKeys,
			"hckms":    m.HCKmsKeys,
			"azure_kv": m.AzureKeyVaultKeys,
			"hc_vault": m.VaultKeys,
			"age":      m.AgeKeys,
			"plugins":  m.PluginKeys,
		}
		internalGroup, err := internalGroupFrom(topLevelGroup)
		if err != nil {
			return nil, err
		}
		internalGroups = append(internalGroups, internalGroup)
		return internalGroups, nil
	} else if len(m.KeyGroups) > 0 {
		for _, group := range m.KeyGroups {
			internalGroup, err := internalGroupFrom(group)
			if err != nil {
				return nil, err
			}
			internalGroups = append(internalGroups, internalGroup)
		}
		return internalGroups, nil
	} else {
		return nil, fmt.Errorf("No keys found in file")
	}
}
// ExampleComplexTree is an example sops.Tree object exhibiting complex relationships
var ExampleComplexTree = sops.Tree{
	Branches: sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key:   "hello",
				Value: `Welcome to SOPS! Edit this file as you please!`,
			},
			sops.TreeItem{
				Key:   "example_key",
				Value: "example_value",
			},
			sops.TreeItem{
				Key:   sops.Comment{Value: " Example comment"},
				Value: nil,
			},
			sops.TreeItem{
				Key: "example_array",
				Value: []interface{}{
					"example_value1",
					"example_value2",
				},
			},
			sops.TreeItem{
				Key:   "example_number",
				Value: 1234.56789,
			},
			sops.TreeItem{
				Key:   "example_booleans",
				Value: []interface{}{true, false},
			},
		},
	},
}

// ExampleSimpleTree is an example sops.Tree object exhibiting only simple relationships
// with only one nested branch and only simple string values
var ExampleSimpleTree = sops.Tree{
	Branches: sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key: "Welcome!",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   sops.Comment{Value: " This is an example file."},
						Value: nil,
					},
					sops.TreeItem{
						Key:   "hello",
						Value: "Welcome to SOPS! Edit this file as you please!",
					},
					sops.TreeItem{
						Key:   "example_key",
						Value: "example_value",
					},
				},
			},
		},
	},
}

// ExampleFlatTree is an example sops.Tree object exhibiting only simple relationships
// with no nested branches and only simple string values
var ExampleFlatTree = sops.Tree{
	Branches: sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key:   sops.Comment{Value: " This is an example file."},
				Value: nil,
			},
			sops.TreeItem{
				Key:   "hello",
				Value: "Welcome to SOPS! Edit this file as you please!",
			},
			sops.TreeItem{
				Key:   "example_key",
				Value: "example_value",
			},
			sops.TreeItem{
				Key:   "example_multiline",
				Value: "foo\nbar\nbaz",
			},
		},
	},
}

// HasSopsTopLevelKey returns true if the given branch has a top-level key called "sops".
func HasSopsTopLevelKey(branch sops.TreeBranch) bool {
	for _, b := range branch {
		if b.Key == SopsMetadataKey {
			return true
		}
	}
	return false
}

// IsComplexValue returns true if the given value is an array or dictionary/hash.
func IsComplexValue(v interface{}) bool {
	switch v.(type) {
	case []interface{}:
		return true
	case sops.TreeBranch:
		return true
	}
	return false
}

// ValToString converts a simple value to a string.
// It does not handle complex values (arrays and mappings).
func ValToString(v interface{}) string {
	switch v := v.(type) {
	case float64:
		result := strconv.FormatFloat(v, 'G', -1, 64)
		// If the result can be confused with an integer, make sure we have at least one decimal digit
		if !strings.ContainsRune(result, '.') && !strings.ContainsRune(result, 'E') {
			result = strconv.FormatFloat(v, 'f', 1, 64)
		}
		return result
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(time.RFC3339)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
