# frozen_string_literal: true

lib = File.expand_path('lib', __dir__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'open_source/version'

Gem::Specification.new do |spec|
  spec.name          = 'opensource'
  spec.version       = OpenSource::VERSION
  spec.authors       = ['Mohnish Thallavajhula']
  spec.email         = ['hi@iam.mt']
  spec.summary       = spec.description = %q{Command line tool that lets you add an open source license to your project by running a simple command}
  spec.homepage      = 'https://github.com/mohnish/opensource'
  spec.license       = 'MIT'
  spec.required_ruby_version = '>= 3.2'
  spec.metadata['rubygems_mfa_required'] = 'true'

  spec.files         = `git ls-files -z`.split("\x0").select { |f| File.file?(f) }.reject { |f| f.match(%r{^(test|spec|features)/}) }
  spec.bindir        = 'exe'
  spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = ['lib']

  spec.add_dependency 'logger', '>= 1.7', '< 2'

  spec.add_development_dependency 'rake', '>= 13', '< 15'
  spec.add_development_dependency 'rspec', '>= 3.13', '< 5'
end
