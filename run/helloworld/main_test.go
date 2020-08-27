package main

import (
	"net/http/httptest"
	"testing"
	"os"
)

func TestHandler(t *testing.T) {
    tests := []struct{
		label string
		want string
		name string
	}{
		{
			label: "default",
			want: "Hello World!\n",
			name: "",
		},
		{
			label: "override",
			want: "Hello Override!\n",
			name: "Override",
		},
	}

	originalName := os.Getenv("NAME")
	defer os.Setenv("NAME", originalName)

	for _, test := range tests {
		os.Setenv("NAME", test.name)

		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("%s: got %q, want %q", test.label, got, test.want)
		}
	}
}
