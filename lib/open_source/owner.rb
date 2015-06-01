module OpenSource
  module Owner
    extend self

    def set_credentials
      puts "Enter full name: "
      name = gets.chomp
      puts "Enter email address: "
      email = gets.chomp

      f = File.new(File.expand_path("~/.osrc"), "w")
      f.write <<-CREDENTIALS
        name: #{name}
        email: #{email}
      CREDENTIALS
      f.close
    end

    def get_credentials
      user = {}

      IO.foreach(File.expand_path("~/.osrc")) do |line|
        current_line = line.strip
        if !current_line.empty?
          key, value = current_line.split(':')
          user[key.strip] = value.strip
        end
      end

      user['escaped_email'] = "&lt;#{user['email']}&gt;"

      user
    end
  end
end
