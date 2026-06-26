# FIPS 140-3 cryptographic module — compliance evidence

This document records the validated cryptographic module that a FIPS build of
Syncthing (see [FIPS.md](FIPS.md)) links against, and how to establish the
provenance of a given binary. It is intended as supporting evidence for
**NIST SP 800-171 control 3.13.11** — *"Employ FIPS-validated cryptography when
used to protect the confidentiality of CUI"* (equivalently NIST SP 800-53
SC-13).

> [!NOTE]
> Certificate numbers below are recorded from public CMVP/vendor sources. The
> authoritative record is the live CMVP validated-modules list. An assessor
> should confirm the **current validation status** (Active vs. Historical) and
> the **tested operating environments** on the CMVP site at assessment time —
> validation status and the module-version-to-certificate mapping can change
> over time. Direct retrieval of the CMVP pages was not possible from the build
> environment (egress policy), so the live status has not been embedded here.

## Module linked by this build

Syncthing's FIPS build does **not** use cgo or BoringCrypto. It uses the native
**Go Cryptographic Module**, selected by the `GOFIPS140` build setting, running
in FIPS 140-3 approved mode.

| Field | Value |
| --- | --- |
| Module name | Go Cryptographic Module |
| Validating vendor | Geomys LLC (maintainer of the Go FIPS module, on behalf of the Go project) |
| Module version | v1.0.0 |
| Standard | FIPS 140-3 |
| **CMVP certificate** | **#5247** |
| Security Policy (NIST CSRC) | `140sp5247.pdf` — *Geomys LLC Go Cryptographic Module FIPS 140-3 Non-Proprietary Security Policy* |
| Selected via | `GOFIPS140=v1.0.0` at build time |
| Approved mode at run time | `GODEBUG=fips140=on` (default for a `GOFIPS140` build) |

Source corroboration: NIST publishes the module's non-proprietary security
policy as `140sp5247.pdf`. NIST's `140sp<NNNN>` naming convention encodes the
CMVP certificate number, i.e. certificate **#5247**. Confirm on the live CMVP
listing:

- CMVP certificate: <https://csrc.nist.gov/projects/cryptographic-module-validation-program/certificate/5247>
- CMVP validated modules search: <https://csrc.nist.gov/projects/cryptographic-module-validation-program/validated-modules/search>
- Go project statement: <https://go.dev/doc/security/fips140>

### Tested vs. vendor-affirmed operating environments

FIPS 140-3 certificates list specific *tested* operating environments. The same
module may be run on other environments under vendor affirmation (FIPS 140-3 IG
2.3.C). Syncthing FIPS builds run on Linux, macOS (Apple Silicon and Intel) and
Windows using the identical module source; whether a given OS/CPU is a tested OE
or a vendor-affirmed one must be read from the security policy (`140sp5247.pdf`)
for the deployment's platform. Record the relevant entry as part of the
assessment.

## Alternative module: BoringCrypto (not used here)

If a BoringCrypto-based build is ever required instead, that is a **different
module with a different certificate**, and the evidence above does not apply to
it. For reference:

| Module | Standard | CMVP certificate (verify on CMVP) |
| --- | --- | --- |
| Google BoringCrypto | FIPS 140-3 | #5104 |
| Google BoringCrypto | FIPS 140-2 (historical) | #2964, #3318, #3678, #3753, #4156, #4735 |

BoringCrypto requires `GOEXPERIMENT=boringcrypto`, cgo, and is limited to
`linux/amd64`; it would not cover the macOS/Windows targets this build supports.
This repository's FIPS build deliberately uses the native Go module instead.

## Binary provenance

A FIPS binary records its build inputs, which can be read back with
`go version -m`. This binds the binary to the toolchain, the FIPS module
version, the build tags, and the source revision.

```
$ go version -m ./syncthing
./syncthing: go1.25.0
        build   -tags=fips,fips140v1.0
        build   GOFIPS140=v1.0.0
        build   GOOS=...
        build   GOARCH=...
        build   vcs.revision=<git commit sha>
        build   vcs.time=<commit timestamp>
        build   vcs.modified=false
```

Evidence points for an assessment:

1. **`GOFIPS140=v1.0.0`** — confirms the validated module version was selected.
2. **`fips140v1.0`** build tag — added automatically by the toolchain when
   `GOFIPS140` is set; its presence is proof the module is linked.
3. **`fips`** build tag — Syncthing's own tag that compiles out all
   non-approved cryptography (see FIPS.md).
4. **`go1.25.0`** — the toolchain that bundles the module.
5. **`vcs.revision` / `vcs.modified=false`** — ties the binary to a specific,
   unmodified commit. Build release binaries from a clean checkout so
   `vcs.modified` is `false` and the revision is meaningful.

### Reproducing a build with provenance

```sh
git checkout <release-commit>
GOFIPS140=v1.0.0 go build -tags fips -o syncthing ./cmd/syncthing
go version -m ./syncthing | grep -E 'GOFIPS140|fips140v1.0|vcs.revision'
```

### Run-time evidence

The module runs its power-on self-tests (KATs/CASTs) at startup when in
approved mode. Syncthing additionally logs its FIPS state:

```
INF FIPS 140-3 approved-mode cryptography is active (log.pkg=fips)
```

A FIPS build that is started without approved mode logs a warning instead, so
the operational state is visible in the service logs.

### Confirming no non-approved cryptography is linked

```sh
go list -tags fips -deps ./cmd/syncthing \
  | grep -E 'chacha20|scrypt|bcrypt|miscreant|x/crypto/hkdf' \
  | grep -v '^vendor/'
```

This prints nothing for a FIPS build. (The Go standard library's own vendored
`golang.org/x/crypto/chacha20poly1305`, under a `vendor/` path, is part of
`crypto/tls`/HPKE and is disabled by the module in approved mode.)

## Algorithms in use (all approved)

| Function | Algorithm |
| --- | --- |
| Device identity / certificates | Ed25519 (FIPS 186-5) |
| Sync transport | TLS 1.3, AES-GCM, ECDHE |
| GUI server certificate | ECDSA P-256 |
| GUI password hashing | PBKDF2-HMAC-SHA256 |
| Device IDs, block hashing | SHA-256 |

The specific approved algorithms and their CAVP mappings for certificate #5247
are enumerated in the module's security policy (`140sp5247.pdf`).
