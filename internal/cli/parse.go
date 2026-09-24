package cli

import (
	"fmt"
	"strings"

	"github.com/mohnish/opensource/internal/license"
)

func parse(args []string) (parsed, error) {
	var (
		helpReq, versionReq, setupReq bool
		licenseVal, appendVal         string
		licenseSet, appendSet         bool
	)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-h" || arg == "--help":
			helpReq = true
		case arg == "-v" || arg == "--version":
			versionReq = true
		case arg == "-s" || arg == "--setup":
			setupReq = true

		case arg == "-l" || arg == "--license":
			val, next, err := takeValue(args, i, arg)
			if err != nil {
				return parsed{}, err
			}
			licenseVal, licenseSet, i = val, true, next
		case strings.HasPrefix(arg, "--license="):
			licenseVal, licenseSet = strings.TrimPrefix(arg, "--license="), true
		case strings.HasPrefix(arg, "-l") && len(arg) > 2:
			licenseVal, licenseSet = arg[2:], true

		case arg == "-a" || arg == "--append":
			val, next, err := takeValue(args, i, arg)
			if err != nil {
				return parsed{}, err
			}
			appendVal, appendSet, i = val, true, next
		case strings.HasPrefix(arg, "--append="):
			appendVal, appendSet = strings.TrimPrefix(arg, "--append="), true
		case strings.HasPrefix(arg, "-a") && len(arg) > 2:
			appendVal, appendSet = arg[2:], true

		default:
			return parsed{}, fmt.Errorf("invalid option: %s", arg)
		}
	}

	if licenseSet && !license.IsSupported(licenseVal) {
		return parsed{}, fmt.Errorf("invalid argument: --license %s", licenseVal)
	}

	result := parsed{license: licenseVal, append: appendVal}
	switch {
	case helpReq:
		result.action = actionHelp
	case versionReq:
		result.action = actionVersion
	case setupReq:
		result.action = actionSetup
	case !licenseSet && !appendSet:
		// No actionable options: show usage, like invoking with no arguments.
		result.action = actionHelp
	default:
		result.action = actionGenerate
	}

	return result, nil
}

// takeValue returns the argument following the flag at index i, along with the
// index of that value (the caller's loop increment then advances past it).
func takeValue(args []string, i int, flag string) (string, int, error) {
	if i+1 >= len(args) {
		return "", i, fmt.Errorf("missing argument: %s", flag)
	}

	return args[i+1], i + 1, nil
}

func usage() string {
	lines := []string{
		"Usage: opensource OPTIONS",
		"",
		"Specific options:",
		optionLine("-s, --setup", "Setup user credentials in ~/.osrc file"),
		optionLine("-l, --license LICENSE", "LICENSE can be "+strings.Join(license.Supported(), ", ")),
		optionLine("-a, --append README", "Append LICENSE content to README file"),
		"",
		"Common options:",
		optionLine("-v, --version", "Print the version"),
		optionLine("-h, --help", "Show this message"),
		"",
	}

	return strings.Join(lines, "\n")
}

func optionLine(flags, summary string) string {
	return fmt.Sprintf("    %-33s%s", flags, summary)
}
