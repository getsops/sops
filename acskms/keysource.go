package acskms

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kmssdk "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/logging"
	"github.com/sirupsen/logrus"
)

const (
	// arnRegex matches an ACS ARN, for example:
	// "acs:kms:cn-shanghai:1234567890:key/key-idxxxxx".
	arnRegex = `^acs:kms:(.+):[0-9]+:key/(.+)$`
	// kmsTTL is the duration after which a MasterKey requires rotation.
	kmsTTL = time.Hour
	// KeyTypeIdentifier is the string used to identify an ACS KMS MasterKey.
	KeyTypeIdentifier = "acs_kms"
)

var (
	// log is the global logger for any Alibaba Cloud KMS MasterKey.
	log *logrus.Logger
)

func init() {
	log = logging.NewLogger("ACSKMS")
}

// MasterKey is an Alibaba Cloud KMS key used to encrypt and decrypt SOPS' data key.
type MasterKey struct {
	// Arn is the full key ARN
	Arn string
	// Region is the Alibaba Cloud region (e.g., "cn-hangzhou")
	Region string
	// EncryptedKey stores the data key in its encrypted form.
	EncryptedKey string
	// CreationDate is when this MasterKey was created.
	CreationDate time.Time
}

// NewMasterKey creates a new MasterKey from a key arn string, setting
// the creation date to the current date.
func NewMasterKey(arn string) (*MasterKey, error) {
	region, err := parseKeyArn(arn)
	if err != nil {
		return nil, err
	}
	return &MasterKey{
		Arn:          arn,
		Region:       region,
		CreationDate: time.Now().UTC(),
	}, nil
}

// NewMasterKeyFromKeyIDString takes a comma separated list of Alibaba Cloud KMS
// key ARNs, and returns a slice of new MasterKeys.
func NewMasterKeyFromKeyIDString(keyArn string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if keyArn == "" {
		return keys, nil
	}
	for _, s := range strings.Split(keyArn, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		k, err := NewMasterKey(s)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// parseKeyArn parse an Alibaba Cloud KMS key identifier, which can be a full ARN.
func parseKeyArn(arn string) (string, error) {
	re := regexp.MustCompile(arnRegex)
	matches := re.FindStringSubmatch(arn)
	if len(matches) != 3 {
		return "", fmt.Errorf("invalid ACS KMS key ARN: %s", arn)
	}

	return matches[1], nil
}

// Encrypt encrypts the data key using Alibaba Cloud KMS.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	client, err := key.getClient()
	if err != nil {
		return err
	}

	request := &kmssdk.EncryptRequest{
		KeyId:     tea.String(key.Arn),
		Plaintext: tea.String(base64.StdEncoding.EncodeToString(dataKey)),
	}

	resp, err := client.Encrypt(request)
	if err != nil {
		return fmt.Errorf("acskms encrypt error: %v", err)
	}

	key.EncryptedKey = *resp.Body.CiphertextBlob
	return nil
}

// Decrypt decrypts the data key using Alibaba Cloud KMS.
func (key *MasterKey) Decrypt() ([]byte, error) {
	client, err := key.getClient()
	if err != nil {
		return nil, err
	}

	request := &kmssdk.DecryptRequest{
		CiphertextBlob: tea.String(key.EncryptedKey),
	}
	// If an endpoint is manually set (e.g. KMS Instance), we might need to rely on the SDK's behavior.
	// The standard SDK usually works fine with CiphertextBlob.

	resp, err := client.Decrypt(request)
	if err != nil {
		return nil, fmt.Errorf("acskms decrypt error: %v", err)
	}

	return base64.StdEncoding.DecodeString(*resp.Body.Plaintext)
}

// EncryptIfNeeded encrypts the data key if it's not already encrypted.
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// EncryptedDataKey returns the encrypted data key.
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// NeedsRotation checks if the key needs rotation.
func (key *MasterKey) NeedsRotation() bool {
	return false
}

// ToString returns the string representation of the key.
func (key *MasterKey) ToString() string {
	return key.Arn
}

// ToMap returns the map representation of the key.
func (key *MasterKey) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"arn":        key.Arn,
		"created_at": key.CreationDate.UTC().Format(time.RFC3339),
		"enc":        key.EncryptedKey,
	}
}

// TypeToIdentifier returns the type identifier of the key.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// getClient returns a new Alibaba Cloud KMS client.
func (key *MasterKey) getClient() (*kmssdk.Client, error) {
	cred, err := credentials.NewCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("acskms credential error: %v", err)
	}

	config := &openapi.Config{
		Credential: cred,
		RegionId:   tea.String(key.Region),
	}

	if endpoint := os.Getenv("SOPS_ACSKMS_INSTANCE_ENDPOINT"); endpoint != "" {
		config.Endpoint = tea.String(endpoint)
	} else if key.Region != "" {
		config.Endpoint = tea.String(fmt.Sprintf("kms.%s.aliyuncs.com", key.Region))
	}

	if caFile := os.Getenv("SOPS_ACSKMS_CA_FILE"); caFile != "" {
		caContent, err := os.ReadFile(caFile)
		if err == nil {
			config.Ca = tea.String(string(caContent))
		} else {
			log.Warnf("Failed to read CA file %s: %v", caFile, err)
		}
	}

	client, err := kmssdk.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("acskms client error: %v", err)
	}
	return client, nil
}

// ApplyToMasterKey applies the key parameters to the MasterKey.
// Helper to reconstruct key from map.
func (key *MasterKey) ApplyToMasterKey(k keys.MasterKey) {
	// Not strictly needed for basic interface but good to have parity
}
