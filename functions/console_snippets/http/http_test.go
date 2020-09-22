package p

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		want     string
		wantCode int
	}{
		{
			name:     "valid",
			data:     `{"message": "Greetings, Ocean!"}`,
			want:     "Greetings, Ocean!",
			wantCode: http.StatusOK,
		},
		{
			name:     "empty",
			data:     "",
			want:     "Hello, World!",
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid",
			data:     "not-valid-JSON",
			want:     http.StatusText(http.StatusBadRequest) + "\n",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/", strings.NewReader(test.data))
		rr := httptest.NewRecorder()
		HelloWorld(rr, req)

		if got := rr.Result().StatusCode; got != test.wantCode {
			t.Errorf("HelloWorld(%s) Status: got '%d', want '%d'", test.name, got, test.wantCode)
		}

		if got := rr.Body.String(); got != test.want {
			t.Errorf("HelloWorld(%s) Body: got %q, want %q", test.name, got, test.want)
		}
	}
}
