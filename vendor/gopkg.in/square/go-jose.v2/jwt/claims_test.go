/*-
 * Copyright 2016 Zbigniew Mandziejewicz
 * Copyright 2016 Square, Inc.
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

package jwt

import (
	"testing"
	"time"

	"gopkg.in/square/go-jose.v2/json"

	"github.com/stretchr/testify/assert"
)

func TestEncodeClaims(t *testing.T) {
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

	c := Claims{
		Issuer:    "issuer",
		Subject:   "subject",
		Audience:  Audience{"a1", "a2"},
		NotBefore: NewNumericDate(time.Time{}),
		IssuedAt:  NewNumericDate(now),
		Expiry:    NewNumericDate(now.Add(1 * time.Hour)),
	}

	b, err := json.Marshal(c)
	assert.NoError(t, err)

	expected := `{"iss":"issuer","sub":"subject","aud":["a1","a2"],"exp":1451610000,"iat":1451606400}`
	assert.Equal(t, expected, string(b))
}

func TestDecodeClaims(t *testing.T) {
	s := []byte(`{"iss":"issuer","sub":"subject","aud":["a1","a2"],"exp":1451610000,"iat":1451606400}`)
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

	c := Claims{}
	if err := json.Unmarshal(s, &c); assert.NoError(t, err) {
		assert.Equal(t, "issuer", c.Issuer)
		assert.Equal(t, "subject", c.Subject)
		assert.Equal(t, Audience{"a1", "a2"}, c.Audience)
		assert.True(t, now.Equal(c.IssuedAt.Time()))
		assert.True(t, now.Add(1*time.Hour).Equal(c.Expiry.Time()))
	}

	s2 := []byte(`{"aud": "a1"}`)
	c2 := Claims{}
	if err := json.Unmarshal(s2, &c2); assert.NoError(t, err) {
		assert.Equal(t, Audience{"a1"}, c2.Audience)
	}

	invalid := []struct {
		Raw string
		Err error
	}{
		{`{"aud": 5}`, ErrUnmarshalAudience},
		{`{"aud": ["foo", 5, "bar"]}`, ErrUnmarshalAudience},
		{`{"exp": "invalid"}`, ErrUnmarshalNumericDate},
	}

	for _, v := range invalid {
		c := Claims{}
		assert.Equal(t, v.Err, json.Unmarshal([]byte(v.Raw), &c))
	}
}

func TestNumericDate(t *testing.T) {
	zeroDate := NewNumericDate(time.Time{})
	assert.True(t, time.Time{}.Equal(zeroDate.Time()), "Expected derived time to be time.Time{}")

	zeroDate2 := (*NumericDate)(nil)
	assert.True(t, time.Time{}.Equal(zeroDate2.Time()), "Expected derived time to be time.Time{}")

	nonZeroDate := NewNumericDate(time.Unix(0, 0))
	expected := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Truef(t, expected.Equal(nonZeroDate.Time()), "Expected derived time to be %s", expected)
}

func TestEncodeClaimsTimeValues(t *testing.T) {
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

	c := Claims{
		NotBefore: NewNumericDate(time.Time{}),
		IssuedAt:  NewNumericDate(time.Unix(0, 0)),
		Expiry:    NewNumericDate(now),
	}

	b, err := json.Marshal(c)
	assert.NoError(t, err)

	expected := `{"exp":1451606400,"iat":0}`
	assert.Equal(t, expected, string(b))

	c2 := Claims{}
	if err := json.Unmarshal(b, &c2); assert.NoError(t, err) {
		assert.True(t, c.NotBefore.Time().Equal(c2.NotBefore.Time()))
		assert.True(t, c.IssuedAt.Time().Equal(c2.IssuedAt.Time()))
		assert.True(t, c.Expiry.Time().Equal(c2.Expiry.Time()))
	}
}
