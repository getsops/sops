package stores

import (
	"testing"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/azkv"
	"github.com/getsops/sops/v3/gcpkms"
	"github.com/getsops/sops/v3/hckms"
	"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/pgp"
	"github.com/stretchr/testify/assert"
)

func TestExtractMetadata(t *testing.T) {
	var (
		broken1 = sops.TreeBranch{
			sops.TreeItem{
				Key:   "foo",
				Value: "bar",
			},
		}
		broken2 = sops.TreeBranch{
			sops.TreeItem{
				Key:   "sops",
				Value: "foo",
			},
		}
		empty = sops.TreeBranch{
			sops.TreeItem{
				Key:   "sops",
				Value: sops.TreeBranch{},
			},
		}

		minimal1 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "LastModified",
						Value: "2025-03-23T10:20:30Z",
					},
				},
			},
		}
		minimal2 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "1234",
								},
							},
						},
					},
				},
			},
		}

		multiple = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "unencrypted_suffix",
						Value: "foo",
					},
					sops.TreeItem{
						Key:   "encrypted_suffix",
						Value: "bar",
					},
					sops.TreeItem{
						Key:   "unencrypted_regex",
						Value: "baz",
					},
					sops.TreeItem{
						Key:   "encrypted_regex",
						Value: "bam",
					},
					sops.TreeItem{
						Key:   "unencrypted_comment_regex",
						Value: "foobar",
					},
					sops.TreeItem{
						Key:   "encrypted_comment_regex",
						Value: "bazbam",
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "1234",
								},
							},
						},
					},
				},
			},
		}

		single1 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "unencrypted_suffix",
						Value: "foo",
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "1234",
								},
							},
						},
					},
				},
			},
		}
		single2 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "encrypted_suffix",
						Value: "bar",
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "1234",
								},
							},
						},
					},
				},
			},
		}
		single3 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "unencrypted_regex",
						Value: "baz",
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "1234",
								},
							},
						},
					},
				},
			},
		}
		single4 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "encrypted_regex",
						Value: "bam",
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "1234",
								},
							},
						},
					},
				},
			},
		}
		single5 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "unencrypted_comment_regex",
						Value: "foobar",
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "1234",
								},
							},
						},
					},
				},
			},
		}
		single6 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "encrypted_comment_regex",
						Value: "bazbam",
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "1234",
								},
							},
						},
					},
				},
			},
		}

		completeKeyGroup = sops.TreeBranch{
			sops.TreeItem{
				Key: "kms",
				Value: []interface{}{
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "created_at",
							Value: "2025-03-23T10:20:29Z",
						},
						sops.TreeItem{
							Key:   "enc",
							Value: "ABCD AWS KMS (inner)",
						},
						sops.TreeItem{
							Key:   "arn",
							Value: "AWS KMS ARN (inner)",
						},
						sops.TreeItem{
							Key:   "role",
							Value: "AWS KMS role (inner)",
						},
						sops.TreeItem{
							Key: "context",
							Value: sops.TreeBranch{
								sops.TreeItem{
									Key:   "foo",
									Value: "bar",
								},
							},
						},
						sops.TreeItem{
							Key:   "aws_profile",
							Value: "AWS KMS profile (inner)",
						},
					},
				},
			},
			sops.TreeItem{
				Key: "gcp_kms",
				Value: []interface{}{
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "created_at",
							Value: "2025-03-23T10:20:29Z",
						},
						sops.TreeItem{
							Key:   "enc",
							Value: "ABCD GCP KMS (inner)",
						},
						sops.TreeItem{
							Key:   "resource_id",
							Value: "GCP KMS resource ID (inner)",
						},
					},
				},
			},
			sops.TreeItem{
				Key: "hckms",
				Value: []interface{}{
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "created_at",
							Value: "2025-03-23T10:20:29Z",
						},
						sops.TreeItem{
							Key:   "enc",
							Value: "ABCD HC KMS (inner)",
						},
						sops.TreeItem{
							Key:   "key_id",
							Value: "HC KMS (inner):key ID (inner)",
						},
					},
				},
			},
			sops.TreeItem{
				Key: "azure_kv",
				Value: []interface{}{
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "created_at",
							Value: "2025-03-23T10:20:29Z",
						},
						sops.TreeItem{
							Key:   "enc",
							Value: "ABCD AZKV (inner)",
						},
						sops.TreeItem{
							Key:   "vault_url",
							Value: "AZKV vault URL (inner)",
						},
						sops.TreeItem{
							Key:   "name",
							Value: "AZKV name (inner)",
						},
						sops.TreeItem{
							Key:   "version",
							Value: "AZKV version (inner)",
						},
					},
				},
			},
			sops.TreeItem{
				Key: "hc_vault",
				Value: []interface{}{
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "created_at",
							Value: "2025-03-23T10:20:29Z",
						},
						sops.TreeItem{
							Key:   "enc",
							Value: "ABCD HC Vault (inner)",
						},
						sops.TreeItem{
							Key:   "vault_address",
							Value: "HC Vault address (inner)",
						},
						sops.TreeItem{
							Key:   "engine_path",
							Value: "HC Vault engine path (inner)",
						},
						sops.TreeItem{
							Key:   "key_name",
							Value: "HC Vault key name (inner)",
						},
					},
				},
			},
			sops.TreeItem{
				Key: "agekey",
				Value: []interface{}{
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "recipient",
							Value: "age recipient (inner)",
						},
						sops.TreeItem{
							Key:   "enc",
							Value: "ABCD age (inner)",
						},
					},
				},
			},
			sops.TreeItem{
				Key: "pgp",
				Value: []interface{}{
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "created_at",
							Value: "2025-03-23T10:20:29Z",
						},
						sops.TreeItem{
							Key:   "enc",
							Value: "ABCD PGP (inner)",
						},
						sops.TreeItem{
							Key:   "fp",
							Value: "PGP fingerprint (inner)",
						},
					},
				},
			},
		}

		everything1 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "shamir_threshold",
						Value: 2,
					},
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "mac",
						Value: "asdf",
					},
					sops.TreeItem{
						Key:   "encrypted_comment_regex",
						Value: "bazbam",
					},
					sops.TreeItem{
						Key:   "mac_only_encrypted",
						Value: true,
					},
					sops.TreeItem{
						Key:   "version",
						Value: "barbaz",
					},
					sops.TreeItem{
						Key: "kms",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD AWS KMS",
								},
								sops.TreeItem{
									Key:   "arn",
									Value: "AWS KMS ARN",
								},
								sops.TreeItem{
									Key:   "role",
									Value: "AWS KMS role",
								},
								sops.TreeItem{
									Key: "context",
									Value: sops.TreeBranch{
										sops.TreeItem{
											Key:   "foo",
											Value: "bar",
										},
									},
								},
								sops.TreeItem{
									Key:   "aws_profile",
									Value: "AWS KMS profile",
								},
							},
						},
					},
					sops.TreeItem{
						Key: "gcp_kms",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD GCP KMS",
								},
								sops.TreeItem{
									Key:   "resource_id",
									Value: "GCP KMS resource ID",
								},
							},
						},
					},
					sops.TreeItem{
						Key: "hckms",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD HC KMS",
								},
								sops.TreeItem{
									Key:   "key_id",
									Value: "HC KMS:key ID",
								},
							},
						},
					},
					sops.TreeItem{
						Key: "azure_kv",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD AZKV",
								},
								sops.TreeItem{
									Key:   "vault_url",
									Value: "AZKV vault URL",
								},
								sops.TreeItem{
									Key:   "name",
									Value: "AZKV name",
								},
								sops.TreeItem{
									Key:   "version",
									Value: "AZKV version",
								},
							},
						},
					},
					sops.TreeItem{
						Key: "hc_vault",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD HC Vault",
								},
								sops.TreeItem{
									Key:   "vault_address",
									Value: "HC Vault address",
								},
								sops.TreeItem{
									Key:   "engine_path",
									Value: "HC Vault engine path",
								},
								sops.TreeItem{
									Key:   "key_name",
									Value: "HC Vault key name",
								},
							},
						},
					},
					sops.TreeItem{
						Key: "agekey",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "recipient",
									Value: "age recipient",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD age",
								},
							},
						},
					},
					sops.TreeItem{
						Key: "pgp",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{
									Key:   "created_at",
									Value: "2025-03-23T10:20:29Z",
								},
								sops.TreeItem{
									Key:   "enc",
									Value: "ABCD PGP",
								},
								sops.TreeItem{
									Key:   "fp",
									Value: "PGP fingerprint",
								},
							},
						},
					},
					// This will be ignored:
					sops.TreeItem{
						Key: "key_groups",
						Value: []interface{}{
							completeKeyGroup,
						},
					},
				},
			},
		}
		everything2 = sops.TreeBranch{
			sops.TreeItem{
				Key: "sops",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "shamir_threshold",
						Value: 2,
					},
					sops.TreeItem{
						Key:   "lastmodified",
						Value: "2025-03-23T10:20:30Z",
					},
					sops.TreeItem{
						Key:   "mac",
						Value: "asdf",
					},
					sops.TreeItem{
						Key:   "encrypted_comment_regex",
						Value: "bazbam",
					},
					sops.TreeItem{
						Key:   "mac_only_encrypted",
						Value: true,
					},
					sops.TreeItem{
						Key:   "version",
						Value: "barbaz",
					},
					sops.TreeItem{
						Key: "key_groups",
						Value: []interface{}{
							completeKeyGroup,
						},
					},
				},
			},
		}
	)

	branches, metadata, err := ExtractMetadata([]sops.TreeBranch{broken1}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.NotNil(t, err)
	assert.Equal(t, sops.MetadataNotFound, err)
	assert.Nil(t, branches)
	assert.Equal(t, sops.Metadata{}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{broken2}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.NotNil(t, err)
	assert.Equal(t, "Found sops entry that is not a mapping", err.Error())
	assert.Nil(t, branches)
	assert.Equal(t, sops.Metadata{}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{empty}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.NotNil(t, err)
	assert.Equal(t, "parsing time \"\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"2006\"", err.Error())
	assert.Nil(t, branches)
	assert.Equal(t, sops.Metadata{}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{minimal1}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.NotNil(t, err)
	assert.Equal(t, "No keys found in file", err.Error())
	assert.Nil(t, branches)
	assert.Equal(t, sops.Metadata{}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{minimal2}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	assert.Equal(t, sops.Metadata{
		LastModified:      time.Unix(1742725230, 0).UTC(),
		UnencryptedSuffix: "_unencrypted",
		KeyGroups: []sops.KeyGroup{
			{
				&pgp.MasterKey{
					Fingerprint:  "1234",
					EncryptedKey: "ABCD",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{multiple}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.NotNil(t, err)
	assert.Equal(t, "Cannot use more than one of encrypted_suffix, unencrypted_suffix, encrypted_regex, unencrypted_regex, encrypted_comment_regex, or unencrypted_comment_regex in the same file", err.Error())
	assert.Nil(t, branches)
	assert.Equal(t, sops.Metadata{}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{single1}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	assert.Equal(t, sops.Metadata{
		LastModified:      time.Unix(1742725230, 0).UTC(),
		UnencryptedSuffix: "foo",
		KeyGroups: []sops.KeyGroup{
			{
				&pgp.MasterKey{
					Fingerprint:  "1234",
					EncryptedKey: "ABCD",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{single2}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	assert.Equal(t, sops.Metadata{
		LastModified:    time.Unix(1742725230, 0).UTC(),
		EncryptedSuffix: "bar",
		KeyGroups: []sops.KeyGroup{
			{
				&pgp.MasterKey{
					Fingerprint:  "1234",
					EncryptedKey: "ABCD",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{single3}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	assert.Equal(t, sops.Metadata{
		LastModified:     time.Unix(1742725230, 0).UTC(),
		UnencryptedRegex: "baz",
		KeyGroups: []sops.KeyGroup{
			{
				&pgp.MasterKey{
					Fingerprint:  "1234",
					EncryptedKey: "ABCD",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{single4}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	assert.Equal(t, sops.Metadata{
		LastModified:   time.Unix(1742725230, 0).UTC(),
		EncryptedRegex: "bam",
		KeyGroups: []sops.KeyGroup{
			{
				&pgp.MasterKey{
					Fingerprint:  "1234",
					EncryptedKey: "ABCD",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{single5}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	assert.Equal(t, sops.Metadata{
		LastModified:            time.Unix(1742725230, 0).UTC(),
		UnencryptedCommentRegex: "foobar",
		KeyGroups: []sops.KeyGroup{
			{
				&pgp.MasterKey{
					Fingerprint:  "1234",
					EncryptedKey: "ABCD",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{single6}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	assert.Equal(t, sops.Metadata{
		LastModified:          time.Unix(1742725230, 0).UTC(),
		EncryptedCommentRegex: "bazbam",
		KeyGroups: []sops.KeyGroup{
			{
				&pgp.MasterKey{
					Fingerprint:  "1234",
					EncryptedKey: "ABCD",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{everything1}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	bar := "bar"
	assert.Equal(t, sops.Metadata{
		ShamirThreshold:           2,
		LastModified:              time.Unix(1742725230, 0).UTC(),
		MessageAuthenticationCode: "asdf",
		EncryptedCommentRegex:     "bazbam",
		MACOnlyEncrypted:          true,
		Version:                   "barbaz",
		KeyGroups: []sops.KeyGroup{
			{
				&kms.MasterKey{
					Arn:          "AWS KMS ARN",
					Role:         "AWS KMS role",
					EncryptedKey: "ABCD AWS KMS",
					CreationDate: time.Unix(1742725229, 0).UTC(),
					EncryptionContext: map[string]*string{
						"foo": &bar,
					},
					AwsProfile: "AWS KMS profile",
				},
				&gcpkms.MasterKey{
					ResourceID:   "GCP KMS resource ID",
					EncryptedKey: "ABCD GCP KMS",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
				&hckms.MasterKey{
					KeyID:        "HC KMS:key ID",
					Region:       "HC KMS",
					KeyUUID:      "key ID",
					EncryptedKey: "ABCD HC KMS",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
				&azkv.MasterKey{
					VaultURL:     "AZKV vault URL",
					Name:         "AZKV name",
					Version:      "AZKV version",
					EncryptedKey: "ABCD AZKV",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
				&hcvault.MasterKey{
					VaultAddress: "HC Vault address",
					EnginePath:   "HC Vault engine path",
					KeyName:      "HC Vault key name",
					EncryptedKey: "ABCD HC Vault",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
				&pgp.MasterKey{
					Fingerprint:  "PGP fingerprint",
					EncryptedKey: "ABCD PGP",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)

	branches, metadata, err = ExtractMetadata([]sops.TreeBranch{everything2}, MetadataOpts{Flatten: MetadataFlattenNone})
	assert.Nil(t, err)
	assert.Equal(t, sops.Metadata{
		ShamirThreshold:           2,
		LastModified:              time.Unix(1742725230, 0).UTC(),
		MessageAuthenticationCode: "asdf",
		EncryptedCommentRegex:     "bazbam",
		MACOnlyEncrypted:          true,
		Version:                   "barbaz",
		KeyGroups: []sops.KeyGroup{
			{
				&kms.MasterKey{
					Arn:          "AWS KMS ARN (inner)",
					Role:         "AWS KMS role (inner)",
					EncryptedKey: "ABCD AWS KMS (inner)",
					CreationDate: time.Unix(1742725229, 0).UTC(),
					EncryptionContext: map[string]*string{
						"foo": &bar,
					},
					AwsProfile: "AWS KMS profile (inner)",
				},
				&gcpkms.MasterKey{
					ResourceID:   "GCP KMS resource ID (inner)",
					EncryptedKey: "ABCD GCP KMS (inner)",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
				&hckms.MasterKey{
					KeyID:        "HC KMS (inner):key ID (inner)",
					Region:       "HC KMS (inner)",
					KeyUUID:      "key ID (inner)",
					EncryptedKey: "ABCD HC KMS (inner)",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
				&azkv.MasterKey{
					VaultURL:     "AZKV vault URL (inner)",
					Name:         "AZKV name (inner)",
					Version:      "AZKV version (inner)",
					EncryptedKey: "ABCD AZKV (inner)",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
				&hcvault.MasterKey{
					VaultAddress: "HC Vault address (inner)",
					EnginePath:   "HC Vault engine path (inner)",
					KeyName:      "HC Vault key name (inner)",
					EncryptedKey: "ABCD HC Vault (inner)",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
				&pgp.MasterKey{
					Fingerprint:  "PGP fingerprint (inner)",
					EncryptedKey: "ABCD PGP (inner)",
					CreationDate: time.Unix(1742725229, 0).UTC(),
				},
			},
		},
	}, metadata)
}

func TestSerializeMetadata(t *testing.T) {
}
