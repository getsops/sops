package filestatus

import (
	"path"
	"testing"

	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/config"
	"github.com/stretchr/testify/require"
)

const repoRoot = "../../../../"

func fromRepoRoot(p string) string {
	return path.Join(repoRoot, p)
}

func TestFileStatus(t *testing.T) {
	tests := []struct {
		name              string
		file              string
		expectedEncrypted bool
	}{
		{
			name:              "encrypted file should be reported as such",
			file:              "example.yaml",
			expectedEncrypted: true,
		},
		{
			name: "plain text file should be reported as cleartext",
			file: "functional-tests/res/plainfile.yaml",
		},
		{
			name: "file without mac should be reported as cleartext",
			file: "functional-tests/res/plainfile.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fromRepoRoot(tt.file)
			s := common.DefaultStoreForPath(config.NewStoresConfig(), f)
			encrypted, err := cfs(s, f)
			require.Nil(t, err, "should not error")
			if tt.expectedEncrypted {
				require.True(t, encrypted, "file should have been reported as encrypted")
			} else {
				require.False(t, encrypted, "file should have been reported as cleartext")
			}
		})
	}
}
