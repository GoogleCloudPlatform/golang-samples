package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

func init() {
	os.Setenv("EDITOR_UPSTREAM_RENDER_URL", "http://testing.local")
}

func TestEditorHandler(t *testing.T) {
	// Using this helper because template parsing logic is included.
	s, err := NewServiceFromEnv()
	if err != nil {
		t.Fatalf("could not prepare service for testing: %v", err)
	}

	req := httptest.NewRequest("GET", "/", strings.NewReader(""))
	rr := httptest.NewRecorder()
	s.editorHandler(rr, req)

	if got := rr.Result().StatusCode; got != http.StatusOK {
		t.Errorf("response status: got %q, want %q", got, http.StatusOK)
	}

	want := `<title>Markdown Editor</title>`
	re := regexp.MustCompile(`<title>.*</title>`)
	got := re.FindString(rr.Body.String())

	if got != want {
		t.Errorf("body: got %q, want %q", got, want)
	}

	if !strings.Contains(rr.Body.String(), "This UI allows a user to write Markdown text") {
		t.Errorf("body: missing default markdown text from go template")
	}
}

func TestRenderHandlerErrors(t *testing.T) {
	tests := []struct {
		label      string
		req        *http.Request
		wantBody   string
		wantStatus int
	}{
		{
			label:      "Invalid Method",
			req:        httptest.NewRequest("GET", "/render", strings.NewReader("")),
			wantBody:   http.StatusText(http.StatusMethodNotAllowed) + "\n",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			label:      "Invalid JSON",
			req:        httptest.NewRequest("POST", "/render", strings.NewReader("**markdown**")),
			wantBody:   http.StatusText(http.StatusBadRequest) + "\n",
			wantStatus: http.StatusBadRequest,
		},
	}

	// Using this helper because template parsing logic is included.
	s, err := NewServiceFromEnv()
	if err != nil {
		t.Fatalf("could not prepare service for testing: %v", err)
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		s.renderHandler(rr, test.req)

		if got := rr.Result().StatusCode; got != test.wantStatus {
			t.Errorf("%s: response status: got %q, want %q", test.label, got, test.wantStatus)
		}

		if got := rr.Body.String(); got != test.wantBody {
			t.Errorf("%s: body: got %q, want %q", test.label, got, test.wantBody)
		}
	}
}
