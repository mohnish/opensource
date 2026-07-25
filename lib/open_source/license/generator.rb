require 'erb'

module OpenSource
  module License
    class Generator
      def initialize(options)
        @options = options
        @owner = Owner.new
        @license = load_license_template
      end

      def generate
        create_license_file
        append_to_file if @options[:append]
      end

    private
      def load_license_template
        unless OpenSource::SUPPORTED_LICENSES.include?(@options[:license])
          raise LicenseError, "Unsupported license #{@options[:license].inspect}. Supported licenses: #{OpenSource::SUPPORTED_LICENSES.join(', ')}"
        end

        ERB.new(File.read(template_path))
      rescue OpenSource::Error
        raise
      rescue SystemCallError => ex
        raise LicenseError, "Unable to read #{@options[:license]} license template: #{ex.message}"
      end

      def template_path
        "#{File.expand_path("../templates", __FILE__)}/#{@options[:license]}.erb"
      end

      def create_license_file
        license_content = render_license

        File.open(license_path, 'w') do |file|
          file.write(license_content)
        end
      rescue OpenSource::Error
        raise
      rescue SystemCallError => ex
        raise FileError, "Unable to write #{license_path}: #{ex.message}"
      end

      def append_to_file
        append_path = File.expand_path(@options[:append])
        license_content = render_license

        File.open(append_path, 'a') do |file|
          file << "\n## License\n\n#{license_content}"
        end
      rescue OpenSource::Error
        raise
      rescue SystemCallError => ex
        raise FileError, "Unable to append license to #{append_path}: #{ex.message}"
      end

      def license_path
        "#{Dir.pwd}/LICENSE"
      end

      def render_license
        @license.result(binding)
      rescue OpenSource::Error
        raise
      rescue StandardError => ex
        raise LicenseError, "Unable to render #{@options[:license]} license: #{ex.message}"
      end
    end
  end
end
