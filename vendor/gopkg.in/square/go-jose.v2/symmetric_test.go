/*-
 * Copyright 2014 Square Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jose

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"testing"

	"golang.org/x/crypto/pbkdf2"
)

func TestInvalidSymmetricAlgorithms(t *testing.T) {
	_, err := newSymmetricRecipient("XYZ", []byte{})
	if err != ErrUnsupportedAlgorithm {
		t.Error("should not accept invalid algorithm")
	}

	enc := &symmetricKeyCipher{}
	_, err = enc.encryptKey([]byte{}, "XYZ")
	if err != ErrUnsupportedAlgorithm {
		t.Error("should not accept invalid algorithm")
	}
}

func TestAeadErrors(t *testing.T) {
	aead := &aeadContentCipher{
		keyBytes:     16,
		authtagBytes: 16,
		getAead: func(key []byte) (cipher.AEAD, error) {
			return nil, ErrCryptoFailure
		},
	}

	parts, err := aead.encrypt([]byte{}, []byte{}, []byte{})
	if err != ErrCryptoFailure {
		t.Error("should handle aead failure")
	}

	_, err = aead.decrypt([]byte{}, []byte{}, parts)
	if err != ErrCryptoFailure {
		t.Error("should handle aead failure")
	}
}

func TestInvalidKey(t *testing.T) {
	gcm := newAESGCM(16).(*aeadContentCipher)
	_, err := gcm.getAead([]byte{})
	if err == nil {
		t.Error("should not accept invalid key")
	}
}

func TestStaticKeyGen(t *testing.T) {
	key := make([]byte, 32)
	io.ReadFull(rand.Reader, key)

	gen := &staticKeyGenerator{key: key}
	if gen.keySize() != len(key) {
		t.Error("static key generator reports incorrect size")
	}

	generated, _, err := gen.genKey()
	if err != nil {
		t.Error("static key generator should always succeed", err)
	}
	if !bytes.Equal(generated, key) {
		t.Error("static key generator returns different data")
	}
}

func TestAeadInvalidInput(t *testing.T) {
	sample := []byte("1234567890123456")
	tt := []aeadParts{
		{},
		{iv: sample, tag: sample},
	}
	for _, tc := range tt {
		aead := newAESGCM(16).(*aeadContentCipher)
		_, err := aead.decrypt(sample, []byte{}, &tc)
		if err != ErrCryptoFailure {
			t.Error("should handle aead failure")
		}
	}
}

func TestVectorsAESGCM(t *testing.T) {
	// Source: http://tools.ietf.org/html/draft-ietf-jose-json-web-encryption-29#appendix-A.1
	plaintext := []byte{
		84, 104, 101, 32, 116, 114, 117, 101, 32, 115, 105, 103, 110, 32,
		111, 102, 32, 105, 110, 116, 101, 108, 108, 105, 103, 101, 110, 99,
		101, 32, 105, 115, 32, 110, 111, 116, 32, 107, 110, 111, 119, 108,
		101, 100, 103, 101, 32, 98, 117, 116, 32, 105, 109, 97, 103, 105,
		110, 97, 116, 105, 111, 110, 46}

	aad := []byte{
		101, 121, 74, 104, 98, 71, 99, 105, 79, 105, 74, 83, 85, 48, 69,
		116, 84, 48, 70, 70, 85, 67, 73, 115, 73, 109, 86, 117, 89, 121, 73,
		54, 73, 107, 69, 121, 78, 84, 90, 72, 81, 48, 48, 105, 102, 81}

	expectedCiphertext := []byte{
		229, 236, 166, 241, 53, 191, 115, 196, 174, 43, 73, 109, 39, 122,
		233, 96, 140, 206, 120, 52, 51, 237, 48, 11, 190, 219, 186, 80, 111,
		104, 50, 142, 47, 167, 59, 61, 181, 127, 196, 21, 40, 82, 242, 32,
		123, 143, 168, 226, 73, 216, 176, 144, 138, 247, 106, 60, 16, 205,
		160, 109, 64, 63, 192}

	expectedAuthtag := []byte{
		92, 80, 104, 49, 133, 25, 161, 215, 173, 101, 219, 211, 136, 91, 210, 145}

	// Mock random reader
	RandReader = bytes.NewReader([]byte{
		177, 161, 244, 128, 84, 143, 225, 115, 63, 180, 3, 255, 107, 154,
		212, 246, 138, 7, 110, 91, 112, 46, 34, 105, 47, 130, 203, 46, 122,
		234, 64, 252, 227, 197, 117, 252, 2, 219, 233, 68, 180, 225, 77, 219})
	defer resetRandReader()

	enc := newAESGCM(32)
	key, _, _ := randomKeyGenerator{size: 32}.genKey()
	out, err := enc.encrypt(key, aad, plaintext)
	if err != nil {
		t.Error("Unable to encrypt:", err)
		return
	}

	if bytes.Compare(out.ciphertext, expectedCiphertext) != 0 {
		t.Error("Ciphertext did not match")
	}
	if bytes.Compare(out.tag, expectedAuthtag) != 0 {
		t.Error("Auth tag did not match")
	}
}

func TestVectorPBES2_HS256A_128KW(t *testing.T) {
	cipher := &symmetricKeyCipher{
		key: []byte("Thus from my lips, by yours, my sin is purged."),
		p2c: 4096,
		p2s: []byte{
			217, 96, 147, 112, 150, 117, 70,
			247, 127, 8, 155, 137, 174, 42, 80, 215,
		},
	}

	cek := []byte{
		111, 27, 25, 52, 66, 29, 20, 78, 92, 176, 56, 240, 65, 208, 82, 112,
		161, 131, 36, 55, 202, 236, 185, 172, 129, 23, 153, 194, 195, 48,
		253, 182,
	}

	// PBES2-HS256+A128KW || 0x00 || p2s
	salt := []byte{
		80, 66, 69, 83, 50, 45, 72, 83, 50, 53, 54, 43, 65, 49, 50, 56, 75,
		87, 0, 217, 96, 147, 112, 150, 117, 70, 247, 127, 8, 155, 137, 174,
		42, 80, 215,
	}

	expectedDerivedKey := []byte{
		110, 171, 169, 92, 129, 92, 109, 117, 233, 242, 116, 233, 170, 14,
		24, 75}

	expectedEncryptedKey := []byte{
		78, 186, 151, 59, 11, 141, 81, 240, 213, 245, 83, 211, 53, 188, 134,
		188, 66, 125, 36, 200, 222, 124, 5, 103, 249, 52, 117, 184, 140, 81,
		246, 158, 161, 177, 20, 33, 245, 57, 59, 4}

	derivedKey := pbkdf2.Key(cipher.key, salt, cipher.p2c, 16, sha256.New)
	if bytes.Compare(derivedKey, expectedDerivedKey) != 0 {
		t.Error("Derived key did not match")
	}

	encryptedKey, err := cipher.encryptKey(cek, PBES2_HS256_A128KW)
	if err != nil {
		t.Fatal("Unable to encrypt:", err)
	}

	if bytes.Compare(encryptedKey.encryptedKey, expectedEncryptedKey) != 0 {
		t.Error("Encrypted key did not match")
	}
}
