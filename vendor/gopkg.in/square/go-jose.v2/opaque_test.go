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

package jose

import (
	"fmt"
	"testing"
)

type signWrapper struct {
	pk      *JSONWebKey
	wrapped payloadSigner
	algs    []SignatureAlgorithm
}

var _ = OpaqueSigner(&signWrapper{})

func (sw *signWrapper) Algs() []SignatureAlgorithm {
	return sw.algs
}

func (sw *signWrapper) Public() *JSONWebKey {
	return sw.pk
}

func (sw *signWrapper) SignPayload(payload []byte, alg SignatureAlgorithm) ([]byte, error) {
	sig, err := sw.wrapped.signPayload(payload, alg)
	if err != nil {
		return nil, err
	}
	return sig.Signature, nil
}

type verifyWrapper struct {
	wrapped []payloadVerifier
}

var _ = OpaqueVerifier(&verifyWrapper{})

func (vw *verifyWrapper) VerifyPayload(payload []byte, signature []byte, alg SignatureAlgorithm) error {
	if len(vw.wrapped) == 0 {
		return fmt.Errorf("error: verifier had no keys")
	}
	var err error
	for _, v := range vw.wrapped {
		err = v.verifyPayload(payload, signature, alg)
		if err == nil {
			return nil
		}
	}
	return err
}

func TestRoundtripsJWSOpaque(t *testing.T) {
	sigAlgs := []SignatureAlgorithm{RS256, RS384, RS512, PS256, PS384, PS512, ES256, ES384, ES512, EdDSA}

	serializers := []func(*JSONWebSignature) (string, error){
		func(obj *JSONWebSignature) (string, error) { return obj.CompactSerialize() },
		func(obj *JSONWebSignature) (string, error) { return obj.FullSerialize(), nil },
	}

	corrupter := func(obj *JSONWebSignature) {}

	for _, alg := range sigAlgs {
		signingKey, verificationKey := GenerateSigningTestKey(alg)

		for i, serializer := range serializers {
			sw := makeOpaqueSigner(t, signingKey, alg)
			vw := makeOpaqueVerifier(t, []interface{}{verificationKey}, alg)

			err := RoundtripJWS(alg, serializer, corrupter, sw, verificationKey, "test_nonce")
			if err != nil {
				t.Error(err, alg, i)
			}

			err = RoundtripJWS(alg, serializer, corrupter, signingKey, vw, "test_nonce")
			if err != nil {
				t.Error(err, alg, i)
			}

			err = RoundtripJWS(alg, serializer, corrupter, sw, vw, "test_nonce")
			if err != nil {
				t.Error(err, alg, i)
			}
		}
	}
}

func makeOpaqueSigner(t *testing.T, signingKey interface{}, alg SignatureAlgorithm) *signWrapper {
	ri, err := makeJWSRecipient(alg, signingKey)
	if err != nil {
		t.Fatal(err)
	}
	return &signWrapper{
		wrapped: ri.signer,
		algs:    []SignatureAlgorithm{alg},
		pk:      &JSONWebKey{Key: ri.publicKey()},
	}
}

func makeOpaqueVerifier(t *testing.T, verificationKey []interface{}, alg SignatureAlgorithm) *verifyWrapper {
	var verifiers []payloadVerifier
	for _, vk := range verificationKey {
		verifier, err := newVerifier(vk)
		if err != nil {
			t.Fatal(err)
		}
		verifiers = append(verifiers, verifier)
	}
	return &verifyWrapper{wrapped: verifiers}
}

func TestOpaqueSignerKeyRotation(t *testing.T) {

	sigAlgs := []SignatureAlgorithm{RS256, RS384, RS512, PS256, PS384, PS512, ES256, ES384, ES512, EdDSA}

	serializers := []func(*JSONWebSignature) (string, error){
		func(obj *JSONWebSignature) (string, error) { return obj.CompactSerialize() },
		func(obj *JSONWebSignature) (string, error) { return obj.FullSerialize(), nil },
	}

	for _, alg := range sigAlgs {
		for i, serializer := range serializers {
			sk1, pk1 := GenerateSigningTestKey(alg)
			sk2, pk2 := GenerateSigningTestKey(alg)

			sw := makeOpaqueSigner(t, sk1, alg)
			sw.pk.KeyID = "first"
			vw := makeOpaqueVerifier(t, []interface{}{pk1, pk2}, alg)

			signer, err := NewSigner(
				SigningKey{Algorithm: alg, Key: sw},
				&SignerOptions{NonceSource: staticNonceSource("test_nonce")},
			)
			if err != nil {
				t.Fatal(err, alg, i)
			}

			jws1, err := signer.Sign([]byte("foo bar baz"))
			if err != nil {
				t.Fatal(err, alg, i)
			}
			jws1 = rtSerialize(t, serializer, jws1, vw)
			if kid := jws1.Signatures[0].Protected.KeyID; kid != "first" {
				t.Errorf("expected kid %q but got %q", "first", kid)
			}

			swNext := makeOpaqueSigner(t, sk2, alg)
			swNext.pk.KeyID = "next"
			sw.wrapped = swNext.wrapped
			sw.pk = swNext.pk

			jws2, err := signer.Sign([]byte("foo bar baz next"))
			if err != nil {
				t.Error(err, alg, i)
			}
			jws2 = rtSerialize(t, serializer, jws2, vw)
			if kid := jws2.Signatures[0].Protected.KeyID; kid != "next" {
				t.Errorf("expected kid %q but got %q", "next", kid)
			}
		}
	}
}

func rtSerialize(t *testing.T, serializer func(*JSONWebSignature) (string, error), sig *JSONWebSignature, vk interface{}) *JSONWebSignature {
	b, err := serializer(sig)
	if err != nil {
		t.Fatal(err)
	}
	sig, err = ParseSigned(b)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sig.Verify(vk); err != nil {
		t.Fatal(err)
	}
	return sig
}
