package sops

import (
	"go.mozilla.org/sops/keys"
	"testing"
	"time"

	"crypto/rand"

	"github.com/stretchr/testify/assert"
)

type PlaintextMasterKey struct {
	Key []byte
}

func (k *PlaintextMasterKey) Encrypt(dataKey []byte) error {
	k.Key = dataKey
	return nil
}

func (k *PlaintextMasterKey) EncryptIfNeeded(dataKey []byte) error {
	k.Key = dataKey
	return nil
}

func (k *PlaintextMasterKey) Decrypt() ([]byte, error) {
	return k.Key, nil
}

func (k *PlaintextMasterKey) NeedsRotation() bool {
	return false
}

func (k *PlaintextMasterKey) EncryptedDataKey() []byte {
	return k.Key
}

func (k *PlaintextMasterKey) SetEncryptedDataKey(key []byte) {
	k.Key = key
}

func (k *PlaintextMasterKey) ToString() string {
	return string(k.Key)
}

func (k *PlaintextMasterKey) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"key": k.Key,
	}
}

func TestShamirRoundtripAllKeysAvailable(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	assert.NoError(t, err)
	m := Metadata{
		Shamir: true,
		KeySources: []KeySource{
			{
				Name: "mock",
				Keys: []keys.MasterKey{
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
				},
			},
		},
		LastModified: time.Now(),
	}
	errs := m.UpdateMasterKeys(key)
	assert.Empty(t, errs)
	dataKey, err := m.GetDataKey()
	assert.NoError(t, err)
	assert.Equal(t, key, dataKey)
}

func TestShamirRoundtripQuorumAvailable(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	assert.NoError(t, err)
	m := Metadata{
		Shamir: true,
		KeySources: []KeySource{
			{
				Name: "mock",
				Keys: []keys.MasterKey{
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
				},
			},
		},
		LastModified: time.Now(),
	}
	errs := m.UpdateMasterKeys(key)
	assert.Empty(t, errs)
	m.KeySources[0].Keys = m.KeySources[0].Keys[:3]
	dataKey, err := m.GetDataKey()
	assert.NoError(t, err)
	assert.Equal(t, key, dataKey)
}

func TestShamirRoundtripNotEnoughKeys(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	assert.NoError(t, err)
	m := Metadata{
		Shamir:       true,
		ShamirQuorum: 4,
		KeySources: []KeySource{
			{
				Name: "mock",
				Keys: []keys.MasterKey{
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
					&PlaintextMasterKey{},
				},
			},
		},
		LastModified: time.Now(),
	}
	errs := m.UpdateMasterKeys(key)
	assert.Empty(t, errs)
	m.KeySources[0].Keys = m.KeySources[0].Keys[:2]
	dataKey, err := m.GetDataKey()
	assert.Error(t, err)
	assert.NotEqual(t, key, dataKey)
}
