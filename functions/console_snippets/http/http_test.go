package p

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	tests := []struct {
		name string
		want string
		data string
	}{
		{
			name: "valid",
			want: "Greetings, Ocean!",
			data: `{"message": "Greetings, Ocean!"}`,
		},
		{
			name: "empty",
			want: "Hello, World!",
			data: "",
		},
		{
			name: "invalid",
			want: "Hello, World!",
			data: "not-valid-JSON",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/", strings.NewReader(test.data))
		rr := httptest.NewRecorder()
		HelloWorld(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("HelloWorld(%s) got %q, want %q", test.name, got, test.want)
		}
	}
}

func TestHelloWorldErrors(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader("not-valid-JSON"))
	rr := httptest.NewRecorder()
	got := runHandler(http.HandlerFunc(HelloWorld), rr, req).String()
	if want := "json.NewDecoder"; !strings.Contains(got, want) {
		t.Errorf("HelloWorld: got %q, want %q", got, want)
	}
}

func runHandler(h http.Handler, rr http.ResponseWriter, req *http.Request) *bytes.Buffer {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	originalWriter := os.Stderr
	log.SetOutput(writer)
	defer log.SetOutput(originalWriter)

	originalFlags := log.Flags()
	log.SetFlags(0)
	defer log.SetFlags(originalFlags)

	h.ServeHTTP(rr, req)
	writer.Flush()
	return &buf
}
