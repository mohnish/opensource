package license

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrMissingCredentials is returned when the config file does not exist yet.
// Callers can detect it with errors.Is to offer an interactive setup.
var ErrMissingCredentials = errors.New("missing credentials")

// Credentials identifies the person a license is generated for.
type Credentials struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

// Owner reads and writes the owner credentials stored in the config file
// (historically ~/.osrc). Credentials are cached after the first successful
// load so repeated access does not re-read the file.
type Owner struct {
	configPath string
	cached     *Credentials
}

// NewOwner returns an Owner backed by the given config file path.
func NewOwner(configPath string) *Owner {
	return &Owner{configPath: configPath}
}

// Save writes the credentials to the config file and caches them.
func (o *Owner) Save(creds Credentials) error {
	data, err := yaml.Marshal(creds)
	if err != nil {
		return fmt.Errorf("Unable to encode credentials: %w", err)
	}

	if err := os.WriteFile(o.configPath, data, 0o600); err != nil {
		return fmt.Errorf("Unable to write %s: %w", o.configPath, err)
	}

	saved := creds
	o.cached = &saved

	return nil
}

// Credentials returns the owner credentials, loading them from disk on first
// use. It returns an error wrapping ErrMissingCredentials when the config file
// does not exist.
func (o *Owner) Credentials() (Credentials, error) {
	if o.cached != nil {
		return *o.cached, nil
	}

	creds, err := o.load()
	if err != nil {
		return Credentials{}, err
	}

	o.cached = &creds

	return creds, nil
}

func (o *Owner) load() (Credentials, error) {
	data, err := os.ReadFile(o.configPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Credentials{}, fmt.Errorf("Missing %s; run `opensource --setup` first: %w", o.configPath, ErrMissingCredentials)
		}

		return Credentials{}, fmt.Errorf("Unable to read %s: %w", o.configPath, err)
	}

	var document any
	if err := yaml.Unmarshal(data, &document); err != nil {
		return Credentials{}, fmt.Errorf("Unable to parse %s: %w", o.configPath, err)
	}

	fields, ok := asStringMap(document)
	if !ok {
		return Credentials{}, fmt.Errorf("Invalid configuration in %s; run `opensource --setup` to recreate it", o.configPath)
	}

	name, hasName := fields["name"]
	email, hasEmail := fields["email"]
	if !hasName || !hasEmail {
		return Credentials{}, fmt.Errorf("Invalid configuration in %s; run `opensource --setup` to recreate it", o.configPath)
	}

	return Credentials{Name: name, Email: email}, nil
}

// asStringMap normalizes a decoded YAML mapping into string keys/values. It
// tolerates the legacy Ruby format, which serialized symbol keys as ":name" and
// ":email", by stripping a single leading colon from each key.
func asStringMap(document any) (map[string]string, bool) {
	raw, ok := document.(map[string]any)
	if !ok {
		return nil, false
	}

	fields := make(map[string]string, len(raw))
	for key, value := range raw {
		fields[strings.TrimPrefix(key, ":")] = fmt.Sprint(value)
	}

	return fields, true
}
