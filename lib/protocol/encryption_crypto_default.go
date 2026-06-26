// Copyright (C) 2019 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build !fips

// This file holds the low-level cryptographic primitives backing the
// encrypted folder feature for the standard (non-FIPS) build. They rely on
// ChaCha20-Poly1305, AES-SIV, scrypt and HKDF, none of which are part of the
// FIPS 140-3 approved set. The FIPS build replaces this file with
// encryption_crypto_fips.go, which disables the feature.

package protocol

import (
	"crypto/sha256"
	"errors"
	"io"

	"github.com/miscreant/miscreant.go"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/scrypt"

	"github.com/syncthing/syncthing/lib/rand"
)

const (
	nonceSize     = 24 // chacha20poly1305.NonceSizeX
	tagSize       = 16 // chacha20poly1305.Overhead()
	miscreantAlgo = "AES-SIV"
)

var hkdfSalt = []byte("syncthing")

// encryptBytes encrypts bytes with a random nonce
func encryptBytes(data []byte, key *[keySize]byte) []byte {
	nonce := randomNonce()
	return encrypt(data, nonce, key)
}

// encryptDeterministic encrypts bytes using AES-SIV
func encryptDeterministic(data []byte, key *[keySize]byte, additionalData []byte) []byte {
	aead, err := miscreant.NewAEAD(miscreantAlgo, key[:], 0)
	if err != nil {
		panic("cipher failure: " + err.Error())
	}
	return aead.Seal(nil, nil, data, additionalData)
}

// decryptDeterministic decrypts bytes using AES-SIV
func decryptDeterministic(data []byte, key *[keySize]byte, additionalData []byte) ([]byte, error) {
	aead, err := miscreant.NewAEAD(miscreantAlgo, key[:], 0)
	if err != nil {
		panic("cipher failure: " + err.Error())
	}
	return aead.Open(nil, nil, data, additionalData)
}

func encrypt(data []byte, nonce *[nonceSize]byte, key *[keySize]byte) []byte {
	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		// Can only fail if the key is the wrong length
		panic("cipher failure: " + err.Error())
	}

	if aead.NonceSize() != nonceSize || aead.Overhead() != tagSize {
		// We want these values to be constant for our type declarations so
		// we don't use the values returned by the GCM, but we verify them
		// here.
		panic("crypto parameter mismatch")
	}

	// Data is appended to the nonce
	return aead.Seal(nonce[:], nonce[:], data, nil)
}

// DecryptBytes returns the decrypted bytes, or an error if decryption
// failed.
func DecryptBytes(data []byte, key *[keySize]byte) ([]byte, error) {
	if len(data) < blockOverhead {
		return nil, errors.New("data too short")
	}

	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		// Can only fail if the key is the wrong length
		panic("cipher failure: " + err.Error())
	}

	if aead.NonceSize() != nonceSize || aead.Overhead() != tagSize {
		// We want these values to be constant for our type declarations so
		// we don't use the values returned by the GCM, but we verify them
		// here.
		panic("crypto parameter mismatch")
	}

	return aead.Open(nil, data[:nonceSize], data[nonceSize:], nil)
}

// randomNonce is a normal, cryptographically random nonce
func randomNonce() *[nonceSize]byte {
	var nonce [nonceSize]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		panic("catastrophic randomness failure: " + err.Error())
	}
	return &nonce
}

// passwordKDF derives a folder key from a password using scrypt.
func passwordKDF(folderID, password string) ([]byte, error) {
	return scrypt.Key([]byte(password), knownBytes(folderID), 32768, 8, 1, keySize)
}

// fileKDF derives a per-file key from the folder key and file name using
// HKDF-SHA256.
func fileKDF(folderKey *[keySize]byte, filename string) ([keySize]byte, error) {
	kdf := hkdf.New(sha256.New, append(folderKey[:], filename...), hkdfSalt, nil)
	var fileKey [keySize]byte
	n, err := io.ReadFull(kdf, fileKey[:])
	if err != nil || n != keySize {
		return fileKey, errors.New("hkdf failure")
	}
	return fileKey, nil
}
