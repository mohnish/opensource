# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'open_source/version'

Gem::Specification.new do |spec|
  spec.name          = "opensource"
  spec.version       = OpenSource::VERSION
  spec.authors       = ["Mohnish Thallavajhula"]
  spec.email         = ["i@mohni.sh"]
  spec.summary       = %q{Command line tool that lets you add an open source license to your project by running a simple command}
  spec.description   = %q{Command line tool that lets you add an open source license to your project by running a simple command}
  spec.homepage      = "https://github.com/mohnish/opensource.rb"
  spec.license       = "MIT"

  spec.files         = `git ls-files -z`.split("\x0")
  spec.executables   = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }
  spec.test_files    = spec.files.grep(%r{^(test|spec|features)/})
  spec.require_paths = ["lib"]

  spec.add_development_dependency "bundler", "~> 1.6"
  spec.add_development_dependency "rake", "~> 10.0"
end
