/*
Package stores acts as a layer between the internal representation of encrypted files and the encrypted files
themselves.

Subpackages implement serialization and deserialization to multiple formats.

This package defines the structure SOPS files should have and conversions to and from the internal representation. Part
of the purpose of this package is to make it easy to change the SOPS file format while remaining backwards-compatible.
*/
package stores

import (
	"reflect"
	"strconv"
	"time"

	"fmt"

	"github.com/mitchellh/mapstructure"
	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/age"
	"go.mozilla.org/sops/v3/azkv"
	"go.mozilla.org/sops/v3/gcpkms"
	"go.mozilla.org/sops/v3/hcvault"
	"go.mozilla.org/sops/v3/kms"
	"go.mozilla.org/sops/v3/pgp"
)

// SopsFile is a struct used by the stores as a helper to unmarshal the SOPS metadata
type SopsFile struct {
	// Metadata is a pointer so we can easily tell when the field is not present
	// in the SOPS file by checking for nil. This way we can show the user a
	// helpful error message indicating that the metadata wasn't found, instead
	// of showing a cryptic parsing error
	Metadata *Metadata `yaml:"sops" json:"sops" ini:"sops"`
}

// Metadata is stored in SOPS encrypted files, and it contains the information necessary to decrypt the file.
// This struct is just used for serialization, and SOPS uses another struct internally, sops.Metadata. It exists
// in order to allow the binary format to stay backwards compatible over time, but at the same time allow the internal
// representation SOPS uses to change over time.
type Metadata struct {
	ShamirThreshold           int         `yaml:"shamir_threshold,omitempty" json:"shamir_threshold,omitempty" mapstructure:"shamir_threshold,omitempty"`
	KeyGroups                 []keygroup  `yaml:"key_groups,omitempty" json:"key_groups,omitempty" mapstructure:"key_groups,omitempty"`
	KMSKeys                   []kmskey    `yaml:"kms" json:"kms" mapstructure:"kms"`
	GCPKMSKeys                []gcpkmskey `yaml:"gcp_kms" json:"gcp_kms" mapstructure:"gcp_kms"`
	AzureKeyVaultKeys         []azkvkey   `yaml:"azure_kv" json:"azure_kv" mapstructure:"azure_kv"`
	VaultKeys                 []vaultkey  `yaml:"hc_vault" json:"hc_vault" mapstructure:"hc_vault"`
	AgeKeys                   []agekey    `yaml:"age" json:"age" mapstructure:"age"`
	LastModified              string      `yaml:"lastmodified" json:"lastmodified" mapstructure:"lastmodified"`
	MessageAuthenticationCode string      `yaml:"mac" json:"mac" mapstructure:"mac"`
	PGPKeys                   []pgpkey    `yaml:"pgp" json:"pgp" mapstructure:"pgp"`
	UnencryptedSuffix         string      `yaml:"unencrypted_suffix,omitempty" json:"unencrypted_suffix,omitempty" mapstructure:"unencrypted_suffix,omitempty"`
	EncryptedSuffix           string      `yaml:"encrypted_suffix,omitempty" json:"encrypted_suffix,omitempty" mapstructure:"encrypted_suffix,omitempty"`
	UnencryptedRegex          string      `yaml:"unencrypted_regex,omitempty" json:"unencrypted_regex,omitempty" mapstructure:"unencrypted_regex,omitempty"`
	EncryptedRegex            string      `yaml:"encrypted_regex,omitempty" json:"encrypted_regex,omitempty" mapstructure:"encrypted_regex,omitempty"`
	Version                   string      `yaml:"version" json:"version" mapstructure:"version"`
}

type keygroup struct {
	PGPKeys           []pgpkey    `yaml:"pgp,omitempty" json:"pgp,omitempty" mapstructure:"pgp,omitempty"`
	KMSKeys           []kmskey    `yaml:"kms,omitempty" json:"kms,omitempty" mapstructure:"kms,omitempty"`
	GCPKMSKeys        []gcpkmskey `yaml:"gcp_kms,omitempty" json:"gcp_kms,omitempty" mapstructure:"gcp_kms,omitempty"`
	AzureKeyVaultKeys []azkvkey   `yaml:"azure_kv,omitempty" json:"azure_kv,omitempty" mapstructure:"azure_kv,omitempty"`
	VaultKeys         []vaultkey  `yaml:"hc_vault" json:"hc_vault" mapstructure:"hc_vault"`
	AgeKeys           []agekey    `yaml:"age" json:"age" mapstructure:"age"`
}

type pgpkey struct {
	CreatedAt        string `yaml:"created_at" json:"created_at" mapstructure:"created_at"`
	EncryptedDataKey string `yaml:"enc" json:"enc" mapstructure:"enc"`
	Fingerprint      string `yaml:"fp" json:"fp" mapstructure:"fp"`
}

type kmskey struct {
	Arn              string             `yaml:"arn" json:"arn" mapstructure:"arn"`
	Role             string             `yaml:"role,omitempty" json:"role,omitempty" mapstructure:"role,omitempty"`
	Context          map[string]*string `yaml:"context,omitempty" json:"context,omitempty" mapstructure:"context,omitempty"`
	CreatedAt        string             `yaml:"created_at" json:"created_at" mapstructure:"created_at"`
	EncryptedDataKey string             `yaml:"enc" json:"enc" mapstructure:"enc"`
	AwsProfile       string             `yaml:"aws_profile" json:"aws_profile" mapstructure:"aws_profile"`
}

type gcpkmskey struct {
	ResourceID       string `yaml:"resource_id" json:"resource_id" mapstructure:"resource_id"`
	CreatedAt        string `yaml:"created_at" json:"created_at" mapstructure:"created_at"`
	EncryptedDataKey string `yaml:"enc" json:"enc" mapstructure:"enc"`
}

type vaultkey struct {
	VaultAddress     string `yaml:"vault_address" json:"vault_address" mapstructure:"vault_address"`
	EnginePath       string `yaml:"engine_path" json:"engine_path" mapstructure:"engine_path"`
	KeyName          string `yaml:"key_name" json:"key_name" mapstructure:"key_name"`
	CreatedAt        string `yaml:"created_at" json:"created_at" mapstructure:"created_at"`
	EncryptedDataKey string `yaml:"enc" json:"enc" mapstructure:"enc"`
}

type azkvkey struct {
	VaultURL         string `yaml:"vault_url" json:"vault_url" mapstructure:"vault_url"`
	Name             string `yaml:"name" json:"name" mapstructure:"name"`
	Version          string `yaml:"version" json:"version" mapstructure:"version"`
	CreatedAt        string `yaml:"created_at" json:"created_at" mapstructure:"created_at"`
	EncryptedDataKey string `yaml:"enc" json:"enc" mapstructure:"enc"`
}

type agekey struct {
	Recipient        string `yaml:"recipient" json:"recipient" mapstructure:"recipient"`
	EncryptedDataKey string `yaml:"enc" json:"enc" mapstructure:"enc"`
}

// MetadataFromInternal converts an internal SOPS metadata representation to a representation appropriate for storage
func MetadataFromInternal(sopsMetadata sops.Metadata) Metadata {
	var m Metadata
	m.LastModified = sopsMetadata.LastModified.Format(time.RFC3339)
	m.UnencryptedSuffix = sopsMetadata.UnencryptedSuffix
	m.EncryptedSuffix = sopsMetadata.EncryptedSuffix
	m.UnencryptedRegex = sopsMetadata.UnencryptedRegex
	m.EncryptedRegex = sopsMetadata.EncryptedRegex
	m.MessageAuthenticationCode = sopsMetadata.MessageAuthenticationCode
	m.Version = sopsMetadata.Version
	m.ShamirThreshold = sopsMetadata.ShamirThreshold
	if len(sopsMetadata.KeyGroups) == 1 {
		group := sopsMetadata.KeyGroups[0]
		m.PGPKeys = pgpKeysFromGroup(group)
		m.KMSKeys = kmsKeysFromGroup(group)
		m.GCPKMSKeys = gcpkmsKeysFromGroup(group)
		m.VaultKeys = vaultKeysFromGroup(group)
		m.AzureKeyVaultKeys = azkvKeysFromGroup(group)
		m.AgeKeys = ageKeysFromGroup(group)
	} else {
		for _, group := range sopsMetadata.KeyGroups {
			m.KeyGroups = append(m.KeyGroups, keygroup{
				KMSKeys:           kmsKeysFromGroup(group),
				PGPKeys:           pgpKeysFromGroup(group),
				GCPKMSKeys:        gcpkmsKeysFromGroup(group),
				VaultKeys:         vaultKeysFromGroup(group),
				AzureKeyVaultKeys: azkvKeysFromGroup(group),
				AgeKeys:           ageKeysFromGroup(group),
			})
		}
	}
	return m
}

func pgpKeysFromGroup(group sops.KeyGroup) (keys []pgpkey) {
	for _, key := range group {
		switch key := key.(type) {
		case *pgp.MasterKey:
			keys = append(keys, pgpkey{
				Fingerprint:      key.Fingerprint,
				EncryptedDataKey: key.EncryptedKey,
				CreatedAt:        key.CreationDate.Format(time.RFC3339),
			})
		}
	}
	return
}

func kmsKeysFromGroup(group sops.KeyGroup) (keys []kmskey) {
	for _, key := range group {
		switch key := key.(type) {
		case *kms.MasterKey:
			keys = append(keys, kmskey{
				Arn:              key.Arn,
				CreatedAt:        key.CreationDate.Format(time.RFC3339),
				EncryptedDataKey: key.EncryptedKey,
				Context:          key.EncryptionContext,
				Role:             key.Role,
				AwsProfile:       key.AwsProfile,
			})
		}
	}
	return
}

func gcpkmsKeysFromGroup(group sops.KeyGroup) (keys []gcpkmskey) {
	for _, key := range group {
		switch key := key.(type) {
		case *gcpkms.MasterKey:
			keys = append(keys, gcpkmskey{
				ResourceID:       key.ResourceID,
				CreatedAt:        key.CreationDate.Format(time.RFC3339),
				EncryptedDataKey: key.EncryptedKey,
			})
		}
	}
	return
}

func vaultKeysFromGroup(group sops.KeyGroup) (keys []vaultkey) {
	for _, key := range group {
		switch key := key.(type) {
		case *hcvault.MasterKey:
			keys = append(keys, vaultkey{
				VaultAddress:     key.VaultAddress,
				EnginePath:       key.EnginePath,
				KeyName:          key.KeyName,
				CreatedAt:        key.CreationDate.Format(time.RFC3339),
				EncryptedDataKey: key.EncryptedKey,
			})
		}
	}
	return
}

func azkvKeysFromGroup(group sops.KeyGroup) (keys []azkvkey) {
	for _, key := range group {
		switch key := key.(type) {
		case *azkv.MasterKey:
			keys = append(keys, azkvkey{
				VaultURL:         key.VaultURL,
				Name:             key.Name,
				Version:          key.Version,
				CreatedAt:        key.CreationDate.Format(time.RFC3339),
				EncryptedDataKey: key.EncryptedKey,
			})
		}
	}
	return
}

func ageKeysFromGroup(group sops.KeyGroup) (keys []agekey) {
	for _, key := range group {
		switch key := key.(type) {
		case *age.MasterKey:
			keys = append(keys, agekey{
				Recipient:        key.Recipient,
				EncryptedDataKey: key.EncryptedKey,
			})
		}
	}
	return
}

// ToInternal converts a storage-appropriate Metadata struct to a SOPS internal representation
func (m *Metadata) ToInternal() (sops.Metadata, error) {
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

	if cryptRuleCount > 1 {
		return sops.Metadata{}, fmt.Errorf("Cannot use more than one of encrypted_suffix, unencrypted_suffix, encrypted_regex or unencrypted_regex in the same file")
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
		LastModified:              lastModified,
	}, nil
}

func internalGroupFrom(kmsKeys []kmskey, pgpKeys []pgpkey, gcpKmsKeys []gcpkmskey, azkvKeys []azkvkey, vaultKeys []vaultkey, ageKeys []agekey) (sops.KeyGroup, error) {
	var internalGroup sops.KeyGroup
	for _, kmsKey := range kmsKeys {
		k, err := kmsKey.toInternal()
		if err != nil {
			return nil, err
		}
		internalGroup = append(internalGroup, k)
	}
	for _, gcpKmsKey := range gcpKmsKeys {
		k, err := gcpKmsKey.toInternal()
		if err != nil {
			return nil, err
		}
		internalGroup = append(internalGroup, k)
	}
	for _, azkvKey := range azkvKeys {
		k, err := azkvKey.toInternal()
		if err != nil {
			return nil, err
		}
		internalGroup = append(internalGroup, k)
	}
	for _, vaultKey := range vaultKeys {
		k, err := vaultKey.toInternal()
		if err != nil {
			return nil, err
		}
		internalGroup = append(internalGroup, k)
	}
	for _, pgpKey := range pgpKeys {
		k, err := pgpKey.toInternal()
		if err != nil {
			return nil, err
		}
		internalGroup = append(internalGroup, k)
	}
	for _, ageKey := range ageKeys {
		k, err := ageKey.toInternal()
		if err != nil {
			return nil, err
		}
		internalGroup = append(internalGroup, k)
	}
	return internalGroup, nil
}

func (m *Metadata) internalKeygroups() ([]sops.KeyGroup, error) {
	var internalGroups []sops.KeyGroup
	if len(m.PGPKeys) > 0 || len(m.KMSKeys) > 0 || len(m.GCPKMSKeys) > 0 || len(m.AzureKeyVaultKeys) > 0 || len(m.VaultKeys) > 0 || len(m.AgeKeys) > 0 {
		internalGroup, err := internalGroupFrom(m.KMSKeys, m.PGPKeys, m.GCPKMSKeys, m.AzureKeyVaultKeys, m.VaultKeys, m.AgeKeys)
		if err != nil {
			return nil, err
		}
		internalGroups = append(internalGroups, internalGroup)
		return internalGroups, nil
	} else if len(m.KeyGroups) > 0 {
		for _, group := range m.KeyGroups {
			internalGroup, err := internalGroupFrom(group.KMSKeys, group.PGPKeys, group.GCPKMSKeys, group.AzureKeyVaultKeys, group.VaultKeys, group.AgeKeys)
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

func (kmsKey *kmskey) toInternal() (*kms.MasterKey, error) {
	creationDate, err := time.Parse(time.RFC3339, kmsKey.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &kms.MasterKey{
		Role:              kmsKey.Role,
		EncryptionContext: kmsKey.Context,
		EncryptedKey:      kmsKey.EncryptedDataKey,
		CreationDate:      creationDate,
		Arn:               kmsKey.Arn,
		AwsProfile:        kmsKey.AwsProfile,
	}, nil
}

func (gcpKmsKey *gcpkmskey) toInternal() (*gcpkms.MasterKey, error) {
	creationDate, err := time.Parse(time.RFC3339, gcpKmsKey.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &gcpkms.MasterKey{
		ResourceID:   gcpKmsKey.ResourceID,
		EncryptedKey: gcpKmsKey.EncryptedDataKey,
		CreationDate: creationDate,
	}, nil
}

func (azkvKey *azkvkey) toInternal() (*azkv.MasterKey, error) {
	creationDate, err := time.Parse(time.RFC3339, azkvKey.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &azkv.MasterKey{
		VaultURL:     azkvKey.VaultURL,
		Name:         azkvKey.Name,
		Version:      azkvKey.Version,
		EncryptedKey: azkvKey.EncryptedDataKey,
		CreationDate: creationDate,
	}, nil
}

func (vaultKey *vaultkey) toInternal() (*hcvault.MasterKey, error) {
	creationDate, err := time.Parse(time.RFC3339, vaultKey.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &hcvault.MasterKey{
		VaultAddress: vaultKey.VaultAddress,
		EnginePath:   vaultKey.EnginePath,
		KeyName:      vaultKey.KeyName,
		CreationDate: creationDate,
		EncryptedKey: vaultKey.EncryptedDataKey,
	}, nil
}

func (pgpKey *pgpkey) toInternal() (*pgp.MasterKey, error) {
	creationDate, err := time.Parse(time.RFC3339, pgpKey.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &pgp.MasterKey{
		EncryptedKey: pgpKey.EncryptedDataKey,
		CreationDate: creationDate,
		Fingerprint:  pgpKey.Fingerprint,
	}, nil
}

func (ageKey *agekey) toInternal() (*age.MasterKey, error) {
	return &age.MasterKey{
		EncryptedKey: ageKey.EncryptedDataKey,
		Recipient:    ageKey.Recipient,
	}, nil
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

// ConvertStructToMap recursively converts a structure to a map[string]interface{} representation while
// respecting all mapstructure tags on the source structure. This is useful when converting complex structures
// to a map suitable for use with the Flatten function.
//
// Note: this will only emit the public fields of a structure, private fields are ignored entirely.
func ConvertStructToMap(input interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := mapstructure.Decode(input, &result)
	if err != nil {
		return nil, fmt.Errorf("decode struct: %w", err)
	}

	// Mapstructure stops when the output interface is satisfied, in our case we need to delve further into
	// any collections and ensure that all structures are converted to their map representations.
	for k, v := range result {
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Array:
		case reflect.Slice:
			elemType := val.Type().Elem()
			// Ignore any elements that are already primitive types
			if elemType.Kind() != reflect.Interface &&
				elemType.Kind() != reflect.Struct {
				continue
			}

			newList := make([]interface{}, val.Len())
			for j := 0; j < val.Len(); j++ {
				newVal, err := ConvertStructToMap(val.Index(j).Interface())
				if err != nil {
					return nil, fmt.Errorf("convert array field to map: %w", err)
				}

				newList[j] = newVal
			}
			result[k] = newList
		case reflect.Map:
			elemType := val.Type().Elem()
			// Ignore any elements that are already primitive types
			if elemType.Kind() != reflect.Interface &&
				elemType.Kind() != reflect.Struct {
				continue
			}

			// Non-string keys
			if val.Type().Key().Kind() != reflect.String {
				return nil, fmt.Errorf("field '%s' is invalid, only map fields with string keys are supported", k)
			}

			newMap := map[string]interface{}{}
			for _, key := range val.MapKeys() {
				newVal, err := ConvertStructToMap(val.MapIndex(key).Interface())
				if err != nil {
					return nil, fmt.Errorf("convert array field to map: %w", err)
				}

				newMap[key.String()] = newVal
			}
			result[k] = newMap
		}
	}

	return result, nil
}

// ValueToString converts the input value to a string representation. This is useful when encoding data to plain
// text formats as is done in the ini and dotenv stores.
func ValueToString(v interface{}) string {
	switch v := v.(type) {
	case fmt.Stringer:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', 6, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
