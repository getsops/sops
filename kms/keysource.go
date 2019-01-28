/*
Package kms contains an implementation of the go.mozilla.org/sops.MasterKey interface that encrypts and decrypts the
data key using AWS KMS with the AWS Go SDK.
*/
package kms //import "go.mozilla.org/sops/kms"

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"go.mozilla.org/sops/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("AWSKMS")
}

// this needs to be a global var for unit tests to work (mockKMS redefines
// it in keysource_test.go)
var kmsSvc kmsiface.KMSAPI
var isMocked bool

// MasterKey is a AWS KMS key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	Arn               string
	Role              string
	EncryptedKey      string
	CreationDate      time.Time
	EncryptionContext map[string]*string
	AwsProfile        string
}

// EncryptedDataKey returns the encrypted data key this master key holds
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Encrypt takes a sops data key, encrypts it with KMS and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {
	// isMocked is set by unit test to indicate that the KMS service
	// has already been initialized. it's ugly, but it works.
	if kmsSvc == nil || !isMocked {
		sess, err := key.createSession()
		if err != nil {
			log.WithField("arn", key.Arn).Info("Encryption failed")
			return fmt.Errorf("Failed to create session: %v", err)
		}
		kmsSvc = kms.New(sess)
	}
	out, err := kmsSvc.Encrypt(&kms.EncryptInput{Plaintext: dataKey, KeyId: &key.Arn, EncryptionContext: key.EncryptionContext})
	if err != nil {
		log.WithField("arn", key.Arn).Info("Encryption failed")
		return fmt.Errorf("Failed to call KMS encryption service: %v", err)
	}
	key.EncryptedKey = base64.StdEncoding.EncodeToString(out.CiphertextBlob)
	log.WithField("arn", key.Arn).Info("Encryption succeeded")
	return nil
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with AWS KMS and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	k, err := base64.StdEncoding.DecodeString(key.EncryptedKey)
	if err != nil {
		log.WithField("arn", key.Arn).Info("Decryption failed")
		return nil, fmt.Errorf("Error base64-decoding encrypted data key: %s", err)
	}
	// isMocked is set by unit test to indicate that the KMS service
	// has already been initialized. it's ugly, but it works.
	if kmsSvc == nil || !isMocked {
		sess, err := key.createSession()
		if err != nil {
			log.WithField("arn", key.Arn).Info("Decryption failed")
			return nil, fmt.Errorf("Error creating AWS session: %v", err)
		}
		kmsSvc = kms.New(sess)
	}
	decrypted, err := kmsSvc.Decrypt(&kms.DecryptInput{CiphertextBlob: k, EncryptionContext: key.EncryptionContext})
	if err != nil {
		log.WithField("arn", key.Arn).Info("Decryption failed")
		return nil, fmt.Errorf("Error decrypting key: %v", err)
	}
	log.WithField("arn", key.Arn).Info("Decryption succeeded")
	return decrypted.Plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return key.Arn
}

// NewMasterKey creates a new MasterKey from an ARN, role and context, setting the creation date to the current date
func NewMasterKey(arn string, role string, context map[string]*string) *MasterKey {
	return &MasterKey{
		Arn:               arn,
		Role:              role,
		EncryptionContext: context,
		CreationDate:      time.Now().UTC(),
	}
}

// NewMasterKeyFromArn takes an ARN string and returns a new MasterKey for that ARN
func NewMasterKeyFromArn(arn string, context map[string]*string, awsProfile string) *MasterKey {
	k := &MasterKey{}
	arn = strings.Replace(arn, " ", "", -1)
	roleIndex := strings.Index(arn, "+arn:aws:iam::")
	if roleIndex > 0 {
		k.Arn = arn[:roleIndex]
		k.Role = arn[roleIndex+1:]
	} else {
		k.Arn = arn
	}
	k.EncryptionContext = context
	k.CreationDate = time.Now().UTC()
	k.AwsProfile = awsProfile
	return k
}

// MasterKeysFromArnString takes a comma separated list of AWS KMS ARNs and returns a slice of new MasterKeys for those ARNs
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

func (key MasterKey) createStsSession(config aws.Config, sess *session.Session) (*session.Session, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	stsService := sts.New(sess)
	name := "sops@" + hostname
	out, err := stsService.AssumeRole(&sts.AssumeRoleInput{
		RoleArn: &key.Role, RoleSessionName: &name})
	if err != nil {
		return nil, fmt.Errorf("Failed to assume role %q: %v", key.Role, err)
	}
	config.Credentials = credentials.NewStaticCredentials(*out.Credentials.AccessKeyId,
		*out.Credentials.SecretAccessKey, *out.Credentials.SessionToken)
	sess, err = session.NewSession(&config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new aws session: %v", err)
	}
	return sess, nil
}

func (key MasterKey) createSession() (*session.Session, error) {
	re := regexp.MustCompile(`^arn:aws[\w-]*:kms:(.+):[0-9]+:(key|alias)/.+$`)
	matches := re.FindStringSubmatch(key.Arn)
	if matches == nil {
		return nil, fmt.Errorf("No valid ARN found in %q", key.Arn)
	}

	config := aws.Config{Region: aws.String(matches[1])}

	if key.AwsProfile != "" {
		config.Credentials = credentials.NewSharedCredentials("", key.AwsProfile)
	}

	opts := session.Options{
		Config:                  config,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}
	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return nil, err
	}
	if key.Role != "" {
		return key.createStsSession(config, sess)
	}
	return sess, nil
}

// ToMap converts the MasterKey to a map for serialization purposes
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

// ParseKMSContext takes either a KMS context map or a comma-separated list of KMS context key:value pairs and returns a map
func ParseKMSContext(in interface{}) map[string]*string {
	nonStringValueWarning := "Encryption context contains a non-string value, context will not be used"
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
