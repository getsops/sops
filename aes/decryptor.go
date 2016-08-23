package aes

import (
	cryptoaes "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type EncryptedValue struct {
	data     []byte
	iv       []byte
	tag      []byte
	datatype string
}

var encre = regexp.MustCompile(`^ENC\[AES256_GCM,data:(.+),iv:(.+),tag:(.+),type:(.+)\]`)

func parse(value string) (*EncryptedValue, error) {
	matches := encre.FindStringSubmatch(value)
	if matches == nil {
		return nil, fmt.Errorf("Input string %s does not match sops' data format", value)
	}
	data, err := base64.StdEncoding.DecodeString(matches[1])
	if err != nil {
		return nil, errors.New("Error base64-decoding data")
	}
	iv, err := base64.StdEncoding.DecodeString(matches[2])
	if err != nil {
		return nil, errors.New("Error base64-decoding iv")
	}
	tag, err := base64.StdEncoding.DecodeString(matches[3])
	if err != nil {
		return nil, errors.New("Error base64-decoding tag")
	}
	datatype := matches[4]

	return &EncryptedValue{data, iv, tag, datatype}, nil
}

// Decrypt takes a sops-format value string and a key and returns the decrypted value.
func Decrypt(value, key string, additionalAuthData []byte) (interface{}, error) {
	encryptedValue, err := parse(value)
	if err != nil {
		return "", err
	}
	aes, err := cryptoaes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCMWithNonceSize(aes, len(encryptedValue.iv))
	if err != nil {
		return "", err
	}

	data := append(encryptedValue.data, encryptedValue.tag...)
	out, err := gcm.Open(nil, encryptedValue.iv, data, additionalAuthData)
	if err != nil {
		return "", fmt.Errorf("Could not decrypt with AES_GCM: %s", err)
	}
	v := string(out)
	switch encryptedValue.datatype {
	case "str":
		return v, nil
	case "int":
		return strconv.Atoi(v)
	case "float":
		return strconv.ParseFloat(v, 64)
	case "bytes":
		return v, nil
	case "bool":
		return strconv.ParseBool(v)
	default:
		return nil, fmt.Errorf("Unknown datatype: %s", encryptedValue.datatype)
	}
}

func Encrypt(value interface{}, key string, additionalAuthData []byte) (string, error) {
	aes, err := cryptoaes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("Could not create AES Cipher: %s", err)
	}
	iv := make([]byte, 32)
	_, err = rand.Read(iv)
	if err != nil {
		return "", fmt.Errorf("Could not generate random bytes for IV: %s", err)
	}
	gcm, err := cipher.NewGCMWithNonceSize(aes, len(iv))
	if err != nil {
		return "", fmt.Errorf("Could not create GCM: %s", err)
	}
	var plaintext []byte
	var t string
	switch value := value.(type) {
	case string:
		t = "str"
		plaintext = []byte(value)
	case int:
		t = "int"
		plaintext = []byte(strconv.Itoa(value))
	case float64:
		t = "float"
		plaintext = []byte(strconv.FormatFloat(value, 'f', 9, 64))
	case bool:
		t = "bool"
		plaintext = []byte(strconv.FormatBool(value))
	default:
		return "", fmt.Errorf("Value to encrypt has unsupported type %T", value)
	}

	out := gcm.Seal(nil, iv, plaintext, additionalAuthData)
	ciphertext := out[:len(out)-16]
	b64ciphertext := base64.StdEncoding.EncodeToString(ciphertext)
	tag := out[len(out)-16:]
	b64tag := base64.StdEncoding.EncodeToString(tag)
	b64iv := base64.StdEncoding.EncodeToString(iv)
	return fmt.Sprintf("ENC[AES256_GCM,data:%s,iv:%s,tag:%s,type:%s]", b64ciphertext, b64iv, b64tag, t), nil
}
