package plugin

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getsops/sops/v3/keys"
)

const KeyTypeIdentifier = "plugin"

// MasterKey delegates wrapping and unwrapping of the SOPS data key to an
// allowlisted SOPS plugin.
type MasterKey struct {
	Name             string
	Config           map[string]interface{}
	ProviderMetadata map[string]interface{}
	EncryptedKey     []byte
	CreationDate     time.Time
}

// NewMasterKey creates a plugin key with validated JSON-object config.
func NewMasterKey(name string, config map[string]interface{}) (*MasterKey, error) {
	if err := ValidatePluginName(name); err != nil {
		return nil, err
	}
	if _, err := json.Marshal(config); err != nil {
		return nil, fmt.Errorf("invalid plugin configuration: %w", err)
	}
	if config == nil {
		config = map[string]interface{}{}
	}
	return &MasterKey{
		Name:             name,
		Config:           config,
		ProviderMetadata: map[string]interface{}{},
	}, nil
}

func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptWithContext(context.Background(), dataKey)
}

// EncryptWithContext wraps dataKey using the supplied request context.
func (key *MasterKey) EncryptWithContext(ctx context.Context, dataKey []byte) error {
	response, err := WrapKey(ctx, key.Name, WrapKeyRequest{
		Config:   key.Config,
		Metadata: key.ProviderMetadata,
		DataKey:  dataKey,
	})
	if err != nil {
		return err
	}
	key.EncryptedKey = append(key.EncryptedKey[:0], response.EncryptedDataKey...)
	key.ProviderMetadata = response.Metadata
	key.CreationDate = time.Now()
	return nil
}

func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if len(key.EncryptedKey) > 0 {
		return nil
	}
	return key.Encrypt(dataKey)
}

func (key *MasterKey) EncryptedDataKey() []byte {
	return append([]byte(nil), key.EncryptedKey...)
}

func (key *MasterKey) SetEncryptedDataKey(dataKey []byte) {
	key.EncryptedKey = append(key.EncryptedKey[:0], dataKey...)
}

func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptWithContext(context.Background())
}

// DecryptWithContext unwraps the stored data key using the supplied context.
func (key *MasterKey) DecryptWithContext(ctx context.Context) ([]byte, error) {
	response, err := UnwrapKey(ctx, key.Name, UnwrapKeyRequest{
		Config:           key.Config,
		Metadata:         key.ProviderMetadata,
		EncryptedDataKey: key.EncryptedKey,
	})
	if err != nil {
		return nil, err
	}
	return response.DataKey, nil
}

func (key *MasterKey) NeedsRotation() bool {
	return false
}

func (key *MasterKey) ToString() string {
	config := key.Config
	if config == nil {
		config = map[string]interface{}{}
	}
	encodedConfig, err := json.Marshal(config)
	if err != nil {
		return fmt.Sprintf("plugin %s (invalid config)", key.Name)
	}
	hash := sha256.Sum256(encodedConfig)
	return fmt.Sprintf("plugin %s (%x)", key.Name, hash[:6])
}

func (key *MasterKey) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       key.Name,
		"config":     key.Config,
		"metadata":   key.ProviderMetadata,
		"enc":        base64.StdEncoding.EncodeToString(key.EncryptedKey),
		"created_at": key.CreationDate.Format(time.RFC3339),
	}
}

func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// SetProviderMetadataJSON validates and stores metadata returned by a key
// service.
func (key *MasterKey) SetProviderMetadataJSON(data []byte) error {
	if len(data) == 0 {
		key.ProviderMetadata = map[string]interface{}{}
		return nil
	}
	var metadata map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decodeJSONValue(decoder, &metadata); err != nil {
		return err
	}
	if metadata == nil {
		return fmt.Errorf("plugin metadata must be a JSON object")
	}
	key.ProviderMetadata = metadata
	return nil
}

var _ keys.MasterKey = (*MasterKey)(nil)
