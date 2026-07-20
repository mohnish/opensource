# OpenSource
[![CI](https://github.com/mohnish/opensource/actions/workflows/ci.yml/badge.svg)](https://github.com/mohnish/opensource/actions/workflows/ci.yml)

> Command line tool that lets you add an open source license to your project by running a simple command.

Supports **Ruby 3.2+**. Local development targets **Ruby 4.0.6**, as declared in `.ruby-version`.

## Supported Licenses

- MIT
- Apache 2
- BSD 3 Clause
- GPL 3

## Installation

Run:

```bash
gem install opensource
```

## Usage

```bash
Usage: opensource OPTIONS

Specific options:
    -s, --setup                      Setup user credentials in ~/.osrc file
    -l, --license LICENSE            LICENSE can be apache2, bsd, gpl3, mit
    -a, --append README              Append LICENSE content to README file

Common options:
    -v, --version                    Print the version
    -h, --help                       Show this message
```

## Development

This project has been updated for the current Ruby toolchain:

- Ruby `4.0.6` for local development
- Ruby `3.2+` as the supported runtime range
- Bundler `4`
- Rake `13`
- RSpec `3.13`
- GitHub Actions CI instead of Travis CI
- Dependabot for Bundler and GitHub Actions updates

Install dependencies:

```bash
make setup
```

Run the test suite:

```bash
make test
```

Build the gem package:

```bash
make build
```

Run the local CI checks:

```bash
make ci
```

See all available shortcuts:

```bash
make
```

CI runs the test suite on Ruby `3.2`, `3.3`, `3.4`, and `4.0`.

### Runtime Dependencies

Ruby 4 no longer loads `logger` as an always-available standard library dependency. The gem declares `logger` explicitly so the CLI works on Ruby 4 and newer RubyGems installations.

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
