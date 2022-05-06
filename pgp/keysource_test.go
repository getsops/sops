package pgp

import (
	"bytes"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	mockPublicKey   = "testdata/public.gpg"
	mockPrivateKey  = "testdata/private.gpg"
	mockFingerprint = "B59DAF469E8C948138901A649732075EA221A7EA"
)

func TestMasterKeyFromFingerprint(t *testing.T) {
	key := NewMasterKeyFromFingerprint(mockFingerprint)
	assert.Equal(t, mockFingerprint, key.Fingerprint)
	assert.NotNil(t, key.CreationDate)

	key = NewMasterKeyFromFingerprint("B59DAF 469E8C94813 8901A 649732075E A221A7EA")
	assert.Equal(t, mockFingerprint, key.Fingerprint)
}

func TestNewGnuPGHome(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)

	assert.NotEmpty(t, gnuPGHome.String())
	assert.DirExists(t, gnuPGHome.String())
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})
	assert.NoError(t, gnuPGHome.Validate())
}

func TestGnuPGHome_Import(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})

	b, err := os.ReadFile(mockPublicKey)
	assert.NoError(t, err)
	assert.NoError(t, gnuPGHome.Import(b))

	err, _, stderr := gpgExec(gnuPGHome.String(), []string{"--list-keys", mockFingerprint}, nil)
	assert.NoErrorf(t, err, stderr.String())

	b, err = os.ReadFile(mockPrivateKey)
	assert.NoError(t, err)
	assert.NoError(t, gnuPGHome.Import(b))

	err, _, stderr = gpgExec(gnuPGHome.String(), []string{"--list-secret-keys", mockFingerprint}, nil)
	assert.NoErrorf(t, err, stderr.String())

	assert.Error(t, gnuPGHome.Import([]byte("invalid armored data")))
	assert.Error(t, GnuPGHome("").Import(b))
}

func TestGnuPGHome_ImportFile(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})

	assert.NoError(t, gnuPGHome.ImportFile(mockPublicKey))
	assert.Error(t, gnuPGHome.ImportFile("invalid"))
}

func TestGnuPGHome_Validate(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		assert.Error(t, GnuPGHome("").Validate())
	})

	t.Run("relative path", func(t *testing.T) {
		assert.Error(t, GnuPGHome("../../.gnupghome").Validate())
	})

	t.Run("file path", func(t *testing.T) {
		tmpDir := t.TempDir()
		f, err := os.CreateTemp(tmpDir, "file")
		assert.NoError(t, err)
		defer f.Close()

		assert.Error(t, GnuPGHome(f.Name()).Validate())
	})

	t.Run("wrong permissions", func(t *testing.T) {
		// Is created with 0755
		tmpDir := t.TempDir()
		assert.Error(t, GnuPGHome(tmpDir).Validate())
	})

	t.Run("valid", func(t *testing.T) {
		gnupgHome, err := NewGnuPGHome()
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = os.RemoveAll(gnupgHome.String())
		})
		assert.NoError(t, gnupgHome.Validate())
	})
}

func TestGnuPGHome_String(t *testing.T) {
	gnuPGHome := GnuPGHome("/some/absolute/path")
	assert.Equal(t, "/some/absolute/path", gnuPGHome.String())
}

func TestGnuPGHome_ApplyToMasterKey(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})

	key := NewMasterKeyFromFingerprint(mockFingerprint)
	gnuPGHome.ApplyToMasterKey(key)
	assert.Equal(t, gnuPGHome.String(), key.gnuPGHomeDir)

	gnuPGHome = "/non/existing/absolute/path/fails/validate"
	gnuPGHome.ApplyToMasterKey(key)
	assert.NotEqual(t, gnuPGHome.String(), key.gnuPGHomeDir)
}

func TestMasterKey_Encrypt(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})
	assert.NoError(t, gnuPGHome.ImportFile(mockPublicKey))

	key := NewMasterKeyFromFingerprint(mockFingerprint)
	gnuPGHome.ApplyToMasterKey(key)
	data := []byte("oh no, my darkest secret")
	assert.NoError(t, key.Encrypt(data))

	assert.NotEmpty(t, key.EncryptedKey)
	assert.NotEqual(t, data, key.EncryptedKey)

	assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

	args := []string{
		"-d",
	}
	err, stdout, stderr := gpgExec(key.gnuPGHome(), args, strings.NewReader(key.EncryptedKey))
	assert.NoError(t, err, stderr.String())
	assert.Equal(t, data, stdout.Bytes())

	key.Fingerprint = "invalid"
	err = key.Encrypt([]byte("invalid"))
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to encrypt sops data key with pgp: gpg: 'invalid' is not a valid long keyID")
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})
	assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

	key := NewMasterKeyFromFingerprint(mockFingerprint)
	gnuPGHome.ApplyToMasterKey(key)
	assert.NoError(t, key.EncryptIfNeeded([]byte("data")))

	encryptedKey := key.EncryptedKey
	assert.Contains(t, encryptedKey, "END PGP MESSAGE")

	assert.NoError(t, key.EncryptIfNeeded([]byte("some other data")))
	assert.Equal(t, encryptedKey, key.EncryptedKey)
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key := &MasterKey{EncryptedKey: "some key"}
	assert.EqualValues(t, key.EncryptedKey, key.EncryptedDataKey())
}

func TestMasterKey_Decrypt(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})
	assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

	fingerprint := shortenFingerprint(mockFingerprint)

	data := []byte("this data is absolutely top secret")
	err, stdout, stderr := gpgExec(gnuPGHome.String(), []string{
		"--no-default-recipient",
		"--yes",
		"--encrypt",
		"-a",
		"-r",
		fingerprint,
		"--trusted-key",
		fingerprint,
		"--no-encrypt-to",
	}, bytes.NewReader(data))
	assert.NoErrorf(t, gnuPGHome.ImportFile(mockPrivateKey), stderr.String())

	encryptedData := stdout.String()
	assert.NotEqualValues(t, data, encryptedData)

	key := NewMasterKeyFromFingerprint(mockFingerprint)
	gnuPGHome.ApplyToMasterKey(key)
	key.EncryptedKey = encryptedData

	got, err := key.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, data, got)

	key.EncryptedKey = "absolute invalid"
	got, err = key.Decrypt()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "gpg: no valid OpenPGP data found")
	assert.Nil(t, got)
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key := NewMasterKeyFromFingerprint("")
	assert.False(t, key.NeedsRotation())

	key.CreationDate = key.CreationDate.Add(-(pgpTTL + time.Second))
	assert.True(t, key.NeedsRotation())
}

func TestMasterKey_ToString(t *testing.T) {
	key := NewMasterKeyFromFingerprint(mockFingerprint)
	assert.Equal(t, mockFingerprint, key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key := NewMasterKeyFromFingerprint(mockFingerprint)
	key.EncryptedKey = "data"
	assert.Equal(t, map[string]interface{}{
		"fp":         mockFingerprint,
		"created_at": key.CreationDate.UTC().Format(time.RFC3339),
		"enc":        key.EncryptedKey,
	}, key.ToMap())
}

func TestMasterKey_gnuPGHome(t *testing.T) {
	key := &MasterKey{}

	usr, err := user.Current()
	if err == nil {
		assert.Equal(t, filepath.Join(usr.HomeDir, ".gnupg"), key.gnuPGHome())
	} else {
		assert.Equal(t, filepath.Join(os.Getenv("HOME"), ".gnupg"), key.gnuPGHome())
	}

	gnupgHome := "/overwrite/home"
	t.Setenv("GNUPGHOME", gnupgHome)
	assert.Equal(t, gnupgHome, key.gnuPGHome())

	key.gnuPGHomeDir = "/home/dir/overwrite"
	assert.Equal(t, key.gnuPGHomeDir, key.gnuPGHome())
}

func Test_gpgBinary(t *testing.T) {
	assert.Equal(t, "gpg", gpgBinary())

	overwrite := "/some/other/gpg"
	t.Setenv(SopsGpgExecEnv, overwrite)
	assert.Equal(t, overwrite, gpgBinary())
}

func Test_shortenFingerprint(t *testing.T) {
	shortId := shortenFingerprint(mockFingerprint)
	assert.Equal(t, "9732075EA221A7EA", shortId)

	assert.Equal(t, shortId, shortenFingerprint(shortId))
}

// TODO(hidde): previous tests kept around for now.

func TestPGP(t *testing.T) {
	key := NewMasterKeyFromFingerprint("FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4")
	f := func(x []byte) bool {
		if x == nil || len(x) == 0 {
			return true
		}
		if err := key.Encrypt(x); err != nil {
			t.Errorf("Failed to encrypt: %#v err: %w", x, err)
			return false
		}
		k, err := key.Decrypt()
		if err != nil {
			t.Errorf("Failed to decrypt: %#v err: %w", x, err)
			return false
		}
		return bytes.Equal(x, k)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestPGPKeySourceFromString(t *testing.T) {
	s := "C8C5 2C0A B2A4 8174 01E8  12C8 F3CC 3233 3FAD 9F1E, C8C5 2C0A B2A4 8174 01E8  12C8 F3CC 3233 3FAD 9F1E"
	ks := MasterKeysFromFingerprintString(s)
	expected := "C8C52C0AB2A4817401E812C8F3CC32333FAD9F1E"
	if ks[0].Fingerprint != expected {
		t.Errorf("Fingerprint does not match. Got %s, expected %s", ks[0].Fingerprint, expected)
	}

	if ks[1].Fingerprint != expected {
		t.Error("Fingerprint does not match")
	}
}
