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
)

type encryptedValue struct {
	data     []byte
	iv       []byte
	tag      []byte
	datatype string
}

const nonceSize int = 32

type stashData struct {
	iv        []byte
	plaintext interface{}
}

// Cipher encrypts and decrypts data keys with AES GCM 256
type Cipher struct{}

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
func (c Cipher) Decrypt(value string, key []byte, path string) (plaintext interface{}, stash interface{}, err error) {
	if value == "" {
		return "", nil, nil
	}
	encryptedValue, err := parse(value)
	if err != nil {
		return "", nil, err
	}
	aescipher, err := cryptoaes.NewCipher(key)
	if err != nil {
		return "", nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(aescipher, len(encryptedValue.iv))
	if err != nil {
		return "", nil, err
	}
	stashValue := stashData{iv: encryptedValue.iv}
	data := append(encryptedValue.data, encryptedValue.tag...)
	decryptedBytes, err := gcm.Open(nil, encryptedValue.iv, data, []byte(path))
	if err != nil {
		return "", nil, fmt.Errorf("Could not decrypt with AES_GCM: %s", err)
	}
	decryptedValue := string(decryptedBytes)
	switch encryptedValue.datatype {
	case "str":
		stashValue.plaintext = decryptedValue
		return decryptedValue, stashValue, nil
	case "int":
		plaintext, err = strconv.Atoi(decryptedValue)
		stashValue.plaintext = plaintext
		return plaintext, stashValue, err
	case "float":
		plaintext, err = strconv.ParseFloat(decryptedValue, 64)
		stashValue.plaintext = plaintext
		return plaintext, stashValue, err
	case "bytes":
		stashValue.plaintext = decryptedBytes
		return decryptedBytes, stashValue, nil
	case "bool":
		plaintext, err = strconv.ParseBool(decryptedValue)
		stashValue.plaintext = plaintext
		return plaintext, stashValue, err
	default:
		return nil, nil, fmt.Errorf("Unknown datatype: %s", encryptedValue.datatype)
	}
}

// Encrypt takes one of (string, int, float, bool) and encrypts it with the provided key and additional auth data, returning a sops-format encrypted string.
func (c Cipher) Encrypt(value interface{}, key []byte, path string, stash interface{}) (string, error) {
	if value == "" {
		return "", nil
	}
	aescipher, err := cryptoaes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Could not initialize AES GCM encryption cipher: %s", err)
	}
	var iv []byte
	if stash, ok := stash.(stashData); !ok || stash.plaintext != value {
		iv = make([]byte, nonceSize)
		_, err = rand.Read(iv)
		if err != nil {
			return "", fmt.Errorf("Could not generate random bytes for IV: %s", err)
		}
	} else {
		iv = stash.iv
	}
	gcm, err := cipher.NewGCMWithNonceSize(aescipher, nonceSize)
	if err != nil {
		return "", fmt.Errorf("Could not create GCM: %s", err)
	}
	var plaintext []byte
	var encryptedType string
	switch value := value.(type) {
	case string:
		encryptedType = "str"
		plaintext = []byte(value)
	case int:
		encryptedType = "int"
		plaintext = []byte(strconv.Itoa(value))
	case float64:
		encryptedType = "float"
		// The Python version encodes floats without padding 0s after the decimal point.
		plaintext = []byte(strconv.FormatFloat(value, 'f', -1, 64))
	case bool:
		encryptedType = "bool"
		// The Python version encodes booleans with Titlecase
		plaintext = []byte(strings.Title(strconv.FormatBool(value)))
	default:
		return "", fmt.Errorf("Value to encrypt has unsupported type %T", value)
	}
	out := gcm.Seal(nil, iv, plaintext, []byte(path))
	return fmt.Sprintf("ENC[AES256_GCM,data:%s,iv:%s,tag:%s,type:%s]",
		base64.StdEncoding.EncodeToString(out[:len(out)-cryptoaes.BlockSize]),
		base64.StdEncoding.EncodeToString(iv),
		base64.StdEncoding.EncodeToString(out[len(out)-cryptoaes.BlockSize:]),
		encryptedType), nil
}
