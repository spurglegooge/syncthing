// Copyright (C) 2024 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build fips

package config

import "fmt"

// checkFIPSConstraints rejects configuration that depends on cryptography not
// available in a FIPS build. Currently that means encrypted folders, whose
// underlying primitives (ChaCha20-Poly1305, AES-SIV, scrypt) are not part of
// the FIPS 140-3 approved set.
func (cfg *Configuration) checkFIPSConstraints() error {
	for _, folder := range cfg.Folders {
		if folder.Type == FolderTypeReceiveEncrypted {
			return fmt.Errorf("folder %q is a receive-encrypted folder, which is not supported in FIPS builds", folder.ID)
		}
		for _, dev := range folder.Devices {
			if dev.EncryptionPassword != "" {
				return fmt.Errorf("folder %q shares with device %v using an encryption password, which is not supported in FIPS builds", folder.ID, dev.DeviceID)
			}
		}
	}
	return nil
}
