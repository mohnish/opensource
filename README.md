# OpenSource
[![Build Status](https://travis-ci.org/mohnish/opensource.svg?branch=master)](https://travis-ci.org/mohnish/opensource)
[![Code Climate](https://codeclimate.com/github/mohnish/opensource/badges/gpa.svg)](https://codeclimate.com/github/mohnish/opensource)

> Command line tool that lets you add an open source license to your project by running a simple command.

Supports **Ruby 1.9.3+**.

## Supported Licenses

- MIT
- Apache 2
- BSD 3 Clause
- GPL 3

## Installation

Run:

```bash
$ gem install opensource
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

## License

(The MIT License)

Copyright (c) 2017 Mohnish Thallavajhula &lt;m@mt.cx&gt;

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
