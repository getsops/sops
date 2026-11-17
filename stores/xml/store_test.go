package xml

import (
	"testing"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/pgp"
	"github.com/stretchr/testify/assert"
)

func newTestStore() *Store {
	return NewStore(&config.XMLStoreConfig{Indent: 2})
}

func TestDecodeXML(t *testing.T) {
	in := `<?xml version="1.0" encoding="UTF-8"?>
<config version="1.0">
  <!-- a comment -->
  <database>
    <host>localhost</host>
    <password>s3cr3t</password>
  </database>
  <user name="alice">token-a</user>
  <user name="bob">token-b</user>
</config>`
	expected := sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key: "config",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "@version",
						Value: "1.0",
					},
					sops.TreeItem{
						Key:   sops.Comment{Value: " a comment "},
						Value: nil,
					},
					sops.TreeItem{
						Key: "database",
						Value: sops.TreeBranch{
							sops.TreeItem{
								Key:   "host",
								Value: "localhost",
							},
							sops.TreeItem{
								Key:   "password",
								Value: "s3cr3t",
							},
						},
					},
					sops.TreeItem{
						Key: "user",
						Value: []interface{}{
							sops.TreeBranch{
								sops.TreeItem{Key: "@name", Value: "alice"},
								sops.TreeItem{Key: textKey, Value: "token-a"},
							},
							sops.TreeBranch{
								sops.TreeItem{Key: "@name", Value: "bob"},
								sops.TreeItem{Key: textKey, Value: "token-b"},
							},
						},
					},
				},
			},
		},
	}
	branches, err := newTestStore().LoadPlainFile([]byte(in))
	assert.Nil(t, err)
	assert.Equal(t, expected, branches)
}

func TestEncodeXML(t *testing.T) {
	branches := sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key: "config",
				Value: sops.TreeBranch{
					sops.TreeItem{Key: "@version", Value: "1.0"},
					sops.TreeItem{Key: sops.Comment{Value: " a comment "}, Value: nil},
					sops.TreeItem{Key: "host", Value: "local&host"},
					sops.TreeItem{Key: "port", Value: 5432},
					sops.TreeItem{Key: "user", Value: []interface{}{"alice", "bob"}},
				},
			},
		},
	}
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<config version="1.0">
  <!-- a comment -->
  <host>local&amp;host</host>
  <port>5432</port>
  <user>alice</user>
  <user>bob</user>
</config>
`
	out, err := newTestStore().EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(out))
}

func TestRoundTripXML(t *testing.T) {
	in := `<?xml version="1.0" encoding="UTF-8"?>
<config>
  <!-- a comment -->
  <secret>hunter2</secret>
  <list>
    <item>one</item>
    <item>two</item>
  </list>
</config>
`
	store := newTestStore()
	branches, err := store.LoadPlainFile([]byte(in))
	assert.Nil(t, err)
	out, err := store.EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, in, string(out))
}

func TestMetadataRoundTrip(t *testing.T) {
	store := newTestStore()
	branches, err := store.LoadPlainFile([]byte(`<config><secret>hunter2</secret></config>`))
	assert.Nil(t, err)
	tree := sops.Tree{Branches: branches}
	tree.Metadata.Version = "3.0.0"
	tree.Metadata.MessageAuthenticationCode = "mac"
	tree.Metadata.LastModified = tree.Metadata.LastModified.UTC()
	tree.Metadata.KeyGroups = []sops.KeyGroup{
		{&pgp.MasterKey{Fingerprint: "1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A", EncryptedKey: "data", CreationDate: time.Now().UTC()}},
	}

	encrypted, err := store.EmitEncryptedFile(tree)
	assert.Nil(t, err)
	assert.True(t, store.HasSopsTopLevelKey(mustLoad(t, store, encrypted)[0]))

	loaded, err := store.LoadEncryptedFile(encrypted)
	assert.Nil(t, err)
	assert.Equal(t, "3.0.0", loaded.Metadata.Version)
	assert.Equal(t, branches, loaded.Branches)
}

func mustLoad(t *testing.T, store *Store, in []byte) sops.TreeBranches {
	branches, err := store.LoadPlainFile(in)
	assert.Nil(t, err)
	return branches
}

func TestSingleRootElement(t *testing.T) {
	_, err := newTestStore().LoadPlainFile([]byte(`<a/><b/>`))
	assert.NotNil(t, err)
}
