module OpenSource
  module Owner
    extend self

    OSRC_PATH = File.expand_path("~/.osrc")

    def set_credentials credentials
      f = File.new(OSRC_PATH, "w")
      f.write <<-CREDENTIALS.gsub /^\s+/, ""
        name: #{credentials[:name]}
        email: #{credentials[:email]}
        escaped_email: &lt;#{credentials[:email]}&gt;
      CREDENTIALS
      f.close
    end

    def get_credentials
      credentials = {}

      IO.foreach(OSRC_PATH) do |line|
        current_line = line.strip
        if !current_line.empty?
          key, value = current_line.split(':')
          credentials[key.strip.to_sym] = value.strip
        end
      end

      credentials
    end
  end
end
