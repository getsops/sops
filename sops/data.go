package sops

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"

	"gopkg.in/yaml.v2"
)

type KMS struct {
	CreatedAt   time.Time `yaml:"created_at"`
	Enc         string
	decodedKey  []byte
	decodeError error
	Role        string
	Arn         string
}

var i = 0

func (k KMS) decodeKey() ([]byte, error) {
	base64Decoded, err := base64.StdEncoding.DecodeString(k.Enc)
	if err != nil {
		return nil, err
	}

	sess, err := k.AWSSession()
	if err != nil {
		return nil, fmt.Errorf("AWSSession: %v", err)
	}

	svc := kms.New(sess)
	params := &kms.DecryptInput{
		CiphertextBlob: base64Decoded,
	}
	out, err := svc.Decrypt(params)
	if err != nil {
		return nil, err
	}
	return out.Plaintext, nil
}

func (k KMS) AWSSession() (*session.Session, error) {
	re := regexp.MustCompile(`^arn:aws:kms:(.+):([0-9]+):key/(.+)$`)
	matches := re.FindStringSubmatch(k.Arn)
	if matches == nil {
		return nil, fmt.Errorf("Could not find valid ARN in %s", k.Arn)
	}

	return session.New(&aws.Config{Region: aws.String(matches[1])}), nil
}

func (k KMS) DecodeKey() ([]byte, error) {
	if k.decodedKey == nil && k.decodeError == nil {
		k.decodedKey, k.decodeError = k.decodeKey()
	}

	return k.decodedKey, k.decodeError
}

func (k KMS) Decrypt(val, iv, tag, additionalData []byte) ([]byte, error) {
	key, err := k.DecodeKey()
	if err != nil {
		return nil, fmt.Errorf("DecodeKey: %v", err)
	}

	aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(aes, len(iv))
	if err != nil {
		return nil, err
	}

	data := append(val, tag...)
	out, err := gcm.Open(nil, iv, data, additionalData)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type PGP struct {
	Fp        string
	CreatedAt time.Time `yaml:"created_at"`
	Enc       string
}

type SopsData struct {
	Mac               string
	Version           string
	KMS               []KMS
	PGP               []map[string]string
	LastModifed       time.Time
	UnencryptedSuffix string `yaml:"unencrypted_suffix"`
}

func NewSopsData(in []byte) (*SopsData, error) {
	sops := new(SopsData)
	err := yaml.Unmarshal(in, sops)
	return sops, err
}

func (s *SopsData) DecryptKMS(val, iv, tag, additionalData []byte) (string, error) {
	for _, kms := range s.KMS {
		out, err := kms.Decrypt(val, iv, tag, additionalData)
		if err != nil {
			log.Printf("DecryptKMS: %v", err)
			continue
		}
		return string(out), err
	}
	return "", errors.New("Decryption failed")
}

func (s *SopsData) DecryptString(in, accKey string) string {
	if s.UnencryptedSuffix != "" && strings.HasSuffix(accKey, fmt.Sprintf("%s:", s.UnencryptedSuffix)) {
		return in
	}
	encRegex := regexp.MustCompile(`^ENC\[AES256_GCM,data:(.+),iv:(.+),tag:(.+),type:(.+)\]`)
	matches := encRegex.FindStringSubmatch(in)
	if matches == nil {
		return in
	}
	data, err := base64.StdEncoding.DecodeString(matches[1])
	if err != nil {
		log.Printf("Error decoding data: %v", err)
		return in
	}
	iv, err := base64.StdEncoding.DecodeString(matches[2])
	if err != nil {
		log.Printf("Error decoding iv: %v", err)
		return in
	}
	tag, err := base64.StdEncoding.DecodeString(matches[3])
	if err != nil {
		log.Printf("Error decoding tag: %v", err)
		return in
	}

	out, err := s.DecryptKMS(data, iv, tag, []byte(accKey))
	if err != nil {
		log.Printf("Error decrypting data: %v", err)
		return in
	}

	return out
}

func (s *SopsData) DecryptValue(in, key interface{}, accKey string) interface{} {
	if key != nil {
		accKey = fmt.Sprintf("%v%v:", accKey, key)
	}
	switch in := in.(type) {
	case string:
		return s.DecryptString(in, accKey)
	case map[interface{}]interface{}:
		return s.DecryptMap(in, accKey)
	case []interface{}:
		return s.DecryptSlice(in, accKey)
	default:
		fmt.Printf("Could not decode type: %v\n", reflect.TypeOf(in))
	}
	return nil
}

func (s *SopsData) DecryptMap(in map[interface{}]interface{}, accKey string) map[interface{}]interface{} {
	branch := make(map[interface{}]interface{})
	for k, v := range in {
		branch[k] = s.DecryptValue(v, k, accKey)
	}

	return branch
}

func (s *SopsData) DecryptSlice(in []interface{}, accKey string) []interface{} {
	list := make([]interface{}, len(in))
	for i, v := range in {
		list[i] = s.DecryptValue(v, nil, accKey)
	}
	return list
}
