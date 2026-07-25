require 'spec_helper'

describe OpenSource::License::Generator do
  let(:config_dir) { Dir.mktmpdir }
  let(:project_dir) { Dir.mktmpdir }
  let(:config_path) { File.join(config_dir, '.osrc') }

  before do
    stub_const('OpenSource::CONFIG_PATH', config_path)
    File.write(config_path, YAML.dump({ name: 'mt', email: 'mt@example.com' }))
  end

  after do
    FileUtils.remove_entry(config_dir)
    FileUtils.remove_entry(project_dir)
  end

  describe '#generate' do
    around do |example|
      Dir.chdir(project_dir) { example.run }
    end

    context 'when append is requested' do
      let(:options) { { license: 'mit', append: 'README.md' } }

      it 'generates a license file and appends it to the requested file' do
        File.write('README.md', '# Example')

        OpenSource::License::Generator.new(options).generate

        license = File.read('LICENSE')
        readme = File.read('README.md')

        expect(license).to include('The MIT License')
        expect(license).to include('<mt@example.com>')
        expect(license).not_to include('&lt;mt@example.com&gt;')
        expect(readme).to include('## License')
        expect(readme).to include('&lt;mt@example.com&gt;')
        expect(readme).not_to include('<mt@example.com>')
      end
    end

    context 'when append is not requested' do
      let(:options) { { license: 'mit' } }

      it 'generates a license file' do
        OpenSource::License::Generator.new(options).generate

        license = File.read('LICENSE')

        expect(license).to include('The MIT License')
        expect(license).to include('<mt@example.com>')
        expect(license).not_to include('&lt;mt@example.com&gt;')
      end
    end

    it 'raises an OpenSource::LicenseError for unsupported licenses' do
      expect do
        OpenSource::License::Generator.new(license: 'mpl')
      end.to raise_error(OpenSource::LicenseError, /Unsupported license/)
    end

    it 'raises an OpenSource::MissingCredentialsError when owner credentials are missing' do
      File.delete(config_path)

      expect do
        OpenSource::License::Generator.new(license: 'mit').generate
      end.to raise_error(OpenSource::MissingCredentialsError, /Missing .*--setup/)

      expect(File).not_to exist('LICENSE')
    end

    it 'raises an OpenSource::FileError when the license file cannot be written' do
      license_path = File.join(Dir.pwd, 'LICENSE')
      allow(File).to receive(:open).and_call_original
      allow(File).to receive(:open).with(license_path, 'w').and_raise(Errno::EACCES.new(license_path))

      expect do
        OpenSource::License::Generator.new(license: 'mit').generate
      end.to raise_error(OpenSource::FileError, /Unable to write/)
    end

    it 'raises an OpenSource::FileError when the append target cannot be written' do
      append_path = File.expand_path('README.md')
      allow(File).to receive(:open).and_call_original
      allow(File).to receive(:open).with(append_path, 'a').and_raise(Errno::EACCES.new(append_path))

      expect do
        OpenSource::License::Generator.new(license: 'mit', append: 'README.md').generate
      end.to raise_error(OpenSource::FileError, /Unable to append/)
    end
  end
end
