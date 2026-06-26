# FIPS 140-3 mode

This repository can be built in a **FIPS variant** that restricts Syncthing to
cryptography provided by the [Go Cryptographic Module][go-fips] operating in
FIPS 140-3 approved mode. It is intended for deployments that must satisfy a
FIPS 140-3 cryptographic-module requirement and that run on their own,
controlled set of nodes.

> [!IMPORTANT]
> The FIPS variant is **not wire-compatible with the public Syncthing
> network** for every feature (see [Restrictions](#restrictions)). It is meant
> for a private cluster of nodes that are all running a FIPS build, or a
> standard build limited to the features that remain enabled. Core
> device-to-device synchronization remains compatible with standard builds.

## Why the native Go module (and not BoringCrypto)

Go ships a pure-Go, cross-platform cryptographic module that carries the
FIPS 140-3 validation. It is selected with the `GOFIPS140` build setting and/or
the `fips140` GODEBUG and requires no cgo. This is preferred over the older
`GOEXPERIMENT=boringcrypto` approach, which is limited to `linux/amd64`,
requires cgo, and is not officially supported. Because the native module is
pure Go, the FIPS build works on **Linux, macOS (Apple Silicon and Intel) and
Windows**, and cross-compiles like any other Go target.

## Building

Add the `fips` build tag and build against the FIPS module:

```sh
GOFIPS140=v1.0.0 go build -tags fips -o syncthing ./cmd/syncthing
```

Or, using the build script:

```sh
GOFIPS140=v1.0.0 go run build.go -tags fips build
```

`GOFIPS140=v1.0.0` bakes the FIPS module into the binary and makes approved
mode the default at run time. The exact module version string tracks the Go
release; use `latest` to follow the toolchain's bundled version.

The `fips` build tag does two things beyond selecting the module:

1. It **compiles out** every non-approved cryptographic primitive that
   Syncthing's own code or its third-party dependencies would otherwise link
   (ChaCha20-Poly1305 for application data, AES-SIV, scrypt, bcrypt, and the
   QUIC transport — see below). A FIPS binary therefore does not contain that
   code at all, not merely "avoids calling it".
2. It swaps approved replacements into the paths that remain (PBKDF2 for GUI
   password hashing, the approved TLS cipher subset).

### Apple Silicon / Intel macOS

No special steps. The commands above produce a working macOS binary:

```sh
GOOS=darwin GOARCH=arm64 GOFIPS140=v1.0.0 go build -tags fips ./cmd/syncthing  # Apple Silicon
GOOS=darwin GOARCH=amd64 GOFIPS140=v1.0.0 go build -tags fips ./cmd/syncthing  # Intel
```

## Running

A binary built with `GOFIPS140` runs in approved mode by default. You can also
force the module on (or off) at run time with the GODEBUG setting:

```sh
GODEBUG=fips140=on  ./syncthing serve   # approved mode
GODEBUG=fips140=only ./syncthing serve  # approved mode; panic on any non-approved call
```

At startup Syncthing logs its FIPS state, e.g.:

```
INF FIPS 140-3 approved-mode cryptography is active (log.pkg=fips)
```

If a `fips` build is started **without** the module active, it logs a warning
instead, because that is almost always a misconfiguration:

```
WRN This is a FIPS build but the Go cryptographic module is not in approved
    mode; build with GOFIPS140=v1.0.0 (or later) or run with GODEBUG=fips140=on
```

## What stays the same

The parts of Syncthing that matter for ordinary synchronization already use
only FIPS-approved primitives and are unchanged:

| Function | Algorithm | Status |
| --- | --- | --- |
| Device certificates / identity | Ed25519 (FIPS 186-5) | approved |
| GUI/browser certificate | ECDSA P-256 | approved |
| Sync transport | TLS 1.3 (AES-GCM, ECDHE) | approved |
| Device IDs | SHA-256 | approved |
| Block hashing (BEP) | SHA-256 | approved |

A FIPS node and a standard node interoperate for normal (unencrypted) folder
sharing.

## Restrictions

To keep the binary free of non-approved cryptography, the FIPS build disables
the following. Each is enforced so you cannot accidentally rely on it:

- **Encrypted folders** ("receive encrypted" / untrusted devices). These rely
  on ChaCha20-Poly1305, AES-SIV, scrypt and HKDF in a construction that is not
  FIPS-approved. A configuration containing a receive-encrypted folder or a
  device shared with an encryption password is **rejected at startup** with a
  clear error.
- **QUIC transport.** QUIC mandates ChaCha20-Poly1305 for packet protection,
  which is not approved. The FIPS build omits the QUIC dialer/listener; use TCP
  (and the relay) instead. `quic://` addresses are unavailable, exactly as in a
  `noquic` build.
- **bcrypt GUI passwords.** GUI/admin passwords are hashed with
  PBKDF2-HMAC-SHA256 instead of bcrypt. The on-disk hash format differs
  (`$pbkdf2-sha256$...`), so a host migrated from a standard build to a FIPS
  build must have its GUI password set again.

The GUI TLS endpoint also advertises only the approved AES-GCM/AES-CBC ECDHE
cipher suites (the ChaCha20 suites are dropped); the Go module additionally
enforces the approved set at the protocol level.

## Verifying a build is clean

The FIPS binary's dependency graph should contain no non-approved cryptography
from application code or third-party modules. You can confirm with:

```sh
go list -tags fips -deps ./cmd/syncthing \
  | grep -E 'chacha20|scrypt|bcrypt|miscreant|x/crypto/hkdf' \
  | grep -v '^vendor/'
```

This should print nothing. (The Go standard library's own *vendored* copy of
`golang.org/x/crypto/chacha20poly1305`, under a `vendor/` path, is part of
`crypto/tls`/HPKE and is disabled by the module in approved mode; it is not
application code and cannot be removed.)

## Compliance evidence

The validated module identity (CMVP certificate), the BoringCrypto alternative,
and how to establish a binary's provenance for an assessment (e.g. NIST
SP 800-171 control 3.13.11) are documented separately in
[FIPS-COMPLIANCE.md](FIPS-COMPLIANCE.md).

## Scope and caveats

Building on the validated Go Cryptographic Module means the cryptography in use
is FIPS 140-3 validated. Formal validation of the *Syncthing binary itself* is
a separate, organizational certification process and is not claimed here. This
build gives you a Syncthing that uses only approved-mode cryptography, which is
what is normally required.

[go-fips]: https://go.dev/doc/security/fips140
