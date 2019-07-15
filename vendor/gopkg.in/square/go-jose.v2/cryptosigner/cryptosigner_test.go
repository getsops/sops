/*-
 * Copyright 2018 Square Inc.
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

package cryptosigner

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"

	"golang.org/x/crypto/ed25519"
	"gopkg.in/square/go-jose.v2"
)

func TestRoundtripsJWSCryptoSigner(t *testing.T) {
	sigAlgs := []jose.SignatureAlgorithm{jose.RS256, jose.RS384, jose.RS512, jose.PS256, jose.PS384, jose.PS512, jose.ES256, jose.ES384, jose.ES512, jose.EdDSA}

	serializers := []func(*jose.JSONWebSignature) (string, error){
		func(obj *jose.JSONWebSignature) (string, error) { return obj.CompactSerialize() },
		func(obj *jose.JSONWebSignature) (string, error) { return obj.FullSerialize(), nil },
	}

	for _, alg := range sigAlgs {
		signingKey, verificationKey := generateSigningTestKey(alg)

		for i, serializer := range serializers {
			err := roundtripJWS(alg, serializer, Opaque(signingKey.(crypto.Signer)), verificationKey)
			if err != nil {
				t.Error(err, alg, i)
			}
		}
	}
}

type staticNonceSource string

func (sns staticNonceSource) Nonce() (string, error) {
	return string(sns), nil
}

func roundtripJWS(sigAlg jose.SignatureAlgorithm, serializer func(*jose.JSONWebSignature) (string, error), signingKey interface{}, verificationKey interface{}) error {
	nonce := "test_nonce"
	opts := &jose.SignerOptions{
		NonceSource: staticNonceSource(nonce),
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: sigAlg, Key: signingKey}, opts)
	if err != nil {
		return fmt.Errorf("error on new signer: %s", err)
	}

	input := []byte("Lorem ipsum dolor sit amet")
	obj, err := signer.Sign(input)
	if err != nil {
		return fmt.Errorf("error on sign: %s", err)
	}

	msg, err := serializer(obj)
	if err != nil {
		return fmt.Errorf("error on serialize: %s", err)
	}

	obj, err = jose.ParseSigned(msg)
	if err != nil {
		return fmt.Errorf("error on parse: %s", err)
	}

	output, err := obj.Verify(verificationKey)
	if err != nil {
		return fmt.Errorf("error on verify: %s", err)
	}

	// Check that verify works with embedded keys (if present)
	for i, sig := range obj.Signatures {
		if sig.Header.JSONWebKey != nil {
			_, err = obj.Verify(sig.Header.JSONWebKey)
			if err != nil {
				return fmt.Errorf("error on verify with embedded key %d: %s", i, err)
			}
		}

		// Check that the nonce correctly round-tripped (if present)
		if sig.Header.Nonce != nonce {
			return fmt.Errorf("Incorrect nonce returned: [%s]", sig.Header.Nonce)
		}
	}

	if bytes.Compare(output, input) != 0 {
		return fmt.Errorf("input/output do not match, got '%s', expected '%s'", output, input)
	}

	return nil
}

func generateSigningTestKey(sigAlg jose.SignatureAlgorithm) (sig, ver interface{}) {
	switch sigAlg {
	case jose.EdDSA:
		ver, sig, _ = ed25519.GenerateKey(rand.Reader)
	case jose.RS256, jose.RS384, jose.RS512, jose.PS256, jose.PS384, jose.PS512:
		rsaTestKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		sig = rsaTestKey
		ver = &rsaTestKey.PublicKey
	case jose.ES256:
		key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		sig = key
		ver = &key.PublicKey
	case jose.ES384:
		key, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		sig = key
		ver = &key.PublicKey
	case jose.ES512:
		key, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		sig = key
		ver = &key.PublicKey
	default:
		panic("Must update test case")
	}
	return
}
