// Copyright (C) 2019 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build fips

// This file is the FIPS-build counterpart to encryption_crypto_default.go.
//
// Syncthing's encrypted folder feature relies on ChaCha20-Poly1305, AES-SIV
// (via miscreant), scrypt and HKDF. None of these are part of the FIPS 140-3
// approved algorithm set provided by the Go cryptographic module, so in a
// FIPS build the feature is disabled rather than re-implemented with weaker
// guarantees. Importantly, this file pulls in none of those packages, so a
// FIPS binary never links non-approved cryptography.
//
// These functions should never be reached at runtime: encrypted folders are
// rejected at configuration load time (see lib/config, build tag fips). They
// are defined only so the rest of lib/protocol compiles, and they fail
// closed.

package protocol

import "errors"

// ErrEncryptedFoldersFIPS is returned by the encrypted-folder primitives in a
// FIPS build, where the feature is unavailable.
var ErrEncryptedFoldersFIPS = errors.New("encrypted folders are not supported in FIPS builds")

const (
	// These mirror the AES-256-GCM parameters that an approved-mode build
	// would use. They only need to be defined for the package to compile;
	// the encrypted-folder code paths that consume them are unreachable in a
	// FIPS build.
	nonceSize = 12
	tagSize   = 16
)

func encryptBytes(_ []byte, _ *[keySize]byte) []byte {
	panic(ErrEncryptedFoldersFIPS)
}

func DecryptBytes(_ []byte, _ *[keySize]byte) ([]byte, error) {
	return nil, ErrEncryptedFoldersFIPS
}

func encryptDeterministic(_ []byte, _ *[keySize]byte, _ []byte) []byte {
	panic(ErrEncryptedFoldersFIPS)
}

func decryptDeterministic(_ []byte, _ *[keySize]byte, _ []byte) ([]byte, error) {
	return nil, ErrEncryptedFoldersFIPS
}

func passwordKDF(_, _ string) ([]byte, error) {
	return nil, ErrEncryptedFoldersFIPS
}

func fileKDF(_ *[keySize]byte, _ string) ([keySize]byte, error) {
	return [keySize]byte{}, ErrEncryptedFoldersFIPS
}
