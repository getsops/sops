package sops

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"regexp"
)

// KeySource provides a way to obtain the symmetric encryption key used by sops
type KeySource interface {
	DecryptKey(encryptedKey string) string
	EncryptKey(key string) string
}

type KMS struct {
	Arn  string
	Role string
}

type KMSKeySource struct {
	KMS []KMS
}

type GPGKeySource struct{}

func (k KMS) createStsSession(config aws.Config, sess *session.Session) (*session.Session, error) {
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

func (k KMS) createSession() (*session.Session, error) {
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

func (k KMS) DecryptKey(encryptedKey string) (string, error) {
	sess, err := k.createSession()
	if err != nil {
		return "", fmt.Errorf("Error creating AWS session: %v", err)
	}

	service := kms.New(sess)
	decrypted, err := service.Decrypt(&kms.DecryptInput{CiphertextBlob: []byte(encryptedKey)})
	if err != nil {
		return "", fmt.Errorf("Error decrypting key: %v", err)
	}
	return string(decrypted.Plaintext), nil
}

func (ks KMSKeySource) DecryptKey(encryptedKey string) (string, error) {
	for _, kms := range ks.KMS {
		key, err := kms.DecryptKey(encryptedKey)
		if err != nil {
			return "", err
		}
		return key, nil
	}
	return "", fmt.Errorf("Could not decrypt key with KMS")
}

func (ks KMSKeySource) EncryptKey(key string) string {
	return key
}

func (gpg GPGKeySource) DecryptKey(encryptedKey string) string {
	return encryptedKey
}

func (gpg GPGKeySource) EncryptKey(key string) string {
	return key
}
