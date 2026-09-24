// Package cli implements the opensource command line interface.
package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mohnish/opensource/internal/license"
)

// IO bundles the input/output streams and terminal state the CLI needs, so the
// behavior can be exercised in tests without touching the real process streams.
type IO struct {
	Stdin io.Reader
	Out   io.Writer
	// Interactive reports whether Stdin is connected to a terminal. When true
	// and credentials are missing, the CLI prompts to create them instead of
	// failing.
	Interactive bool
}

type action int

const (
	actionGenerate action = iota
	actionSetup
	actionHelp
	actionVersion
)

type parsed struct {
	action  action
	license string
	append  string
}

// Run parses args and executes the requested command, writing all output to
// io.Out. It returns the process exit code: 0 on success, 1 on any error.
func Run(args []string, io IO, version, configPath string) int {
	opts, err := parse(args)
	if err != nil {
		return fail(io, err)
	}

	switch opts.action {
	case actionHelp:
		fmt.Fprint(io.Out, usage())
		return 0
	case actionVersion:
		fmt.Fprintln(io.Out, version)
		return 0
	case actionSetup:
		if _, err := setupCredentials(io, configPath); err != nil {
			return fail(io, err)
		}
		return 0
	default:
		if err := generate(opts, io, configPath); err != nil {
			return fail(io, err)
		}
		return 0
	}
}

func generate(opts parsed, io IO, configPath string) error {
	owner := license.NewOwner(configPath)
	generator := license.NewGenerator(license.Options{License: opts.license, Append: opts.append}, owner)

	err := generator.Generate()
	if err == nil {
		return nil
	}

	// When credentials are missing and we're attached to a terminal, set them
	// up interactively and retry, mirroring the original tool's behavior.
	if errors.Is(err, license.ErrMissingCredentials) && io.Interactive {
		fmt.Fprintln(io.Out, "Owner credentials are not set. Let's set them up now.")
		if _, err := setupCredentials(io, configPath); err != nil {
			return err
		}
		return generator.Generate()
	}

	return err
}

func setupCredentials(io IO, configPath string) (license.Credentials, error) {
	creds, err := promptCredentials(io)
	if err != nil {
		return license.Credentials{}, err
	}

	if err := license.NewOwner(configPath).Save(creds); err != nil {
		return license.Credentials{}, err
	}

	return creds, nil
}

func promptCredentials(io IO) (license.Credentials, error) {
	reader := bufio.NewReader(io.Stdin)

	fmt.Fprintln(io.Out, "Enter full name: ")
	name, err := readLine(reader)
	if err != nil {
		return license.Credentials{}, err
	}

	fmt.Fprintln(io.Out, "Enter email address: ")
	email, err := readLine(reader)
	if err != nil {
		return license.Credentials{}, err
	}

	return license.Credentials{Name: name, Email: email}, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("Unable to read input: %w", err)
	}

	return strings.TrimSpace(line), nil
}

func fail(io IO, err error) int {
	fmt.Fprintf(io.Out, "Error: %s\n", err)
	return 1
}
