require 'open_source/version'
require 'open_source/error'
require 'open_source/owner'
require 'open_source/license'

module OpenSource
  LICENSES = Dir.entries(File.expand_path("../../templates", __FILE__)).map do |filename|
    File.basename(filename, '.erb') if !['.', '..'].include?(filename)
  end.compact

  def self.request_owner_credentials
    puts "Enter full name: "
    name = gets.chomp
    puts "Enter email address: "
    email = gets.chomp

    { name: name, email: email }
  end

  def self.setup_owner_credentails
    owner_credentials = request_owner_credentials
    Owner.set_credentials(owner_credentials)
  end
end
