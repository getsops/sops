package plugin

import (
	"time"

	"github.com/getsops/sops/v3/keys"
)

func init() {
	keys.RegisterProvider(&Provider{})
}

type Provider struct{}

func (p *Provider) Type() string {
	return "plugins"
}

func (p *Provider) MarshalKey(key keys.MasterKey) (map[string]any, error) {
	k, ok := key.(*MasterKey)
	if !ok {
		return nil, nil
	}
	return map[string]any{
		"binary_name": k.BinaryName,
		"instance_id": k.InstanceID,
		"config":      k.PluginConfig,
		"enc":         k.EncryptedKey,
		"created_at":  k.CreationDate.Format(time.RFC3339),
		"timeout":     k.Timeout,
	}, nil
}

func (p *Provider) UnmarshalKey(data map[string]any) (keys.MasterKey, error) {
	var binaryName, instanceID, enc, timeout, createdAtStr string
	if v, ok := data["binary_name"].(string); ok {
		binaryName = v
	}
	if v, ok := data["instance_id"].(string); ok {
		instanceID = v
	}
	if v, ok := data["enc"].(string); ok {
		enc = v
	}
	if v, ok := data["timeout"].(string); ok {
		timeout = v
	}
	if v, ok := data["created_at"].(string); ok {
		createdAtStr = v
	}

	var config map[string]any
	if v, ok := data["config"].(map[string]any); ok {
		config = v
	} else if v, ok := data["config"].(map[interface{}]interface{}); ok {
		config = make(map[string]any)
		for k, val := range v {
			if s, ok := k.(string); ok {
				config[s] = val
			}
		}
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}
	return &MasterKey{
		BinaryName:   binaryName,
		InstanceID:   instanceID,
		PluginConfig: config,
		EncryptedKey: enc,
		CreationDate: createdAt,
		Timeout:      timeout,
	}, nil
}

func (p *Provider) KeysFromConfig(config any, opts keys.CreationOptions) ([]keys.MasterKey, error) {
	maps, ok := config.([]interface{})
	if !ok {
		return nil, nil
	}

	var globalTimeout string
	if opts.GlobalConfig != nil {
		if t, ok := opts.GlobalConfig["timeout"].(string); ok {
			globalTimeout = t
		}
	}

	var res []keys.MasterKey
	for _, item := range maps {
		m, ok := item.(map[string]interface{})
		if !ok {
			// yaml.v3 might decode to map[string]interface{}
			continue
		}
		var binaryName, instanceID, timeout string
		if v, ok := m["binary_name"].(string); ok {
			binaryName = v
		}
		if v, ok := m["instance_id"].(string); ok {
			instanceID = v
		}
		if v, ok := m["timeout"].(string); ok {
			timeout = v
		}

		if timeout == "" {
			timeout = globalTimeout
		}

		var pluginConfig map[string]interface{}
		if v, ok := m["config"].(map[string]interface{}); ok {
			pluginConfig = v
		}
		res = append(res, NewMasterKey(binaryName, pluginConfig, timeout, instanceID))
	}
	return res, nil
}
