package license

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// chdir switches into dir for the duration of the test.
func chdir(t *testing.T, dir string) {
	t.Helper()

	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
}

func ownerWithCredentials(t *testing.T) *Owner {
	t.Helper()

	owner := NewOwner(filepath.Join(t.TempDir(), ".osrc"))
	if err := owner.Save(Credentials{Name: "mt", Email: "mt@example.com"}); err != nil {
		t.Fatalf("seeding credentials: %v", err)
	}

	return owner
}

func TestGenerateWritesLicense(t *testing.T) {
	project := t.TempDir()
	chdir(t, project)

	gen := NewGenerator(Options{License: "mit"}, ownerWithCredentials(t))
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(project, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "The MIT License") {
		t.Errorf("LICENSE missing title:\n%s", text)
	}
	if !strings.Contains(text, "<mt@example.com>") {
		t.Errorf("LICENSE missing raw email:\n%s", text)
	}
	if strings.Contains(text, "&lt;mt@example.com&gt;") {
		t.Errorf("LICENSE should not contain escaped email:\n%s", text)
	}
}

func TestGenerateAppendsToReadme(t *testing.T) {
	project := t.TempDir()
	chdir(t, project)
	readme := filepath.Join(project, "README.md")
	if err := os.WriteFile(readme, []byte("# Example\n"), 0o644); err != nil {
		t.Fatalf("writing README: %v", err)
	}

	gen := NewGenerator(Options{License: "mit", Append: "README.md"}, ownerWithCredentials(t))
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	license, err := os.ReadFile(filepath.Join(project, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	if !strings.Contains(string(license), "<mt@example.com>") {
		t.Errorf("LICENSE missing raw email")
	}

	appended, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("reading README: %v", err)
	}
	text := string(appended)
	if !strings.HasPrefix(text, "# Example\n") {
		t.Errorf("README lost its original content:\n%s", text)
	}
	if !strings.Contains(text, "\n## License\n\n") {
		t.Errorf("README missing License section:\n%s", text)
	}
	if !strings.Contains(text, "&lt;mt@example.com&gt;") {
		t.Errorf("README missing escaped email:\n%s", text)
	}
	if strings.Contains(text, "<mt@example.com>") {
		t.Errorf("README should use the escaped email, not the raw form:\n%s", text)
	}
}

func TestGenerateAppendCreatesMissingReadme(t *testing.T) {
	project := t.TempDir()
	chdir(t, project)

	gen := NewGenerator(Options{License: "mit", Append: "README.md"}, ownerWithCredentials(t))
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if _, err := os.Stat(filepath.Join(project, "README.md")); err != nil {
		t.Fatalf("expected README.md to be created: %v", err)
	}
}

func TestGenerateUnsupportedLicense(t *testing.T) {
	chdir(t, t.TempDir())

	gen := NewGenerator(Options{License: "mpl"}, ownerWithCredentials(t))
	err := gen.Generate()
	if !errors.Is(err, ErrUnsupportedLicense) {
		t.Fatalf("got %v, want ErrUnsupportedLicense", err)
	}
}

func TestGenerateMissingCredentialsLeavesNoLicense(t *testing.T) {
	project := t.TempDir()
	chdir(t, project)

	owner := NewOwner(filepath.Join(t.TempDir(), ".osrc")) // never saved
	gen := NewGenerator(Options{License: "mit"}, owner)

	err := gen.Generate()
	if !errors.Is(err, ErrMissingCredentials) {
		t.Fatalf("got %v, want ErrMissingCredentials", err)
	}
	if _, statErr := os.Stat(filepath.Join(project, "LICENSE")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("LICENSE should not exist, stat err = %v", statErr)
	}
}

func TestGenerateUsesCurrentYear(t *testing.T) {
	project := t.TempDir()
	chdir(t, project)

	gen := NewGenerator(Options{License: "mit"}, ownerWithCredentials(t))
	gen.now = func() time.Time { return time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC) }
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(project, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	if !strings.Contains(string(content), "Copyright (c) 1999 mt") {
		t.Errorf("LICENSE missing expected year:\n%s", content)
	}
}

func TestGenerateEverySupportedLicense(t *testing.T) {
	for _, name := range Supported() {
		t.Run(name, func(t *testing.T) {
			project := t.TempDir()
			chdir(t, project)

			gen := NewGenerator(Options{License: name}, ownerWithCredentials(t))
			if err := gen.Generate(); err != nil {
				t.Fatalf("Generate(%s): %v", name, err)
			}

			content, err := os.ReadFile(filepath.Join(project, "LICENSE"))
			if err != nil {
				t.Fatalf("reading LICENSE: %v", err)
			}
			if len(content) == 0 {
				t.Fatalf("LICENSE for %s is empty", name)
			}
			// The owner name must always be substituted into the output.
			if !strings.Contains(string(content), "mt") {
				t.Errorf("LICENSE for %s did not include the owner name", name)
			}
		})
	}
}

func TestBSDUppercasesName(t *testing.T) {
	project := t.TempDir()
	chdir(t, project)

	owner := NewOwner(filepath.Join(t.TempDir(), ".osrc"))
	if err := owner.Save(Credentials{Name: "Ada Lovelace", Email: "ada@example.com"}); err != nil {
		t.Fatalf("seeding credentials: %v", err)
	}

	gen := NewGenerator(Options{License: "bsd"}, owner)
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(project, "LICENSE"))
	if err != nil {
		t.Fatalf("reading LICENSE: %v", err)
	}
	if !strings.Contains(string(content), "SHALL ADA LOVELACE BE LIABLE") {
		t.Errorf("BSD LICENSE missing uppercased name:\n%s", content)
	}
}
