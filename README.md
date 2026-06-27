[![AirinSync][14]][15]

---

[![MPLv2 License](https://img.shields.io/badge/license-MPLv2-blue.svg?style=flat-square)](https://www.mozilla.org/MPL/2.0/)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/88/badge)](https://bestpractices.coreinfrastructure.org/projects/88)
[![Go Report Card](https://goreportcard.com/badge/github.com/syncthing/syncthing)](https://goreportcard.com/report/github.com/syncthing/syncthing)

AirinSync is a FIPS-hardened distribution of [Syncthing][15], built to use only
FIPS 140-3 approved cryptography (see [FIPS.md](FIPS.md)). It tracks the
upstream Syncthing codebase; the underlying module path and wire protocol are
unchanged, so it interoperates with standard Syncthing for ordinary folder
sharing. Upstream resources (documentation, forum) referenced below apply
except where AirinSync intentionally differs.

## Goals

AirinSync is a **continuous file synchronization program**. It synchronizes
files between two or more computers. We strive to fulfill the goals below.
The goals are listed in order of importance, the most important ones first.
This is the summary version of the goal list - for more
commentary, see the full [Goals document][13].

AirinSync should be:

1. **Safe From Data Loss**

   Protecting the user's data is paramount. We take every reasonable
   precaution to avoid corrupting the user's files.

2. **Secure Against Attackers**

   Again, protecting the user's data is paramount. Regardless of our other
   goals, we must never allow the user's data to be susceptible to
   eavesdropping or modification by unauthorized parties.

3. **Easy to Use**

   AirinSync should be approachable, understandable, and inclusive.

4. **Automatic**

   User interaction should be required only when absolutely necessary.

5. **Universally Available**

   AirinSync should run on every common computer. We are mindful that the
   latest technology is not always available to every individual.

6. **For Individuals**

   AirinSync is primarily about empowering the individual user with safe,
   secure, and easy to use file synchronization.

7. **Everything Else**

   There are many things we care about that don't make it on to the list. It
   is fine to optimize for these values, as long as they are not in conflict
   with the stated goals above.

## Getting Started

Take a look at the [getting started guide][2].

There are a few examples for keeping AirinSync running in the background
on your system in [the etc directory][3]. There are also several [GUI
implementations][11] for Windows, Mac, and Linux.

## Docker

To run AirinSync in Docker, see [the Docker README][16].

## Getting in Touch

For questions about AirinSync specifically, contact your AirinSync maintainer.
For upstream Syncthing, the first and best point of contact is the [Forum][8].
If you've found something that is clearly a bug in the upstream code, feel free
to report it in the upstream [GitHub issue tracker][10].

If you believe that you’ve found a security vulnerability in the upstream
Syncthing code, please report it by emailing security@syncthing.net. Do not
report it in the Forum or issue tracker.

## Building

Building AirinSync from source is easy. After extracting the source bundle from
a release or checking out git, you just need to run `go run build.go` and the
binaries are created in `./bin`. There's [a guide][5] with more details on the
build process.

### FIPS 140-3 mode

AirinSync can be built in a FIPS variant that restricts it to cryptography from
the FIPS 140-3 validated Go cryptographic module
(`GOFIPS140=v1.0.0 go build -tags fips -o AirinSync ./cmd/syncthing`). This is intended for
deployments with a FIPS requirement running on their own set of nodes. See
[FIPS.md](FIPS.md) for how to build, run, and verify it, and for the feature
restrictions it imposes (encrypted folders and QUIC are disabled). Compliance
evidence (CMVP certificate, binary provenance) is in
[FIPS-COMPLIANCE.md](FIPS-COMPLIANCE.md).

## Signed Releases

Release binaries are GPG signed with the key available from
https://syncthing.net/security/. There is also a built-in automatic
upgrade mechanism (disabled in some distribution channels) which uses a
compiled in ECDSA signature. macOS and Windows binaries are also
code-signed.

## Documentation

AirinSync shares the upstream Syncthing [documentation site][6] [[source]][17],
except for the AirinSync-specific behavior documented in [FIPS.md](FIPS.md).

All code is licensed under the [MPLv2 License][7].

[1]: https://docs.syncthing.net/specs/bep-v1.html
[2]: https://docs.syncthing.net/intro/getting-started.html
[3]: https://github.com/syncthing/syncthing/blob/main/etc
[5]: https://docs.syncthing.net/dev/building.html
[6]: https://docs.syncthing.net/
[7]: https://github.com/syncthing/syncthing/blob/main/LICENSE
[8]: https://forum.syncthing.net/
[10]: https://github.com/syncthing/syncthing/issues
[11]: https://docs.syncthing.net/users/contrib.html#gui-wrappers
[13]: https://github.com/syncthing/syncthing/blob/main/GOALS.md
[14]: assets/logo-text-128.png
[15]: https://syncthing.net/
[16]: https://github.com/syncthing/syncthing/blob/main/README-Docker.md
[17]: https://github.com/syncthing/docs
