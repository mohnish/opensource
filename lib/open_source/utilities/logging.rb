require 'logger'

module OpenSource
  module Utilities
    module Logging
      class Simple < Logger::Formatter
        def call(severity, time, progname, msg)
          "OpenSource: #{msg}\n"
        end
      end

      def logger
        @logger ||= initialize_logger
      end

      private
        def initialize_logger
          logger = Logger.new(STDOUT)
          logger.level = Logger::INFO
          logger.formatter = Simple.new
          logger
        end
    end
  end
end
