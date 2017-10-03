/*
Package stores acts as a layer between the internal representation of encrypted files and the encrypted files
themselves.

Subpackages implement serialization and deserialization to multiple formats.

This package defines the structure SOPS files should have and conversions to and from the internal representation. Part
of the purpose of this package is to make it easy to change the SOPS file format while remaining backwards-compatible.
*/
package stores

import (
	"time"

	"fmt"

	"go.mozilla.org/sops"
	"go.mozilla.org/sops/gcpkms"
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
)

// SopsFile is a struct used by the stores as a helper to unmarshal the SOPS metadata
type SopsFile struct {
	// Metadata is a pointer so we can easily tell when the field is not present
	// in the SOPS file by checking for nil. This way we can show the user a
	// helpful error message indicating that the metadata wasn't found, instead
	// of showing a cryptic parsing error
	Metadata *Metadata `yaml:"sops" json:"sops"`
}

// Metadata is stored in SOPS encrypted files, and it contains the information necessary to decrypt the file.
// This struct is just used for serialization, and SOPS uses another struct internally, sops.Metadata. It exists
// in order to allow the binary format to stay backwards compatible over time, but at the same time allow the internal
// representation SOPS uses to change over time.
type Metadata struct {
	ShamirThreshold           int         `yaml:"shamir_threshold,omitempty" json:"shamir_threshold,omitempty"`
	KeyGroups                 []keygroup  `yaml:"key_groups,omitempty" json:"key_groups,omitempty"`
	KMSKeys                   []kmskey    `yaml:"kms" json:"kms"`
	GCPKMSKeys                []gcpkmskey `yaml:"gcp_kms" json:"gcp_kms"`
	LastModified              string      `yaml:"lastmodified" json:"lastmodified"`
	MessageAuthenticationCode string      `yaml:"mac" json:"mac"`
	PGPKeys                   []pgpkey    `yaml:"pgp" json:"pgp"`
	UnencryptedSuffix         string      `yaml:"unencrypted_suffix" json:"unencrypted_suffix"`
	Version                   string      `yaml:"version" json:"version"`
}

type keygroup struct {
	PGPKeys    []pgpkey    `yaml:"pgp,omitempty" json:"pgp,omitempty"`
	KMSKeys    []kmskey    `yaml:"kms,omitempty" json:"kms,omitempty"`
	GCPKMSKeys []gcpkmskey `yaml:"gcp_kms,omitempty" json:"gcp_kms,omitempty"`
}

type pgpkey struct {
	CreatedAt        string `yaml:"created_at" json:"created_at"`
	EncryptedDataKey string `yaml:"enc" json:"enc"`
	Fingerprint      string `yaml:"fp" json:"fp"`
}

type kmskey struct {
	Arn              string             `yaml:"arn" json:"arn"`
	Role             string             `yaml:"role,omitempty" json:"role,omitempty"`
	Context          map[string]*string `yaml:"context,omitempty" json:"context,omitempty"`
	CreatedAt        string             `yaml:"created_at" json:"created_at"`
	EncryptedDataKey string             `yaml:"enc" json:"enc"`
}

type gcpkmskey struct {
	ResourceID       string `yaml:"resource_id" json:"resource_id"`
	CreatedAt        string `yaml:"created_at" json:"created_at"`
	EncryptedDataKey string `yaml:"enc" json:"enc"`
}

// MetadataFromInternal converts an internal SOPS metadata representation to a representation appropriate for storage
func MetadataFromInternal(sopsMetadata sops.Metadata) Metadata {
	var m Metadata
	m.LastModified = sopsMetadata.LastModified.Format(time.RFC3339)
	m.UnencryptedSuffix = sopsMetadata.UnencryptedSuffix
	m.MessageAuthenticationCode = sopsMetadata.MessageAuthenticationCode
	m.Version = sopsMetadata.Version
	m.ShamirThreshold = sopsMetadata.ShamirThreshold
	if len(sopsMetadata.KeyGroups) == 1 {
		group := sopsMetadata.KeyGroups[0]
		m.PGPKeys = pgpKeysFromGroup(group)
		m.KMSKeys = kmsKeysFromGroup(group)
		m.GCPKMSKeys = gcpkmsKeysFromGroup(group)
	} else {
		for _, group := range sopsMetadata.KeyGroups {
			m.KeyGroups = append(m.KeyGroups, keygroup{
				KMSKeys:    kmsKeysFromGroup(group),
				PGPKeys:    pgpKeysFromGroup(group),
				GCPKMSKeys: gcpkmsKeysFromGroup(group),
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
	if m.UnencryptedSuffix == "" {
		m.UnencryptedSuffix = sops.DefaultUnencryptedSuffix
	}
	return sops.Metadata{
		KeyGroups:                 groups,
		ShamirThreshold:           m.ShamirThreshold,
		Version:                   m.Version,
		MessageAuthenticationCode: m.MessageAuthenticationCode,
		UnencryptedSuffix:         m.UnencryptedSuffix,
		LastModified:              lastModified,
	}, nil
}

func internalGroupFrom(kmsKeys []kmskey, pgpKeys []pgpkey, gcpKmsKeys []gcpkmskey) (sops.KeyGroup, error) {
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
	for _, pgpKey := range pgpKeys {
		k, err := pgpKey.toInternal()
		if err != nil {
			return nil, err
		}
		internalGroup = append(internalGroup, k)
	}
	return internalGroup, nil
}

func (m *Metadata) internalKeygroups() ([]sops.KeyGroup, error) {
	var internalGroups []sops.KeyGroup
	if len(m.PGPKeys) > 0 || len(m.KMSKeys) > 0 || len(m.GCPKMSKeys) > 0 {
		internalGroup, err := internalGroupFrom(m.KMSKeys, m.PGPKeys, m.GCPKMSKeys)
		if err != nil {
			return nil, err
		}
		internalGroups = append(internalGroups, internalGroup)
		return internalGroups, nil
	} else if len(m.KeyGroups) > 0 {
		for _, group := range m.KeyGroups {
			internalGroup, err := internalGroupFrom(group.KMSKeys, group.PGPKeys, group.GCPKMSKeys)
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
