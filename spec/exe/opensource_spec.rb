require 'spec_helper'
require 'open3'
require 'pty'
require 'rbconfig'

describe 'exe/opensource' do
  let(:root_dir) { File.expand_path('../..', __dir__) }
  let(:executable) { File.join(root_dir, 'exe', 'opensource') }
  let(:env) { { 'RUBYLIB' => File.join(root_dir, 'lib'), 'HOME' => home_dir } }
  let(:home_dir) { Dir.mktmpdir }
  let(:project_dir) { Dir.mktmpdir }

  after do
    FileUtils.remove_entry(home_dir)
    FileUtils.remove_entry(project_dir)
  end

  def run_cli(*args, stdin_data: nil)
    Open3.capture3(env, RbConfig.ruby, executable, *args, chdir: project_dir, stdin_data: stdin_data)
  end

  def run_cli_in_terminal(*args, stdin_data:)
    output = +''
    status = nil

    PTY.spawn(env, RbConfig.ruby, executable, *args, chdir: project_dir) do |reader, writer, pid|
      writer.write(stdin_data)
      writer.close

      begin
        loop { output << reader.readpartial(1024) }
      rescue EOFError, Errno::EIO
        nil
      end

      _pid, status = Process.wait2(pid)
    end

    [output, status]
  end

  it 'reports option parser failures through OpenSource::Error' do
    stdout, stderr, status = run_cli('--license', 'mpl')

    expect(status.exitstatus).to eq(1)
    expect(stdout).to include('Error:')
    expect(stdout).to include('invalid argument')
    expect(stderr).to eq('')
  end

  it 'reports missing credentials through OpenSource::Error' do
    stdout, stderr, status = run_cli('--license', 'mit')

    expect(status.exitstatus).to eq(1)
    expect(stdout).to include('Error: Missing')
    expect(stdout).to include('--setup')
    expect(stderr).to eq('')
    expect(File).not_to exist(File.join(project_dir, 'LICENSE'))
  end

  it 'prompts for credentials and continues when missing credentials are hit interactively' do
    output, status = run_cli_in_terminal('--license', 'mit', stdin_data: "mt\nmt@example.com\n")

    expect(status.exitstatus).to eq(0)
    expect(output).to include('Owner credentials are not set')
    expect(output).to include('Enter full name:')
    expect(output).to include('Enter email address:')
    expect(File.read(File.join(project_dir, 'LICENSE'))).to include('<mt@example.com>')
    expect(YAML.load_file(File.join(home_dir, '.osrc'))).to eq({ name: 'mt', email: 'mt@example.com' })
  end

  it 'sets up credentials, generates a license, and appends it to a README' do
    File.write(File.join(project_dir, 'README.md'), "# Example\n")

    setup_stdout, setup_stderr, setup_status = run_cli('--setup', stdin_data: "mt\nmt@example.com\n")
    generate_stdout, generate_stderr, generate_status = run_cli('--license', 'mit', '--append', 'README.md')

    expect(setup_status.exitstatus).to eq(0)
    expect(setup_stdout).to include('Enter full name:')
    expect(setup_stderr).to eq('')
    expect(generate_status.exitstatus).to eq(0)
    expect(generate_stdout).to eq('')
    expect(generate_stderr).to eq('')
    license = File.read(File.join(project_dir, 'LICENSE'))
    readme = File.read(File.join(project_dir, 'README.md'))

    expect(license).to include('<mt@example.com>')
    expect(license).not_to include('&lt;mt@example.com&gt;')
    expect(readme).to include('## License')
    expect(readme).to include('&lt;mt@example.com&gt;')
  end
end
