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

	"filippo.io/age"
	"filippo.io/age/agessh"
	"golang.org/x/crypto/ssh"
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

// parseSSHIdentityFromPrivateKeyFile returns an age.Identity from the given
// private key file. If the private key file is encrypted, it will configure
// the identity to prompt for a passphrase.
func parseSSHIdentityFromPrivateKeyFile(keyPath string) (age.Identity, error) {
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
			pass, err := readPassphrase(fmt.Sprintf("Enter passphrase for %q:", keyPath))
			if err != nil {
				return nil, fmt.Errorf("could not read passphrase for %q: %v", keyPath, err)
			}
			return pass, nil
		}
		i, err := agessh.NewEncryptedSSHIdentity(pubKey, contents, passphrasePrompt)
		if err != nil {
			return nil, fmt.Errorf("could not create encrypted SSH identity: %w", err)
		}
		return i, nil
	}
	if err != nil {
		return nil, fmt.Errorf("malformed SSH identity in %q: %w", keyPath, err)
	}
	return id, nil
}
