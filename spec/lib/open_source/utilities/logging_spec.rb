require 'spec_helper'

describe OpenSource::Utilities::Logging do
  let(:klass) { Class.new.extend(OpenSource::Utilities::Logging) }

  describe '.logger' do
    it 'returns the logger' do
      expect(klass.logger.info('woot')).to be_truthy
    end
  end
end
