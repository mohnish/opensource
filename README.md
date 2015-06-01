# Opensource.rb

[![Code Climate](https://codeclimate.com/github/mohnish/opensource.rb/badges/gpa.svg)](https://codeclimate.com/github/mohnish/opensource.rb)

Add an Opensource License to your project. Helps get rid of the annoyance of adding a new license to the project and optionally append the content of the license to the README.md file.

## Supported Licenses

- MIT
- Apache 2
- BSD 3 Clause
- GPL 3

## Installation

Add this line to your application's Gemfile:

```ruby
gem 'opensource'
```

And then execute:

    $ bundle

Or install it yourself as:

    $ gem install opensource

## Usage
```
Usage: opensource OPTIONS

Specific options:
    -s, --setup                      Setup user credentials in ~/.osrc file
    -l, --license LICENSE            LICENSE can be apache2|bsd|gpl3|mit
    -a, --append README              Append LICENSE content to README file

Common options:
    -v, --version                    Print the version
    -h, --help                       Show this message
```

## License

(The MIT License)

Copyright (c) 2015 Mohnish Thallavajhula &lt;i@mohni.sh&gt;

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
