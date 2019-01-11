package ini

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops"
)

func TestDecodeIni(t *testing.T) {
	in := `
; last modified 1 April 2001 by John Doe
[owner]
name=John Doe
organization=Acme Widgets Inc.

[database]
; use IP address in case network name resolution is not working
server=192.0.2.62     
port=143
file="payroll.dat"
`
	expected := sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key:   "DEFAULT",
				Value: sops.TreeBranch(nil),
			},
			sops.TreeItem{
				Key: "owner",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   sops.Comment{Value: "last modified 1 April 2001 by John Doe"},
						Value: nil,
					},
					sops.TreeItem{
						Key:   "name",
						Value: "John Doe",
					},
					sops.TreeItem{
						Key:   "organization",
						Value: "Acme Widgets Inc.",
					},
				},
			},
			sops.TreeItem{
				Key: "database",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "server",
						Value: "192.0.2.62",
					},
					sops.TreeItem{
						Key:   sops.Comment{Value: "use IP address in case network name resolution is not working"},
						Value: nil,
					},
					sops.TreeItem{
						Key:   "port",
						Value: "143",
					},
					sops.TreeItem{
						Key:   "file",
						Value: "payroll.dat",
					},
				},
			},
		},
	}
	branch, err := Store{}.treeBranchesFromIni([]byte(in))
	assert.Nil(t, err)
	assert.Equal(t, expected, branch)
}

func TestEncodeSimpleIni(t *testing.T) {
	branches := sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key: "DEFAULT",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "foo",
						Value: "bar",
					},
					sops.TreeItem{
						Key:   "baz",
						Value: "3.0",
					},
					sops.TreeItem{
						Key:   "qux",
						Value: "false",
					},
				},
			},
		},
	}
	out, err := Store{}.iniFromTreeBranches(branches)
	assert.Nil(t, err)
	expected, _ := Store{}.treeBranchesFromIni(out)
	assert.Equal(t, expected, branches)
}

func TestEncodeIniWithEscaping(t *testing.T) {
	branches := sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key: "DEFAULT",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "foo\\bar",
						Value: "value",
					},
					sops.TreeItem{
						Key:   "a_key_with\"quotes\"",
						Value: "4.0",
					},
					sops.TreeItem{
						Key:   "baz\\\\foo",
						Value: "2.0",
					},
				},
			},
		},
	}
	out, err := Store{}.iniFromTreeBranches(branches)
	assert.Nil(t, err)
	expected, _ := Store{}.treeBranchesFromIni(out)
	assert.Equal(t, expected, branches)
}

func TestUnmarshalMetadataFromNonSOPSFile(t *testing.T) {
	data := []byte(`hello=2`)
	store := Store{}
	_, err := store.LoadEncryptedFile(data)
	assert.Equal(t, sops.MetadataNotFound, err)
}
