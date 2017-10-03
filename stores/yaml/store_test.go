package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops"
)

func TestUnmarshalMetadataFromNonSOPSFile(t *testing.T) {
	data := []byte(`hello: 2`)
	_, err := (&Store{}).UnmarshalMetadata(data)
	assert.Equal(t, sops.MetadataNotFound, err)
}
