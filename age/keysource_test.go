package age

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAge(t *testing.T) {
	assert := assert.New(t)

	key, err := MasterKeyFromRecipient("age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw")

	assert.NoError(err)
	assert.Equal("age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw", key.ToString())

	dataKey := []byte("abcdefghijklmnopqrstuvwxyz123456")

	err = key.Encrypt(dataKey)
	assert.NoError(err)

	_, filename, _, _ := runtime.Caller(0)
	err = os.Setenv("SOPS_AGE_KEY_FILE", path.Join(path.Dir(filename), "keys.txt"))
	assert.NoError(err)

	decryptedKey, err := key.Decrypt()
	assert.NoError(err)
	assert.Equal(dataKey, decryptedKey)
}
