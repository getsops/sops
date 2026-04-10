/*
Package acskms contains an implementation of the github.com/getsops/sops/v3.MasterKey
interface that encrypts and decrypts the data key using Alibaba Cloud (Aliyun) KMS.

Authentication credential resolution order:
 1. Environment variables: ALIBABA_CLOUD_ACCESS_KEY_ID / ALIBABA_CLOUD_ACCESS_KEY_SECRET
    Optional STS token:    ALIBABA_CLOUD_SECURITY_TOKEN
 2. Aliyun CLI config file (~/.aliyun/config.json)
    Profile is selected via ALIBABA_CLOUD_PROFILE env var, or the current profile.
    Supported modes: AK, StsToken, CloudSSO (CloudSSO uses the pre-resolved STS creds).
    Config file path can be overridden with ALIBABA_CLOUD_CONFIG_FILE.

The ARN format used in .sops.yaml and encrypted file metadata is:

	acs:kms:{region}:{account-id}:key/{key-id}
	acs:kms:{region}:{account-id}:alias/{alias-name}
*/
package acskms // import "github.com/getsops/sops/v3/acskms"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms "github.com/alibabacloud-go/kms-20160120/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify an Alibaba Cloud KMS MasterKey.
	// This matches the key used in .sops.yaml and in the sops metadata block of encrypted files.
	KeyTypeIdentifier = "acs_kms"
	// acskmsKeyTTL is the duration after which a MasterKey requires rotation.
	acskmsKeyTTL = time.Hour * 24 * 30 * 6 // 6 months

	// SopsACSAccessKeyIDEnv is the environment variable for the Alibaba Cloud access key ID.
	SopsACSAccessKeyIDEnv = "ALIBABA_CLOUD_ACCESS_KEY_ID"
	// SopsACSAccessKeySecretEnv is the environment variable for the Alibaba Cloud access key secret.
	SopsACSAccessKeySecretEnv = "ALIBABA_CLOUD_ACCESS_KEY_SECRET"
	// SopsACSSecurityTokenEnv is the environment variable for the Alibaba Cloud STS security token.
	SopsACSSecurityTokenEnv = "ALIBABA_CLOUD_SECURITY_TOKEN"
)

// arnRegex matches an Alibaba Cloud KMS ARN, for example:
// "acs:kms:cn-shanghai:123456789:key/key-abc123"
// "acs:kms:cn-shanghai:123456789:alias/my-alias"
var arnRegex = regexp.MustCompile(`^acs:kms:([a-z0-9-]+):\d+:(key|alias)/.+$`)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("ACSKMS")
}

// MasterKey is an Alibaba Cloud KMS key used to encrypt and decrypt SOPS' data key.
type MasterKey struct {
	// Arn is the full ARN of the KMS key, e.g.:
	// "acs:kms:cn-shanghai:123456789:key/key-abc123"
	Arn string
	// EncryptedKey stores the data key in its encrypted form (CiphertextBlob, base64 encoded).
	EncryptedKey string
	// CreationDate is when this MasterKey was created.
	CreationDate time.Time

	// region is extracted from the ARN and used to create the KMS client.
	region string
}

// NewMasterKey creates a new MasterKey from an ARN string, setting the creation date to now.
func NewMasterKey(arn string) (*MasterKey, error) {
	arn = strings.TrimSpace(arn)
	region, err := regionFromARN(arn)
	if err != nil {
		return nil, err
	}
	return &MasterKey{
		Arn:          arn,
		CreationDate: time.Now().UTC(),
		region:       region,
	}, nil
}

// MasterKeysFromARNString takes a comma-separated list of ARNs and returns a slice of MasterKeys.
func MasterKeysFromARNString(arnList string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if arnList == "" {
		return keys, nil
	}
	for _, arn := range strings.Split(arnList, ",") {
		key, err := NewMasterKey(strings.TrimSpace(arn))
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

// Encrypt takes a SOPS data key, encrypts it with Alibaba Cloud KMS, and stores
// the result (base64 CiphertextBlob) in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	client, err := key.newKMSClient()
	if err != nil {
		log.WithField("arn", key.Arn).Info("Encryption failed")
		return fmt.Errorf("cannot create Alibaba Cloud KMS client: %w", err)
	}

	// The KMS Encrypt API expects the plaintext as base64-encoded bytes.
	encodedPlaintext := base64.StdEncoding.EncodeToString(dataKey)

	req := &kms.EncryptRequest{
		KeyId:     tea.String(key.keyID()),
		Plaintext: tea.String(encodedPlaintext),
	}
	runtime := &util.RuntimeOptions{}
	resp, err := client.EncryptWithOptions(req, runtime)
	if err != nil {
		log.WithField("arn", key.Arn).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with Alibaba Cloud KMS key %q: %w", key.Arn, err)
	}

	key.EncryptedKey = tea.StringValue(resp.Body.CiphertextBlob)
	log.WithField("arn", key.Arn).Info("Encryption succeeded")
	return nil
}

// EncryptIfNeeded encrypts the provided SOPS data key, if it has not been encrypted yet.
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// EncryptedDataKey returns the encrypted data key this master key holds.
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Decrypt decrypts the EncryptedKey field with Alibaba Cloud KMS and returns the plaintext data key.
func (key *MasterKey) Decrypt() ([]byte, error) {
	client, err := key.newKMSClient()
	if err != nil {
		log.WithField("arn", key.Arn).Info("Decryption failed")
		return nil, fmt.Errorf("cannot create Alibaba Cloud KMS client: %w", err)
	}

	req := &kms.DecryptRequest{
		CiphertextBlob: tea.String(key.EncryptedKey),
	}
	runtime := &util.RuntimeOptions{}
	resp, err := client.DecryptWithOptions(req, runtime)
	if err != nil {
		log.WithField("arn", key.Arn).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with Alibaba Cloud KMS key %q: %w", key.Arn, err)
	}

	// The KMS Decrypt API returns the plaintext as base64-encoded bytes.
	plaintext, err := base64.StdEncoding.DecodeString(tea.StringValue(resp.Body.Plaintext))
	if err != nil {
		log.WithField("arn", key.Arn).Info("Decryption failed")
		return nil, fmt.Errorf("failed to base64-decode decrypted plaintext for key %q: %w", key.Arn, err)
	}

	log.WithField("arn", key.Arn).Info("Decryption succeeded")
	return plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > acskmsKeyTTL
}

// ToString converts the key to a string representation (the ARN).
func (key *MasterKey) ToString() string {
	return key.Arn
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key *MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["arn"] = key.Arn
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// newKMSClient creates a new Alibaba Cloud KMS client for this key's region.
// Credentials are read from environment variables.
// aliyunCredentials holds resolved Alibaba Cloud credentials.
type aliyunCredentials struct {
	AccessKeyID     string
	AccessKeySecret string
	SecurityToken   string
}

// aliyunProfile mirrors fields from ~/.aliyun/config.json used for credential resolution.
type aliyunProfile struct {
	Name            string `json:"name"`
	Mode            string `json:"mode"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SecurityToken   string `json:"sts_token"`
}

type aliyunConfig struct {
	Current  string          `json:"current"`
	Profiles []aliyunProfile `json:"profiles"`
}

// loadAliyunConfigCredentials reads ~/.aliyun/config.json (or ALIBABA_CLOUD_CONFIG_FILE)
// and returns credentials for the requested profile (ALIBABA_CLOUD_PROFILE or current).
// Supports modes: AK, StsToken, CloudSSO (already resolved STS creds in profile).
func loadAliyunConfigCredentials() (*aliyunCredentials, error) {
	cfgFile := os.Getenv("ALIBABA_CLOUD_CONFIG_FILE")
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		cfgFile = filepath.Join(home, ".aliyun", "config.json")
	}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", cfgFile, err)
	}

	var cfg aliyunConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse %s: %w", cfgFile, err)
	}

	profileName := os.Getenv("ALIBABA_CLOUD_PROFILE")
	if profileName == "" {
		profileName = cfg.Current
	}

	for _, p := range cfg.Profiles {
		if p.Name != profileName {
			continue
		}
		switch p.Mode {
		case "AK", "StsToken", "CloudSSO":
			if p.AccessKeyID == "" || p.AccessKeySecret == "" {
				return nil, fmt.Errorf("profile %q has empty credentials (mode: %s)", profileName, p.Mode)
			}
			return &aliyunCredentials{
				AccessKeyID:     p.AccessKeyID,
				AccessKeySecret: p.AccessKeySecret,
				SecurityToken:   p.SecurityToken,
			}, nil
		default:
			return nil, fmt.Errorf("profile %q mode %q is not supported for direct credential loading; use env vars instead", profileName, p.Mode)
		}
	}

	return nil, fmt.Errorf("profile %q not found in %s", profileName, cfgFile)
}

func (key *MasterKey) newKMSClient() (*kms.Client, error) {
	if key.region == "" {
		region, err := regionFromARN(key.Arn)
		if err != nil {
			return nil, err
		}
		key.region = region
	}

	// Credential resolution chain:
	// 1. Environment variables (ALIBABA_CLOUD_ACCESS_KEY_ID / ALIBABA_CLOUD_ACCESS_KEY_SECRET)
	// 2. Aliyun CLI config file (~/.aliyun/config.json), respecting ALIBABA_CLOUD_PROFILE
	accessKeyID := os.Getenv(SopsACSAccessKeyIDEnv)
	accessKeySecret := os.Getenv(SopsACSAccessKeySecretEnv)
	securityToken := os.Getenv(SopsACSSecurityTokenEnv)

	if accessKeyID == "" || accessKeySecret == "" {
		creds, err := loadAliyunConfigCredentials()
		if err != nil {
			return nil, fmt.Errorf("Alibaba Cloud credentials not found: set %s and %s environment variables, or configure an aliyun CLI profile (set ALIBABA_CLOUD_PROFILE or use the current profile in ~/.aliyun/config.json): %w",
				SopsACSAccessKeyIDEnv, SopsACSAccessKeySecretEnv, err)
		}
		accessKeyID = creds.AccessKeyID
		accessKeySecret = creds.AccessKeySecret
		securityToken = creds.SecurityToken
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(key.region),
	}
	if securityToken != "" {
		config.SecurityToken = tea.String(securityToken)
	}

	return kms.NewClient(config)
}

// keyID extracts the KeyId portion from the ARN to pass to the KMS API.
// The KMS Encrypt API accepts the full ARN directly as the KeyId parameter.
func (key *MasterKey) keyID() string {
	return key.Arn
}

// regionFromARN extracts the region from an Alibaba Cloud KMS ARN.
func regionFromARN(arn string) (string, error) {
	matches := arnRegex.FindStringSubmatch(arn)
	if matches == nil {
		return "", fmt.Errorf("invalid Alibaba Cloud KMS ARN %q: expected format acs:kms:{region}:{account}:key/{key-id}", arn)
	}
	return matches[1], nil
}
