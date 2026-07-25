module OpenSource
  class Error < StandardError; end
  class OptionError < Error; end
  class ConfigError < Error; end
  class LicenseError < Error; end
  class FileError < Error; end
end
