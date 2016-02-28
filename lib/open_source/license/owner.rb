require 'yaml'

module OpenSource
  module License
    class Owner
      def credentials=(credentials)
        config_file = File.new(CONFIG_PATH, 'w')
        file_contents = YAML.dump(credentials)
        config_file.write(file_contents)
      rescue StandardError => ex
        OpenSource.logger.fatal("Unable to access #{CONFIG_PATH}")
      end

      def credentials
        @credentials ||= YAML.load_file(CONFIG_PATH)
      end

      def email
        credentials[:email]
      end

      def markdown_supported_email
        "&lt;#{credentials[:email]}&gt;"
      end

      def name
        credentials[:name]
      end
    end
  end
end
