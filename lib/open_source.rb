require 'open_source/version'
require 'open_source/error'
require 'open_source/owner'
require 'open_source/license'

module OpenSource
  LICENSES = Dir.entries(File.expand_path("../../templates", __FILE__)).map do |filename|
    File.basename(filename, '.erb') if !['.', '..'].include?(filename)
  end.compact
end
