package stores_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	sopsplugin "github.com/getsops/sops/v3/plugin"
	"github.com/getsops/sops/v3/stores"
	jsonstore "github.com/getsops/sops/v3/stores/json"
	yamlstore "github.com/getsops/sops/v3/stores/yaml"
	"github.com/stretchr/testify/require"
)

func TestPluginMetadataRoundTripAcrossLayouts(t *testing.T) {
	first, err := sopsplugin.NewMasterKey("first", map[string]interface{}{
		"key_id":  "key-one",
		"options": map[string]interface{}{"region": "west"},
	})
	require.NoError(t, err)
	first.EncryptedKey = []byte{0, 1, 2, 255}
	first.ProviderMetadata = map[string]interface{}{"key_version": "v1"}
	first.CreationDate = time.Unix(1_700_000_000, 0).UTC()

	second, err := sopsplugin.NewMasterKey("second", map[string]interface{}{"key_id": "key-two"})
	require.NoError(t, err)
	second.EncryptedKey = []byte{4, 5, 6}
	second.ProviderMetadata = map[string]interface{}{"operation": "offline-wrap"}
	second.CreationDate = time.Unix(1_700_000_001, 0).UTC()

	for _, test := range []struct {
		name   string
		groups []sops.KeyGroup
	}{
		{name: "single legacy group", groups: []sops.KeyGroup{{first}}},
		{name: "multiple key groups", groups: []sops.KeyGroup{{first}, {second}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			original := sops.Tree{
				Branches: sops.TreeBranches{sops.TreeBranch{}},
				Metadata: sops.Metadata{
					LastModified: time.Unix(1_700_000_100, 0).UTC(),
					Version:      "3.0.0",
					KeyGroups:    test.groups,
				},
			}
			for _, flatten := range []stores.MetadataFlatten{
				stores.MetadataFlattenNone,
				stores.MetadataFlattenBelowTop,
				stores.MetadataFlattenFull,
			} {
				branches, err := stores.SerializeMetadata(original, stores.MetadataOpts{Flatten: flatten})
				require.NoError(t, err)
				_, decoded, err := stores.ExtractMetadata(branches, stores.MetadataOpts{Flatten: flatten})
				require.NoError(t, err)
				require.Len(t, decoded.KeyGroups, len(test.groups))
				for groupIndex, expectedGroup := range test.groups {
					require.Len(t, decoded.KeyGroups[groupIndex], len(expectedGroup))
					for keyIndex, expected := range expectedGroup {
						got, ok := decoded.KeyGroups[groupIndex][keyIndex].(*sopsplugin.MasterKey)
						require.True(t, ok)
						want := expected.(*sopsplugin.MasterKey)
						require.Equal(t, want.Name, got.Name)
						require.Equal(t, want.Config, got.Config)
						require.Equal(t, want.ProviderMetadata, got.ProviderMetadata)
						require.Equal(t, want.EncryptedKey, got.EncryptedKey)
						require.Equal(t, want.CreationDate, got.CreationDate)
					}
				}
			}
		})
	}
}

func TestPluginConfigAndVersionedMetadataKeepLargeNumbersAcrossStores(t *testing.T) {
	const largeInteger = "9007199254740993"
	const preciseDecimal = "0.12345678901234567890123456789"
	key, err := sopsplugin.NewMasterKey("number-test", map[string]interface{}{
		"key_version":   json.Number(largeInteger),
		"precise_value": json.Number(preciseDecimal),
	})
	require.NoError(t, err)
	key.ProviderMetadata = map[string]interface{}{
		"metadata_version": json.Number(largeInteger),
		"precise_value":    json.Number(preciseDecimal),
	}
	key.EncryptedKey = []byte{1, 2, 3}
	key.CreationDate = time.Unix(1_700_000_000, 0).UTC()
	tree := sops.Tree{
		Branches: sops.TreeBranches{sops.TreeBranch{{Key: "secret", Value: "encrypted"}}},
		Metadata: sops.Metadata{
			LastModified: time.Unix(1_700_000_100, 0).UTC(),
			Version:      "3.0.0",
			KeyGroups:    []sops.KeyGroup{{key}},
		},
	}

	for _, test := range []struct {
		name  string
		store interface {
			EmitEncryptedFile(sops.Tree) ([]byte, error)
			LoadEncryptedFile([]byte) (sops.Tree, error)
		}
	}{
		{name: "JSON", store: jsonstore.NewStore(&config.JSONStoreConfig{Indent: -1})},
		{name: "YAML", store: yamlstore.NewStore(&config.YAMLStoreConfig{})},
	} {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := test.store.EmitEncryptedFile(tree)
			require.NoError(t, err)
			decoded, err := test.store.LoadEncryptedFile(encoded)
			require.NoError(t, err)
			pluginKey, ok := decoded.Metadata.KeyGroups[0][0].(*sopsplugin.MasterKey)
			require.True(t, ok)

			var decodedConfig struct {
				KeyVersion   int64       `json:"key_version"`
				PreciseValue json.Number `json:"precise_value"`
			}
			require.NoError(t, sopsplugin.DecodeConfig(pluginKey.Config, &decodedConfig))
			require.Equal(t, int64(9007199254740993), decodedConfig.KeyVersion)
			require.Equal(t, json.Number(preciseDecimal), decodedConfig.PreciseValue)

			var decodedMetadata struct {
				MetadataVersion int64       `json:"metadata_version"`
				PreciseValue    json.Number `json:"precise_value"`
			}
			require.NoError(t, sopsplugin.DecodeMetadata(pluginKey.ProviderMetadata, &decodedMetadata))
			require.Equal(t, int64(9007199254740993), decodedMetadata.MetadataVersion)
			require.Equal(t, json.Number(preciseDecimal), decodedMetadata.PreciseValue)
		})
	}
}
