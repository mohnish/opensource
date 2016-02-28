require 'open_source/version'
require 'open_source/error'
require 'open_source/utilities/logging'
require 'open_source/license/owner'
require 'open_source/license/generator'

module OpenSource
  extend Utilities::Logging

  SUPPORTED_LICENSES = Dir.entries(File.expand_path('../open_source/license/templates', __FILE__)).map do |filename|
    File.basename(filename, '.erb') if !['.', '..'].include?(filename)
  end.compact

  CONFIG_PATH = File.expand_path('~/.osrc')

  def self.request_owner_credentials
    OpenSource.logger.info("Enter full name: ")
    name = gets.chomp
    OpenSource.logger.info("Enter email address: ")
    email = gets.chomp

    { name: name, email: email }
  end

  def self.setup_owner_credentials
    owner_credentials = request_owner_credentials
    owner = License::Owner.new
    owner.credentials = owner_credentials
  end
end
