package ocikms

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/stretchr/testify/require"
)

// writeTempRSAKey writes an unencrypted PKCS#1 RSA private key to a temp file.
func writeTempRSAKey(t *testing.T, dir string) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}
	pemData := pem.EncodeToMemory(pemBlock)
	path := filepath.Join(dir, "oci-test-private-key.pem")
	if err := os.WriteFile(path, pemData, 0600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	return path
}

// writeOCIConfig writes a minimal ~/.oci/config style file.
func writeOCIConfig(t *testing.T, path string, profile string, user string, tenancy string, region string, fingerprint string, keyFile string) {
	t.Helper()
	content := strings.Join([]string{
		"[" + profile + "]",
		"user=" + user,
		"fingerprint=" + fingerprint,
		"key_file=" + keyFile,
		"tenancy=" + tenancy,
		"region=" + region,
		"",
	}, "\n")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

// clearOCIEnv clears OCI SDK environment variables to prevent interference
func clearOCIEnv(t *testing.T) {
	t.Helper()
	envVars := []string{
		"OCI_tenancy_ocid",
		"OCI_user_ocid",
		"OCI_region",
		"OCI_fingerprint",
		"OCI_private_key_path",
	}
	for _, env := range envVars {
		t.Setenv(env, "")
	}
}

// clearCLIOCIEnv clears OCI CLI environment variables to prevent interference
func clearCLIOCIEnv() {
	envVars := []string{
		OCICLITenancy,
		OCICLIUser,
		OCICLIRegion,
		OCICLIFingerprint,
		OCICLIKeyFile,
	}
	for _, env := range envVars {
		os.Unsetenv(env)
	}
}

// disableIPProvider disables Instance Principal provider in tests
func disableIPProvider(t *testing.T) {
	old := newIPProvider
	t.Cleanup(func() { newIPProvider = old })
	newIPProvider = func() (common.ConfigurationProvider, error) {
		return nil, fmt.Errorf("ip disabled in tests")
	}
}

func TestConfigurationProvider_OCI_CLI_Env(t *testing.T) {
	// Disable IP network path in tests by overriding factory
	disableIPProvider(t)
	// Isolate HOME to avoid default file provider interference
	t.Setenv(HomeEnv, t.TempDir())

	// Generate key
	keyDir := t.TempDir()
	keyPath := writeTempRSAKey(t, keyDir)

	// Set OCI_CLI_* envs
	t.Setenv(OCICLITenancy, "ocid1.tenancy.oc1..exampletenancy")
	t.Setenv(OCICLIUser, "ocid1.user.oc1..exampleuser")
	t.Setenv(OCICLIRegion, "us-ashburn-1")
	t.Setenv(OCICLIFingerprint, "aa:bb:cc:dd")
	t.Setenv(OCICLIKeyFile, keyPath)

	// Ensure other providers are not set by accident
	// Native SDK env provider uses lower-case suffixes with prefix OCI_
	clearOCIEnv(t)

	prov, err := configurationProvider()
	require.NoError(t, err)

	tenancy, err := prov.TenancyOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.tenancy.oc1..exampletenancy", tenancy)

	user, err := prov.UserOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.user.oc1..exampleuser", user)

	region, err := prov.Region()
	require.NoError(t, err)
	require.Equal(t, "us-ashburn-1", region)

	fp, err := prov.KeyFingerprint()
	require.NoError(t, err)
	require.Equal(t, "aa:bb:cc:dd", fp)
}

func TestConfigurationProvider_OCI_Env(t *testing.T) {
	disableIPProvider(t)
	// Isolate HOME
	t.Setenv(HomeEnv, t.TempDir())

	keyDir := t.TempDir()
	keyPath := writeTempRSAKey(t, keyDir)

	// SDK env provider expects lower-case suffixes
	t.Setenv(OCITenancyOCID, "ocid1.tenancy.oc1..ten")
	t.Setenv(OCIUserOCID, "ocid1.user.oc1..usr")
	t.Setenv(OCIRegion, "eu-frankfurt-1")
	t.Setenv(OCIFingerprint, "11:22:33:44")
	t.Setenv(OCIPrivateKeyPath, keyPath)

	// Ensure CLI envs are not set (unset, not empty strings)
	clearCLIOCIEnv()

	prov, err := configurationProvider()
	require.NoError(t, err)

	tenancy, err := prov.TenancyOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.tenancy.oc1..ten", tenancy)

	user, err := prov.UserOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.user.oc1..usr", user)

	region, err := prov.Region()
	require.NoError(t, err)
	require.Equal(t, "eu-frankfurt-1", region)

	fp, err := prov.KeyFingerprint()
	require.NoError(t, err)
	require.Equal(t, "11:22:33:44", fp)
}

func TestConfigurationProvider_FileViaEnv(t *testing.T) {
	disableIPProvider(t)
	// Isolate HOME
	t.Setenv(HomeEnv, t.TempDir())

	d := t.TempDir()
	keyPath := writeTempRSAKey(t, d)
	cfgPath := filepath.Join(d, "config")
	writeOCIConfig(t, cfgPath, "DEFAULT", "ocid1.user.oc1..fileusr", "ocid1.tenancy.oc1..fileten", "uk-london-1", "ff:ee:dd:cc", keyPath)

	// Point to config via env
	t.Setenv(OCICLIConfigFile, cfgPath)
	// Explicit profile not required; default is DEFAULT

	// Ensure env-based providers are not set
	clearCLIOCIEnv()

	clearOCIEnv(t)

	prov, err := configurationProvider()
	require.NoError(t, err)

	tenancy, err := prov.TenancyOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.tenancy.oc1..fileten", tenancy)

	user, err := prov.UserOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.user.oc1..fileusr", user)

	region, err := prov.Region()
	require.NoError(t, err)
	require.Equal(t, "uk-london-1", region)

	fp, err := prov.KeyFingerprint()
	require.NoError(t, err)
	require.Equal(t, "ff:ee:dd:cc", fp)
}

func TestConfigurationProvider_DefaultFileFallback(t *testing.T) {
	disableIPProvider(t)
	// Set HOME to a temp dir and create ~/.oci/config
	home := t.TempDir()
	if runtime.GOOS == "windows" {
		// USERPROFILE is also consulted on Windows
		t.Setenv(UserProfileEnv, home)
	}
	t.Setenv(HomeEnv, home)

	ociDir := filepath.Join(home, ".oci")
	if err := os.MkdirAll(ociDir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	keyPath := writeTempRSAKey(t, ociDir)
	cfgPath := filepath.Join(ociDir, "config")
	writeOCIConfig(t, cfgPath, "DEFAULT", "ocid1.user.oc1..defusr", "ocid1.tenancy.oc1..deften", "ap-tokyo-1", "00:aa:bb:cc", keyPath)

	// Ensure no env points to explicit file and env providers are empty
	os.Unsetenv(OCICLIConfigFile)

	clearCLIOCIEnv()

	clearOCIEnv(t)

	prov, err := common.ConfigurationProviderFromFile(cfgPath, "")
	require.NoError(t, err)

	tenancy, err := prov.TenancyOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.tenancy.oc1..deften", tenancy)

	user, err := prov.UserOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.user.oc1..defusr", user)

	region, err := prov.Region()
	require.NoError(t, err)
	require.Equal(t, "ap-tokyo-1", region)

	fp, err := prov.KeyFingerprint()
	require.NoError(t, err)
	require.Equal(t, "00:aa:bb:cc", fp)
}

// ipStubProvider implements common.ConfigurationProvider to stub Instance Principal in tests
type ipStubProvider struct{}

func (ipStubProvider) TenancyOCID() (string, error)    { return "ocid1.tenancy.oc1..ipstub", nil }
func (ipStubProvider) UserOCID() (string, error)       { return "", nil }
func (ipStubProvider) KeyFingerprint() (string, error) { return "ip:stub:fp", nil }
func (ipStubProvider) Region() (string, error)         { return "me-dubai-1", nil }
func (ipStubProvider) KeyID() (string, error)          { return "ST$ipstub", nil }
func (ipStubProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	// generate a small key for completeness
	k, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	return k, nil
}
func (ipStubProvider) AuthType() (common.AuthConfig, error) { return common.AuthConfig{}, nil }

func TestConfigurationProvider_InstancePrincipal_Stubbed(t *testing.T) {
	// Override IP factory to return stub, no network
	old := newIPProvider
	t.Cleanup(func() { newIPProvider = old })
	newIPProvider = func() (common.ConfigurationProvider, error) { return ipStubProvider{}, nil }

	// Isolate environment so that only IP path is viable
	t.Setenv(HomeEnv, t.TempDir())
	os.Unsetenv(OCICLIConfigFile)
	os.Unsetenv(OCICLIProfile)

	// Clear CLI envs
	clearCLIOCIEnv()

	// Clear native SDK envs
	clearOCIEnv(t)

	prov, err := configurationProvider()
	require.NoError(t, err)

	tenancy, err := prov.TenancyOCID()
	require.NoError(t, err)
	require.Equal(t, "ocid1.tenancy.oc1..ipstub", tenancy)

	region, err := prov.Region()
	require.NoError(t, err)
	require.Equal(t, "me-dubai-1", region)

	fp, err := prov.KeyFingerprint()
	require.NoError(t, err)
	require.Equal(t, "ip:stub:fp", fp)
}
