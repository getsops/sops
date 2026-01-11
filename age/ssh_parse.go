// These functions are similar to those in the age project
// https://github.com/FiloSottile/age/blob/v1.0.0/cmd/age/parse.go
//
// Copyright 2021 The age Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in age's LICENSE file at
// https://github.com/FiloSottile/age/blob/v1.0.0/LICENSE
//
// SPDX-License-Identifier: BSD-3-Clause

package age

import (
	"fmt"
	"io"
	"os"
	"sync"

	"filippo.io/age"
	"filippo.io/age/agessh"
	agesshconv "github.com/Mic92/ssh-to-age"
	"golang.org/x/crypto/ssh"
)

const (
	sshEd25519KeyType   = "ssh-ed25519"
	ageX25519StanzaType = "X25519"
)

// readPublicKeyFile attempts to read a public key based on the given private
// key path. It assumes the public key is in the same directory, with the same
// name, but with a ".pub" extension. If the public key cannot be read, an
// error is returned.
func readPublicKeyFile(privateKeyPath string) (ssh.PublicKey, error) {
	publicKeyPath := privateKeyPath + ".pub"
	f, err := os.Open(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain public %q key for %q SSH key: %w", publicKeyPath, privateKeyPath, err)
	}
	defer f.Close()
	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", publicKeyPath, err)
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(contents)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q: %w", publicKeyPath, err)
	}
	return pubKey, nil
}

// lazyEd25519AgeIdentity wraps an encrypted SSH ed25519 key and lazily converts
// it to an age X25519 identity only when decryption is attempted.
type lazyEd25519AgeIdentity struct {
	contents      []byte
	expectedRecip string // age recipient derived from SSH public key
	getPassphrase func() ([]byte, error)

	mutex   sync.Mutex
	wrapped age.Identity // nil until successfully initialized
}

// matchesRecipient checks if this identity matches the recipient.
func (l *lazyEd25519AgeIdentity) matchesRecipient(recipient string) bool {
	return l.expectedRecip == recipient
}

func (l *lazyEd25519AgeIdentity) Unwrap(stanzas []*age.Stanza) ([]byte, error) {
	// X25519 identities only handle X25519 stanzas. Check before prompting for passphrase.
	hasX25519 := false
	for _, s := range stanzas {
		if s.Type == ageX25519StanzaType {
			hasX25519 = true
			break
		}
	}
	if !hasX25519 {
		return nil, age.ErrIncorrectIdentity
	}

	wrapped, err := l.getOrInitWrapped()
	if err != nil {
		return nil, err
	}
	return wrapped.Unwrap(stanzas)
}

// getOrInitWrapped lazily initializes the wrapped age identity, prompting for
// passphrase if needed. Returns the cached identity on subsequent calls.
func (l *lazyEd25519AgeIdentity) getOrInitWrapped() (age.Identity, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.wrapped != nil {
		return l.wrapped, nil
	}

	passphrase, err := l.getPassphrase()
	if err != nil {
		return nil, fmt.Errorf("could not read passphrase: %w", err)
	}

	ageIdentityStr, _, err := agesshconv.SSHPrivateKeyToAge(l.contents, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to convert SSH key to age identity: %w", err)
	}
	if ageIdentityStr == nil {
		return nil, fmt.Errorf("failed to convert SSH key to age identity: no identity returned")
	}

	l.wrapped, err = age.ParseX25519Identity(*ageIdentityStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse age identity: %w", err)
	}

	return l.wrapped, nil
}

// parseSSHIdentitiesFromPrivateKeyFile returns age identities from the given
// private key file. For ed25519 keys (encrypted or unencrypted), it returns:
//   - An SSH identity (for decrypting data encrypted to SSH recipients)
//   - An age X25519 identity (for decrypting data encrypted to age recipients
//     derived from the same SSH key via ssh-to-age)
//
// For non-ed25519 keys, only the SSH identity is returned.
func parseSSHIdentitiesFromPrivateKeyFile(keyPath string) ([]age.Identity, error) {
	keyFile, err := os.Open(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer keyFile.Close()
	contents, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	id, err := agessh.ParseIdentity(contents)
	if sshErr, ok := err.(*ssh.PassphraseMissingError); ok {
		pubKey := sshErr.PublicKey
		if pubKey == nil {
			pubKey, err = readPublicKeyFile(keyPath)
			if err != nil {
				return nil, err
			}
		}

		passphrasePrompt := func() ([]byte, error) {
			pass, err := readSecret(fmt.Sprintf("Enter passphrase for %q:", keyPath))
			if err != nil {
				return nil, fmt.Errorf("could not read passphrase for %q: %v", keyPath, err)
			}
			return pass, nil
		}

		sshIdentity, err := agessh.NewEncryptedSSHIdentity(pubKey, contents, passphrasePrompt)
		if err != nil {
			return nil, fmt.Errorf("could not create encrypted SSH identity: %w", err)
		}

		identities := []age.Identity{sshIdentity}

		// For ed25519 keys, also create a lazy age X25519 identity
		if pubKey.Type() == sshEd25519KeyType {
			pubKeyBytes := ssh.MarshalAuthorizedKey(pubKey)
			ageRecip, err := agesshconv.SSHPublicKeyToAge(pubKeyBytes)
			if err != nil {
				log.WithField("path", keyPath).Debug("Failed to derive age recipient from SSH public key, skipping age identity: " + err.Error())
			} else if ageRecip != nil {
				identities = append(identities, &lazyEd25519AgeIdentity{
					contents:      contents,
					expectedRecip: *ageRecip,
					getPassphrase: passphrasePrompt,
				})
			}
		}

		return identities, nil
	}
	if err != nil {
		return nil, fmt.Errorf("malformed SSH identity in %q: %w", keyPath, err)
	}

	identities := []age.Identity{id}

	// For ed25519 keys, also create an age X25519 identity so we can decrypt
	// data encrypted to age recipients derived from this SSH key (via ssh-to-age).
	ageIdentityStr, _, err := agesshconv.SSHPrivateKeyToAge(contents, nil)
	if err != nil {
		log.WithField("path", keyPath).Debug("Failed to convert SSH key to age identity, skipping: " + err.Error())
	} else if ageIdentityStr != nil {
		ageIdentity, err := age.ParseX25519Identity(*ageIdentityStr)
		if err != nil {
			log.WithField("path", keyPath).Debug("Failed to parse age identity from converted SSH key: " + err.Error())
		} else {
			identities = append(identities, ageIdentity)
		}
	}

	return identities, nil
}
