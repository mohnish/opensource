package license

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), ".osrc")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	return path
}

func TestOwnerSaveThenLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".osrc")
	owner := NewOwner(path)

	want := Credentials{Name: "mt", Email: "mt@example.com"}
	if err := owner.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// A fresh Owner reads the file written above.
	got, err := NewOwner(path).Credentials()
	if err != nil {
		t.Fatalf("Credentials: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestOwnerCredentialsStringKeys(t *testing.T) {
	path := writeConfig(t, "name: mt\nemail: mt@example.com\n")

	got, err := NewOwner(path).Credentials()
	if err != nil {
		t.Fatalf("Credentials: %v", err)
	}
	if want := (Credentials{Name: "mt", Email: "mt@example.com"}); got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestOwnerCredentialsLegacySymbolKeys(t *testing.T) {
	// The original Ruby tool serialized symbol keys as ":name" / ":email".
	path := writeConfig(t, ":name: mt\n:email: mt@example.com\n")

	got, err := NewOwner(path).Credentials()
	if err != nil {
		t.Fatalf("Credentials: %v", err)
	}
	if want := (Credentials{Name: "mt", Email: "mt@example.com"}); got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestOwnerCredentialsMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".osrc")

	_, err := NewOwner(path).Credentials()
	if !errors.Is(err, ErrMissingCredentials) {
		t.Fatalf("got %v, want ErrMissingCredentials", err)
	}
}

func TestOwnerCredentialsUnparseable(t *testing.T) {
	path := writeConfig(t, ": [")

	_, err := NewOwner(path).Credentials()
	if err == nil || errors.Is(err, ErrMissingCredentials) {
		t.Fatalf("got %v, want a parse error", err)
	}
}

func TestOwnerCredentialsInvalidShape(t *testing.T) {
	// A YAML sequence is not a valid credentials mapping.
	path := writeConfig(t, "- mt\n- mt@example.com\n")

	_, err := NewOwner(path).Credentials()
	if err == nil {
		t.Fatalf("expected an error for invalid configuration")
	}
}

func TestOwnerCredentialsMissingFields(t *testing.T) {
	path := writeConfig(t, "name: mt\n")

	_, err := NewOwner(path).Credentials()
	if err == nil {
		t.Fatalf("expected an error when email is missing")
	}
}

func TestOwnerCredentialsCached(t *testing.T) {
	path := writeConfig(t, "name: mt\nemail: mt@example.com\n")
	owner := NewOwner(path)

	if _, err := owner.Credentials(); err != nil {
		t.Fatalf("Credentials: %v", err)
	}

	// Removing the file after the first load must not affect the cached value.
	if err := os.Remove(path); err != nil {
		t.Fatalf("removing config: %v", err)
	}

	got, err := owner.Credentials()
	if err != nil {
		t.Fatalf("Credentials (cached): %v", err)
	}
	if want := (Credentials{Name: "mt", Email: "mt@example.com"}); got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestOwnerSaveWritesPrivateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".osrc")
	if err := NewOwner(path).Save(Credentials{Name: "mt", Email: "mt@example.com"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("config permissions = %o, want 600", perm)
	}
}
