require 'spec_helper'

describe OpenSource do
  describe '::CONFIG_PATH' do
    it 'returns the configuration path of osrc file' do
      expect(OpenSource::CONFIG_PATH).to eql(File.expand_path('~/.osrc'))
    end
  end

  describe '::SUPPORTED_LICENSES' do
    it 'returns a list of supported licenses' do
      expect(OpenSource::SUPPORTED_LICENSES).to include('apache2', 'bsd', 'gpl3', 'mit')
    end
  end

  describe '.request_owner_credentials' do
    let(:stdin_mock) { double('stdin') }
    before { allow(OpenSource).to receive(:gets).and_return(stdin_mock) }

    it 'requests the owner credentials' do
      expect(stdin_mock).to receive(:chomp).twice
      OpenSource.request_owner_credentials
    end
  end

  describe '.setup_owner_credentials' do
    it 'writes owner credentials to the osrc file' do
      owner_mock = double('owner')
      expect(OpenSource).to receive(:request_owner_credentials)
      expect(OpenSource::License::Owner).to receive(:new).and_return(owner_mock)
      expect(owner_mock).to receive('credentials=')
      OpenSource.setup_owner_credentials
    end
  end
end
