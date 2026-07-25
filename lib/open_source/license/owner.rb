require 'yaml'

module OpenSource
  module License
    class Owner
      def credentials=(credentials)
        File.open(CONFIG_PATH, 'w') do |config_file|
          config_file.write(YAML.dump(credentials))
        end
      rescue SystemCallError => ex
        raise ConfigError, "Unable to write #{CONFIG_PATH}: #{ex.message}"
      end

      def credentials
        @credentials ||= load_credentials
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

    private
      def load_credentials
        credentials = YAML.load_file(CONFIG_PATH)
        credentials = normalize_credentials(credentials)

        unless credentials.key?(:name) && credentials.key?(:email)
          raise ConfigError, "Invalid configuration in #{CONFIG_PATH}; run `opensource --setup` to recreate it"
        end

        credentials
      rescue OpenSource::Error
        raise
      rescue Errno::ENOENT
        raise ConfigError, "Missing #{CONFIG_PATH}; run `opensource --setup` first"
      rescue Psych::Exception => ex
        raise ConfigError, "Unable to parse #{CONFIG_PATH}: #{ex.message}"
      rescue SystemCallError => ex
        raise ConfigError, "Unable to read #{CONFIG_PATH}: #{ex.message}"
      end

      def normalize_credentials(credentials)
        unless credentials.is_a?(Hash)
          raise ConfigError, "Invalid configuration in #{CONFIG_PATH}; run `opensource --setup` to recreate it"
        end

        credentials.each_with_object({}) do |(key, value), normalized|
          normalized[key.to_sym] = value if key.respond_to?(:to_sym)
        end
      end
    end
  end
end
