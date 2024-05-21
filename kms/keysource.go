/*
Package kms contains an implementation of the github.com/getsops/sops/v3.MasterKey
interface that encrypts and decrypts the data key using AWS KMS with the SDK
for Go V2.
*/
package kms // import "github.com/getsops/sops/v3/kms"

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// arnRegex matches an AWS ARN, for example:
	// "arn:aws:kms:us-west-2:107501996527:key/612d5f0p-p1l3-45e6-aca6-a5b005693a48".
	arnRegex = `^arn:aws[\w-]*:kms:(.+):[0-9]+:(key|alias)/.+$`
	// stsSessionRegex matches an AWS STS session name, for example:
	// "john_s", "sops@42WQm042".
	stsSessionRegex = "[^a-zA-Z0-9=,.@-_]+"
	// roleSessionNameLengthLimit is the AWS role session name length limit.
	roleSessionNameLengthLimit = 64
	// kmsTTL is the duration after which a MasterKey requires rotation.
	kmsTTL = time.Hour * 24 * 30 * 6
	// KeyTypeIdentifier is the string used to identify an AWS KMS MasterKey.
	KeyTypeIdentifier = "kms"
)

var (
	// log is the global logger for any AWS KMS MasterKey.
	log *logrus.Logger
	// osHostname returns the hostname as reported by the kernel.
	osHostname = os.Hostname
)

func init() {
	log = logging.NewLogger("AWSKMS")
}

// MasterKey is an AWS KMS key used to encrypt and decrypt SOPS' data key using
// AWS SDK for Go V2.
type MasterKey struct {
	// Arn associated with the AWS KMS key.
	Arn string
	// Role ARN used to assume a role through AWS STS.
	Role string
	// EncryptedKey stores the data key in it's encrypted form.
	EncryptedKey string
	// CreationDate is when this MasterKey was created.
	CreationDate time.Time
	// EncryptionContext provides additional context about the data key.
	// Ref: https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#encrypt_context
	EncryptionContext map[string]*string
	// AwsProfile is the profile to use for loading configuration and credentials.
	// Ref: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-profiles
	AwsProfile string

	// credentialsProvider is used to configure the AWS client config with
	// credentials. It can be injected by a (local) keyservice.KeyServiceServer
	// using CredentialsProvider.ApplyToMasterKey. If nil, the default client is used
	// which utilizes runtime environmental values.
	credentialsProvider aws.CredentialsProvider
	// baseEndpoint can be used to override the endpoint the AWS client resolves
	// to by default. This is mostly used for testing purposes as it can not be
	// injected using e.g. an environment variable. The field is not publicly
	// exposed, nor configurable.
	baseEndpoint string
}

// NewMasterKey creates a new MasterKey from an ARN, role and context, setting
// the creation date to the current date.
func NewMasterKey(arn string, role string, context map[string]*string) *MasterKey {
	return &MasterKey{
		Arn:               arn,
		Role:              role,
		EncryptionContext: context,
		CreationDate:      time.Now().UTC(),
	}
}

// NewMasterKeyWithProfile creates a new MasterKey from an ARN, role, context
// and awsProfile, setting the creation date to the current date.
func NewMasterKeyWithProfile(arn string, role string, context map[string]*string, awsProfile string) *MasterKey {
	k := NewMasterKey(arn, role, context)
	k.AwsProfile = awsProfile
	return k
}

// NewMasterKeyFromArn takes an ARN string and returns a new MasterKey for that
// ARN.
func NewMasterKeyFromArn(arn string, context map[string]*string, awsProfile string) *MasterKey {
	key := &MasterKey{}
	arn = strings.Replace(arn, " ", "", -1)
	key.Arn = arn
	roleIndex := strings.Index(arn, "+arn:aws:iam::")
	if roleIndex > 0 {
		// Overwrite ARN
		key.Arn = arn[:roleIndex]
		key.Role = arn[roleIndex+1:]
	}
	key.EncryptionContext = context
	key.CreationDate = time.Now().UTC()
	key.AwsProfile = awsProfile
	return key
}

// MasterKeysFromArnString takes a comma separated list of AWS KMS ARNs, and
// returns a slice of new MasterKeys for those ARNs.
func MasterKeysFromArnString(arn string, context map[string]*string, awsProfile string) []*MasterKey {
	var keys []*MasterKey
	if arn == "" {
		return keys
	}
	for _, s := range strings.Split(arn, ",") {
		keys = append(keys, NewMasterKeyFromArn(s, context, awsProfile))
	}
	return keys
}

// ParseKMSContext takes either a KMS context map or a comma-separated list of
// KMS context key:value pairs, and returns a map.
func ParseKMSContext(in interface{}) map[string]*string {
	const nonStringValueWarning = "Encryption context contains a non-string value, context will not be used"
	out := make(map[string]*string)
	switch in := in.(type) {
	case map[string]interface{}:
		if len(in) == 0 {
			return nil
		}
		for k, v := range in {
			value, ok := v.(string)
			if !ok {
				log.Warn(nonStringValueWarning)
				return nil
			}
			out[k] = &value
		}
	case map[interface{}]interface{}:
		if len(in) == 0 {
			return nil
		}
		for k, v := range in {
			key, ok := k.(string)
			if !ok {
				log.Warn(nonStringValueWarning)
				return nil
			}
			value, ok := v.(string)
			if !ok {
				log.Warn(nonStringValueWarning)
				return nil
			}
			out[key] = &value
		}
	case string:
		if in == "" {
			return nil
		}
		for _, kv := range strings.Split(in, ",") {
			kv := strings.Split(kv, ":")
			if len(kv) != 2 {
				log.Warn(nonStringValueWarning)
				return nil
			}
			out[kv[0]] = &kv[1]
		}
	}
	return out
}

// CredentialsProvider is a wrapper around aws.CredentialsProvider used for
// authentication towards AWS KMS.
type CredentialsProvider struct {
	provider aws.CredentialsProvider
}

// NewCredentialsProvider returns a CredentialsProvider object with the provided
// aws.CredentialsProvider.
func NewCredentialsProvider(cp aws.CredentialsProvider) *CredentialsProvider {
	return &CredentialsProvider{
		provider: cp,
	}
}

// ApplyToMasterKey configures the credentials on the provided key.
func (c CredentialsProvider) ApplyToMasterKey(key *MasterKey) {
	key.credentialsProvider = c.provider
}

// Encrypt takes a SOPS data key, encrypts it with KMS and stores the result
// in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	cfg, err := key.createKMSConfig()
	if err != nil {
		log.WithField("arn", key.Arn).Info("Encryption failed")
		return err
	}
	client := key.createClient(cfg)
	input := &kms.EncryptInput{
		KeyId:             &key.Arn,
		Plaintext:         dataKey,
		EncryptionContext: stringPointerToStringMap(key.EncryptionContext),
	}
	out, err := client.Encrypt(context.TODO(), input)
	if err != nil {
		log.WithField("arn", key.Arn).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with AWS KMS: %w", err)
	}
	key.EncryptedKey = base64.StdEncoding.EncodeToString(out.CiphertextBlob)
	log.WithField("arn", key.Arn).Info("Encryption succeeded")
	return nil
}

// EncryptIfNeeded encrypts the provided SOPS data key, if it has not been
// encrypted yet.
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

// Decrypt decrypts the EncryptedKey with a newly created AWS KMS config, and
// returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	k, err := base64.StdEncoding.DecodeString(key.EncryptedKey)
	if err != nil {
		log.WithField("arn", key.Arn).Info("Decryption failed")
		return nil, fmt.Errorf("error base64-decoding encrypted data key: %s", err)
	}
	cfg, err := key.createKMSConfig()
	if err != nil {
		log.WithField("arn", key.Arn).Info("Decryption failed")
		return nil, err
	}
	client := key.createClient(cfg)
	input := &kms.DecryptInput{
		KeyId:             &key.Arn,
		CiphertextBlob:    k,
		EncryptionContext: stringPointerToStringMap(key.EncryptionContext),
	}
	decrypted, err := client.Decrypt(context.TODO(), input)
	if err != nil {
		log.WithField("arn", key.Arn).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with AWS KMS: %w", err)
	}
	log.WithField("arn", key.Arn).Info("Decryption succeeded")
	return decrypted.Plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > kmsTTL
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return key.Arn
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["arn"] = key.Arn
	if key.Role != "" {
		out["role"] = key.Role
	}
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	if key.EncryptionContext != nil {
		outcontext := make(map[string]string)
		for k, v := range key.EncryptionContext {
			outcontext[k] = *v
		}
		out["context"] = outcontext
	}
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// createKMSConfig returns an AWS config with the credentialsProvider of the
// MasterKey, or the default configuration sources.
func (key MasterKey) createKMSConfig() (*aws.Config, error) {
	re := regexp.MustCompile(arnRegex)
	matches := re.FindStringSubmatch(key.Arn)
	if matches == nil {
		return nil, fmt.Errorf("no valid ARN found in '%s'", key.Arn)
	}
	region := matches[1]

	cfg, err := config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
		// Use the credentialsProvider if present, otherwise default to reading credentials
		// from the environment.
		if key.credentialsProvider != nil {
			lo.Credentials = key.credentialsProvider
		}
		if key.AwsProfile != "" {
			lo.SharedConfigProfile = key.AwsProfile
		}
		lo.Region = region
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not load AWS config: %w", err)
	}

	if key.Role != "" {
		return key.createSTSConfig(&cfg)
	}
	return &cfg, nil
}

// createClient creates a new AWS KMS client with the provided config.
func (key MasterKey) createClient(config *aws.Config) *kms.Client {
	return kms.NewFromConfig(*config, func(o *kms.Options) {
		if key.baseEndpoint != "" {
			o.BaseEndpoint = aws.String(key.baseEndpoint)
		}
	})
}

// createSTSConfig uses AWS STS to assume a role and returns a config
// configured with that role's credentials. It returns an error if
// it fails to construct a session name, or assume the role.
func (key MasterKey) createSTSConfig(config *aws.Config) (*aws.Config, error) {
	name, err := stsSessionName()
	if err != nil {
		return nil, err
	}
	input := &sts.AssumeRoleInput{
		RoleArn:         &key.Role,
		RoleSessionName: &name,
	}

	client := sts.NewFromConfig(*config)
	out, err := client.AssumeRole(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to assume role '%s': %w", key.Role, err)
	}

	config.Credentials = credentials.NewStaticCredentialsProvider(*out.Credentials.AccessKeyId,
		*out.Credentials.SecretAccessKey, *out.Credentials.SessionToken,
	)
	return config, nil
}

// stsSessionName returns the name for the STS session in the format of
// `sops@<hostname>`. It sanitizes the hostname with stsSessionRegex, and
// truncates to roleSessionNameLengthLimit when it exceeds the limit.
func stsSessionName() (string, error) {
	hostname, err := osHostname()
	if err != nil {
		return "", fmt.Errorf("failed to construct STS session name: %w", err)
	}

	re := regexp.MustCompile(stsSessionRegex)
	sanitizedHostname := re.ReplaceAllString(hostname, "")

	name := "sops@" + sanitizedHostname
	if len(name) >= roleSessionNameLengthLimit {
		name = name[:roleSessionNameLengthLimit]
	}
	return name, nil
}

func stringPointerToStringMap(in map[string]*string) map[string]string {
	var out = make(map[string]string)
	for k, v := range in {
		if v == nil {
			continue
		}
		out[k] = *v
	}
	return out
}
