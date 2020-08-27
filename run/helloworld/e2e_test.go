package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestHelloworldService(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	service := cloudrunci.NewService("helloworld", tc.ProjectID)
	service.Env = cloudrunci.EnvVars{"NAME": "Override"}
	if err := service.Deploy(); err != nil {
		t.Fatalf("service.Deploy %q: %v", service.Name, err)
	}
	defer service.Clean()

	req, err := service.NewRequest("GET", "/")
	if err != nil {
		t.Fatalf("service.NewRequest: %v", err)
	}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	fmt.Printf("client.Do: %s %s\n", req.Method, req.URL)

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %v", err)
	}

	if got, want := string(out), "Hello Override!\n"; got != want {
		t.Errorf("body: got %q, want %q", got, want)
	}

	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("response status: got %d, want %d", got, http.StatusOK)
	}

}
