package sops

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/howeyc/gopass"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
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

type GPG struct {
	Fingerprint  string
	EncryptedKey string
}

type GPGKeySource struct {
	GPG []GPG
}

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
	errors := make([]error, 1)
	for _, kms := range ks.KMS {
		encKey, err := base64.StdEncoding.DecodeString(kms.EncryptedKey)
		if err != nil {
			continue
		}
		key, err := kms.DecryptKey(string(encKey))
		if err == nil {
			return key, nil
		}
		errors = append(errors, err)
	}
	return "", fmt.Errorf("The key could not be decrypted with any KMS entries", errors)
}

func (ks KMSKeySource) EncryptKeys(plaintext string) error {
	for i, _ := range ks.KMS {
		sess, err := ks.KMS[i].createSession()
		if err != nil {
			return err
		}
		service := kms.New(sess)
		out, err := service.Encrypt(&kms.EncryptInput{Plaintext: []byte(plaintext), KeyId: &ks.KMS[i].Arn})
		if err != nil {
			return err
		}
		ks.KMS[i].EncryptedKey = base64.StdEncoding.EncodeToString(out.CiphertextBlob)
	}
	return nil
}

func (gpg GPGKeySource) gpgHome() string {
	dir := os.Getenv("GNUPGHOME")
	if dir == "" {
		usr, err := user.Current()
		if err != nil {
			return "~/.gnupg"
		}
		return path.Join(usr.HomeDir, ".gnupg")
	}
	return dir
}

func (gpg GPGKeySource) loadRing(path string) (openpgp.EntityList, error) {
	f, err := os.Open(path)
	if err != nil {
		return openpgp.EntityList{}, err
	}
	defer f.Close()
	keyring, err := openpgp.ReadKeyRing(f)
	if err != nil {
		return keyring, err
	}
	return keyring, nil
}

func (gpg GPGKeySource) secRing() (openpgp.EntityList, error) {
	return gpg.loadRing(gpg.gpgHome() + "/secring.gpg")
}

func (gpg GPGKeySource) pubRing() (openpgp.EntityList, error) {
	return gpg.loadRing(gpg.gpgHome() + "/pubring.gpg")
}

func (gpg GPGKeySource) fingerprintMap(ring openpgp.EntityList) map[string]openpgp.Entity {
	fps := make(map[string]openpgp.Entity)
	for _, entity := range ring {
		fp := strings.ToUpper(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
		if entity != nil {
			fps[fp] = *entity
		}
	}
	return fps
}

func (gpg GPGKeySource) passphrasePrompt(keys []openpgp.Key, symmetric bool) ([]byte, error) {
	fmt.Print("Enter PGP key passphrase: ")
	psswd, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println(err)
	}
	return psswd, err
}

func (gpg GPGKeySource) DecryptKeys() (string, error) {
	ring, err := gpg.secRing()
	if err != nil {
		return "", fmt.Errorf("Could not load secring: %s", err)
	}
	for _, g := range gpg.GPG {
		block, err := armor.Decode(strings.NewReader(g.EncryptedKey))
		if err != nil {
			fmt.Println("Decode failed", err)
			continue
		}
		md, err := openpgp.ReadMessage(block.Body, ring, gpg.passphrasePrompt, nil)
		if err != nil {
			fmt.Println("ReadMessage failed", err)
			continue
		}
		if b, err := ioutil.ReadAll(md.UnverifiedBody); err == nil {
			return string(b), nil
		}
	}
	return "", fmt.Errorf("The key could not be decrypted with any of the GPG entries")
}

func (gpg GPGKeySource) EncryptKeys(plaintext string) error {
	ring, err := gpg.pubRing()
	if err != nil {
		return err
	}
	fingerprints := gpg.fingerprintMap(ring)
	for i, _ := range gpg.GPG {
		entity, ok := fingerprints[gpg.GPG[i].Fingerprint]
		if !ok {
			return fmt.Errorf("Key with fingerprint %s is not available in keyring.", gpg.GPG[i].Fingerprint)
		}
		encbuf := new(bytes.Buffer)
		armorbuf, err := armor.Encode(encbuf, "PGP MESSAGE", nil)
		if err != nil {
			return err
		}
		plaintextbuf, err := openpgp.Encrypt(armorbuf, []*openpgp.Entity{&entity}, nil, nil, nil)
		if err != nil {
			return err
		}
		_, err = plaintextbuf.Write([]byte(plaintext))
		if err != nil {
			return err
		}
		err = plaintextbuf.Close()
		if err != nil {
			return err
		}
		err = armorbuf.Close()
		if err != nil {
			return err
		}
		bytes, err := ioutil.ReadAll(encbuf)
		if err != nil {
			return err
		}
		gpg.GPG[i].EncryptedKey = string(bytes)
	}
	return nil
}
