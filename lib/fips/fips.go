// Copyright (C) 2024 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

// Package fips exposes the FIPS 140-3 status of the running binary.
//
// Syncthing can be built in a FIPS variant (with the "fips" build tag) that
// links only FIPS 140-3 approved cryptography. Such a build is intended to be
// compiled against and run with the Go cryptographic module in approved mode,
// i.e. built with GOFIPS140=v1.0.0 (or later) and/or run with
// GODEBUG=fips140=on. See FIPS.md for details.
package fips

import (
	"crypto/fips140"
	"log/slog"
)

// Enabled reports whether the Go cryptographic module is currently operating
// in FIPS 140-3 approved mode (as selected by GOFIPS140 at build time or
// GODEBUG=fips140=on at run time).
func Enabled() bool {
	return fips140.Enabled()
}

// RequiredByBuild reports whether this binary was built with the "fips" build
// tag, meaning non-approved cryptography has been compiled out and approved
// mode is expected to be active at run time.
//
// Its value is set in the build-tagged files required_fips.go and
// required_default.go.
var RequiredByBuild = requiredByBuild

// LogStartupStatus emits a log line describing the FIPS 140-3 state of the
// process. A FIPS build that is not running in approved mode is flagged as a
// warning, since that is almost certainly a misconfiguration.
func LogStartupStatus() {
	switch {
	case Enabled():
		slog.Info("FIPS 140-3 approved-mode cryptography is active")
	case RequiredByBuild:
		slog.Warn("This is a FIPS build but the Go cryptographic module is not in approved mode; build with GOFIPS140=v1.0.0 (or later) or run with GODEBUG=fips140=on")
	}
}
