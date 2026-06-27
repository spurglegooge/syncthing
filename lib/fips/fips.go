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
	"runtime/debug"
)

// Enabled reports whether the Go cryptographic module is currently operating
// in FIPS 140-3 approved mode (as selected by GOFIPS140 at build time or
// GODEBUG=fips140=on at run time).
func Enabled() bool {
	return fips140.Enabled()
}

// ModuleVersion returns the version of the Go Cryptographic Module that this
// binary was built against, e.g. "v1.0.0". The value is the GOFIPS140 build
// setting, which the toolchain records in the binary when the module is
// selected at build time (GOFIPS140=v1.0.0). It is read back here from the
// embedded build information, so the running binary self-reports the exact
// validated module version it links. Returns "" if no module version was
// pinned at build time.
func ModuleVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, s := range info.Settings {
		if s.Key == "GOFIPS140" {
			if s.Value == "" || s.Value == "off" {
				return ""
			}
			return s.Value
		}
	}
	return ""
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
		if v := ModuleVersion(); v != "" {
			slog.Info("FIPS 140-3 approved-mode cryptography is active (Go Cryptographic Module " + v + ", CMVP #5247)")
		} else {
			slog.Info("FIPS 140-3 approved-mode cryptography is active")
		}
	case RequiredByBuild:
		slog.Warn("This is a FIPS build but the Go cryptographic module is not in approved mode; build with GOFIPS140=v1.0.0 (or later) or run with GODEBUG=fips140=on")
	}
}
