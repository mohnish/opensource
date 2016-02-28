$LOAD_PATH.unshift File.expand_path('../../lib', __FILE__)
require 'open_source'

RSpec.configure do |config|
  config.order = 'random'
end
