package common

import (
	"testing"

	"github.com/getsops/sops/v3"
	sopsplugin "github.com/getsops/sops/v3/plugin"
	"github.com/stretchr/testify/require"
)

func TestDiffKeyGroupsIncludesPluginConfigurationIdentity(t *testing.T) {
	before, err := sopsplugin.NewMasterKey("kms", map[string]interface{}{"key_id": "old"})
	require.NoError(t, err)
	before.EncryptedKey = []byte("old ciphertext")

	sameIdentity, err := sopsplugin.NewMasterKey("kms", map[string]interface{}{"key_id": "old"})
	require.NoError(t, err)
	sameIdentity.EncryptedKey = []byte("new ciphertext")

	changedConfig, err := sopsplugin.NewMasterKey("kms", map[string]interface{}{"key_id": "new"})
	require.NoError(t, err)

	unchanged := DiffKeyGroups(
		[]sops.KeyGroup{{before}},
		[]sops.KeyGroup{{sameIdentity}},
	)
	require.Len(t, unchanged, 1)
	require.Len(t, unchanged[0].Common, 1)
	require.Same(t, sameIdentity, unchanged[0].Common[0])
	require.Empty(t, unchanged[0].Added)
	require.Empty(t, unchanged[0].Removed)

	changed := DiffKeyGroups(
		[]sops.KeyGroup{{before}},
		[]sops.KeyGroup{{changedConfig}},
	)
	require.Len(t, changed, 1)
	require.Empty(t, changed[0].Common)
	require.Len(t, changed[0].Added, 1)
	require.Same(t, changedConfig, changed[0].Added[0])
	require.Len(t, changed[0].Removed, 1)
	require.Same(t, before, changed[0].Removed[0])
	require.NotContains(t, before.ToString(), "old")
}
