package pgp

import (
	"bytes"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/stretchr/testify/assert"
)

var (
	mockPublicKey   = "testdata/public.gpg"
	mockPrivateKey  = "testdata/private.gpg"
	mockPubRing     = "testdata/ring/pubring.gpg"
	mockSecRing     = "testdata/ring/secring.gpg"
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

	_, stderr, err := gpgExec(gnuPGHome.String(), []string{"--list-keys", mockFingerprint}, nil)
	assert.NoErrorf(t, err, stderr.String())

	b, err = os.ReadFile(mockPrivateKey)
	assert.NoError(t, err)
	assert.NoError(t, gnuPGHome.Import(b))

	_, stderr, err = gpgExec(gnuPGHome.String(), []string{"--list-secret-keys", mockFingerprint}, nil)
	assert.NoErrorf(t, err, stderr.String())

	err = gnuPGHome.Import([]byte("invalid armored data"))
	assert.Error(t, err)
	assert.ErrorContains(t, err, "(exit status 2): gpg: no valid OpenPGP data found.\ngpg: Total number processed: 0")
	assert.Error(t, GnuPGHome("").Import(b))
}

func TestGnuPGHome_Import_With_Missing_Binary(t *testing.T) {
	t.Setenv(SopsGpgExecEnv, "/does/not/exist")

	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})

	b, err := os.ReadFile(mockPublicKey)
	assert.NoError(t, err)
	err = gnuPGHome.Import(b)
	assert.ErrorContains(t, err, "failed to import armored key data into GnuPG keyring: fork/exec /does/not/exist: no such file or directory")
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

func TestGnuPGHome_Cleanup(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)

	err = gnuPGHome.Cleanup()
	assert.NoError(t, err)
	assert.Error(t, gnuPGHome.Validate())

	gnuPGHome = "/an/absolute/invalid/path"
	assert.Error(t, gnuPGHome.Cleanup())
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
		tmpDir := t.TempDir()

		err := os.Chmod(tmpDir, 0o755)
		assert.NoError(t, err)

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

func TestDisableOpenPGP_ApplyToMasterKey(t *testing.T) {
	key := NewMasterKeyFromFingerprint(mockFingerprint)
	DisableOpenPGP{}.ApplyToMasterKey(key)
	assert.True(t, key.disableOpenPGP)
}

func TestPubRing_ApplyToMasterKey(t *testing.T) {
	key := NewMasterKeyFromFingerprint(mockFingerprint)
	pubring := PubRing("/some/path.pgp")
	pubring.ApplyToMasterKey(key)
	assert.Equal(t, string(pubring), key.pubRing)
}

func TestSecRing_ApplyToMasterKey(t *testing.T) {
	key := NewMasterKeyFromFingerprint(mockFingerprint)
	secring := SecRing("/some/path.pgp")
	secring.ApplyToMasterKey(key)
	assert.Equal(t, string(secring), key.secRing)
}

func TestMasterKey_Encrypt(t *testing.T) {
	t.Run("with OpenPGP", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		PubRing(mockPubRing).ApplyToMasterKey(key)

		data := []byte("oh no, my darkest secret")
		assert.NoError(t, key.Encrypt(data))
		assert.NotEqual(t, data, key.EncryptedKey)
		// Detailed testing is done by TestMasterKey_encryptWithOpenPGP
	})

	t.Run("with GnuPG", func(t *testing.T) {
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
		// Detailed testing is done by TestMasterKey_encryptWithGnuPG
	})

	t.Run("with error", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)

		data := []byte("oh no, my darkest secret")
		err := key.Encrypt(data)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "GnuPG binary error")
		assert.ErrorContains(t, err, "github.com/ProtonMail/go-crypto/openpgp error")
	})

	t.Run("with OpenPGP disabled", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		DisableOpenPGP{}.ApplyToMasterKey(key)

		data := []byte("oh no, my darkest secret")
		err := key.Encrypt(data)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "GnuPG binary error")
		assert.NotContains(t, err.Error(), "github.com/ProtonMail/go-crypto/openpgp error")
	})
}

func TestMasterKey_encryptWithOpenPGP(t *testing.T) {
	t.Run("encrypt", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		PubRing(mockPubRing).ApplyToMasterKey(key)

		data := []byte("oh no, my darkest secret")
		assert.NoError(t, key.encryptWithOpenPGP(data))

		assert.NotEmpty(t, key.EncryptedKey)
		assert.NotEqual(t, data, key.EncryptedKey)

		secRing, err := loadRing(mockSecRing)
		assert.NoError(t, err)
		block, err := armor.Decode(strings.NewReader(key.EncryptedKey))
		assert.NoError(t, err)
		md, err := openpgp.ReadMessage(block.Body, secRing, nil, nil)
		assert.NoError(t, err)
		b, err := io.ReadAll(md.UnverifiedBody)
		assert.NoError(t, err)

		assert.Equal(t, data, b)
	})

	t.Run("invalid fingerprint error", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint("invalid")
		err := key.encryptWithOpenPGP([]byte("invalid"))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "key with fingerprint 'invalid' is not available in keyring")
	})
}

func TestMasterKey_encryptWithGnuPG(t *testing.T) {
	t.Run("encrypt", func(t *testing.T) {
		gnuPGHome, err := NewGnuPGHome()
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = os.RemoveAll(gnuPGHome.String())
		})
		assert.NoError(t, gnuPGHome.ImportFile(mockPublicKey))

		key := NewMasterKeyFromFingerprint(mockFingerprint)
		gnuPGHome.ApplyToMasterKey(key)
		data := []byte("oh no, my darkest secret")
		assert.NoError(t, key.encryptWithGnuPG(data))

		assert.NotEmpty(t, key.EncryptedKey)
		assert.NotEqual(t, data, key.EncryptedKey)

		assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

		args := []string{
			"-d",
		}
		stdout, stderr, err := gpgExec(key.gnuPGHomeDir, args, strings.NewReader(key.EncryptedKey))
		assert.NoError(t, err, stderr.String())
		assert.Equal(t, data, stdout.Bytes())
	})

	t.Run("invalid fingerprint error", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint("invalid")
		err := key.encryptWithGnuPG([]byte("invalid"))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to encrypt sops data key with pgp: gpg: 'invalid' is not a valid long keyID")
	})

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
	// Mock encrypted data
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})
	assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

	fingerprint := shortenFingerprint(mockFingerprint)

	data := []byte("this data is absolutely top secret")
	stdout, stderr, err := gpgExec(gnuPGHome.String(), []string{
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
	assert.Nil(t, err)
	assert.NoErrorf(t, gnuPGHome.ImportFile(mockPrivateKey), stderr.String())

	encryptedData := stdout.String()
	assert.NotEqualValues(t, data, encryptedData)

	// Actual tests
	t.Run("with OpenPGP", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.EncryptedKey = encryptedData
		SecRing(mockSecRing).ApplyToMasterKey(key)

		got, err := key.Decrypt()
		assert.NoError(t, err)
		assert.Equal(t, data, got)
		// Detailed testing is done by TestMasterKey_decryptWithOpenPGP
	})

	t.Run("with GnuPG", func(t *testing.T) {
		gnuPGHome, err := NewGnuPGHome()
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = os.RemoveAll(gnuPGHome.String())
		})
		assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.EncryptedKey = encryptedData
		gnuPGHome.ApplyToMasterKey(key)

		got, err := key.Decrypt()
		assert.NoError(t, err)
		assert.Equal(t, data, got)
		// Detailed testing is done by TestMasterKey_decryptWithGnuPG
	})

	t.Run("with error", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.EncryptedKey = encryptedData

		data, err := key.Decrypt()
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.ErrorContains(t, err, "GnuPG binary error")
		assert.ErrorContains(t, err, "github.com/ProtonMail/go-crypto/openpgp error")
	})

	t.Run("with OpenPGP disabled", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.EncryptedKey = encryptedData
		DisableOpenPGP{}.ApplyToMasterKey(key)

		data, err := key.Decrypt()
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.ErrorContains(t, err, "GnuPG binary error")
		assert.NotContains(t, err.Error(), "github.com/ProtonMail/go-crypto/openpgp error")
	})
}

func TestMasterKey_decryptWithOpenPGP(t *testing.T) {
	t.Run("decrypt", func(t *testing.T) {
		gnuPGHome, err := NewGnuPGHome()
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = os.RemoveAll(gnuPGHome.String())
		})
		assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

		fingerprint := shortenFingerprint(mockFingerprint)

		data := []byte("this data is absolutely top secret")
		stdout, stderr, err := gpgExec(gnuPGHome.String(), []string{
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
		assert.Nil(t, err)
		assert.NoErrorf(t, gnuPGHome.ImportFile(mockPrivateKey), stderr.String())

		encryptedData := stdout.String()
		assert.NotEqualValues(t, data, encryptedData)

		key := NewMasterKeyFromFingerprint(mockFingerprint)
		SecRing(mockSecRing).ApplyToMasterKey(key)
		key.EncryptedKey = encryptedData

		got, err := key.decryptWithOpenPGP()
		assert.NoError(t, err)
		assert.Equal(t, data, got)
	})

	t.Run("invalid data error", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.EncryptedKey = "absolute invalid"
		SecRing(mockSecRing).ApplyToMasterKey(key)
		got, err := key.decryptWithOpenPGP()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "armor decoding failed: EOF")
		assert.Nil(t, got)
	})
}

func TestMasterKey_decryptWithGnuPG(t *testing.T) {
	t.Run("decrypt", func(t *testing.T) {
		gnuPGHome, err := NewGnuPGHome()
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = os.RemoveAll(gnuPGHome.String())
		})
		assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

		fingerprint := shortenFingerprint(mockFingerprint)

		data := []byte("this data is absolutely top secret")
		stdout, stderr, err := gpgExec(gnuPGHome.String(), []string{
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
		assert.Nil(t, err)
		assert.NoErrorf(t, gnuPGHome.ImportFile(mockPrivateKey), stderr.String())

		encryptedData := stdout.String()
		assert.NotEqualValues(t, data, encryptedData)

		key := NewMasterKeyFromFingerprint(mockFingerprint)
		gnuPGHome.ApplyToMasterKey(key)
		key.EncryptedKey = encryptedData

		got, err := key.decryptWithGnuPG()
		assert.NoError(t, err)
		assert.Equal(t, data, got)
	})

	t.Run("invalid data error", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.EncryptedKey = "absolute invalid"
		got, err := key.decryptWithGnuPG()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "gpg: no valid OpenPGP data found")
		assert.Nil(t, got)
	})
}

func TestMasterKey_EncryptDecrypt_RoundTrip(t *testing.T) {
	gnuPGHome, err := NewGnuPGHome()
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(gnuPGHome.String())
	})
	assert.NoError(t, gnuPGHome.ImportFile(mockPrivateKey))

	key := NewMasterKeyFromFingerprint(mockFingerprint)
	gnuPGHome.ApplyToMasterKey(key)

	data := []byte("some secret data")
	assert.NoError(t, key.Encrypt(data))
	assert.NotEmpty(t, key.EncryptedKey)

	decryptedData, err := key.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, data, decryptedData)
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

func TestMasterKey_retrievePubKey(t *testing.T) {
	t.Run("existing fingerprint", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		PubRing(mockPubRing).ApplyToMasterKey(key)

		got, err := key.retrievePubKey()
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
	})

	t.Run("non-existing fingerprint", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint("invalid")
		PubRing(mockPubRing).ApplyToMasterKey(key)

		got, err := key.retrievePubKey()
		assert.Error(t, err)
		assert.Empty(t, got)
	})
}

func TestMasterKey_getPubRing(t *testing.T) {
	t.Run("default pub ring", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.gnuPGHomeDir = "testdata/ring"

		got, err := key.getPubRing()
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
	})

	t.Run("key pub ring", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		PubRing(mockPubRing).ApplyToMasterKey(key)

		got, err := key.getPubRing()
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
	})

	t.Run("no pub ring", func(t *testing.T) {
		tmpDir := t.TempDir()

		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.gnuPGHomeDir = tmpDir

		got, err := key.getPubRing()
		assert.Error(t, err)
		assert.Empty(t, got)
	})
}

func TestMasterKey_getSecRing(t *testing.T) {
	t.Run("default sec ring", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		key.gnuPGHomeDir = "testdata/ring"

		got, err := key.getSecRing()
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
	})

	t.Run("key sec ring", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		SecRing(mockSecRing).ApplyToMasterKey(key)

		got, err := key.getSecRing()
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
	})

	t.Run("pub ring fallback", func(t *testing.T) {
		key := NewMasterKeyFromFingerprint(mockFingerprint)
		PubRing(mockSecRing).ApplyToMasterKey(key)

		got, err := key.getSecRing()
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
	})
}

func Test_fingerprintIndex(t *testing.T) {
	r, err := loadRing(mockPubRing)
	assert.NoError(t, err)

	got := fingerprintIndex(r)
	assert.Len(t, got, 1)
	_, ok := got[mockFingerprint]
	assert.True(t, ok)
}

func Test_loadRing(t *testing.T) {
	t.Run("pub ring", func(t *testing.T) {
		r, err := loadRing(mockPubRing)
		assert.NoError(t, err)
		assert.Len(t, r, 1)
	})

	t.Run("sec ring", func(t *testing.T) {
		r, err := loadRing(mockSecRing)
		assert.NoError(t, err)
		assert.Len(t, r, 1)
	})

	t.Run("read error", func(t *testing.T) {
		r, err := loadRing(mockPublicKey)
		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("not found error", func(t *testing.T) {
		r, err := loadRing("/an/absolute/invalid/path")
		assert.Error(t, err)
		assert.Nil(t, r)
	})
}

func Test_gpgBinary(t *testing.T) {
	assert.Equal(t, "gpg", gpgBinary())

	overwrite := "/some/other/gpg"
	t.Setenv(SopsGpgExecEnv, overwrite)
	assert.Equal(t, overwrite, gpgBinary())

	overwrite = "not_abs_path"
	t.Setenv(SopsGpgExecEnv, overwrite)
	assert.Equal(t, overwrite, gpgBinary())
}

func Test_gnuPGHome(t *testing.T) {
	usr, err := user.Current()
	if err == nil {
		assert.Equal(t, filepath.Join(usr.HomeDir, ".gnupg"), gnuPGHome(""))
	} else {
		assert.Equal(t, filepath.Join(os.Getenv("HOME"), ".gnupg"), gnuPGHome(""))
	}

	gnupgHome := "/overwrite/home"
	t.Setenv("GNUPGHOME", gnupgHome)
	assert.Equal(t, gnupgHome, gnuPGHome(""))

	customP := "/home/dir/overwrite"
	assert.Equal(t, customP, gnuPGHome(customP))
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
		if len(x) == 0 {
			return true
		}
		if err := key.Encrypt(x); err != nil {
			t.Errorf("Failed to encrypt: %#v err: %s", x, err)
			return false
		}
		k, err := key.Decrypt()
		if err != nil {
			t.Errorf("Failed to decrypt: %#v err: %s", x, err)
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
