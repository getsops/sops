package hcvault

import (
	"fmt"
	logger "log"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("vault", "1.2.2", []string{"VAULT_DEV_ROOT_TOKEN_ID=secret"})
	if err != nil {
		logger.Fatalf("Could not start resource: %s", err)
	}

	os.Setenv("VAULT_ADDR", fmt.Sprintf("http://127.0.0.1:%v", resource.GetPort("8200/tcp")))
	os.Setenv("VAULT_TOKEN", "secret")
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		cli, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			return fmt.Errorf("Cannot create Vault Client: %v", err)
		}
		status, err := cli.Sys().InitStatus()
		if err != nil {
			return err
		}
		if status != true {
			return fmt.Errorf("Vault not ready yet")
		}
		return nil
	}); err != nil {
		logger.Fatalf("Could not connect to docker: %s", err)
	}

	key := NewMasterKey(fmt.Sprintf("http://127.0.0.1:%v", resource.GetPort("8200/tcp")), "sops", "main")
	err = key.createVaultTransitAndKey()
	if err != nil {
		logger.Fatal(err)
	}
	code := 0
	if err == nil {
		code = m.Run()
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		logger.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestKeyToMap(t *testing.T) {
	key := MasterKey{
		CreationDate: time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		VaultAddress: "http://127.0.0.1:8200",
		EnginePath:   "foo",
		KeyName:      "bar",
		EncryptedKey: "this is encrypted",
	}
	assert.Equal(t, map[string]interface{}{
		"vault_address": "http://127.0.0.1:8200",
		"engine_path":   "foo",
		"key_name":      "bar",
		"enc":           "this is encrypted",
		"created_at":    "2016-10-31T10:00:00Z",
	}, key.ToMap())
}

func TestEncryptionDecryption(t *testing.T) {
	dataKey := []byte("super very Secret Key!!!")
	key := MasterKey{
		VaultAddress: os.Getenv("VAULT_ADDR"),
		EnginePath:   "sops",
		KeyName:      "main",
	}
	err := key.Encrypt(dataKey)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	decrypted, err := key.Decrypt()
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	assert.Equal(t, dataKey, decrypted)
}

func TestNewMasterKeyFromURI(t *testing.T) {
	uri1 := "https://vault.example.com:8200/v1/transit/keys/keyName"
	uri2 := "https://vault.me.com/v1/super42/bestmarket/keys/slig"
	uri3 := "http://127.0.0.1:12121/v1/transit/keys/dev"

	mk1 := &MasterKey{
		VaultAddress: "https://vault.example.com:8200",
		EnginePath:   "transit",
		KeyName:      "keyName",
	}
	mk2 := &MasterKey{
		VaultAddress: "https://vault.me.com",
		EnginePath:   "super42/bestmarket",
		KeyName:      "slig",
	}
	mk3 := &MasterKey{
		VaultAddress: "http://127.0.0.1:12121",
		EnginePath:   "transit",
		KeyName:      "dev",
	}
	genMk1, err := NewMasterKeyFromURI(uri1)
	if err != nil {
		log.Errorln(err)
		t.Fail()
	}

	genMk2, err := NewMasterKeyFromURI(uri2)
	if err != nil {
		log.Errorln(err)
		t.Fail()
	}

	genMk3, err := NewMasterKeyFromURI(uri3)
	if err != nil {
		log.Errorln(err)
		t.Fail()
	}

	if assert.NotNil(t, genMk1) {
		mk1.CreationDate = genMk1.CreationDate
		assert.Equal(t, mk1, genMk1)
	}
	if assert.NotNil(t, genMk2) {
		mk2.CreationDate = genMk2.CreationDate
		assert.Equal(t, mk2, genMk2)
	}
	if assert.NotNil(t, genMk3) {
		mk3.CreationDate = genMk3.CreationDate
		assert.Equal(t, mk3, genMk3)
	}

	badURIs := []string{
		"vault.me/keys/dev/mykey",
		"http://127.0.0.1:12121/v1/keys/dev",
		"tcp://127.0.0.1:12121/v1/keys/dev",
	}
	for _, uri := range badURIs {
		if _, err = NewMasterKeyFromURI(uri); err == nil {
			log.Errorf("Should be a invalid uri: %s", uri)
			t.Fail()
		}
	}

}
