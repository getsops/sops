package decryptor

import (
	"crypto/aes"
	"crypto/cipher"
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

func parse(value string) (*EncryptedValue, error) {
	re := regexp.MustCompile(`^ENC\[AES256_GCM,data:(.+),iv:(.+),tag:(.+),type:(.+)\]`)
	matches := re.FindStringSubmatch(value)
	if matches == nil {
		return nil, errors.New("Input string does not match sops' data format")
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
	aes, err := aes.NewCipher([]byte(key))
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
		return strconv.ParseInt(v, 10, 64)
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
