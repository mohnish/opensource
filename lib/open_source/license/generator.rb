require 'erb'

module OpenSource
  module License
    class Generator
      def initialize(options)
        @options = options
        @owner = Owner.new
        @license = ERB.new(File.read("#{File.expand_path("../templates", __FILE__)}/#{@options[:license]}.erb"))
      end

      def generate
        create_license_file
        append_to_file if @options[:append]
      end

    private
      def create_license_file
        f = File.new("#{Dir.pwd}/LICENSE", 'w')
        f.write(@license.result(binding))
        f.close
      end

      def append_to_file
        File.open(File.expand_path(@options[:append]), 'a') do |f|
          f << "\n## License\n\n#{@license.result(binding)}"
        end
      end
    end
  end
end
