/*
Package aes defines a Cipher that uses 256-bit AES-GCM authenticated encryption to encrypt values the SOPS tree.
*/
package aes //import "go.mozilla.org/sops/aes"

import (
	cryptoaes "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops"
	"go.mozilla.org/sops/logging"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("AES")
}

type encryptedValue struct {
	data     []byte
	iv       []byte
	tag      []byte
	datatype string
}

const nonceSize int = 32

type stashKey struct {
	additionalData string
	plaintext      interface{}
}

// Cipher encrypts and decrypts data keys with AES GCM 256
type Cipher struct {
	// stash is a map that stores IVs for reuse, so that the ciphertext doesn't change when decrypting and reencrypting
	// the same values.
	stash map[stashKey][]byte
}

// NewCipher is the constructor for a new Cipher object
func NewCipher() Cipher {
	return Cipher{
		stash: make(map[stashKey][]byte),
	}
}

var encre = regexp.MustCompile(`^ENC\[AES256_GCM,data:(.+),iv:(.+),tag:(.+),type:(.+)\]`)

func parse(value string) (*encryptedValue, error) {
	matches := encre.FindStringSubmatch(value)
	if matches == nil {
		return nil, fmt.Errorf("Input string %s does not match sops' data format", value)
	}
	data, err := base64.StdEncoding.DecodeString(matches[1])
	if err != nil {
		return nil, fmt.Errorf("Error base64-decoding data: %s", err)
	}
	iv, err := base64.StdEncoding.DecodeString(matches[2])
	if err != nil {
		return nil, fmt.Errorf("Error base64-decoding iv: %s", err)
	}
	tag, err := base64.StdEncoding.DecodeString(matches[3])
	if err != nil {
		return nil, fmt.Errorf("Error base64-decoding tag: %s", err)
	}
	datatype := string(matches[4])

	return &encryptedValue{data, iv, tag, datatype}, nil
}

// Decrypt takes a sops-format value string and a key and returns the decrypted value and a stash value
func (c Cipher) Decrypt(ciphertext string, key []byte, additionalData string) (plaintext interface{}, err error) {
	if isEmpty(ciphertext) {
		return "", nil
	}
	encryptedValue, err := parse(ciphertext)
	if err != nil {
		return nil, err
	}
	aescipher, err := cryptoaes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(aescipher, len(encryptedValue.iv))
	if err != nil {
		return nil, err
	}
	data := append(encryptedValue.data, encryptedValue.tag...)
	decryptedBytes, err := gcm.Open(nil, encryptedValue.iv, data, []byte(additionalData))
	if err != nil {
		return nil, fmt.Errorf("Could not decrypt with AES_GCM: %s", err)
	}
	decryptedValue := string(decryptedBytes)
	switch encryptedValue.datatype {
	case "str":
		plaintext = decryptedValue
	case "int":
		plaintext, err = strconv.Atoi(decryptedValue)
	case "float":
		plaintext, err = strconv.ParseFloat(decryptedValue, 64)
	case "bytes":
		plaintext = decryptedBytes
	case "bool":
		plaintext, err = strconv.ParseBool(decryptedValue)
	case "comment":
		plaintext = sops.Comment{Value: decryptedValue}
	default:
		return nil, fmt.Errorf("Unknown datatype: %s", encryptedValue.datatype)
	}
	c.stash[stashKey{plaintext: plaintext, additionalData: additionalData}] = encryptedValue.iv
	return plaintext, err
}

func isEmpty(value interface{}) bool {
	switch value := value.(type) {
	case string:
		return value == ""
	case []byte:
		return len(value) == 0
	case sops.Comment:
		return isEmpty(value.Value)
	default:
		return false
	}
}

// Encrypt takes one of (string, int, float, bool) and encrypts it with the provided key and additional auth data, returning a sops-format encrypted string.
func (c Cipher) Encrypt(plaintext interface{}, key []byte, additionalData string) (ciphertext string, err error) {
	if isEmpty(plaintext) {
		return "", nil
	}
	aescipher, err := cryptoaes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Could not initialize AES GCM encryption cipher: %s", err)
	}
	var iv []byte
	if stash, ok := c.stash[stashKey{plaintext: plaintext, additionalData: additionalData}]; !ok {
		iv = make([]byte, nonceSize)
		_, err = rand.Read(iv)
		if err != nil {
			return "", fmt.Errorf("Could not generate random bytes for IV: %s", err)
		}
	} else {
		iv = stash
	}
	gcm, err := cipher.NewGCMWithNonceSize(aescipher, nonceSize)
	if err != nil {
		return "", fmt.Errorf("Could not create GCM: %s", err)
	}
	var plainBytes []byte
	var encryptedType string
	switch value := plaintext.(type) {
	case string:
		encryptedType = "str"
		plainBytes = []byte(value)
	case int:
		encryptedType = "int"
		plainBytes = []byte(strconv.Itoa(value))
	case float64:
		encryptedType = "float"
		// The Python version encodes floats without padding 0s after the decimal point.
		plainBytes = []byte(strconv.FormatFloat(value, 'f', -1, 64))
	case bool:
		encryptedType = "bool"
		// The Python version encodes booleans with Titlecase
		plainBytes = []byte(strings.Title(strconv.FormatBool(value)))
	case sops.Comment:
		encryptedType = "comment"
		plainBytes = []byte(value.Value)
	default:
		return "", fmt.Errorf("Value to encrypt has unsupported type %T", value)
	}
	out := gcm.Seal(nil, iv, plainBytes, []byte(additionalData))
	return fmt.Sprintf("ENC[AES256_GCM,data:%s,iv:%s,tag:%s,type:%s]",
		base64.StdEncoding.EncodeToString(out[:len(out)-cryptoaes.BlockSize]),
		base64.StdEncoding.EncodeToString(iv),
		base64.StdEncoding.EncodeToString(out[len(out)-cryptoaes.BlockSize:]),
		encryptedType), nil
}
