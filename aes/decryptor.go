package aes

import (
	cryptoaes "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
)

type encryptedValue struct {
	data     []byte
	iv       []byte
	tag      []byte
	datatype string
}

const nonceSize int = 32

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

// Decrypt takes a sops-format value string and a key and returns the decrypted value.
func (c Cipher) Decrypt(value string, key []byte, additionalAuthData []byte) (interface{}, error) {
	encryptedValue, err := parse(value)
	if err != nil {
		return "", err
	}
	aescipher, err := cryptoaes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCMWithNonceSize(aescipher, len(encryptedValue.iv))
	if err != nil {
		return "", err
	}

	data := append(encryptedValue.data, encryptedValue.tag...)
	decryptedBytes, err := gcm.Open(nil, encryptedValue.iv, data, additionalAuthData)
	if err != nil {
		return "", fmt.Errorf("Could not decrypt with AES_GCM: %s", err)
	}
	decryptedValue := string(decryptedBytes)
	switch encryptedValue.datatype {
	case "str":
		return decryptedValue, nil
	case "int":
		return strconv.Atoi(decryptedValue)
	case "float":
		return strconv.ParseFloat(decryptedValue, 64)
	case "bytes":
		return decryptedValue, nil
	case "bool":
		return strconv.ParseBool(decryptedValue)
	default:
		return nil, fmt.Errorf("Unknown datatype: %s", encryptedValue.datatype)
	}
}

// Encrypt takes one of (string, int, float, bool) and encrypts it with the provided key and additional auth data, returning a sops-format encrypted string.
func (c Cipher) Encrypt(value interface{}, key []byte, additionalAuthData []byte) (string, error) {
	aescipher, err := cryptoaes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Could not initialize AES GCM encryption cipher: %s", err)
	}
	iv := make([]byte, nonceSize)
	_, err = rand.Read(iv)
	if err != nil {
		return "", fmt.Errorf("Could not generate random bytes for IV: %s", err)
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
		plaintext = []byte(strconv.FormatFloat(value, 'f', 9, 64))
	case bool:
		encryptedType = "bool"
		plaintext = []byte(strconv.FormatBool(value))
	default:
		return "", fmt.Errorf("Value to encrypt has unsupported type %T", value)
	}
	out := gcm.Seal(nil, iv, plaintext, additionalAuthData)
	return fmt.Sprintf("ENC[AES256_GCM,data:%s,iv:%s,tag:%s,type:%s]",
		base64.StdEncoding.EncodeToString(out[:len(out)-cryptoaes.BlockSize]),
		base64.StdEncoding.EncodeToString(iv),
		base64.StdEncoding.EncodeToString(out[len(out)-cryptoaes.BlockSize:]),
		encryptedType), nil
}
