require 'spec_helper'

describe OpenSource::License::Owner do
  let(:config_path) { File.expand_path('~/.osrc') }

  describe '#credentials=' do
    let(:credentials) { { name: 'mt', email: 'mt@example.com' } }
    let(:config_file_mock) { double('config_file') }
    let(:file_contents_mock) { double('file_contents') }

    it 'creates an osrc file' do
      expect(File).to receive(:new).with(config_path, 'w').and_return(config_file_mock)
      expect(YAML).to receive(:dump).with(credentials).and_return(file_contents_mock)
      expect(config_file_mock).to receive(:write).with(file_contents_mock)
      subject.credentials = credentials
    end
  end

  describe '#credentials' do
    it 'retrieves the credentials from the config file' do
      expect(YAML).to receive(:load_file).with(config_path)
      subject.credentials
    end
  end

  describe '#email' do
    before { allow(subject).to receive(:credentials).and_return({ name: 'mt', email: 'mt@example.com' }) }

    it 'returns the email address of the owner' do
      expect(subject.email).to eql('mt@example.com')
    end
  end

  describe '#markdown_supported_email' do
    before { allow(subject).to receive(:credentials).and_return({ name: 'mt', email: 'mt@example.com' }) }

    it 'returns the mark down supported email address of the owner' do
      expect(subject.markdown_supported_email).to eql('&lt;mt@example.com&gt;')
    end
  end

  describe '#name' do
    before { allow(subject).to receive(:credentials).and_return({ name: 'mt', email: 'mt@example.com' }) }

    it 'returns the name of the owner' do
      expect(subject.name).to eql('mt')
    end
  end
end
