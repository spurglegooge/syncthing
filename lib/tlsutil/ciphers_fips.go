// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build fips

package tlsutil

import "crypto/tls"

// The list of cipher suites we will use / suggest for TLS 1.2 connections in
// a FIPS build. This drops the ChaCha20-Poly1305 suites, which are not in the
// FIPS 140-3 approved set, leaving only AES-GCM and AES-CBC suites with
// ephemeral (ECDHE) key exchange. The Go cryptographic module additionally
// enforces the approved set at runtime when FIPS mode is active.
var cipherSuites = []uint16{
	// AES-GCM suites, 256-bit before 128-bit.
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,

	// AES-CBC fallbacks.
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
}
