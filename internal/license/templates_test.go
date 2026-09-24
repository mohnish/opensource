package license

import (
	"reflect"
	"testing"
)

func TestSupportedIsSortedAndComplete(t *testing.T) {
	got := Supported()
	want := []string{"apache2", "bsd", "gpl3", "isc", "mit"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Supported() = %v, want %v", got, want)
	}
}

func TestIsSupported(t *testing.T) {
	cases := map[string]bool{
		"mit":     true,
		"apache2": true,
		"MIT":     false,
		"mpl":     false,
		"":        false,
	}

	for name, want := range cases {
		if got := IsSupported(name); got != want {
			t.Errorf("IsSupported(%q) = %v, want %v", name, got, want)
		}
	}
}
