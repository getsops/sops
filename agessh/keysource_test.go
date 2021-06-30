package agessh

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasterKeysFromFilesEmpty(t *testing.T) {
	assert := assert.New(t)

	commaSeparatedRecipients := ""
	recipients, err := MasterKeysFromPublicKeyFiles(commaSeparatedRecipients)

	assert.NoError(err)

	assert.Equal(recipients, make([]*MasterKey, 0))
}

func TestAgeSSH(t *testing.T) {
	assert := assert.New(t)

	key, err := MasterKeyFromFile("key.pub")

	assert.NoError(err)
	assert.Equal("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBI1+Q+ntzANebtnicCiAY29uqLGWOaOGLEa8bUqOVmS", key.ToString())

	dataKey := []byte("abcdefghijklmnopqrstuvwxyz123456")

	err = key.Encrypt(dataKey)
	assert.NoError(err)

	_, filename, _, _ := runtime.Caller(0)
	err = os.Setenv(fileEnv, path.Join(path.Dir(filename), "key.priv"))
	assert.NoError(err)

	decryptedKey, err := key.Decrypt()
	assert.NoError(err)
	assert.Equal(dataKey, decryptedKey)
}

func TestAgeDotEnv(t *testing.T) {
	assert := assert.New(t)

	key, err := MasterKeyFromFile("key.pub")

	assert.NoError(err)
	assert.Equal("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBI1+Q+ntzANebtnicCiAY29uqLGWOaOGLEa8bUqOVmS", key.ToString())

	dotenv := `IMAGE_PREFIX=repo/service-
APPLICATION_KEY=K6pfAWuUVND9Fz5SC7jmA6pfAWuUVND9Fz5SC7jmA
KEY_ID=003683d721f2ae683d721f2a1
DOMAIN=files.127.0.0.1.nip.io`
	dataKey := []byte(dotenv)

	err = key.Encrypt(dataKey)
	assert.NoError(err)

	_, filename, _, _ := runtime.Caller(0)
	err = os.Setenv(fileEnv, path.Join(path.Dir(filename), "key.priv"))
	assert.NoError(err)

	decryptedKey, err := key.Decrypt()
	assert.NoError(err)
	assert.Equal(dataKey, decryptedKey)
}
