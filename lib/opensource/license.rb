require 'erb'

module Opensource
  class License
    def initialize(options)
      @options = options
      @user = Owner.get_credentials
      @user['escaped_email'] = "&lt;#{@user['email']}&gt;"
      @license = ERB.new(File.read(File.expand_path("./templates/#{@options[:license]}.erb"))).result(binding)
    end

    def process
      generate
      append if @options[:append]
    end

    private
    def generate
      f = File.new("#{Dir.pwd}/LICENSE", "w")
      f.write(@license)
      f.close
    end

    def append
      File.open(File.expand_path(@options[:append]), "a") do |f|
        f << "\n## License\n\n#{@license}"
      end
    end
  end
end