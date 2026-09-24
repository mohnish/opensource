package license

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// ErrUnsupportedLicense is returned when a license identifier has no template.
var ErrUnsupportedLicense = errors.New("unsupported license")

// Options controls what the generator produces.
type Options struct {
	// License is the identifier of the license to generate (e.g. "mit").
	License string
	// Append, when non-empty, is the path of a README-style file to which a
	// "## License" section is appended.
	Append string
}

// Generator writes a LICENSE file (and optionally appends a license section to
// a README) for the configured owner.
type Generator struct {
	opts  Options
	owner *Owner
	now   func() time.Time
}

// NewGenerator returns a Generator for the given options and owner.
func NewGenerator(opts Options, owner *Owner) *Generator {
	return &Generator{opts: opts, owner: owner, now: time.Now}
}

type templateData struct {
	Year      int
	Name      string
	NameUpper string
	Email     string
}

// Generate writes the LICENSE file in the current working directory and, when
// requested, appends a Markdown license section to the append target. Owner
// credentials are resolved before anything is written, so a missing-credentials
// error leaves the filesystem untouched.
func (g *Generator) Generate() error {
	if !IsSupported(g.opts.License) {
		return fmt.Errorf("Unsupported license %q. Supported licenses: %s: %w",
			g.opts.License, strings.Join(Supported(), ", "), ErrUnsupportedLicense)
	}

	creds, err := g.owner.Credentials()
	if err != nil {
		return err
	}

	tmpl, err := g.loadTemplate()
	if err != nil {
		return err
	}

	if err := g.writeLicense(tmpl, creds); err != nil {
		return err
	}

	if g.opts.Append != "" {
		if err := g.appendLicense(tmpl, creds); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) loadTemplate() (*template.Template, error) {
	raw, err := templateFS.ReadFile("templates/" + g.opts.License + templateExt)
	if err != nil {
		return nil, fmt.Errorf("Unable to read %s license template: %w", g.opts.License, err)
	}

	tmpl, err := template.New(g.opts.License).Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse %s license template: %w", g.opts.License, err)
	}

	return tmpl, nil
}

func (g *Generator) writeLicense(tmpl *template.Template, creds Credentials) error {
	content, err := g.render(tmpl, creds, false)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Unable to determine current directory: %w", err)
	}

	path := filepath.Join(cwd, "LICENSE")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("Unable to write %s: %w", path, err)
	}

	return nil
}

func (g *Generator) appendLicense(tmpl *template.Template, creds Credentials) error {
	path, err := filepath.Abs(g.opts.Append)
	if err != nil {
		return fmt.Errorf("Unable to resolve %s: %w", g.opts.Append, err)
	}

	content, err := g.render(tmpl, creds, true)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("Unable to append license to %s: %w", path, err)
	}
	defer file.Close()

	if _, err := file.WriteString("\n## License\n\n" + content); err != nil {
		return fmt.Errorf("Unable to append license to %s: %w", path, err)
	}

	return nil
}

// render executes the license template. When markdown is true the owner email
// is HTML-escaped (<a@b> becomes &lt;a@b&gt;) so it renders correctly inside a
// Markdown document; otherwise the raw angle-bracket form is used.
func (g *Generator) render(tmpl *template.Template, creds Credentials, markdown bool) (string, error) {
	email := "<" + creds.Email + ">"
	if markdown {
		email = "&lt;" + creds.Email + "&gt;"
	}

	data := templateData{
		Year:      g.now().Year(),
		Name:      creds.Name,
		NameUpper: strings.ToUpper(creds.Name),
		Email:     email,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("Unable to render %s license: %w", g.opts.License, err)
	}

	return buf.String(), nil
}
