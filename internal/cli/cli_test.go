package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runContext sets up an isolated project dir (cwd) and config path for a Run.
type runContext struct {
	project    string
	configPath string
}

func newRunContext(t *testing.T) runContext {
	t.Helper()

	project := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(project); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	return runContext{
		project:    project,
		configPath: filepath.Join(t.TempDir(), ".osrc"),
	}
}

func TestRunNoArgsShowsUsage(t *testing.T) {
	rc := newRunContext(t)
	out := &bytes.Buffer{}

	code := Run(nil, IO{Stdin: strings.NewReader(""), Out: out}, "2.1.0", rc.configPath)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "Usage: opensource OPTIONS") {
		t.Errorf("expected usage, got:\n%s", out.String())
	}
}

func TestRunVersion(t *testing.T) {
	rc := newRunContext(t)
	out := &bytes.Buffer{}

	code := Run([]string{"--version"}, IO{Stdin: strings.NewReader(""), Out: out}, "2.1.0", rc.configPath)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if got := strings.TrimSpace(out.String()); got != "2.1.0" {
		t.Errorf("version output = %q, want %q", got, "2.1.0")
	}
}

func TestRunInvalidLicense(t *testing.T) {
	rc := newRunContext(t)
	out := &bytes.Buffer{}

	code := Run([]string{"--license", "mpl"}, IO{Stdin: strings.NewReader(""), Out: out}, "2.1.0", rc.configPath)

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	text := out.String()
	if !strings.Contains(text, "Error:") || !strings.Contains(text, "invalid argument") {
		t.Errorf("unexpected error output:\n%s", text)
	}
}

func TestRunMissingCredentialsNonInteractive(t *testing.T) {
	rc := newRunContext(t)
	out := &bytes.Buffer{}

	code := Run([]string{"--license", "mit"}, IO{Stdin: strings.NewReader(""), Out: out, Interactive: false}, "2.1.0", rc.configPath)

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	text := out.String()
	if !strings.Contains(text, "Error: Missing") || !strings.Contains(text, "--setup") {
		t.Errorf("unexpected error output:\n%s", text)
	}
	if _, err := os.Stat(filepath.Join(rc.project, "LICENSE")); !os.IsNotExist(err) {
		t.Errorf("LICENSE should not have been created")
	}
}

func TestRunSetupThenGenerate(t *testing.T) {
	rc := newRunContext(t)
	if err := os.WriteFile(filepath.Join(rc.project, "README.md"), []byte("# Example\n"), 0o644); err != nil {
		t.Fatalf("writing README: %v", err)
	}

	setupOut := &bytes.Buffer{}
	setupCode := Run([]string{"--setup"},
		IO{Stdin: strings.NewReader("mt\nmt@example.com\n"), Out: setupOut}, "2.1.0", rc.configPath)
	if setupCode != 0 {
		t.Fatalf("setup exit code = %d, want 0", setupCode)
	}
	if !strings.Contains(setupOut.String(), "Enter full name:") {
		t.Errorf("setup missing prompt:\n%s", setupOut.String())
	}

	genOut := &bytes.Buffer{}
	genCode := Run([]string{"--license", "mit", "--append", "README.md"},
		IO{Stdin: strings.NewReader(""), Out: genOut}, "2.1.0", rc.configPath)
	if genCode != 0 {
		t.Fatalf("generate exit code = %d, want 0", genCode)
	}
	if genOut.String() != "" {
		t.Errorf("generate should be silent, got:\n%s", genOut.String())
	}

	license, err := os.ReadFile(filepath.Join(rc.project, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	if !strings.Contains(string(license), "<mt@example.com>") {
		t.Errorf("LICENSE missing raw email:\n%s", license)
	}

	readme, err := os.ReadFile(filepath.Join(rc.project, "README.md"))
	if err != nil {
		t.Fatalf("reading README: %v", err)
	}
	if !strings.Contains(string(readme), "## License") || !strings.Contains(string(readme), "&lt;mt@example.com&gt;") {
		t.Errorf("README not updated correctly:\n%s", readme)
	}
}

func TestRunInteractivePromptsWhenCredentialsMissing(t *testing.T) {
	rc := newRunContext(t)
	out := &bytes.Buffer{}

	code := Run([]string{"--license", "mit"},
		IO{Stdin: strings.NewReader("mt\nmt@example.com\n"), Out: out, Interactive: true}, "2.1.0", rc.configPath)

	if code != 0 {
		t.Fatalf("exit code = %d, want 0\n%s", code, out.String())
	}
	text := out.String()
	if !strings.Contains(text, "Owner credentials are not set") {
		t.Errorf("missing setup notice:\n%s", text)
	}
	if !strings.Contains(text, "Enter full name:") || !strings.Contains(text, "Enter email address:") {
		t.Errorf("missing prompts:\n%s", text)
	}

	license, err := os.ReadFile(filepath.Join(rc.project, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	if !strings.Contains(string(license), "<mt@example.com>") {
		t.Errorf("LICENSE missing email:\n%s", license)
	}

	// Credentials should have been persisted for next time.
	if _, err := os.Stat(rc.configPath); err != nil {
		t.Errorf("config not written: %v", err)
	}
}
