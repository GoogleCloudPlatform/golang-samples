// Package cloudrunci_test holds tests for the cloudrunci package.
package cloudrunci_test

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/cloudrunci"
)

func TestManagedPlatformRequest(t *testing.T) {
	p := cloudrunci.ManagedPlatform{}

	req, err := p.NewRequest("GET", "http://example.com")
	if err != nil {
		t.Errorf("ManagedPlatform.Request: %q", err)
	}
	authzHeader := req.Header.Get("Authorization")
	if authzHeader == "" {
		t.Errorf("ManagedPlatform.Request: missing authentication header: %q", err)
	}
}

func TestGKEPlatformRequest(t *testing.T) {
	p := cloudrunci.GKEPlatform{}

	req, err := p.NewRequest("GET", "http://example.com")
	if err != nil {
		t.Errorf("GKEPlatform.Request: %q", err)
	}
	authzHeader := req.Header.Get("Authorization")
	if authzHeader != "" {
		t.Errorf("GKEPlatform.Request: unexpected authentication header: %q", err)
	}
}

func TestKubernetesPlatformRequest(t *testing.T) {
	p := cloudrunci.KubernetesPlatform{}

	req, err := p.NewRequest("GET", "http://example.com")
	if err != nil {
		t.Errorf("KubernetesPlatform.Request: %q", err)
	}
	authzHeader := req.Header.Get("Authorization")
	if authzHeader != "" {
		t.Errorf("KubernetesPlatform.Request: unexpected authentication header: %q", err)
	}
}
