package functest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/functest"
)

func TestURL(t *testing.T) {
	fn := functest.NewCloudFunction("testfn", "testproject")
	got := fn.URL()

	if want := "https://us-central1-testproject.cloudfunctions.net/" + fn.DeployName; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTTPClient(t *testing.T) {
	fn := functest.NewCloudFunction("testfn", "testproject")
	client, err := fn.HTTPClient()
	if err != nil {
		t.Fatalf("CloudFunction.HTTPClient: %v", err)
	}

	helloServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("FUNCTEST_TEST_DEBUG") != "" {
			fmt.Println("[DEBUG] Authoriziation Header:", r.Header.Get("Authorization"))
		}
		//w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Hello, World!"))
	}))
	defer helloServer.Close()

	resp, err := client.Get(helloServer.URL)
	if err != nil {
		t.Fatalf("http.Client.Get: %v", err)
	}

	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
