require 'spec_helper'

describe OpenSource::License::Owner do
  let(:config_dir) { Dir.mktmpdir }
  let(:config_path) { File.join(config_dir, '.osrc') }

  before do
    stub_const('OpenSource::CONFIG_PATH', config_path)
  end

  after do
    FileUtils.remove_entry(config_dir)
  end

  describe '#credentials=' do
    let(:credentials) { { name: 'mt', email: 'mt@example.com' } }

    it 'creates an osrc file' do
      subject.credentials = credentials

      expect(YAML.load_file(config_path)).to eql(credentials)
    end

    it 'raises an OpenSource::ConfigError when the config cannot be written' do
      allow(File).to receive(:open).with(config_path, 'w').and_raise(Errno::EACCES.new(config_path))

      expect { subject.credentials = credentials }.to raise_error(OpenSource::ConfigError, /Unable to write/)
    end
  end

  describe '#credentials' do
    it 'retrieves the credentials from the config file' do
      File.write(config_path, YAML.dump({ name: 'mt', email: 'mt@example.com' }))

      expect(subject.credentials).to eql({ name: 'mt', email: 'mt@example.com' })
    end

    it 'accepts string keys in the config file' do
      File.write(config_path, YAML.dump({ 'name' => 'mt', 'email' => 'mt@example.com' }))

      expect(subject.credentials).to eql({ name: 'mt', email: 'mt@example.com' })
    end

    it 'raises an OpenSource::ConfigError when the config file is missing' do
      expect { subject.credentials }.to raise_error(OpenSource::ConfigError, /Missing .*--setup/)
    end

    it 'raises an OpenSource::ConfigError when the config file cannot be parsed' do
      File.write(config_path, ': [')

      expect { subject.credentials }.to raise_error(OpenSource::ConfigError, /Unable to parse/)
    end

    it 'raises an OpenSource::ConfigError when the config file is invalid' do
      File.write(config_path, YAML.dump(['mt', 'mt@example.com']))

      expect { subject.credentials }.to raise_error(OpenSource::ConfigError, /Invalid configuration/)
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
