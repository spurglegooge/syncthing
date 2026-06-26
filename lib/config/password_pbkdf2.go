// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build fips

package config

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// FIPS builds cannot use bcrypt (it is not a FIPS 140-3 approved algorithm),
// so GUI/admin passwords are hashed with PBKDF2-HMAC-SHA256 instead. The
// stored representation is self-describing:
//
//	$pbkdf2-sha256$<iterations>$<base64-salt>$<base64-hash>
//
// This is intentionally not compatible with the bcrypt hashes produced by a
// standard build; a host migrated to a FIPS build needs its GUI password set
// again.

const (
	pbkdf2Prefix     = "$pbkdf2-sha256$"
	pbkdf2Iterations = 600000 // OWASP-recommended minimum for PBKDF2-HMAC-SHA256
	pbkdf2SaltBytes  = 16
	pbkdf2KeyBytes   = 32
)

// SetPassword takes an already-hashed value (in the format above) or a
// plaintext password and stores it. Plaintext passwords are hashed.
func (c *GUIConfiguration) SetPassword(password string) error {
	if strings.HasPrefix(password, pbkdf2Prefix) {
		// Already hashed
		c.Password = password
		return nil
	}

	salt := make([]byte, pbkdf2SaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	dk, err := pbkdf2.Key(sha256.New, password, salt, pbkdf2Iterations, pbkdf2KeyBytes)
	if err != nil {
		return err
	}
	c.Password = fmt.Sprintf("%s%d$%s$%s", pbkdf2Prefix, pbkdf2Iterations,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(dk))
	return nil
}

// CompareHashedPassword returns nil when the given plaintext password matches
// the stored hash.
func (c GUIConfiguration) CompareHashedPassword(password string) error {
	iter, salt, want, err := parsePBKDF2(c.Password)
	if err != nil {
		return err
	}
	got, err := pbkdf2.Key(sha256.New, password, salt, iter, len(want))
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return errors.New("password mismatch")
	}
	return nil
}

func parsePBKDF2(stored string) (iter int, salt, hash []byte, err error) {
	if !strings.HasPrefix(stored, pbkdf2Prefix) {
		return 0, nil, nil, errors.New("unrecognized password hash format")
	}
	parts := strings.Split(strings.TrimPrefix(stored, pbkdf2Prefix), "$")
	if len(parts) != 3 {
		return 0, nil, nil, errors.New("malformed password hash")
	}
	iter, err = strconv.Atoi(parts[0])
	if err != nil || iter <= 0 {
		return 0, nil, nil, errors.New("malformed password hash iterations")
	}
	salt, err = base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, nil, nil, errors.New("malformed password hash salt")
	}
	hash, err = base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return 0, nil, nil, errors.New("malformed password hash value")
	}
	return iter, salt, hash, nil
}
