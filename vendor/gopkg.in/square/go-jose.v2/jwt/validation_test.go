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

	"github.com/stretchr/testify/assert"
)

func TestFieldsMatch(t *testing.T) {
	c := Claims{
		Issuer:   "issuer",
		Subject:  "subject",
		Audience: []string{"a1", "a2"},
		ID:       "42",
	}

	valid := []Expected{
		{Issuer: "issuer"},
		{Subject: "subject"},
		{Audience: Audience{"a1", "a2"}},
		{Audience: Audience{"a2", "a1"}},
		{ID: "42"},
	}

	for _, v := range valid {
		assert.NoError(t, c.Validate(v))
	}

	invalid := []struct {
		Expected Expected
		Error    error
	}{
		{Expected{Issuer: "invalid-issuer"}, ErrInvalidIssuer},
		{Expected{Subject: "invalid-subject"}, ErrInvalidSubject},
		{Expected{Audience: Audience{"invalid-audience"}}, ErrInvalidAudience},
		{Expected{ID: "invalid-id"}, ErrInvalidID},
	}

	for _, v := range invalid {
		assert.Equal(t, v.Error, c.Validate(v.Expected))
	}
}

func TestExpiryAndNotBefore(t *testing.T) {
	now := time.Date(2016, 1, 1, 12, 0, 0, 0, time.UTC)
	twelveHoursAgo := now.Add(-12 * time.Hour)

	c := Claims{
		IssuedAt:  NewNumericDate(twelveHoursAgo),
		NotBefore: NewNumericDate(twelveHoursAgo),
		Expiry:    NewNumericDate(now),
	}

	// expired - default leeway (1 minute)
	assert.NoError(t, c.Validate(Expected{Time: now}))
	err := c.Validate(Expected{Time: now.Add(2 * DefaultLeeway)})
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrExpired)
	}

	// expired - no leeway
	assert.NoError(t, c.ValidateWithLeeway(Expected{Time: now}, 0))
	err = c.ValidateWithLeeway(Expected{Time: now.Add(1 * time.Second)}, 0)
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrExpired)
	}

	// not before - default leeway (1 minute)
	assert.NoError(t, c.Validate(Expected{Time: twelveHoursAgo}))
	err = c.Validate(Expected{Time: twelveHoursAgo.Add(-2 * DefaultLeeway)})
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrNotValidYet)
	}
}

func TestIssuedInFuture(t *testing.T) {
	now := time.Date(2016, 1, 1, 12, 0, 0, 0, time.UTC)
	oneHourInThePast := now.Add(-time.Hour)
	oneHourInTheFuture := now.Add(time.Hour)

	c := Claims{
		IssuedAt:  NewNumericDate(now),
		NotBefore: NewNumericDate(oneHourInThePast),
		Expiry:    NewNumericDate(oneHourInTheFuture),
	}

	// defaults: valid right now, or in the future
	assert.NoError(t, c.Validate(Expected{Time: now}))
	assert.NoError(t, c.Validate(Expected{Time: now.Add(2 * DefaultLeeway)}))

	// cannot be issued in the future
	err := c.Validate(Expected{Time: now.Add(-2 * DefaultLeeway)})
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrIssuedInTheFuture)
	}
	err = c.Validate(Expected{Time: now.Add(-DefaultLeeway - 1)})
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrIssuedInTheFuture)
	}

	// valid when we are within the default leeway
	assert.NoError(t, c.Validate(Expected{Time: now.Add(-DefaultLeeway)}))

	// no leeway: valid up to the exact time; expired if issued one ns in the future
	assert.NoError(t, c.ValidateWithLeeway(Expected{Time: now}, 0))
	err = c.ValidateWithLeeway(Expected{Time: now.Add(-1)}, 0)
	if assert.Error(t, err) {
		assert.Equal(t, err, ErrIssuedInTheFuture)
	}
}

func TestOptionalDateClaims(t *testing.T) {
	var epoch time.Time

	testCases := []struct {
		name  string
		claim Claims
		want  error
	}{
		{
			"no claims",
			Claims{},
			nil,
		},
		{
			"fail nbf",
			Claims{NotBefore: NewNumericDate(time.Now())},
			ErrNotValidYet,
		},
		{
			"fail exp",
			Claims{Expiry: NewNumericDate(epoch.Add(-7 * 24 * time.Hour))},
			ErrExpired,
		},
		{
			"fail iat",
			Claims{IssuedAt: NewNumericDate(time.Now())},
			ErrIssuedInTheFuture,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expect := Expected{}.WithTime(epoch.Add(-24 * time.Hour))
			err := tc.claim.Validate(expect)
			assert.Equal(t, tc.want, err)
		})
	}
}
