module OpenSource
  module Error
    def self.exception(*args)
      RuntimeError.new(*args).extend(self)
    end
  end
end
