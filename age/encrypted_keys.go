// These functions have been copied from the age project
// https://github.com/FiloSottile/age/blob/101cc8676386b0503571a929a88618cae2f0b1cd/cmd/age/encrypted_keys.go
// https://github.com/FiloSottile/age/blob/101cc8676386b0503571a929a88618cae2f0b1cd/cmd/age/parse.go
//
// Copyright 2021 The age Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in age's LICENSE file at
// https://github.com/FiloSottile/age/blob/v1.0.0/LICENSE
//
// SPDX-License-Identifier: BSD-3-Clause

package age

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"filippo.io/age"
	"filippo.io/age/armor"

	gpgagent "github.com/getsops/gopgagent"
)

type EncryptedIdentity struct {
	Contents            []byte
	Passphrase          func() (string, error)
	NoMatchWarning      func()
	IncorrectPassphrase func()

	identities []age.Identity
}

var _ age.Identity = &EncryptedIdentity{}

func (i *EncryptedIdentity) Unwrap(stanzas []*age.Stanza) (fileKey []byte, err error) {
	if i.identities == nil {
		if err := i.decrypt(); err != nil {
			return nil, err
		}
	}

	for _, id := range i.identities {
		fileKey, err = id.Unwrap(stanzas)
		if errors.Is(err, age.ErrIncorrectIdentity) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return fileKey, nil
	}
	i.NoMatchWarning()
	return nil, age.ErrIncorrectIdentity
}

func (i *EncryptedIdentity) decrypt() error {
	d, err := age.Decrypt(bytes.NewReader(i.Contents), &LazyScryptIdentity{i.Passphrase})
	if e := new(age.NoIdentityMatchError); errors.As(err, &e) {
		// ScryptIdentity returns ErrIncorrectIdentity for an incorrect
		// passphrase, which would lead Decrypt to returning "no identity
		// matched any recipient". That makes sense in the API, where there
		// might be multiple configured ScryptIdentity. Since in cmd/age there
		// can be only one, return a better error message.
		i.IncorrectPassphrase()
		return fmt.Errorf("incorrect passphrase")
	}
	if err != nil {
		return fmt.Errorf("failed to decrypt identity file: %v", err)
	}
	i.identities, err = age.ParseIdentities(d)
	return err
}

// LazyScryptIdentity is an age.Identity that requests a passphrase only if it
// encounters an scrypt stanza. After obtaining a passphrase, it delegates to
// ScryptIdentity.
type LazyScryptIdentity struct {
	Passphrase func() (string, error)
}

var _ age.Identity = &LazyScryptIdentity{}

func (i *LazyScryptIdentity) Unwrap(stanzas []*age.Stanza) (fileKey []byte, err error) {
	for _, s := range stanzas {
		if s.Type == "scrypt" && len(stanzas) != 1 {
			return nil, errors.New("an scrypt recipient must be the only one")
		}
	}
	if len(stanzas) != 1 || stanzas[0].Type != "scrypt" {
		return nil, age.ErrIncorrectIdentity
	}
	pass, err := i.Passphrase()
	if err != nil {
		return nil, fmt.Errorf("could not read passphrase: %v", err)
	}
	ii, err := age.NewScryptIdentity(pass)
	if err != nil {
		return nil, err
	}
	fileKey, err = ii.Unwrap(stanzas)
	return fileKey, err
}

func unwrapIdentities(key string, reader io.Reader) (ParsedIdentities, error) {
	b := bufio.NewReader(reader)
	p, _ := b.Peek(14) // length of "age-encryption" and "-----BEGIN AGE"
	peeked := string(p)

	switch {
	// An age encrypted file, plain or armored.
	case peeked == "age-encryption" || peeked == "-----BEGIN AGE":
		var r io.Reader = b
		if peeked == "-----BEGIN AGE" {
			r = armor.NewReader(r)
		}
		const privateKeySizeLimit = 1 << 24 // 16 MiB
		contents, err := io.ReadAll(io.LimitReader(r, privateKeySizeLimit))
		if err != nil {
			return nil, fmt.Errorf("failed to read '%s': %w", key, err)
		}
		if len(contents) == privateKeySizeLimit {
			return nil, fmt.Errorf("failed to read '%s': file too long", key)
		}
		IncorrectPassphrase := func() {
			conn, err := gpgagent.NewConn()
			if err != nil {
				return
			}
			defer func(conn *gpgagent.Conn) {
				if err := conn.Close(); err != nil {
					log.Errorf("failed to close connection with gpg-agent: %s", err)
				}
			}(conn)
			err = conn.RemoveFromCache(key)
			if err != nil {
				log.Warnf("gpg-agent remove cache request errored: %s", err)
				return
			}
		}
		ids := []age.Identity{&EncryptedIdentity{
			Contents: contents,
			Passphrase: func() (string, error) {
				conn, err := gpgagent.NewConn()
				if err != nil {
					passphrase, err := readPassphrase("Enter passphrase for identity " + key + ":")
					if err != nil {
						return "", err
					}
					return string(passphrase), nil
				}
				defer func(conn *gpgagent.Conn) {
					if err := conn.Close(); err != nil {
						log.Errorf("failed to close connection with gpg-agent: %s", err)
					}
				}(conn)

				req := gpgagent.PassphraseRequest{
					// TODO is the cachekey good enough?
					CacheKey: key,
					Prompt:   "Passphrase",
					Desc:     fmt.Sprintf("Enter passphrase for identity '%s':", key),
				}
				pass, err := conn.GetPassphrase(&req)
				if err != nil {
					return "", fmt.Errorf("gpg-agent passphrase request errored: %s", err)
				}
				//make sure that we won't store empty pass
				if len(pass) == 0 {
					IncorrectPassphrase()
				}
				return pass, nil
			},
			IncorrectPassphrase: IncorrectPassphrase,
			NoMatchWarning: func() {
				log.Warnf("encrypted identity '%s' didn't match file's recipients", key)
			},
		}}
		return ids, nil
	// An unencrypted age identity file.
	default:
		ids, err := parseIdentities(b)
		if err != nil {
			return nil, fmt.Errorf("failed to parse '%s' age identities: %w", key, err)
		}
		return ids, nil
	}
}
