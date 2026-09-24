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

### Homebrew

```bash
brew install mohnish/tap/opensource
```

To upgrade later:

```bash
brew upgrade opensource
```

### go install

```bash
go install github.com/mohnish/opensource@latest
```

### Prebuilt binaries

Download a binary for your platform from the [releases page](https://github.com/mohnish/opensource/releases) and place it on your `PATH`.

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

Build a local release snapshot (no publishing) with [GoReleaser](https://goreleaser.com):

```bash
make snapshot
```

## Releasing

Releases are automated with GoReleaser and GitHub Actions.

1. One-time setup: create a `homebrew-tap` repository under your account (i.e. `mohnish/homebrew-tap`) and add a repository secret named `HOMEBREW_TAP_GITHUB_TOKEN` to this repo — a fine-grained personal access token with `contents: write` on the tap repo.
2. Cut a release by pushing a semver tag:

   ```bash
   git tag 3.0.0
   git push origin 3.0.0
   ```

The [Release workflow](.github/workflows/release.yml) builds cross-platform binaries, publishes a GitHub Release with checksums, and updates the Homebrew formula in the tap so `brew install mohnish/tap/opensource` picks up the new version.

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
