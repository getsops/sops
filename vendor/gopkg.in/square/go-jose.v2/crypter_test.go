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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"regexp"
	"testing"

	"golang.org/x/crypto/ed25519"
)

// We generate only a single RSA and EC key for testing, speeds up tests.
var rsaTestKey, _ = rsa.GenerateKey(rand.Reader, 2048)

var ecTestKey256, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
var ecTestKey384, _ = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
var ecTestKey521, _ = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

var ed25519PublicKey, ed25519PrivateKey, _ = ed25519.GenerateKey(rand.Reader)

func RoundtripJWE(keyAlg KeyAlgorithm, encAlg ContentEncryption, compressionAlg CompressionAlgorithm, serializer func(*JSONWebEncryption) (string, error), corrupter func(*JSONWebEncryption) bool, aad []byte, encryptionKey interface{}, decryptionKey interface{}) error {
	var rcpt Recipient
	switch keyAlg {
	case PBES2_HS256_A128KW, PBES2_HS384_A192KW, PBES2_HS512_A256KW:
		// use 1k iterations instead of 100k to reduce computational cost
		rcpt = Recipient{Algorithm: keyAlg, Key: encryptionKey, PBES2Count: 1000}
	default:
		rcpt = Recipient{Algorithm: keyAlg, Key: encryptionKey}
	}
	enc, err := NewEncrypter(encAlg, rcpt, &EncrypterOptions{Compression: compressionAlg})
	if err != nil {
		return fmt.Errorf("error on new encrypter: %s", err)
	}

	input := []byte("Lorem ipsum dolor sit amet")
	obj, err := enc.EncryptWithAuthData(input, aad)
	if err != nil {
		return fmt.Errorf("error in encrypt: %s", err)
	}

	msg, err := serializer(obj)
	if err != nil {
		return fmt.Errorf("error in serializer: %s", err)
	}

	parsed, err := ParseEncrypted(msg)
	if err != nil {
		return fmt.Errorf("error in parse: %s, on msg '%s'", err, msg)
	}

	// (Maybe) mangle object
	skip := corrupter(parsed)
	if skip {
		return fmt.Errorf("corrupter indicated message should be skipped")
	}

	if bytes.Compare(parsed.GetAuthData(), aad) != 0 {
		return fmt.Errorf("auth data in parsed object does not match")
	}

	output, err := parsed.Decrypt(decryptionKey)
	if err != nil {
		return fmt.Errorf("error on decrypt: %s", err)
	}

	if bytes.Compare(input, output) != 0 {
		return fmt.Errorf("Decrypted output does not match input, got '%s' but wanted '%s'", output, input)
	}

	return nil
}

func TestRoundtripsJWE(t *testing.T) {
	// Test matrix
	keyAlgs := []KeyAlgorithm{
		DIRECT, ECDH_ES, ECDH_ES_A128KW, ECDH_ES_A192KW, ECDH_ES_A256KW, A128KW, A192KW, A256KW,
		RSA1_5, RSA_OAEP, RSA_OAEP_256, A128GCMKW, A192GCMKW, A256GCMKW,
		PBES2_HS256_A128KW, PBES2_HS384_A192KW, PBES2_HS512_A256KW,
	}
	encAlgs := []ContentEncryption{A128GCM, A192GCM, A256GCM, A128CBC_HS256, A192CBC_HS384, A256CBC_HS512}
	zipAlgs := []CompressionAlgorithm{NONE, DEFLATE}

	serializers := []func(*JSONWebEncryption) (string, error){
		func(obj *JSONWebEncryption) (string, error) { return obj.CompactSerialize() },
		func(obj *JSONWebEncryption) (string, error) { return obj.FullSerialize(), nil },
	}

	corrupter := func(obj *JSONWebEncryption) bool { return false }

	// Note: can't use AAD with compact serialization
	aads := [][]byte{
		nil,
		[]byte("Ut enim ad minim veniam"),
	}

	// Test all different configurations
	for _, alg := range keyAlgs {
		for _, enc := range encAlgs {
			for _, key := range generateTestKeys(alg, enc) {
				for _, zip := range zipAlgs {
					for i, serializer := range serializers {
						err := RoundtripJWE(alg, enc, zip, serializer, corrupter, aads[i], key.enc, key.dec)
						if err != nil {
							t.Error(err, alg, enc, zip, i)
						}
					}
				}
			}
		}
	}
}

func TestRoundtripsJWECorrupted(t *testing.T) {
	// Test matrix
	keyAlgs := []KeyAlgorithm{DIRECT, ECDH_ES, ECDH_ES_A128KW, A128KW, RSA1_5, RSA_OAEP, RSA_OAEP_256, A128GCMKW, PBES2_HS256_A128KW}
	encAlgs := []ContentEncryption{A128GCM, A192GCM, A256GCM, A128CBC_HS256, A192CBC_HS384, A256CBC_HS512}
	zipAlgs := []CompressionAlgorithm{NONE, DEFLATE}

	serializers := []func(*JSONWebEncryption) (string, error){
		func(obj *JSONWebEncryption) (string, error) { return obj.CompactSerialize() },
		func(obj *JSONWebEncryption) (string, error) { return obj.FullSerialize(), nil },
	}

	bitflip := func(slice []byte) bool {
		if len(slice) > 0 {
			slice[0] ^= 0xFF
			return false
		}
		return true
	}

	corrupters := []func(*JSONWebEncryption) bool{
		func(obj *JSONWebEncryption) bool {
			// Set invalid ciphertext
			return bitflip(obj.ciphertext)
		},
		func(obj *JSONWebEncryption) bool {
			// Set invalid auth tag
			return bitflip(obj.tag)
		},
		func(obj *JSONWebEncryption) bool {
			// Set invalid AAD
			return bitflip(obj.aad)
		},
		func(obj *JSONWebEncryption) bool {
			// Mess with encrypted key
			return bitflip(obj.recipients[0].encryptedKey)
		},
		func(obj *JSONWebEncryption) bool {
			// Mess with GCM-KW auth tag
			tag, _ := obj.protected.getTag()
			skip := bitflip(tag.bytes())
			if skip {
				return true
			}
			obj.protected.set(headerTag, tag)
			return false
		},
	}

	// Note: can't use AAD with compact serialization
	aads := [][]byte{
		nil,
		[]byte("Ut enim ad minim veniam"),
	}

	// Test all different configurations
	for _, alg := range keyAlgs {
		for _, enc := range encAlgs {
			for _, key := range generateTestKeys(alg, enc) {
				for _, zip := range zipAlgs {
					for i, serializer := range serializers {
						for j, corrupter := range corrupters {
							err := RoundtripJWE(alg, enc, zip, serializer, corrupter, aads[i], key.enc, key.dec)
							if err == nil {
								t.Error("failed to detect corrupt data", err, alg, enc, zip, i, j)
							}
						}
					}
				}
			}
		}
	}
}

func TestEncrypterWithJWKAndKeyID(t *testing.T) {
	enc, err := NewEncrypter(A128GCM, Recipient{Algorithm: A128KW, Key: &JSONWebKey{
		KeyID: "test-id",
		Key:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	}}, nil)
	if err != nil {
		t.Error(err)
	}

	ciphertext, _ := enc.Encrypt([]byte("Lorem ipsum dolor sit amet"))

	serialized1, _ := ciphertext.CompactSerialize()
	serialized2 := ciphertext.FullSerialize()

	parsed1, _ := ParseEncrypted(serialized1)
	parsed2, _ := ParseEncrypted(serialized2)

	if parsed1.Header.KeyID != "test-id" {
		t.Errorf("expected message to have key id from JWK, but found '%s' instead", parsed1.Header.KeyID)
	}
	if parsed2.Header.KeyID != "test-id" {
		t.Errorf("expected message to have key id from JWK, but found '%s' instead", parsed2.Header.KeyID)
	}
}

func TestEncrypterWithBrokenRand(t *testing.T) {
	keyAlgs := []KeyAlgorithm{ECDH_ES_A128KW, A128KW, RSA1_5, RSA_OAEP, RSA_OAEP_256, A128GCMKW, PBES2_HS256_A128KW}
	encAlgs := []ContentEncryption{A128GCM, A192GCM, A256GCM, A128CBC_HS256, A192CBC_HS384, A256CBC_HS512}

	serializer := func(obj *JSONWebEncryption) (string, error) { return obj.CompactSerialize() }
	corrupter := func(obj *JSONWebEncryption) bool { return false }

	// Break rand reader
	readers := []func() io.Reader{
		// Totally broken
		func() io.Reader { return bytes.NewReader([]byte{}) },
		// Not enough bytes
		func() io.Reader { return io.LimitReader(rand.Reader, 20) },
	}

	defer resetRandReader()

	for _, alg := range keyAlgs {
		for _, enc := range encAlgs {
			for _, key := range generateTestKeys(alg, enc) {
				for i, getReader := range readers {
					RandReader = getReader()
					err := RoundtripJWE(alg, enc, NONE, serializer, corrupter, nil, key.enc, key.dec)
					if err == nil {
						t.Error("encrypter should fail if rand is broken", i)
					}
				}
			}
		}
	}
}

func TestNewEncrypterErrors(t *testing.T) {
	_, err := NewEncrypter("XYZ", Recipient{}, nil)
	if err == nil {
		t.Error("was able to instantiate encrypter with invalid cipher")
	}

	_, err = NewMultiEncrypter("XYZ", []Recipient{}, nil)
	if err == nil {
		t.Error("was able to instantiate multi-encrypter with invalid cipher")
	}

	_, err = NewEncrypter(A128GCM, Recipient{Algorithm: DIRECT, Key: nil}, nil)
	if err == nil {
		t.Error("was able to instantiate encrypter with invalid direct key")
	}

	_, err = NewEncrypter(A128GCM, Recipient{Algorithm: ECDH_ES, Key: nil}, nil)
	if err == nil {
		t.Error("was able to instantiate encrypter with invalid EC key")
	}
}

func TestMultiRecipientJWE(t *testing.T) {
	sharedKey := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	}

	enc, err := NewMultiEncrypter(A128GCM, []Recipient{
		{Algorithm: RSA_OAEP, Key: &rsaTestKey.PublicKey},
		{Algorithm: A256GCMKW, Key: sharedKey},
	}, nil)
	if err != nil {
		panic(err)
	}

	input := []byte("Lorem ipsum dolor sit amet")
	obj, err := enc.Encrypt(input)
	if err != nil {
		t.Fatal("error in encrypt: ", err)
	}

	msg := obj.FullSerialize()

	parsed, err := ParseEncrypted(msg)
	if err != nil {
		t.Fatal("error in parse: ", err)
	}

	i, _, output, err := parsed.DecryptMulti(rsaTestKey)
	if err != nil {
		t.Fatal("error on decrypt with RSA: ", err)
	}

	if i != 0 {
		t.Fatal("recipient index should be 0 for RSA key")
	}

	if bytes.Compare(input, output) != 0 {
		t.Fatal("Decrypted output does not match input: ", output, input)
	}

	i, _, output, err = parsed.DecryptMulti(sharedKey)
	if err != nil {
		t.Fatal("error on decrypt with AES: ", err)
	}

	if i != 1 {
		t.Fatal("recipient index should be 1 for shared key")
	}

	if bytes.Compare(input, output) != 0 {
		t.Fatal("Decrypted output does not match input", output, input)
	}
}

func TestMultiRecipientErrors(t *testing.T) {
	_, err := NewMultiEncrypter(A128GCM, []Recipient{}, nil)
	if err == nil {
		t.Error("should fail to instantiate with zero recipients")
	}
}

func TestEncrypterOptions(t *testing.T) {
	sharedKey := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	}

	opts := &EncrypterOptions{
		Compression: DEFLATE,
	}
	opts.WithType("JWT")
	opts.WithContentType("JWT")
	enc, err := NewEncrypter(A256GCM, Recipient{Algorithm: A256GCMKW, Key: sharedKey}, opts)
	if err != nil {
		fmt.Println(err)
		t.Error("Failed to create encrypter")
	}

	if !reflect.DeepEqual(*opts, enc.Options()) {
		t.Error("Encrypter options do not match")
	}
}

// Test that extra headers are generated and parsed in a round trip.
func TestEncrypterExtraHeaderInclusion(t *testing.T) {
	sharedKey := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	}

	opts := &EncrypterOptions{
		Compression: DEFLATE,
	}
	opts.WithType("JWT")
	opts.WithContentType("JWT")
	opts.WithHeader(HeaderKey("myCustomHeader"), "xyz")
	enc, err := NewEncrypter(A256GCM, Recipient{Algorithm: A256GCMKW, Key: sharedKey}, opts)
	if err != nil {
		fmt.Println(err)
		t.Error("Failed to create encrypter")
	}

	if !reflect.DeepEqual(*opts, enc.Options()) {
		t.Error("Encrypter options do not match")
	}

	input := []byte("Lorem ipsum dolor sit amet")
	obj, err := enc.Encrypt(input)
	if err != nil {
		t.Fatal("error in encrypt: ", err)
	}

	parsed, err := ParseEncrypted(obj.FullSerialize())
	if err != nil {
		t.Fatal("error in parse: ", err)
	}

	output, err := parsed.Decrypt(sharedKey)
	if err != nil {
		t.Fatal("error on decrypt: ", err)
	}

	if bytes.Compare(input, output) != 0 {
		t.Fatal("Decrypted output does not match input: ", output, input)
	}

	if parsed.Header.ExtraHeaders[HeaderType] != "JWT" ||
		parsed.Header.ExtraHeaders[HeaderContentType] != "JWT" ||
		parsed.Header.ExtraHeaders[HeaderKey("myCustomHeader")] != "xyz" {
		t.Fatalf("Mismatch in extra headers: %#v", parsed.Header.ExtraHeaders)
	}
}

// TestPBES2JWKEncryption uses the plaintext and serialization reference of
// JWK RFC https://tools.ietf.org/html/rfc7517#appendix-C.4
func TestPBES2JWKEncryption(t *testing.T) {
	passphrase := []byte("Thus from my lips, by yours, my sin is purged.")

	plaintext := []byte(`{
      "kty":"RSA",
      "kid":"juliet@capulet.lit",
      "use":"enc",
      "n":"t6Q8PWSi1dkJj9hTP8hNYFlvadM7DflW9mWepOJhJ66w7nyoK1gPNqFMSQRy
           O125Gp-TEkodhWr0iujjHVx7BcV0llS4w5ACGgPrcAd6ZcSR0-Iqom-QFcNP
           8Sjg086MwoqQU_LYywlAGZ21WSdS_PERyGFiNnj3QQlO8Yns5jCtLCRwLHL0
           Pb1fEv45AuRIuUfVcPySBWYnDyGxvjYGDSM-AqWS9zIQ2ZilgT-GqUmipg0X
           OC0Cc20rgLe2ymLHjpHciCKVAbY5-L32-lSeZO-Os6U15_aXrk9Gw8cPUaX1
           _I8sLGuSiVdt3C_Fn2PZ3Z8i744FPFGGcG1qs2Wz-Q",
      "e":"AQAB",
      "d":"GRtbIQmhOZtyszfgKdg4u_N-R_mZGU_9k7JQ_jn1DnfTuMdSNprTeaSTyWfS
           NkuaAwnOEbIQVy1IQbWVV25NY3ybc_IhUJtfri7bAXYEReWaCl3hdlPKXy9U
           vqPYGR0kIXTQRqns-dVJ7jahlI7LyckrpTmrM8dWBo4_PMaenNnPiQgO0xnu
           ToxutRZJfJvG4Ox4ka3GORQd9CsCZ2vsUDmsXOfUENOyMqADC6p1M3h33tsu
           rY15k9qMSpG9OX_IJAXmxzAh_tWiZOwk2K4yxH9tS3Lq1yX8C1EWmeRDkK2a
           hecG85-oLKQt5VEpWHKmjOi_gJSdSgqcN96X52esAQ",
      "p":"2rnSOV4hKSN8sS4CgcQHFbs08XboFDqKum3sc4h3GRxrTmQdl1ZK9uw-PIHf
           QP0FkxXVrx-WE-ZEbrqivH_2iCLUS7wAl6XvARt1KkIaUxPPSYB9yk31s0Q8
           UK96E3_OrADAYtAJs-M3JxCLfNgqh56HDnETTQhH3rCT5T3yJws",
      "q":"1u_RiFDP7LBYh3N4GXLT9OpSKYP0uQZyiaZwBtOCBNJgQxaj10RWjsZu0c6I
           edis4S7B_coSKB0Kj9PaPaBzg-IySRvvcQuPamQu66riMhjVtG6TlV8CLCYK
           rYl52ziqK0E_ym2QnkwsUX7eYTB7LbAHRK9GqocDE5B0f808I4s",
      "dp":"KkMTWqBUefVwZ2_Dbj1pPQqyHSHjj90L5x_MOzqYAJMcLMZtbUtwKqvVDq3
           tbEo3ZIcohbDtt6SbfmWzggabpQxNxuBpoOOf_a_HgMXK_lhqigI4y_kqS1w
           Y52IwjUn5rgRrJ-yYo1h41KR-vz2pYhEAeYrhttWtxVqLCRViD6c",
      "dq":"AvfS0-gRxvn0bwJoMSnFxYcK1WnuEjQFluMGfwGitQBWtfZ1Er7t1xDkbN9
           GQTB9yqpDoYaN06H7CFtrkxhJIBQaj6nkF5KKS3TQtQ5qCzkOkmxIe3KRbBy
           mXxkb5qwUpX5ELD5xFc6FeiafWYY63TmmEAu_lRFCOJ3xDea-ots",
      "qi":"lSQi-w9CpyUReMErP1RsBLk7wNtOvs5EQpPqmuMvqW57NBUczScEoPwmUqq
           abu9V0-Py4dQ57_bapoKRu1R90bvuFnU63SHWEFglZQvJDMeAvmj4sm-Fp0o
           Yu_neotgQ0hzbI5gry7ajdYy9-2lNx_76aBZoOUu9HCJ-UsfSOI8"
     }`)

	serializationReference := `
	  eyJhbGciOiJQQkVTMi1IUzI1NitBMTI4S1ciLCJwMnMiOiIyV0NUY0paMVJ2ZF9DSn
	  VKcmlwUTF3IiwicDJjIjo0MDk2LCJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwiY3R5Ijoi
	  andrK2pzb24ifQ.
	  TrqXOwuNUfDV9VPTNbyGvEJ9JMjefAVn-TR1uIxR9p6hsRQh9Tk7BA.
	  Ye9j1qs22DmRSAddIh-VnA.
	  AwhB8lxrlKjFn02LGWEqg27H4Tg9fyZAbFv3p5ZicHpj64QyHC44qqlZ3JEmnZTgQo
	  wIqZJ13jbyHB8LgePiqUJ1hf6M2HPLgzw8L-mEeQ0jvDUTrE07NtOerBk8bwBQyZ6g
	  0kQ3DEOIglfYxV8-FJvNBYwbqN1Bck6d_i7OtjSHV-8DIrp-3JcRIe05YKy3Oi34Z_
	  GOiAc1EK21B11c_AE11PII_wvvtRiUiG8YofQXakWd1_O98Kap-UgmyWPfreUJ3lJP
	  nbD4Ve95owEfMGLOPflo2MnjaTDCwQokoJ_xplQ2vNPz8iguLcHBoKllyQFJL2mOWB
	  wqhBo9Oj-O800as5mmLsvQMTflIrIEbbTMzHMBZ8EFW9fWwwFu0DWQJGkMNhmBZQ-3
	  lvqTc-M6-gWA6D8PDhONfP2Oib2HGizwG1iEaX8GRyUpfLuljCLIe1DkGOewhKuKkZ
	  h04DKNM5Nbugf2atmU9OP0Ldx5peCUtRG1gMVl7Qup5ZXHTjgPDr5b2N731UooCGAU
	  qHdgGhg0JVJ_ObCTdjsH4CF1SJsdUhrXvYx3HJh2Xd7CwJRzU_3Y1GxYU6-s3GFPbi
	  rfqqEipJDBTHpcoCmyrwYjYHFgnlqBZRotRrS95g8F95bRXqsaDY7UgQGwBQBwy665
	  d0zpvTasvfXf_c0MWAl-neFaKOW_Px6g4EUDjG1GWSXV9cLStLw_0ovdApDIFLHYHe
	  PyagyHjouQUuGiq7BsYwYrwaF06tgB8hV8omLNfMEmDPJaZUzMuHw6tBDwGkzD-tS_
	  ub9hxrpJ4UsOWnt5rGUyoN2N_c1-TQlXxm5oto14MxnoAyBQBpwIEgSH3Y4ZhwKBhH
	  PjSo0cdwuNdYbGPpb-YUvF-2NZzODiQ1OvWQBRHSbPWYz_xbGkgD504LRtqRwCO7CC
	  _CyyURi1sEssPVsMJRX_U4LFEOc82TiDdqjKOjRUfKK5rqLi8nBE9soQ0DSaOoFQZi
	  GrBrqxDsNYiAYAmxxkos-i3nX4qtByVx85sCE5U_0MqG7COxZWMOPEFrDaepUV-cOy
	  rvoUIng8i8ljKBKxETY2BgPegKBYCxsAUcAkKamSCC9AiBxA0UOHyhTqtlvMksO7AE
	  hNC2-YzPyx1FkhMoS4LLe6E_pFsMlmjA6P1NSge9C5G5tETYXGAn6b1xZbHtmwrPSc
	  ro9LWhVmAaA7_bxYObnFUxgWtK4vzzQBjZJ36UTk4OTB-JvKWgfVWCFsaw5WCHj6Oo
	  4jpO7d2yN7WMfAj2hTEabz9wumQ0TMhBduZ-QON3pYObSy7TSC1vVme0NJrwF_cJRe
	  hKTFmdlXGVldPxZCplr7ZQqRQhF8JP-l4mEQVnCaWGn9ONHlemczGOS-A-wwtnmwjI
	  B1V_vgJRf4FdpV-4hUk4-QLpu3-1lWFxrtZKcggq3tWTduRo5_QebQbUUT_VSCgsFc
	  OmyWKoj56lbxthN19hq1XGWbLGfrrR6MWh23vk01zn8FVwi7uFwEnRYSafsnWLa1Z5
	  TpBj9GvAdl2H9NHwzpB5NqHpZNkQ3NMDj13Fn8fzO0JB83Etbm_tnFQfcb13X3bJ15
	  Cz-Ww1MGhvIpGGnMBT_ADp9xSIyAM9dQ1yeVXk-AIgWBUlN5uyWSGyCxp0cJwx7HxM
	  38z0UIeBu-MytL-eqndM7LxytsVzCbjOTSVRmhYEMIzUAnS1gs7uMQAGRdgRIElTJE
	  SGMjb_4bZq9s6Ve1LKkSi0_QDsrABaLe55UY0zF4ZSfOV5PMyPtocwV_dcNPlxLgNA
	  D1BFX_Z9kAdMZQW6fAmsfFle0zAoMe4l9pMESH0JB4sJGdCKtQXj1cXNydDYozF7l8
	  H00BV_Er7zd6VtIw0MxwkFCTatsv_R-GsBCH218RgVPsfYhwVuT8R4HarpzsDBufC4
	  r8_c8fc9Z278sQ081jFjOja6L2x0N_ImzFNXU6xwO-Ska-QeuvYZ3X_L31ZOX4Llp-
	  7QSfgDoHnOxFv1Xws-D5mDHD3zxOup2b2TppdKTZb9eW2vxUVviM8OI9atBfPKMGAO
	  v9omA-6vv5IxUH0-lWMiHLQ_g8vnswp-Jav0c4t6URVUzujNOoNd_CBGGVnHiJTCHl
	  88LQxsqLHHIu4Fz-U2SGnlxGTj0-ihit2ELGRv4vO8E1BosTmf0cx3qgG0Pq0eOLBD
	  IHsrdZ_CCAiTc0HVkMbyq1M6qEhM-q5P6y1QCIrwg.
	  0HFmhOzsQ98nNWJjIHkR7A`

	// remove white spaces and line breaks
	r := regexp.MustCompile(`\s`)
	plaintext = r.ReplaceAll(plaintext, []byte(""))
	serializationReference = r.ReplaceAllString(serializationReference, "")

	rcpt := Recipient{
		Algorithm:  PBES2_HS256_A128KW,
		Key:        passphrase,
		PBES2Count: 4096,
		PBES2Salt: []byte{
			217, 96, 147, 112, 150, 117, 70,
			247, 127, 8, 155, 137, 174, 42, 80, 215,
		},
	}

	enc, err := NewEncrypter(A128CBC_HS256, rcpt, nil)
	if err != nil {
		t.Fatal("error on NewEncrypter:", err)
	}

	obj, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatal("error on new Encrypt:", err)
	}

	serialized, err := obj.CompactSerialize()
	if err != nil {
		t.Fatal("error on CompactSerialize")
	}

	jwe1, err := ParseEncrypted(serialized)
	if err != nil {
		t.Fatal("error in ParseEncrypted")
	}

	jwe2, err := ParseEncrypted(serializationReference)
	if err != nil {
		t.Fatal("error in ParseEncrypted")
	}

	original1, err := jwe1.Decrypt(passphrase)
	if err != nil {
		t.Fatal("error in Decrypt:", err)
	}

	original2, err := jwe2.Decrypt(passphrase)
	if err != nil {
		t.Fatal("error in Decrypt reference:", err)
	}

	if bytes.Compare(original1, original2) != 0 {
		t.Error("decryption does not match reference decryption")
	}

	if bytes.Compare(plaintext, original1) != 0 {
		t.Error("decryption does not match plaintext")
	}

	if bytes.Compare(plaintext, original2) != 0 {
		t.Error("reference decryption does not match plaintext")
	}
}

func TestEncrypterWithPBES2(t *testing.T) {
	expected := []byte("Lorem ipsum dolor sit amet")
	algs := []KeyAlgorithm{
		PBES2_HS256_A128KW, PBES2_HS384_A192KW, PBES2_HS512_A256KW,
	}

	// Check with both strings and []byte
	recipientKeys := []interface{}{"password", []byte("password")}
	for _, key := range recipientKeys {
		for _, alg := range algs {
			enc, err := NewEncrypter(A128GCM, Recipient{Algorithm: alg, Key: &JSONWebKey{
				KeyID: "test-id",
				Key:   key,
			}}, nil)
			if err != nil {
				t.Error(err)
			}

			ciphertext, _ := enc.Encrypt(expected)

			serialized1, _ := ciphertext.CompactSerialize()
			serialized2 := ciphertext.FullSerialize()

			parsed1, _ := ParseEncrypted(serialized1)
			parsed2, _ := ParseEncrypted(serialized2)

			actual1, err := parsed1.Decrypt("password")
			if err != nil {
				t.Fatal("error on Decrypt:", err)
			}

			actual2, err := parsed2.Decrypt([]byte("password"))
			if err != nil {
				t.Fatal("error on Decrypt:", err)
			}

			if bytes.Compare(actual1, expected) != 0 {
				t.Errorf("error comparing decrypted message (%s) and expected (%s)", actual1, expected)
			}

			if bytes.Compare(actual2, expected) != 0 {
				t.Errorf("error comparing decrypted message (%s) and expected (%s)", actual2, expected)
			}
		}
	}
}

type testKey struct {
	enc, dec interface{}
}

func symmetricTestKey(size int) []testKey {
	key, _, _ := randomKeyGenerator{size: size}.genKey()

	return []testKey{
		{
			enc: key,
			dec: key,
		},
		{
			enc: &JSONWebKey{KeyID: "test", Key: key},
			dec: &JSONWebKey{KeyID: "test", Key: key},
		},
	}
}

func TestDirectEncryptionKeySizeCheck(t *testing.T) {
	// 16-byte key
	key16 := []byte("0123456789ABCDEF")

	// 32-byte key
	key32 := []byte("0123456789ABCDEF0123456789ABCDEF")

	// AES-128 with 32-byte key should reject
	_, err := NewEncrypter(A128GCM, Recipient{Algorithm: DIRECT, Key: key32}, nil)
	if err != ErrInvalidKeySize {
		t.Error("Should reject AES-128 with 32-byte key")
	}

	// AES-256 with 16-byte key should reject
	_, err = NewEncrypter(A256GCM, Recipient{Algorithm: DIRECT, Key: key16}, nil)
	if err != ErrInvalidKeySize {
		t.Error("Should reject AES-256 with 16-byte key")
	}
}

func generateTestKeys(keyAlg KeyAlgorithm, encAlg ContentEncryption) []testKey {
	switch keyAlg {
	case DIRECT:
		return symmetricTestKey(getContentCipher(encAlg).keySize())
	case ECDH_ES, ECDH_ES_A128KW, ECDH_ES_A192KW, ECDH_ES_A256KW:
		return []testKey{
			{
				dec: ecTestKey256,
				enc: &ecTestKey256.PublicKey,
			},
			{
				dec: ecTestKey384,
				enc: &ecTestKey384.PublicKey,
			},
			{
				dec: ecTestKey521,
				enc: &ecTestKey521.PublicKey,
			},
			{
				dec: &JSONWebKey{KeyID: "test", Key: ecTestKey256},
				enc: &JSONWebKey{KeyID: "test", Key: &ecTestKey256.PublicKey},
			},
		}
	case A128GCMKW, A128KW:
		return symmetricTestKey(16)
	case A192GCMKW, A192KW:
		return symmetricTestKey(24)
	case A256GCMKW, A256KW:
		return symmetricTestKey(32)
	case RSA1_5, RSA_OAEP, RSA_OAEP_256:
		return []testKey{{
			dec: rsaTestKey,
			enc: &rsaTestKey.PublicKey,
		}}
	case PBES2_HS256_A128KW, PBES2_HS384_A192KW, PBES2_HS512_A256KW:
		// size does not matter, use random integer
		i, err := rand.Int(rand.Reader, big.NewInt(64))
		if err != nil {
			panic(err)
		}
		return symmetricTestKey(int(i.Int64()))
	}

	panic("Must update test case")
}

func RunRoundtripsJWE(b *testing.B, alg KeyAlgorithm, enc ContentEncryption, zip CompressionAlgorithm, priv, pub interface{}) {
	serializer := func(obj *JSONWebEncryption) (string, error) {
		return obj.CompactSerialize()
	}

	corrupter := func(obj *JSONWebEncryption) bool { return false }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := RoundtripJWE(alg, enc, zip, serializer, corrupter, nil, pub, priv)
		if err != nil {
			b.Error(err)
		}
	}
}

var (
	chunks = map[string][]byte{
		"1B":   make([]byte, 1),
		"64B":  make([]byte, 64),
		"1KB":  make([]byte, 1024),
		"64KB": make([]byte, 65536),
		"1MB":  make([]byte, 1048576),
		"64MB": make([]byte, 67108864),
	}

	symKey16, _, _ = randomKeyGenerator{size: 16}.genKey()
	symKey32, _, _ = randomKeyGenerator{size: 32}.genKey()
	symKey64, _, _ = randomKeyGenerator{size: 64}.genKey()

	encrypters = map[string]Encrypter{
		"OAEPAndGCM":          mustEncrypter(RSA_OAEP, A128GCM, &rsaTestKey.PublicKey),
		"PKCSAndGCM":          mustEncrypter(RSA1_5, A128GCM, &rsaTestKey.PublicKey),
		"OAEPAndCBC":          mustEncrypter(RSA_OAEP, A128CBC_HS256, &rsaTestKey.PublicKey),
		"PKCSAndCBC":          mustEncrypter(RSA1_5, A128CBC_HS256, &rsaTestKey.PublicKey),
		"DirectGCM128":        mustEncrypter(DIRECT, A128GCM, symKey16),
		"DirectCBC128":        mustEncrypter(DIRECT, A128CBC_HS256, symKey32),
		"DirectGCM256":        mustEncrypter(DIRECT, A256GCM, symKey32),
		"DirectCBC256":        mustEncrypter(DIRECT, A256CBC_HS512, symKey64),
		"AESKWAndGCM128":      mustEncrypter(A128KW, A128GCM, symKey16),
		"AESKWAndCBC256":      mustEncrypter(A256KW, A256GCM, symKey32),
		"ECDHOnP256AndGCM128": mustEncrypter(ECDH_ES, A128GCM, &ecTestKey256.PublicKey),
		"ECDHOnP384AndGCM128": mustEncrypter(ECDH_ES, A128GCM, &ecTestKey384.PublicKey),
		"ECDHOnP521AndGCM128": mustEncrypter(ECDH_ES, A128GCM, &ecTestKey521.PublicKey),
	}
)

func BenchmarkEncrypt1BWithOAEPAndGCM(b *testing.B)   { benchEncrypt("1B", "OAEPAndGCM", b) }
func BenchmarkEncrypt64BWithOAEPAndGCM(b *testing.B)  { benchEncrypt("64B", "OAEPAndGCM", b) }
func BenchmarkEncrypt1KBWithOAEPAndGCM(b *testing.B)  { benchEncrypt("1KB", "OAEPAndGCM", b) }
func BenchmarkEncrypt64KBWithOAEPAndGCM(b *testing.B) { benchEncrypt("64KB", "OAEPAndGCM", b) }
func BenchmarkEncrypt1MBWithOAEPAndGCM(b *testing.B)  { benchEncrypt("1MB", "OAEPAndGCM", b) }
func BenchmarkEncrypt64MBWithOAEPAndGCM(b *testing.B) { benchEncrypt("64MB", "OAEPAndGCM", b) }

func BenchmarkEncrypt1BWithPKCSAndGCM(b *testing.B)   { benchEncrypt("1B", "PKCSAndGCM", b) }
func BenchmarkEncrypt64BWithPKCSAndGCM(b *testing.B)  { benchEncrypt("64B", "PKCSAndGCM", b) }
func BenchmarkEncrypt1KBWithPKCSAndGCM(b *testing.B)  { benchEncrypt("1KB", "PKCSAndGCM", b) }
func BenchmarkEncrypt64KBWithPKCSAndGCM(b *testing.B) { benchEncrypt("64KB", "PKCSAndGCM", b) }
func BenchmarkEncrypt1MBWithPKCSAndGCM(b *testing.B)  { benchEncrypt("1MB", "PKCSAndGCM", b) }
func BenchmarkEncrypt64MBWithPKCSAndGCM(b *testing.B) { benchEncrypt("64MB", "PKCSAndGCM", b) }

func BenchmarkEncrypt1BWithOAEPAndCBC(b *testing.B)   { benchEncrypt("1B", "OAEPAndCBC", b) }
func BenchmarkEncrypt64BWithOAEPAndCBC(b *testing.B)  { benchEncrypt("64B", "OAEPAndCBC", b) }
func BenchmarkEncrypt1KBWithOAEPAndCBC(b *testing.B)  { benchEncrypt("1KB", "OAEPAndCBC", b) }
func BenchmarkEncrypt64KBWithOAEPAndCBC(b *testing.B) { benchEncrypt("64KB", "OAEPAndCBC", b) }
func BenchmarkEncrypt1MBWithOAEPAndCBC(b *testing.B)  { benchEncrypt("1MB", "OAEPAndCBC", b) }
func BenchmarkEncrypt64MBWithOAEPAndCBC(b *testing.B) { benchEncrypt("64MB", "OAEPAndCBC", b) }

func BenchmarkEncrypt1BWithPKCSAndCBC(b *testing.B)   { benchEncrypt("1B", "PKCSAndCBC", b) }
func BenchmarkEncrypt64BWithPKCSAndCBC(b *testing.B)  { benchEncrypt("64B", "PKCSAndCBC", b) }
func BenchmarkEncrypt1KBWithPKCSAndCBC(b *testing.B)  { benchEncrypt("1KB", "PKCSAndCBC", b) }
func BenchmarkEncrypt64KBWithPKCSAndCBC(b *testing.B) { benchEncrypt("64KB", "PKCSAndCBC", b) }
func BenchmarkEncrypt1MBWithPKCSAndCBC(b *testing.B)  { benchEncrypt("1MB", "PKCSAndCBC", b) }
func BenchmarkEncrypt64MBWithPKCSAndCBC(b *testing.B) { benchEncrypt("64MB", "PKCSAndCBC", b) }

func BenchmarkEncrypt1BWithDirectGCM128(b *testing.B)   { benchEncrypt("1B", "DirectGCM128", b) }
func BenchmarkEncrypt64BWithDirectGCM128(b *testing.B)  { benchEncrypt("64B", "DirectGCM128", b) }
func BenchmarkEncrypt1KBWithDirectGCM128(b *testing.B)  { benchEncrypt("1KB", "DirectGCM128", b) }
func BenchmarkEncrypt64KBWithDirectGCM128(b *testing.B) { benchEncrypt("64KB", "DirectGCM128", b) }
func BenchmarkEncrypt1MBWithDirectGCM128(b *testing.B)  { benchEncrypt("1MB", "DirectGCM128", b) }
func BenchmarkEncrypt64MBWithDirectGCM128(b *testing.B) { benchEncrypt("64MB", "DirectGCM128", b) }

func BenchmarkEncrypt1BWithDirectCBC128(b *testing.B)   { benchEncrypt("1B", "DirectCBC128", b) }
func BenchmarkEncrypt64BWithDirectCBC128(b *testing.B)  { benchEncrypt("64B", "DirectCBC128", b) }
func BenchmarkEncrypt1KBWithDirectCBC128(b *testing.B)  { benchEncrypt("1KB", "DirectCBC128", b) }
func BenchmarkEncrypt64KBWithDirectCBC128(b *testing.B) { benchEncrypt("64KB", "DirectCBC128", b) }
func BenchmarkEncrypt1MBWithDirectCBC128(b *testing.B)  { benchEncrypt("1MB", "DirectCBC128", b) }
func BenchmarkEncrypt64MBWithDirectCBC128(b *testing.B) { benchEncrypt("64MB", "DirectCBC128", b) }

func BenchmarkEncrypt1BWithDirectGCM256(b *testing.B)   { benchEncrypt("1B", "DirectGCM256", b) }
func BenchmarkEncrypt64BWithDirectGCM256(b *testing.B)  { benchEncrypt("64B", "DirectGCM256", b) }
func BenchmarkEncrypt1KBWithDirectGCM256(b *testing.B)  { benchEncrypt("1KB", "DirectGCM256", b) }
func BenchmarkEncrypt64KBWithDirectGCM256(b *testing.B) { benchEncrypt("64KB", "DirectGCM256", b) }
func BenchmarkEncrypt1MBWithDirectGCM256(b *testing.B)  { benchEncrypt("1MB", "DirectGCM256", b) }
func BenchmarkEncrypt64MBWithDirectGCM256(b *testing.B) { benchEncrypt("64MB", "DirectGCM256", b) }

func BenchmarkEncrypt1BWithDirectCBC256(b *testing.B)   { benchEncrypt("1B", "DirectCBC256", b) }
func BenchmarkEncrypt64BWithDirectCBC256(b *testing.B)  { benchEncrypt("64B", "DirectCBC256", b) }
func BenchmarkEncrypt1KBWithDirectCBC256(b *testing.B)  { benchEncrypt("1KB", "DirectCBC256", b) }
func BenchmarkEncrypt64KBWithDirectCBC256(b *testing.B) { benchEncrypt("64KB", "DirectCBC256", b) }
func BenchmarkEncrypt1MBWithDirectCBC256(b *testing.B)  { benchEncrypt("1MB", "DirectCBC256", b) }
func BenchmarkEncrypt64MBWithDirectCBC256(b *testing.B) { benchEncrypt("64MB", "DirectCBC256", b) }

func BenchmarkEncrypt1BWithAESKWAndGCM128(b *testing.B)   { benchEncrypt("1B", "AESKWAndGCM128", b) }
func BenchmarkEncrypt64BWithAESKWAndGCM128(b *testing.B)  { benchEncrypt("64B", "AESKWAndGCM128", b) }
func BenchmarkEncrypt1KBWithAESKWAndGCM128(b *testing.B)  { benchEncrypt("1KB", "AESKWAndGCM128", b) }
func BenchmarkEncrypt64KBWithAESKWAndGCM128(b *testing.B) { benchEncrypt("64KB", "AESKWAndGCM128", b) }
func BenchmarkEncrypt1MBWithAESKWAndGCM128(b *testing.B)  { benchEncrypt("1MB", "AESKWAndGCM128", b) }
func BenchmarkEncrypt64MBWithAESKWAndGCM128(b *testing.B) { benchEncrypt("64MB", "AESKWAndGCM128", b) }

func BenchmarkEncrypt1BWithAESKWAndCBC256(b *testing.B)   { benchEncrypt("1B", "AESKWAndCBC256", b) }
func BenchmarkEncrypt64BWithAESKWAndCBC256(b *testing.B)  { benchEncrypt("64B", "AESKWAndCBC256", b) }
func BenchmarkEncrypt1KBWithAESKWAndCBC256(b *testing.B)  { benchEncrypt("1KB", "AESKWAndCBC256", b) }
func BenchmarkEncrypt64KBWithAESKWAndCBC256(b *testing.B) { benchEncrypt("64KB", "AESKWAndCBC256", b) }
func BenchmarkEncrypt1MBWithAESKWAndCBC256(b *testing.B)  { benchEncrypt("1MB", "AESKWAndCBC256", b) }
func BenchmarkEncrypt64MBWithAESKWAndCBC256(b *testing.B) { benchEncrypt("64MB", "AESKWAndCBC256", b) }

func BenchmarkEncrypt1BWithECDHOnP256AndGCM128(b *testing.B) {
	benchEncrypt("1B", "ECDHOnP256AndGCM128", b)
}
func BenchmarkEncrypt64BWithECDHOnP256AndGCM128(b *testing.B) {
	benchEncrypt("64B", "ECDHOnP256AndGCM128", b)
}
func BenchmarkEncrypt1KBWithECDHOnP256AndGCM128(b *testing.B) {
	benchEncrypt("1KB", "ECDHOnP256AndGCM128", b)
}
func BenchmarkEncrypt64KBWithECDHOnP256AndGCM128(b *testing.B) {
	benchEncrypt("64KB", "ECDHOnP256AndGCM128", b)
}
func BenchmarkEncrypt1MBWithECDHOnP256AndGCM128(b *testing.B) {
	benchEncrypt("1MB", "ECDHOnP256AndGCM128", b)
}
func BenchmarkEncrypt64MBWithECDHOnP256AndGCM128(b *testing.B) {
	benchEncrypt("64MB", "ECDHOnP256AndGCM128", b)
}

func BenchmarkEncrypt1BWithECDHOnP384AndGCM128(b *testing.B) {
	benchEncrypt("1B", "ECDHOnP384AndGCM128", b)
}
func BenchmarkEncrypt64BWithECDHOnP384AndGCM128(b *testing.B) {
	benchEncrypt("64B", "ECDHOnP384AndGCM128", b)
}
func BenchmarkEncrypt1KBWithECDHOnP384AndGCM128(b *testing.B) {
	benchEncrypt("1KB", "ECDHOnP384AndGCM128", b)
}
func BenchmarkEncrypt64KBWithECDHOnP384AndGCM128(b *testing.B) {
	benchEncrypt("64KB", "ECDHOnP384AndGCM128", b)
}
func BenchmarkEncrypt1MBWithECDHOnP384AndGCM128(b *testing.B) {
	benchEncrypt("1MB", "ECDHOnP384AndGCM128", b)
}
func BenchmarkEncrypt64MBWithECDHOnP384AndGCM128(b *testing.B) {
	benchEncrypt("64MB", "ECDHOnP384AndGCM128", b)
}

func BenchmarkEncrypt1BWithECDHOnP521AndGCM128(b *testing.B) {
	benchEncrypt("1B", "ECDHOnP521AndGCM128", b)
}
func BenchmarkEncrypt64BWithECDHOnP521AndGCM128(b *testing.B) {
	benchEncrypt("64B", "ECDHOnP521AndGCM128", b)
}
func BenchmarkEncrypt1KBWithECDHOnP521AndGCM128(b *testing.B) {
	benchEncrypt("1KB", "ECDHOnP521AndGCM128", b)
}
func BenchmarkEncrypt64KBWithECDHOnP521AndGCM128(b *testing.B) {
	benchEncrypt("64KB", "ECDHOnP521AndGCM128", b)
}
func BenchmarkEncrypt1MBWithECDHOnP521AndGCM128(b *testing.B) {
	benchEncrypt("1MB", "ECDHOnP521AndGCM128", b)
}
func BenchmarkEncrypt64MBWithECDHOnP521AndGCM128(b *testing.B) {
	benchEncrypt("64MB", "ECDHOnP521AndGCM128", b)
}

func benchEncrypt(chunkKey, primKey string, b *testing.B) {
	data, ok := chunks[chunkKey]
	if !ok {
		b.Fatalf("unknown chunk size %s", chunkKey)
	}

	enc, ok := encrypters[primKey]
	if !ok {
		b.Fatalf("unknown encrypter %s", primKey)
	}

	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		enc.Encrypt(data)
	}
}

var (
	decryptionKeys = map[string]interface{}{
		"OAEPAndGCM": rsaTestKey,
		"PKCSAndGCM": rsaTestKey,
		"OAEPAndCBC": rsaTestKey,
		"PKCSAndCBC": rsaTestKey,

		"DirectGCM128": symKey16,
		"DirectCBC128": symKey32,
		"DirectGCM256": symKey32,
		"DirectCBC256": symKey64,

		"AESKWAndGCM128": symKey16,
		"AESKWAndCBC256": symKey32,

		"ECDHOnP256AndGCM128": ecTestKey256,
		"ECDHOnP384AndGCM128": ecTestKey384,
		"ECDHOnP521AndGCM128": ecTestKey521,
	}
)

func BenchmarkDecrypt1BWithOAEPAndGCM(b *testing.B)   { benchDecrypt("1B", "OAEPAndGCM", b) }
func BenchmarkDecrypt64BWithOAEPAndGCM(b *testing.B)  { benchDecrypt("64B", "OAEPAndGCM", b) }
func BenchmarkDecrypt1KBWithOAEPAndGCM(b *testing.B)  { benchDecrypt("1KB", "OAEPAndGCM", b) }
func BenchmarkDecrypt64KBWithOAEPAndGCM(b *testing.B) { benchDecrypt("64KB", "OAEPAndGCM", b) }
func BenchmarkDecrypt1MBWithOAEPAndGCM(b *testing.B)  { benchDecrypt("1MB", "OAEPAndGCM", b) }
func BenchmarkDecrypt64MBWithOAEPAndGCM(b *testing.B) { benchDecrypt("64MB", "OAEPAndGCM", b) }

func BenchmarkDecrypt1BWithPKCSAndGCM(b *testing.B)   { benchDecrypt("1B", "PKCSAndGCM", b) }
func BenchmarkDecrypt64BWithPKCSAndGCM(b *testing.B)  { benchDecrypt("64B", "PKCSAndGCM", b) }
func BenchmarkDecrypt1KBWithPKCSAndGCM(b *testing.B)  { benchDecrypt("1KB", "PKCSAndGCM", b) }
func BenchmarkDecrypt64KBWithPKCSAndGCM(b *testing.B) { benchDecrypt("64KB", "PKCSAndGCM", b) }
func BenchmarkDecrypt1MBWithPKCSAndGCM(b *testing.B)  { benchDecrypt("1MB", "PKCSAndGCM", b) }
func BenchmarkDecrypt64MBWithPKCSAndGCM(b *testing.B) { benchDecrypt("64MB", "PKCSAndGCM", b) }

func BenchmarkDecrypt1BWithOAEPAndCBC(b *testing.B)   { benchDecrypt("1B", "OAEPAndCBC", b) }
func BenchmarkDecrypt64BWithOAEPAndCBC(b *testing.B)  { benchDecrypt("64B", "OAEPAndCBC", b) }
func BenchmarkDecrypt1KBWithOAEPAndCBC(b *testing.B)  { benchDecrypt("1KB", "OAEPAndCBC", b) }
func BenchmarkDecrypt64KBWithOAEPAndCBC(b *testing.B) { benchDecrypt("64KB", "OAEPAndCBC", b) }
func BenchmarkDecrypt1MBWithOAEPAndCBC(b *testing.B)  { benchDecrypt("1MB", "OAEPAndCBC", b) }
func BenchmarkDecrypt64MBWithOAEPAndCBC(b *testing.B) { benchDecrypt("64MB", "OAEPAndCBC", b) }

func BenchmarkDecrypt1BWithPKCSAndCBC(b *testing.B)   { benchDecrypt("1B", "PKCSAndCBC", b) }
func BenchmarkDecrypt64BWithPKCSAndCBC(b *testing.B)  { benchDecrypt("64B", "PKCSAndCBC", b) }
func BenchmarkDecrypt1KBWithPKCSAndCBC(b *testing.B)  { benchDecrypt("1KB", "PKCSAndCBC", b) }
func BenchmarkDecrypt64KBWithPKCSAndCBC(b *testing.B) { benchDecrypt("64KB", "PKCSAndCBC", b) }
func BenchmarkDecrypt1MBWithPKCSAndCBC(b *testing.B)  { benchDecrypt("1MB", "PKCSAndCBC", b) }
func BenchmarkDecrypt64MBWithPKCSAndCBC(b *testing.B) { benchDecrypt("64MB", "PKCSAndCBC", b) }

func BenchmarkDecrypt1BWithDirectGCM128(b *testing.B)   { benchDecrypt("1B", "DirectGCM128", b) }
func BenchmarkDecrypt64BWithDirectGCM128(b *testing.B)  { benchDecrypt("64B", "DirectGCM128", b) }
func BenchmarkDecrypt1KBWithDirectGCM128(b *testing.B)  { benchDecrypt("1KB", "DirectGCM128", b) }
func BenchmarkDecrypt64KBWithDirectGCM128(b *testing.B) { benchDecrypt("64KB", "DirectGCM128", b) }
func BenchmarkDecrypt1MBWithDirectGCM128(b *testing.B)  { benchDecrypt("1MB", "DirectGCM128", b) }
func BenchmarkDecrypt64MBWithDirectGCM128(b *testing.B) { benchDecrypt("64MB", "DirectGCM128", b) }

func BenchmarkDecrypt1BWithDirectCBC128(b *testing.B)   { benchDecrypt("1B", "DirectCBC128", b) }
func BenchmarkDecrypt64BWithDirectCBC128(b *testing.B)  { benchDecrypt("64B", "DirectCBC128", b) }
func BenchmarkDecrypt1KBWithDirectCBC128(b *testing.B)  { benchDecrypt("1KB", "DirectCBC128", b) }
func BenchmarkDecrypt64KBWithDirectCBC128(b *testing.B) { benchDecrypt("64KB", "DirectCBC128", b) }
func BenchmarkDecrypt1MBWithDirectCBC128(b *testing.B)  { benchDecrypt("1MB", "DirectCBC128", b) }
func BenchmarkDecrypt64MBWithDirectCBC128(b *testing.B) { benchDecrypt("64MB", "DirectCBC128", b) }

func BenchmarkDecrypt1BWithDirectGCM256(b *testing.B)   { benchDecrypt("1B", "DirectGCM256", b) }
func BenchmarkDecrypt64BWithDirectGCM256(b *testing.B)  { benchDecrypt("64B", "DirectGCM256", b) }
func BenchmarkDecrypt1KBWithDirectGCM256(b *testing.B)  { benchDecrypt("1KB", "DirectGCM256", b) }
func BenchmarkDecrypt64KBWithDirectGCM256(b *testing.B) { benchDecrypt("64KB", "DirectGCM256", b) }
func BenchmarkDecrypt1MBWithDirectGCM256(b *testing.B)  { benchDecrypt("1MB", "DirectGCM256", b) }
func BenchmarkDecrypt64MBWithDirectGCM256(b *testing.B) { benchDecrypt("64MB", "DirectGCM256", b) }

func BenchmarkDecrypt1BWithDirectCBC256(b *testing.B)   { benchDecrypt("1B", "DirectCBC256", b) }
func BenchmarkDecrypt64BWithDirectCBC256(b *testing.B)  { benchDecrypt("64B", "DirectCBC256", b) }
func BenchmarkDecrypt1KBWithDirectCBC256(b *testing.B)  { benchDecrypt("1KB", "DirectCBC256", b) }
func BenchmarkDecrypt64KBWithDirectCBC256(b *testing.B) { benchDecrypt("64KB", "DirectCBC256", b) }
func BenchmarkDecrypt1MBWithDirectCBC256(b *testing.B)  { benchDecrypt("1MB", "DirectCBC256", b) }
func BenchmarkDecrypt64MBWithDirectCBC256(b *testing.B) { benchDecrypt("64MB", "DirectCBC256", b) }

func BenchmarkDecrypt1BWithAESKWAndGCM128(b *testing.B)   { benchDecrypt("1B", "AESKWAndGCM128", b) }
func BenchmarkDecrypt64BWithAESKWAndGCM128(b *testing.B)  { benchDecrypt("64B", "AESKWAndGCM128", b) }
func BenchmarkDecrypt1KBWithAESKWAndGCM128(b *testing.B)  { benchDecrypt("1KB", "AESKWAndGCM128", b) }
func BenchmarkDecrypt64KBWithAESKWAndGCM128(b *testing.B) { benchDecrypt("64KB", "AESKWAndGCM128", b) }
func BenchmarkDecrypt1MBWithAESKWAndGCM128(b *testing.B)  { benchDecrypt("1MB", "AESKWAndGCM128", b) }
func BenchmarkDecrypt64MBWithAESKWAndGCM128(b *testing.B) { benchDecrypt("64MB", "AESKWAndGCM128", b) }

func BenchmarkDecrypt1BWithAESKWAndCBC256(b *testing.B)   { benchDecrypt("1B", "AESKWAndCBC256", b) }
func BenchmarkDecrypt64BWithAESKWAndCBC256(b *testing.B)  { benchDecrypt("64B", "AESKWAndCBC256", b) }
func BenchmarkDecrypt1KBWithAESKWAndCBC256(b *testing.B)  { benchDecrypt("1KB", "AESKWAndCBC256", b) }
func BenchmarkDecrypt64KBWithAESKWAndCBC256(b *testing.B) { benchDecrypt("64KB", "AESKWAndCBC256", b) }
func BenchmarkDecrypt1MBWithAESKWAndCBC256(b *testing.B)  { benchDecrypt("1MB", "AESKWAndCBC256", b) }
func BenchmarkDecrypt64MBWithAESKWAndCBC256(b *testing.B) { benchDecrypt("64MB", "AESKWAndCBC256", b) }

func BenchmarkDecrypt1BWithECDHOnP256AndGCM128(b *testing.B) {
	benchDecrypt("1B", "ECDHOnP256AndGCM128", b)
}
func BenchmarkDecrypt64BWithECDHOnP256AndGCM128(b *testing.B) {
	benchDecrypt("64B", "ECDHOnP256AndGCM128", b)
}
func BenchmarkDecrypt1KBWithECDHOnP256AndGCM128(b *testing.B) {
	benchDecrypt("1KB", "ECDHOnP256AndGCM128", b)
}
func BenchmarkDecrypt64KBWithECDHOnP256AndGCM128(b *testing.B) {
	benchDecrypt("64KB", "ECDHOnP256AndGCM128", b)
}
func BenchmarkDecrypt1MBWithECDHOnP256AndGCM128(b *testing.B) {
	benchDecrypt("1MB", "ECDHOnP256AndGCM128", b)
}
func BenchmarkDecrypt64MBWithECDHOnP256AndGCM128(b *testing.B) {
	benchDecrypt("64MB", "ECDHOnP256AndGCM128", b)
}

func BenchmarkDecrypt1BWithECDHOnP384AndGCM128(b *testing.B) {
	benchDecrypt("1B", "ECDHOnP384AndGCM128", b)
}
func BenchmarkDecrypt64BWithECDHOnP384AndGCM128(b *testing.B) {
	benchDecrypt("64B", "ECDHOnP384AndGCM128", b)
}
func BenchmarkDecrypt1KBWithECDHOnP384AndGCM128(b *testing.B) {
	benchDecrypt("1KB", "ECDHOnP384AndGCM128", b)
}
func BenchmarkDecrypt64KBWithECDHOnP384AndGCM128(b *testing.B) {
	benchDecrypt("64KB", "ECDHOnP384AndGCM128", b)
}
func BenchmarkDecrypt1MBWithECDHOnP384AndGCM128(b *testing.B) {
	benchDecrypt("1MB", "ECDHOnP384AndGCM128", b)
}
func BenchmarkDecrypt64MBWithECDHOnP384AndGCM128(b *testing.B) {
	benchDecrypt("64MB", "ECDHOnP384AndGCM128", b)
}

func BenchmarkDecrypt1BWithECDHOnP521AndGCM128(b *testing.B) {
	benchDecrypt("1B", "ECDHOnP521AndGCM128", b)
}
func BenchmarkDecrypt64BWithECDHOnP521AndGCM128(b *testing.B) {
	benchDecrypt("64B", "ECDHOnP521AndGCM128", b)
}
func BenchmarkDecrypt1KBWithECDHOnP521AndGCM128(b *testing.B) {
	benchDecrypt("1KB", "ECDHOnP521AndGCM128", b)
}
func BenchmarkDecrypt64KBWithECDHOnP521AndGCM128(b *testing.B) {
	benchDecrypt("64KB", "ECDHOnP521AndGCM128", b)
}
func BenchmarkDecrypt1MBWithECDHOnP521AndGCM128(b *testing.B) {
	benchDecrypt("1MB", "ECDHOnP521AndGCM128", b)
}
func BenchmarkDecrypt64MBWithECDHOnP521AndGCM128(b *testing.B) {
	benchDecrypt("64MB", "ECDHOnP521AndGCM128", b)
}

func benchDecrypt(chunkKey, primKey string, b *testing.B) {
	chunk, ok := chunks[chunkKey]
	if !ok {
		b.Fatalf("unknown chunk size %s", chunkKey)
	}

	enc, ok := encrypters[primKey]
	if !ok {
		b.Fatalf("unknown encrypter %s", primKey)
	}

	dec, ok := decryptionKeys[primKey]
	if !ok {
		b.Fatalf("unknown decryption key %s", primKey)
	}

	data, err := enc.Encrypt(chunk)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(chunk)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Decrypt(dec)
	}
}

func mustEncrypter(keyAlg KeyAlgorithm, encAlg ContentEncryption, encryptionKey interface{}) Encrypter {
	enc, err := NewEncrypter(encAlg, Recipient{Algorithm: keyAlg, Key: encryptionKey}, nil)
	if err != nil {
		panic(err)
	}
	return enc
}
