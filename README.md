# OpenSource
[![CI](https://github.com/mohnish/opensource/actions/workflows/ci.yml/badge.svg)](https://github.com/mohnish/opensource/actions/workflows/ci.yml)

> Command line tool that lets you add an open source license to your project by running a simple command.

A single, dependency-free binary written in Go.

## Supported Licenses

- MIT
- Apache 2
- BSD 3 Clause
- GPL 3
- ISC

## Installation

### Homebrew cask (macOS)

Install the signed and notarized binary for Apple Silicon or Intel Macs from
the [`mohnish/opensource` tap](https://github.com/mohnish/homebrew-opensource).
Go is not required:

```bash
brew install --cask mohnish/opensource/opensource
```

The fully qualified name selects the project's tap and cask explicitly. See
[Homebrew's tap trust documentation](https://docs.brew.sh/Tap-Trust) for how
Homebrew handles installation from third-party taps.

Verify the installation:

```bash
opensource --version
```

To upgrade later:

```bash
brew update
brew upgrade --cask mohnish/opensource/opensource
```

To uninstall:

```bash
brew uninstall --cask opensource
```

If you previously installed the Homebrew formula, remove it before installing
the cask so both installations do not compete for the `opensource` command:

```bash
brew uninstall --formula opensource
brew install --cask mohnish/opensource/opensource
```

Your saved name and email in `~/.osrc` are retained.

### go install

```bash
go install github.com/mohnish/opensource@latest
```

### Prebuilt binaries

Download an archive for your operating system and architecture from the
[releases page](https://github.com/mohnish/opensource/releases):

| Operating system | Architectures | Archive |
| --- | --- | --- |
| macOS | arm64 (Apple Silicon), amd64 (Intel) | `.tar.gz` |
| Linux | arm64, amd64 | `.tar.gz` |
| Windows | arm64, amd64 | `.zip` |

Each release includes `checksums.txt` with SHA-256 checksums. Extract the archive
and put `opensource` (or `opensource.exe` on Windows) in a directory on your
`PATH`. The macOS binaries are signed and notarized.

## Usage

```bash
Usage: opensource OPTIONS

Specific options:
    -s, --setup                      Setup user credentials in ~/.osrc file
    -l, --license LICENSE            LICENSE can be apache2, bsd, gpl3, isc, mit
    -a, --append README              Append LICENSE content to README file

Common options:
    -v, --version                    Print the version
    -h, --help                       Show this message
```

First, store your name and email (written to `~/.osrc`):

```bash
opensource --setup
```

Then, from any project directory, generate a `LICENSE` file:

```bash
opensource --license mit
```

Optionally append a `## License` section to a README:

```bash
opensource --license mit --append README.md
```

If you run `--license` before setting up credentials in an interactive terminal, you'll be prompted to enter them and the command continues automatically.

## Development

Requires Go (see the `go` directive in [`go.mod`](go.mod) for the minimum version).

Build the binary:

```bash
make build
```

Run the test suite:

```bash
make test
```

Format, vet, and see all shortcuts:

```bash
make fmt
make vet
make
```

Build an unsigned local release snapshot (no credentials or publishing) with
[GoReleaser](https://goreleaser.com):

```bash
make snapshot
```

Run the full release readiness check, including Go checks, all six archives,
checksums, and a smoke test of the packaged binary on your machine:

```bash
make release-check
```

These targets use GoReleaser v2.18.2, matching CI. Go downloads it on first use;
Python 3 is also required for `release-check`. To use an installed GoReleaser:
`make release-check GORELEASER=goreleaser`.

## Releasing

The [Release workflow](.github/workflows/release.yml) tests the code, builds
release binaries, signs and notarizes the macOS binaries, publishes GitHub
Release archives and checksums, then updates `Casks/opensource.rb` in
[`mohnish/homebrew-opensource`](https://github.com/mohnish/homebrew-opensource).
It also publishes Linux and Windows binaries for arm64 and amd64. Binaries
remain GitHub Release assets; neither repository stores them in Git history.

### Maintainer setup

The public `mohnish/homebrew-opensource` repository is initialized on `main`.
GoReleaser creates and updates the cask on each stable release; no manual cask
or checksum edits are needed.

Configure these repository secrets in **mohnish/opensource → Settings → Secrets
and variables → Actions**:

| Secret | Value |
| --- | --- |
| `HOMEBREW_TAP_TOKEN` | Fine-grained GitHub PAT with **Contents: Read and write** for `mohnish/homebrew-opensource`. |
| `MACOS_SIGN_P12` | Base64-encoded Developer ID Application certificate and private key exported as a `.p12`. |
| `MACOS_SIGN_PASSWORD` | Password for the `.p12` export. |
| `MACOS_NOTARY_KEY` | Base64-encoded App Store Connect team API key (`.p8`). |
| `MACOS_NOTARY_KEY_ID` | API key ID. |
| `MACOS_NOTARY_ISSUER_ID` | API issuer ID. |

The tap token must grant access to `mohnish/homebrew-opensource`, and the Apple
credentials must belong to the team issuing the Developer ID certificate.
`GITHUB_TOKEN` is provided automatically by Actions for publishing this project's
release.

Encode each credential file on macOS without line wrapping:

```bash
base64 -i DeveloperIDApplication.p12 | tr -d '\n' | pbcopy
base64 -i AuthKey_KEYID.p8 | tr -d '\n' | pbcopy
```

Run each command separately and save its clipboard value to the corresponding
secret. See [GoReleaser's signing and notarization documentation](https://goreleaser.com/customization/sign/notarize/)
for credential setup.

### Publish a release

Run the local checks before committing:

```bash
make release-check
```

Commit and push the release changes, including any updated release notes in
[`History.md`](History.md). Then create and push the next version tag (for
example, `v3.0.0`):

```bash
git tag -a v3.0.0 -m "opensource v3.0.0"
git push origin v3.0.0
```

Both `v3.0.0` and historical unprefixed tags such as `3.0.0` trigger the workflow;
use one tag per version. Prerelease tags such as `v3.0.0-rc.1` publish GitHub
prereleases without updating the stable Homebrew cask.

For stable releases, a macOS job audits and installs the published cask, checks
the binary's signature and version, and tests setup and license generation in a
temporary directory. Wait for both jobs in the
[Release workflow run](https://github.com/mohnish/opensource/actions/workflows/release.yml)
to succeed before announcing the release.

Local `make release-check` builds unsigned snapshots without publishing and
cannot verify Apple credentials or tap permissions. The tagged workflow
validates those integrations.

## License

(The MIT License)

Copyright (c) 2026 Mohnish Thallavajhula &lt;hi@iam.mt&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
