package kms

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"regexp"
	"strings"
	"time"
)

type KMSMasterKey struct {
	Arn          string
	Role         string
	EncryptedKey string
	CreationDate time.Time
}

func (key *KMSMasterKey) Encrypt(dataKey []byte) error {
	sess, err := key.createSession()
	if err != nil {
		return err
	}
	service := kms.New(sess)
	out, err := service.Encrypt(&kms.EncryptInput{Plaintext: dataKey, KeyId: &key.Arn})
	if err != nil {
		return err
	}
	key.EncryptedKey = base64.StdEncoding.EncodeToString(out.CiphertextBlob)
	return nil
}

func (key *KMSMasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

func (key *KMSMasterKey) Decrypt() ([]byte, error) {
	k, err := base64.StdEncoding.DecodeString(key.EncryptedKey)
	if err != nil {
		return nil, fmt.Errorf("Error base64-decoding encrypted data key: %s", err)
	}
	sess, err := key.createSession()
	if err != nil {
		return nil, fmt.Errorf("Error creating AWS session: %v", err)
	}

	service := kms.New(sess)
	decrypted, err := service.Decrypt(&kms.DecryptInput{CiphertextBlob: k})
	if err != nil {
		return nil, fmt.Errorf("Error decrypting key: %v", err)
	}
	return decrypted.Plaintext, nil
}

func (key *KMSMasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

func (key *KMSMasterKey) ToString() string {
	return key.Arn
}

func NewKMSMasterKeyFromArn(arn string) KMSMasterKey {
	k := KMSMasterKey{}
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

func KMSMasterKeysFromArnString(arn string) []KMSMasterKey {
	var keys []KMSMasterKey
	if arn == "" {
		return keys
	}
	for _, s := range strings.Split(arn, ",") {
		keys = append(keys, NewKMSMasterKeyFromArn(s))
	}
	return keys
}

func (k KMSMasterKey) createStsSession(config aws.Config, sess *session.Session) (*session.Session, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	stsService := sts.New(sess)
	name := "sops@" + hostname
	out, err := stsService.AssumeRole(&sts.AssumeRoleInput{
		RoleArn: &k.Role, RoleSessionName: &name})
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

func (k KMSMasterKey) createSession() (*session.Session, error) {
	re := regexp.MustCompile(`^arn:aws:kms:(.+):([0-9]+):key/(.+)$`)
	matches := re.FindStringSubmatch(k.Arn)
	if matches == nil {
		return nil, fmt.Errorf("No valid ARN found in %s", k.Arn)
	}
	config := aws.Config{Region: aws.String(matches[1])}
	sess, err := session.NewSession(&config)
	if err != nil {
		return nil, err
	}
	if k.Role != "" {
		return k.createStsSession(config, sess)
	}
	return sess, nil
}

func (k KMSMasterKey) ToMap() map[string]string {
	out := make(map[string]string)
	out["arn"] = k.Arn
	if k.Role != "" {
		out["role"] = k.Role
	}
	out["created_at"] = k.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = k.EncryptedKey
	return out
}
