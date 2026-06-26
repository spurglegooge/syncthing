// Copyright (C) 2024 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build !fips

package config

// checkFIPSConstraints is a no-op in non-FIPS builds. The FIPS build replaces
// it with a version that rejects configuration relying on non-approved
// cryptography (see fips_constraints.go).
func (cfg *Configuration) checkFIPSConstraints() error {
	return nil
}
