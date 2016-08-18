package sops

import (
	"fmt"
	"time"
)

const DateFormat = "2006-01-02T15:04:05Z"

const DefaultUnencryptedSuffix = "_unencrypted"

type Error string

func (e Error) Error() string { return string(e) }

const MacMismatch = Error("MAC mismatch")

type Metadata struct {
	LastModified              time.Time
	UnencryptedSuffix         string
	MessageAuthenticationCode string
	Version                   string
	KeySources                []KeySource
}

type KeySource struct {
	Name string
	Keys []MasterKey
}

type MasterKey interface {
	Encrypt(dataKey string) error
	EncryptIfNeeded(dataKey string) error
	Decrypt() (string, error)
	NeedsRotation() bool
	ToString() string
	ToMap() map[string]string
}

type Store interface {
	LoadUnencrypted(data string) error
	Load(data, key string) error
	Dump(key string) (string, error)
	DumpUnencrypted() (string, error)
	Metadata() Metadata
	LoadMetadata(in string) error
	SetMetadata(Metadata)
}

func (m *Metadata) MasterKeyCount() int {
	count := 0
	for _, ks := range m.KeySources {
		count += len(ks.Keys)
	}
	return count
}

func (m *Metadata) RemoveMasterKeys(keys []MasterKey) {
	for _, ks := range m.KeySources {
		for i, k := range ks.Keys {
			for _, k2 := range keys {
				if k.ToString() == k2.ToString() {
					ks.Keys = append(ks.Keys[:i], ks.Keys[i+1:]...)
				}
			}
		}
	}
}

func (m *Metadata) UpdateMasterKeys(dataKey string) {
	for _, ks := range m.KeySources {
		for _, k := range ks.Keys {
			err := k.EncryptIfNeeded(dataKey)
			if err != nil {
				fmt.Println("[WARNING]: could not encrypt data key with master key ", k.ToString())
			}
		}
	}
}

func (m *Metadata) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["lastmodified"] = m.LastModified.Format("2006-01-02T15:04:05Z")
	out["unencrypted_suffix"] = m.UnencryptedSuffix
	out["mac"] = m.MessageAuthenticationCode
	out["version"] = m.Version
	for _, ks := range m.KeySources {
		keys := make([]map[string]string, 0)
		for _, k := range ks.Keys {
			keys = append(keys, k.ToMap())
		}
		out[ks.Name] = keys
	}
	return out
}
