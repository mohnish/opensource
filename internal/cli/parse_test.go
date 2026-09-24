package cli

import (
	"strings"
	"testing"
)

func TestParseActions(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantAction action
		wantLic    string
		wantAppend string
	}{
		{"no args shows help", nil, actionHelp, "", ""},
		{"help long", []string{"--help"}, actionHelp, "", ""},
		{"help short", []string{"-h"}, actionHelp, "", ""},
		{"version long", []string{"--version"}, actionVersion, "", ""},
		{"version short", []string{"-v"}, actionVersion, "", ""},
		{"setup long", []string{"--setup"}, actionSetup, "", ""},
		{"setup short", []string{"-s"}, actionSetup, "", ""},
		{"license space", []string{"--license", "mit"}, actionGenerate, "mit", ""},
		{"license equals", []string{"--license=mit"}, actionGenerate, "mit", ""},
		{"license short space", []string{"-l", "mit"}, actionGenerate, "mit", ""},
		{"license short attached", []string{"-lmit"}, actionGenerate, "mit", ""},
		{"license and append", []string{"-l", "mit", "-a", "README.md"}, actionGenerate, "mit", "README.md"},
		{"append equals", []string{"--append=DOC.md", "-l", "gpl3"}, actionGenerate, "gpl3", "DOC.md"},
		{"help wins over generate", []string{"-l", "mit", "-h"}, actionHelp, "mit", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse(tc.args)
			if err != nil {
				t.Fatalf("parse(%v): %v", tc.args, err)
			}
			if got.action != tc.wantAction {
				t.Errorf("action = %v, want %v", got.action, tc.wantAction)
			}
			if got.license != tc.wantLic {
				t.Errorf("license = %q, want %q", got.license, tc.wantLic)
			}
			if got.append != tc.wantAppend {
				t.Errorf("append = %q, want %q", got.append, tc.wantAppend)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"unsupported license", []string{"--license", "mpl"}, "invalid argument: --license mpl"},
		{"unknown option", []string{"--nope"}, "invalid option: --nope"},
		{"missing license value", []string{"--license"}, "missing argument: --license"},
		{"missing append value", []string{"-l", "mit", "-a"}, "missing argument: -a"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parse(tc.args)
			if err == nil {
				t.Fatalf("parse(%v) = nil error, want %q", tc.args, tc.want)
			}
			if err.Error() != tc.want {
				t.Errorf("error = %q, want %q", err.Error(), tc.want)
			}
		})
	}
}

func TestUsageMentionsAllLicenses(t *testing.T) {
	out := usage()
	for _, want := range []string{"apache2", "bsd", "gpl3", "isc", "mit", "--setup", "--license", "--append", "--version", "--help"} {
		if !strings.Contains(out, want) {
			t.Errorf("usage missing %q:\n%s", want, out)
		}
	}
}
