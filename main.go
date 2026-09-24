// Command opensource adds an open source license to your project by running a
// simple command.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/mohnish/opensource/internal/cli"
	"golang.org/x/term"
)

// version is set at release time via -ldflags "-X main.version=...". When empty
// it falls back to the module version embedded by the Go toolchain.
var version = ""

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: unable to determine home directory: %s\n", err)
		os.Exit(1)
	}

	stdio := cli.IO{
		Stdin:       os.Stdin,
		Out:         os.Stdout,
		Interactive: term.IsTerminal(int(os.Stdin.Fd())),
	}

	code := cli.Run(os.Args[1:], stdio, resolveVersion(), filepath.Join(home, ".osrc"))
	os.Exit(code)
}

func resolveVersion() string {
	if version != "" {
		return version
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			return strings.TrimPrefix(v, "v")
		}
	}

	return "dev"
}
