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
	DecryptKeys() (string, error)
	EncryptKeys(plaintext string) error
}

type KMS struct {
	Arn          string
	Role         string
	EncryptedKey string
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

func (ks KMSKeySource) DecryptKeys() (string, error) {
	for _, kms := range ks.KMS {
		key, err := kms.DecryptKey(kms.EncryptedKey)
		if err == nil {
			return key, nil
		}
	}
	return "", fmt.Errorf("The key could not be decrypted with any KMS entries")
}

func (ks KMSKeySource) EncryptKeys(plaintext string) error {
	for _, k := range ks.KMS {
		sess, err := k.createSession()
		if err != nil {
			return err
		}
		service := kms.New(sess)
		out, err := service.Encrypt(&kms.EncryptInput{Plaintext: []byte(plaintext)})
		if err != nil {
			return err
		}
		k.EncryptedKey = string(out.CiphertextBlob)
	}
	return nil
}

func (gpg GPGKeySource) DecryptKeys() (string, error) {
	return "", nil
}

func (gpg GPGKeySource) EncryptKeys(plaintext string) error {
	return nil
}
