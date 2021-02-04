package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops/v3"
)

var PLAIN = []byte(`---
# comment 0
key1: value
key1_a: value
# ^ comment 1
---
key2: value2`)

var PLAIN_0 = []byte(`# comment 0
key1: value
key1_a: value
# ^ comment 1
`)

var BRANCHES = sops.TreeBranches{
	sops.TreeBranch{
		sops.TreeItem{
			Key:   sops.Comment{" comment 0"},
			Value: nil,
		},
		sops.TreeItem{
			Key:   "key1",
			Value: "value",
		},
		sops.TreeItem{
			Key:   "key1_a",
			Value: "value",
		},
		sops.TreeItem{
			Key:   sops.Comment{" ^ comment 1"},
			Value: nil,
		},
	},
	sops.TreeBranch{
		sops.TreeItem{
			Key:   "key2",
			Value: "value2",
		},
	},
}

var COMMENT_1 = []byte(`# test
a:
    b: null
    # foo
`)

var COMMENT_2 = []byte(`a:
    # foo
    b: null
`)

var COMMENT_3_IN = []byte(`## Configuration for prometheus-node-exporter subchart
##
prometheus-node-exporter:
  podLabels:
    ## Add the 'node-exporter' label to be used by serviceMonitor to match standard common usage in rules and grafana dashboards
    ##

    jobLabel: node-exporter
  extraArgs:
    - --collector.filesystem.ignored-mount-points=^/(dev|proc|sys|var/lib/docker/.+)($|/)
    - --collector.filesystem.ignored-fs-types=^(autofs|binfmt_misc|cgroup|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|sysfs|tracefs)$
`)
var COMMENT_3_OUT = []byte(`## Configuration for prometheus-node-exporter subchart
##
prometheus-node-exporter:
    podLabels:
        ## Add the 'node-exporter' label to be used by serviceMonitor to match standard common usage in rules and grafana dashboards
        ##
        jobLabel: node-exporter
    extraArgs:
        - --collector.filesystem.ignored-mount-points=^/(dev|proc|sys|var/lib/docker/.+)($|/)
        - --collector.filesystem.ignored-fs-types=^(autofs|binfmt_misc|cgroup|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|sysfs|tracefs)$
`)

var COMMENT_4 = []byte(`# foo
`)

var COMMENT_5 = []byte(`# foo
---
key: value
`)

func TestUnmarshalMetadataFromNonSOPSFile(t *testing.T) {
	data := []byte(`hello: 2`)
	_, err := (&Store{}).LoadEncryptedFile(data)
	assert.Equal(t, sops.MetadataNotFound, err)
}

func TestLoadPlainFile(t *testing.T) {
	branches, err := (&Store{}).LoadPlainFile(PLAIN)
	assert.Nil(t, err)
	assert.Equal(t, BRANCHES, branches)
}

func TestComment1(t *testing.T) {
	// First iteration: load and store
	branches, err := (&Store{}).LoadPlainFile(COMMENT_1)
	assert.Nil(t, err)
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, string(COMMENT_1), string(bytes))
	assert.Equal(t, COMMENT_1, bytes)
}

func TestComment2(t *testing.T) {
	// First iteration: load and store
	branches, err := (&Store{}).LoadPlainFile(COMMENT_2)
	assert.Nil(t, err)
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, string(COMMENT_2), string(bytes))
	assert.Equal(t, COMMENT_2, bytes)
}

func TestComment3(t *testing.T) {
	// First iteration: load and store
	branches, err := (&Store{}).LoadPlainFile(COMMENT_3_IN)
	assert.Nil(t, err)
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, string(COMMENT_3_OUT), string(bytes))
	assert.Equal(t, COMMENT_3_OUT, bytes)
}

/* TODO: re-enable once https://github.com/go-yaml/yaml/pull/690 is merged
func TestComment4(t *testing.T) {
	// First iteration: load and store
	branches, err := (&Store{}).LoadPlainFile(COMMENT_4)
	assert.Nil(t, err)
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, string(COMMENT_4), string(bytes))
	assert.Equal(t, COMMENT_4, bytes)
}

func TestComment5(t *testing.T) {
	// First iteration: load and store
	branches, err := (&Store{}).LoadPlainFile(COMMENT_5)
	assert.Nil(t, err)
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, string(COMMENT_5), string(bytes))
	assert.Equal(t, COMMENT_5, bytes)
}
*/

func TestEmpty(t *testing.T) {
	// First iteration: load and store
	branches, err := (&Store{}).LoadPlainFile([]byte(``))
	assert.Nil(t, err)
	assert.Equal(t, len(branches), 0)
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, ``, string(bytes))
}

/* TODO: re-enable once https://github.com/go-yaml/yaml/pull/690 is merged
func TestEmpty2(t *testing.T) {
	// First iteration: load and store
	branches, err := (&Store{}).LoadPlainFile([]byte(`---`))
	assert.Nil(t, err)
	assert.Equal(t, len(branches), 1)
	assert.Equal(t, len(branches[0]), 0)
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, ``, string(bytes))
}
*/

func TestEmitValue(t *testing.T) {
	// First iteration: load and store
	bytes, err := (&Store{}).EmitValue(BRANCHES[0])
	assert.Nil(t, err)
	assert.Equal(t, string(PLAIN_0), string(bytes))
	assert.Equal(t, PLAIN_0, bytes)
}
