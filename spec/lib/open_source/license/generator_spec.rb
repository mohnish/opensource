require 'spec_helper'

describe OpenSource::License::Generator do
  let(:generator) { OpenSource::License::Generator.new(options) }
  let(:license_mock) { double('license') }

  before { allow(ERB).to receive(:new).and_return(license_mock) }

  describe '#generate' do
    let(:file_mock) { double('file') }

    context 'when append is requested' do
      let(:options) { { license: 'mit', append: 'README.md' } }

      it 'generates a license file and appends it to the requested file' do
        expect(File).to receive(:new).with("#{Dir.pwd}/LICENSE", 'w').and_return(file_mock)
        expect(license_mock).to receive(:result).once # since we are stubbing the 2nd call
        expect(file_mock).to receive(:write)
        expect(file_mock).to receive(:close)
        expect(generator).to receive(:append_to_file)
        generator.generate
      end
    end

    context 'when append is not requested' do
      let(:options) { { license: 'mit' } }

      it 'generates a license file' do
        expect(File).to receive(:new).with("#{Dir.pwd}/LICENSE", 'w').and_return(file_mock)
        expect(license_mock).to receive(:result).once
        expect(file_mock).to receive(:write)
        expect(file_mock).to receive(:close)
        generator.generate
      end
    end
  end
end
