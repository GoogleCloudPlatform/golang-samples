package cloudrunci_test

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
)

func TestEnvVars(t *testing.T) {
	vars := cloudrunci.EnvVars{
		"a": "1",
	}

	want := "a=1"
	if got := vars.String(); got != want {
		t.Errorf("EnvVars.String: want %s, got %s", want, got)
	}

	vars["b"] = "2"
	vars["c"] = "3"
	want = "a=1,b=2,c=3"
	if got := vars.String(); got != want {
		t.Errorf("EnvVars.String: want %s, got %s", want, got)
	}

	vars["c"] = "7"
	delete(vars, "b")
	want = "a=1,c=7"
	if got := vars.String(); got != want {
		t.Errorf("EnvVars.String: want %s, got %s", want, got)
	}
}

func TestEnvVarsErrors(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{
			name: "empty",
			key:  "",
		},
		{
			name: "hyphenated",
			key:  "-x",
		},
		{
			name: "whitespace",
			key:  " ",
		},
		{
			name: "leading digit",
			key:  "9KEY",
		},
	}

	for _, test := range tests {
		vars := cloudrunci.EnvVars{}
		vars[test.key] = ""

		want := fmt.Sprintf("invalid environment variable names: %s", test.key)
		if err := vars.Validate(); err.Error() != want {
			t.Errorf("envvar key(%s): error expected '%s', got '%s'", test.name, want, err.Error())
		}
	}
}
