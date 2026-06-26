// Copyright (C) 2019 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package protocol

// testKeyGen is shared by tests across the package. It lives here, untagged,
// rather than in encryption_test.go (which is excluded from FIPS builds)
// because non-encryption tests use it too. Constructing a KeyGenerator never
// invokes the (non-approved) key-derivation primitives; those only run when a
// folder key is actually requested.
var testKeyGen = NewKeyGenerator()
