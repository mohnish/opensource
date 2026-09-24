package license

import (
	"embed"
	"sort"
	"strings"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

const templateExt = ".tmpl"

// Supported returns the sorted list of license identifiers that can be
// generated. The list is derived from the embedded template files, so adding a
// new template automatically makes a new license available.
func Supported() []string {
	entries, err := templateFS.ReadDir("templates")
	if err != nil {
		// The templates are embedded at build time, so this cannot fail in a
		// correctly built binary.
		panic(err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, strings.TrimSuffix(entry.Name(), templateExt))
	}
	sort.Strings(names)

	return names
}

// IsSupported reports whether name is a license identifier we can generate.
func IsSupported(name string) bool {
	for _, supported := range Supported() {
		if supported == name {
			return true
		}
	}

	return false
}
