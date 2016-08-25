package kms

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"regexp"
	"strings"
	"time"
)

var kmsSvc kmsiface.KMSAPI

// MasterKey is a AWS KMS key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	Arn          string
	Role         string
	EncryptedKey string
	CreationDate time.Time
}

// Encrypt takes a sops data key, encrypts it with KMS and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {
	if kmsSvc == nil {

		sess, err := key.createSession()
		if err != nil {
			return err
		}
		kmsSvc = kms.New(sess)
	}
	out, err := kmsSvc.Encrypt(&kms.EncryptInput{Plaintext: dataKey, KeyId: &key.Arn})
	if err != nil {
		return err
	}
	key.EncryptedKey = base64.StdEncoding.EncodeToString(out.CiphertextBlob)
	return nil
}

// EncryptIfNeeded encrypts the provided sops' data ket and encrypts it if it hasn't been encrypted yet
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
		return nil, fmt.Errorf("Error base64-decoding encrypted data key: %s", err)
	}
	if kmsSvc == nil {
		sess, err := key.createSession()
		if err != nil {
			return nil, fmt.Errorf("Error creating AWS session: %v", err)
		}
		kmsSvc = kms.New(sess)
	}
	decrypted, err := kmsSvc.Decrypt(&kms.DecryptInput{CiphertextBlob: k})
	if err != nil {
		return nil, fmt.Errorf("Error decrypting key: %v", err)
	}
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

// NewMasterKeyFromArn takes an ARN string and returns a new MasterKey for that ARN
func NewMasterKeyFromArn(arn string) MasterKey {
	k := MasterKey{}
	arn = strings.Replace(arn, " ", "", -1)
	roleIndex := strings.Index(arn, "+arn:aws:iam::")
	if roleIndex > 0 {
		k.Arn = arn[:roleIndex]
		k.Role = arn[roleIndex+1:]
	} else {
		k.Arn = arn
	}
	k.CreationDate = time.Now().UTC()
	return k
}

// MasterKeysFromArnString takes a comma separated list of AWS KMS ARNs and returns a slice of new MasterKeys for those ARNs
func MasterKeysFromArnString(arn string) []MasterKey {
	var keys []MasterKey
	if arn == "" {
		return keys
	}
	for _, s := range strings.Split(arn, ",") {
		keys = append(keys, NewMasterKeyFromArn(s))
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
		return nil, err
	}
	config.Credentials = credentials.NewStaticCredentials(*out.Credentials.AccessKeyId,
		*out.Credentials.SecretAccessKey, *out.Credentials.SessionToken)
	sess, err = session.NewSession(&config)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (key MasterKey) createSession() (*session.Session, error) {
	re := regexp.MustCompile(`^arn:aws:kms:(.+):([0-9]+):key/(.+)$`)
	matches := re.FindStringSubmatch(key.Arn)
	if matches == nil {
		return nil, fmt.Errorf("No valid ARN found in %s", key.Arn)
	}
	config := aws.Config{Region: aws.String(matches[1])}
	sess, err := session.NewSession(&config)
	if err != nil {
		return nil, err
	}
	if key.Role != "" {
		return key.createStsSession(config, sess)
	}
	return sess, nil
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]string {
	out := make(map[string]string)
	out["arn"] = key.Arn
	if key.Role != "" {
		out["role"] = key.Role
	}
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}
