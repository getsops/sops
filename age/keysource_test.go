package age

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// mockRecipient is a mock age recipient, it matches mockIdentity.
	mockRecipient string = "age1lzd99uklcjnc0e7d860axevet2cz99ce9pq6tzuzd05l5nr28ams36nvun"
	// mockIdentity is a mock age identity.
	mockIdentity string = "AGE-SECRET-KEY-1G0Q5K9TV4REQ3ZSQRMTMG8NSWQGYT0T7TZ33RAZEE0GZYVZN0APSU24RK7"
	// mockOtherIdentity is another mock age identity.
	mockOtherIdentity string = "AGE-SECRET-KEY-1432K5YRNSC44GC4986NXMX6GVZ52WTMT9C79CLUVWYY4DKDHD5JSNDP4MC"
	// mockEncryptedKey equals to mockEncryptedKeyPlain when decrypted with mockIdentity.
	mockEncryptedKey string = `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBvY2t2NkdLUGRvY3l2OGNy
MVJWcUhCOEZrUG8yeCtnRnhxL0I5NFk4YjJFCmE4SVQ3MEdyZkFqRWpSa2F0NVhF
VDUybzBxdS9nSGpHSVRVMUI0UEVqZkkKLS0tIGJjeGhNQ0Y5L2VZRVVYSm90djFF
bzdnQ3UwTGljMmtrbWNMV1MxYkFzUFUK4xjOZOTGdcbzuwUY/zeBXhcF+Md3e5PQ
EylloI7MNGbadPGb
-----END AGE ENCRYPTED FILE-----`
	// mockEncryptedKeyPlain is the plain value of mockEncryptedKey.
	mockEncryptedKeyPlain string = "data"
	// passphrase used to encrypt age identity.
	mockIdentityPassphrase string = "passphrase"
	mockEncryptedIdentity  string = `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHNjcnlwdCBMN2FXZW9xSFViYjdNeW5D
dy9iSHFnIDE4Ck9zV0ZoNldmci9rL3VXd3BtZmQvK3VZWEpBQjdhZ0UrcmhqR2lF
YThFMzAKLS0tIGVEQ0xwODI1TlNYeHNHaHZKWHoyLzYwMTMvTGhaZG1oa203cSs0
VUpBL1kKsaTnt+H/z8mkL21UYKIt3YMpWSV/oYqTm1cSSUnF9InZEYU9HndK9rc8
ni+MTJCmYf4mgvvGPMf7oIQvs6ijaTdlQb+zeQsL4eif20w+CWgvPNrS6iXUIs8W
w5/fHsxwmrkG96nDkMErJKhmjmLpC+YdbiMe6P/KIpas09m08RTIqcz7ua0Xm3ey
ndU+8ILJOhcnWV55W43nTw/UUFse7f+qY61n7kcd1sGd7ZfSEdEIqS3K2vEtA3ER
fn0s3cyXVEBxL9OZqcAk45bCFVOl13Fp/DBfquHEjvAyeg0=
-----END AGE ENCRYPTED FILE-----`
	// mockSshRecipient is a mock age ssh recipient, it matches mockSshIdentity
	mockSshRecipient string = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAID+Wi8WZw2bXfBpcs/WECttCzP39OkenS6pHWHWGFJvN Test"
	// mockSshIdentity is a mock age identity based on an OpenSSH private key (ed25519)
	mockSshIdentity string = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACA/lovFmcNm13waXLP1hArbQsz9/TpHp0uqR1h1hhSbzQAAAIgCXDMIAlwz
CAAAAAtzc2gtZWQyNTUxOQAAACA/lovFmcNm13waXLP1hArbQsz9/TpHp0uqR1h1hhSbzQ
AAAEBJdWTJ8dC0OnMcwy4gQ96sp6KG8GE9EiyhFGhKldKiST+Wi8WZw2bXfBpcs/WECttC
zP39OkenS6pHWHWGFJvNAAAABFRlc3QB
-----END OPENSSH PRIVATE KEY-----`
	mockEncryptedSshKey string = `-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHNzaC1lZDI1NTE5IDJjd0R4dyB2R3Ns
VUNHaXBiTEJaNU5BMFFQZUpCYWJqODFyTTZ4WWZoRVpUd2M2aTBFCkduUFJHb1U2
K3RqWVQrLzE4anZKZ3h2T3c2MFpZTHlGaHprcElXenByWTAKLS0tIG56MHFSZERl
em9PWmRMMTY4aytYTnVZN04yeER5Z2E3TWxWT3JTZWR2ekUKp/HZLy4MzQqoszGk
+P0hSPPNhOhvFwv4AqCw1+A+WyeHGQPq
-----END AGE ENCRYPTED FILE-----`
)

func TestMasterKeysFromRecipients(t *testing.T) {
	const otherRecipient = "age1tmaae3ld5vpevmsh5yacsauzx8jetg300mpvc4ugp5zr5l6ssq9sla97ep"

	t.Run("recipient", func(t *testing.T) {
		got, err := MasterKeysFromRecipients(mockRecipient)
		assert.NoError(t, err)

		assert.Len(t, got, 1)
		assert.Equal(t, got[0].Recipient, mockRecipient)
	})

	t.Run("recipient-ssh", func(t *testing.T) {
		got, err := MasterKeysFromRecipients(mockSshRecipient)
		assert.NoError(t, err)

		assert.Len(t, got, 1)
		assert.Equal(t, got[0].Recipient, mockSshRecipient)
	})

	t.Run("recipients", func(t *testing.T) {
		got, err := MasterKeysFromRecipients(mockRecipient + "," + otherRecipient + "," + mockSshRecipient)
		assert.NoError(t, err)

		assert.Len(t, got, 3)
		assert.Equal(t, got[0].Recipient, mockRecipient)
		assert.Equal(t, got[1].Recipient, otherRecipient)
		assert.Equal(t, got[2].Recipient, mockSshRecipient)
	})

	t.Run("leading and trailing spaces", func(t *testing.T) {
		got, err := MasterKeysFromRecipients("   " + mockRecipient + "   ,   " + otherRecipient + " ,  " + mockSshRecipient + "     ")
		assert.NoError(t, err)

		assert.Len(t, got, 3)
		assert.Equal(t, got[0].Recipient, mockRecipient)
		assert.Equal(t, got[1].Recipient, otherRecipient)
		assert.Equal(t, got[2].Recipient, mockSshRecipient)
	})

	t.Run("empty", func(t *testing.T) {
		got, err := MasterKeysFromRecipients("")
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})
}

func TestMasterKeyFromRecipient(t *testing.T) {
	t.Run("recipient", func(t *testing.T) {
		got, err := MasterKeyFromRecipient(mockRecipient)
		assert.NoError(t, err)
		assert.EqualValues(t, mockRecipient, got.Recipient)
		assert.NotNil(t, got.parsedRecipient)
		assert.Nil(t, got.parsedIdentities)
	})

	t.Run("recipient-ssh", func(t *testing.T) {
		got, err := MasterKeyFromRecipient(mockSshRecipient)
		assert.NoError(t, err)
		assert.EqualValues(t, mockSshRecipient, got.Recipient)
		assert.NotNil(t, got.parsedRecipient)
		assert.Nil(t, got.parsedIdentities)
	})

	t.Run("leading and trailing spaces", func(t *testing.T) {
		got, err := MasterKeyFromRecipient("   " + mockRecipient + "   ")
		assert.NoError(t, err)
		assert.EqualValues(t, mockRecipient, got.Recipient)
		assert.NotNil(t, got.parsedRecipient)
		assert.Nil(t, got.parsedIdentities)
	})

	t.Run("leading and trailing spaces - ssh", func(t *testing.T) {
		got, err := MasterKeyFromRecipient("   " + mockSshRecipient + "   ")
		assert.NoError(t, err)
		assert.EqualValues(t, mockSshRecipient, got.Recipient)
		assert.NotNil(t, got.parsedRecipient)
		assert.Nil(t, got.parsedIdentities)
	})

	t.Run("invalid recipient", func(t *testing.T) {
		got, err := MasterKeyFromRecipient("invalid")
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestParsedIdentities_Import(t *testing.T) {
	i := make(ParsedIdentities, 0)
	assert.NoError(t, i.Import(mockIdentity, mockOtherIdentity))
	assert.Len(t, i, 2)

	assert.Error(t, i.Import("invalid"))
	assert.Len(t, i, 2)
}

func TestParsedIdentities_ApplyToMasterKey(t *testing.T) {
	i := make(ParsedIdentities, 0)
	assert.NoError(t, i.Import(mockIdentity, mockOtherIdentity))

	key := &MasterKey{}
	i.ApplyToMasterKey(key)
	assert.EqualValues(t, key.parsedIdentities, i)
}

func TestMasterKey_Encrypt(t *testing.T) {
	mockParsedRecipient, err := parseRecipient(mockRecipient)
	assert.NoError(t, err)
	mockSshParsedRecipient, err := parseRecipient(mockSshRecipient)
	assert.NoError(t, err)

	t.Run("recipient", func(t *testing.T) {
		key := &MasterKey{
			Recipient: mockRecipient,
		}
		assert.NoError(t, key.Encrypt([]byte(mockEncryptedKeyPlain)))
		assert.NotEmpty(t, key.EncryptedKey)
	})

	t.Run("recipient ssh", func(t *testing.T) {
		key := &MasterKey{
			Recipient: mockSshRecipient,
		}
		assert.NoError(t, key.Encrypt([]byte(mockEncryptedKeyPlain)))
		assert.NotEmpty(t, key.EncryptedKey)
	})

	t.Run("parsed recipient", func(t *testing.T) {
		key := &MasterKey{
			parsedRecipient: mockParsedRecipient,
		}
		assert.NoError(t, key.Encrypt([]byte(mockEncryptedKeyPlain)))
		assert.NotEmpty(t, key.EncryptedKey)
	})

	t.Run("parsed recipient ssh", func(t *testing.T) {
		key := &MasterKey{
			parsedRecipient: mockSshParsedRecipient,
		}
		assert.NoError(t, key.Encrypt([]byte(mockEncryptedKeyPlain)))
		assert.NotEmpty(t, key.EncryptedKey)
	})

	t.Run("invalid recipient", func(t *testing.T) {
		key := &MasterKey{
			Recipient: "invalid",
		}
		err := key.Encrypt([]byte(mockEncryptedKeyPlain))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse input, unknown recipient type:")
		assert.Empty(t, key.EncryptedKey)
	})

	t.Run("parsed recipient and invalid recipient", func(t *testing.T) {
		key := &MasterKey{
			Recipient:       "invalid",
			parsedRecipient: mockParsedRecipient,
		}
		// Validates mockParsedRecipient > Recipient
		assert.NoError(t, key.Encrypt([]byte(mockEncryptedKeyPlain)))
		assert.NotEmpty(t, key.EncryptedKey)
	})
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	key, err := MasterKeyFromRecipient(mockRecipient)
	assert.NoError(t, err)

	assert.NoError(t, key.EncryptIfNeeded([]byte(mockEncryptedKeyPlain)))

	encryptedKey := key.EncryptedKey
	assert.Contains(t, encryptedKey, "AGE ENCRYPTED FILE")

	assert.NoError(t, key.EncryptIfNeeded([]byte("some other data")))
	assert.Equal(t, encryptedKey, key.EncryptedKey)
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key := &MasterKey{EncryptedKey: "some key"}
	assert.EqualValues(t, key.EncryptedKey, key.EncryptedDataKey())
}

func TestMasterKey_Decrypt(t *testing.T) {
	t.Run("parsed identities", func(t *testing.T) {
		key := &MasterKey{EncryptedKey: mockEncryptedKey}
		var ids ParsedIdentities
		assert.NoError(t, ids.Import(mockOtherIdentity, mockIdentity))
		ids.ApplyToMasterKey(key)

		got, err := key.Decrypt()
		assert.NoError(t, err)
		assert.EqualValues(t, mockEncryptedKeyPlain, got)
	})

	t.Run("loaded identities", func(t *testing.T) {
		overwriteUserConfigDir(t, t.TempDir())
		key := &MasterKey{EncryptedKey: mockEncryptedKey}
		t.Setenv(SopsAgeKeyEnv, mockIdentity)

		got, err := key.Decrypt()
		assert.NoError(t, err)
		assert.EqualValues(t, mockEncryptedKeyPlain, got)
	})

	t.Run("loaded identities ssh", func(t *testing.T) {
		key := &MasterKey{EncryptedKey: mockEncryptedSshKey}
		tmp := t.TempDir()
		overwriteUserConfigDir(t, tmp)

		homeDir, err := os.UserHomeDir()
		assert.NoError(t, err)
		keyPath := filepath.Join(homeDir, ".ssh/id_25519")
		assert.True(t, strings.HasPrefix(keyPath, homeDir))

		assert.NoError(t, os.MkdirAll(filepath.Dir(keyPath), 0o700))
		assert.NoError(t, os.WriteFile(keyPath, []byte(mockSshIdentity), 0o644))
		t.Setenv(SopsAgeSshPrivateKeyFileEnv, keyPath)

		got, err := key.Decrypt()
		assert.NoError(t, err)
		assert.EqualValues(t, mockEncryptedKeyPlain, got)
	})

	t.Run("no identities", func(t *testing.T) {
		tmpDir := t.TempDir()
		overwriteUserConfigDir(t, tmpDir)

		key := &MasterKey{EncryptedKey: mockEncryptedKey}
		got, err := key.Decrypt()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to load age identities")
		assert.Nil(t, got)
	})

	t.Run("no matching identity", func(t *testing.T) {
		key := &MasterKey{EncryptedKey: mockEncryptedKey}
		var ids ParsedIdentities
		assert.NoError(t, ids.Import(mockOtherIdentity))
		ids.ApplyToMasterKey(key)

		// This confirms lazy-loading works as intended
		t.Setenv(SopsAgeKeyEnv, mockIdentity)

		got, err := key.Decrypt()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no identity matched any of the recipients")
		assert.Nil(t, got)
	})

	t.Run("invalid encrypted key", func(t *testing.T) {
		overwriteUserConfigDir(t, t.TempDir())
		key := &MasterKey{EncryptedKey: "invalid"}
		t.Setenv(SopsAgeKeyEnv, mockIdentity)

		got, err := key.Decrypt()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to create reader for decrypting sops data key with age")
		assert.Nil(t, got)
	})
}

func TestMasterKey_EncryptDecrypt_RoundTrip(t *testing.T) {
	encryptKey, err := MasterKeyFromRecipient(mockRecipient)
	assert.NoError(t, err)

	data := []byte("some secret data")
	assert.NoError(t, encryptKey.Encrypt(data))
	assert.NotEmpty(t, encryptKey.EncryptedKey)

	var ids ParsedIdentities
	assert.NoError(t, ids.Import(mockIdentity))

	decryptKey := &MasterKey{}
	decryptKey.EncryptedKey = encryptKey.EncryptedKey
	ids.ApplyToMasterKey(decryptKey)

	decryptedData, err := decryptKey.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, data, decryptedData)
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key := &MasterKey{Recipient: mockRecipient}
	assert.False(t, key.NeedsRotation())
}

func TestMasterKey_ToString(t *testing.T) {
	key := &MasterKey{Recipient: mockRecipient}
	assert.Equal(t, key.Recipient, key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key := &MasterKey{
		Recipient:    mockRecipient,
		EncryptedKey: "some-encrypted-key",
	}
	assert.Equal(t, map[string]interface{}{
		"recipient": mockRecipient,
		"enc":       key.EncryptedKey,
	}, key.ToMap())
}

func TestMasterKey_loadIdentities(t *testing.T) {
	t.Run(SopsAgeKeyEnv, func(t *testing.T) {
		tmpDir := t.TempDir()
		// Overwrite to ensure local config is not picked up by tests
		overwriteUserConfigDir(t, tmpDir)

		t.Setenv(SopsAgeKeyEnv, mockIdentity)

		key := &MasterKey{}
		got, err := key.loadIdentities()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run(SopsAgeKeyEnv+" multiple", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Overwrite to ensure local config is not picked up by tests
		overwriteUserConfigDir(t, tmpDir)

		t.Setenv(SopsAgeKeyEnv, mockIdentity+"\n"+mockOtherIdentity)

		key := &MasterKey{}
		got, err := key.loadIdentities()
		assert.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run(SopsAgeKeyFileEnv, func(t *testing.T) {
		tmpDir := t.TempDir()
		// Overwrite to ensure local config is not picked up by tests
		overwriteUserConfigDir(t, tmpDir)

		keyPath := filepath.Join(tmpDir, "keys.txt")
		assert.NoError(t, os.WriteFile(keyPath, []byte(mockIdentity), 0o644))

		t.Setenv(SopsAgeKeyFileEnv, keyPath)

		key := &MasterKey{}
		got, err := key.loadIdentities()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run(SopsAgeKeyUserConfigPath, func(t *testing.T) {
		tmpDir := t.TempDir()
		overwriteUserConfigDir(t, tmpDir)

		// We need to use getUserConfigDir and not tmpDir as it may add a suffix
		cfgDir, err := getUserConfigDir()
		assert.NoError(t, err)
		keyPath := filepath.Join(cfgDir, SopsAgeKeyUserConfigPath)
		assert.True(t, strings.HasPrefix(keyPath, cfgDir))

		assert.NoError(t, os.MkdirAll(filepath.Dir(keyPath), 0o700))
		assert.NoError(t, os.WriteFile(keyPath, []byte(mockIdentity), 0o644))

		got, err := (&MasterKey{}).loadIdentities()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run(SopsAgeSshPrivateKeyFileEnv, func(t *testing.T) {
		tmpDir := t.TempDir()
		overwriteUserConfigDir(t, tmpDir)

		homeDir, err := os.UserHomeDir()
		assert.NoError(t, err)
		keyPath := filepath.Join(homeDir, ".ssh/id_25519")
		assert.True(t, strings.HasPrefix(keyPath, homeDir))

		assert.NoError(t, os.MkdirAll(filepath.Dir(keyPath), 0o700))
		assert.NoError(t, os.WriteFile(keyPath, []byte(mockSshIdentity), 0o644))
		t.Setenv(SopsAgeSshPrivateKeyFileEnv, keyPath)

		key := &MasterKey{}
		got, err := key.loadIdentities()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("no identity", func(t *testing.T) {
		tmpDir := t.TempDir()
		overwriteUserConfigDir(t, tmpDir)

		got, err := (&MasterKey{}).loadIdentities()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to open file")
		assert.Nil(t, got)
	})

	t.Run("multiple identities", func(t *testing.T) {
		tmpDir := t.TempDir()
		overwriteUserConfigDir(t, tmpDir)

		// We need to use getUserConfigDir and not tmpDir as it may add a suffix
		cfgDir, err := getUserConfigDir()
		assert.NoError(t, err)
		keyPath1 := filepath.Join(cfgDir, SopsAgeKeyUserConfigPath)
		assert.True(t, strings.HasPrefix(keyPath1, cfgDir))

		assert.NoError(t, os.MkdirAll(filepath.Dir(keyPath1), 0o700))
		assert.NoError(t, os.WriteFile(keyPath1, []byte(mockIdentity), 0o644))

		keyPath2 := filepath.Join(tmpDir, "keys.txt")
		assert.NoError(t, os.WriteFile(keyPath2, []byte(mockOtherIdentity), 0o644))
		t.Setenv(SopsAgeKeyFileEnv, keyPath2)

		got, err := (&MasterKey{}).loadIdentities()
		assert.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("parsing error", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Overwrite to ensure local config is not picked up by tests
		overwriteUserConfigDir(t, tmpDir)

		t.Setenv(SopsAgeKeyEnv, "invalid")

		key := &MasterKey{}
		got, err := key.loadIdentities()
		assert.Error(t, err)
		assert.ErrorContains(t, err, fmt.Sprintf("failed to parse '%s' age identities", SopsAgeKeyEnv))
		assert.Nil(t, got)
	})

	t.Run(SopsAgeKeyCmdEnv, func(t *testing.T) {
		tmpDir := t.TempDir()
		// Overwrite to ensure local config is not picked up by tests
		overwriteUserConfigDir(t, tmpDir)

		t.Setenv(SopsAgeKeyCmdEnv, "echo '"+mockIdentity+"'")

		key := &MasterKey{}
		got, err := key.loadIdentities()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("cmd error", func(t *testing.T) {
		tmpDir := t.TempDir()
		// Overwrite to ensure local config is not picked up by tests
		overwriteUserConfigDir(t, tmpDir)

		t.Setenv(SopsAgeKeyCmdEnv, "meow")

		key := &MasterKey{}
		got, err := key.loadIdentities()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to execute command meow")
		assert.Nil(t, got)
	})
}

// overwriteUserConfigDir sets the user config directory and the user home directory
// based on the os.UserConfigDir logic.
func overwriteUserConfigDir(t *testing.T, path string) {
	switch runtime.GOOS {
	case "windows":
		t.Setenv("AppData", path)
	case "plan9": // This adds "/lib" as a suffix to $home
		t.Setenv("home", path)
	default: // Unix
		t.Setenv("XDG_CONFIG_HOME", path)
		t.Setenv("HOME", path)
	}
}

// Make sure that on all supported platforms but Windows, XDG_CONFIG_HOME
// can be used to specify the user's home directory. For most platforms
// this is handled by Go's os.UserConfigDir(), but for Darwin our code
// in getUserConfigDir() handles this explicitly.
func TestUserConfigDir(t *testing.T) {
	if runtime.GOOS != "windows" {
		const dir = "/test/home/dir"
		t.Setenv("XDG_CONFIG_HOME", dir)
		home, err := getUserConfigDir()
		assert.Nil(t, err)
		assert.Equal(t, home, dir)
	}
}

func TestMasterKey_Identities_Passphrase(t *testing.T) {
	t.Run(SopsAgeKeyEnv, func(t *testing.T) {
		key := &MasterKey{EncryptedKey: mockEncryptedKey}
		t.Setenv(SopsAgeKeyEnv, mockEncryptedIdentity)
		//blocks calling gpg-agent
		os.Unsetenv("XDG_RUNTIME_DIR")
		testOnlyAgePassword = mockIdentityPassphrase
		got, err := key.Decrypt()
		testOnlyAgePassword = ""

		assert.NoError(t, err)
		assert.EqualValues(t, mockEncryptedKeyPlain, got)
	})

	t.Run(SopsAgeKeyFileEnv, func(t *testing.T) {
		tmpDir := t.TempDir()
		// Overwrite to ensure local config is not picked up by tests
		overwriteUserConfigDir(t, tmpDir)

		keyPath := filepath.Join(tmpDir, "keys.txt")
		assert.NoError(t, os.WriteFile(keyPath, []byte(mockEncryptedIdentity), 0o644))

		key := &MasterKey{EncryptedKey: mockEncryptedKey}
		t.Setenv(SopsAgeKeyFileEnv, keyPath)
		//blocks calling gpg-agent
		os.Unsetenv("XDG_RUNTIME_DIR")
		testOnlyAgePassword = mockIdentityPassphrase

		got, err := key.Decrypt()
		testOnlyAgePassword = ""

		assert.NoError(t, err)
		assert.EqualValues(t, mockEncryptedKeyPlain, got)
	})

	t.Run("invalid encrypted key", func(t *testing.T) {
		key := &MasterKey{EncryptedKey: "invalid"}
		t.Setenv(SopsAgeKeyEnv, mockEncryptedIdentity)
		//blocks calling gpg-agent
		os.Unsetenv("XDG_RUNTIME_DIR")
		testOnlyAgePassword = mockIdentityPassphrase

		got, err := key.Decrypt()
		testOnlyAgePassword = ""

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to create reader for decrypting sops data key with age")
		assert.Nil(t, got)
	})
}
