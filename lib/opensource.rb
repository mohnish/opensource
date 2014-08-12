require_relative 'opensource/version'
require_relative 'opensource/error'
require_relative 'opensource/owner'
require_relative 'opensource/license'

module Opensource
  LICENSES = Dir.entries(File.expand_path("../../templates", __FILE__)).map do |filename|
    File.basename(filename, '.erb') if !['.', '..'].include?(filename)
  end.compact
end
