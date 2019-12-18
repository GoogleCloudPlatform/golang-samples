// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudrunci

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

// TestServiceValidateErrors checks for errors in the Service definition.
func TestServiceValidateErrors(t *testing.T) {
	service := Service{Name: "my-serivce"}
	if err := service.validate(); err == nil {
		t.Errorf("service.validate: expected error 'Project ID missing', got success")
	}

	service.ProjectID = "my-project"
	if err := service.validate(); err == nil {
		t.Errorf("service.validate: expected error 'Platform configuration missing', got success")
	}
}

// TestServiceStateErrors checks that a service in the wrong state will be blocked from the requested operation.
func TestServiceStateErrors(t *testing.T) {
	service := NewService("my-service", "my-project")

	want := "Request called before Deploy"
	if _, err := service.Request("GET", "/"); !strings.Contains(err.Error(), want) {
		t.Errorf("service.Request: error expected '%s', got %s", want, err.Error())
	}

	want = "NewRequest called before Deploy"
	if _, err := service.NewRequest("GET", "/"); !strings.Contains(err.Error(), want) {
		t.Errorf("service.NewRequest: error expected '%s', got %s", want, err.Error())
	}

	want = "URL called before Deploy"
	if _, err := service.URL("/"); !strings.Contains(err.Error(), want) {
		t.Errorf("service.URL: error expected '%s', got %s", want, err.Error())
	}

	want = "container image already built"
	service.built = true
	if err := service.Build(); !strings.Contains(err.Error(), want) {
		t.Errorf("service.Build: error expected '%s', got %s", want, err.Error())
	}
}

func TestServiceURL(t *testing.T) {
	want := "https://example.com"
	service := NewService("my-serivce", "my-project")
	mock, err := url.Parse(want)
	if err != nil {
		t.Fatalf("url.Parse: %v", err)
	}
	service.url = mock
	service.deployed = true

	u, err := service.ParsedURL()
	if err != nil {
		t.Errorf("service.ParsedURL: %v", err)
	}
	if u.String() != service.url.String() {
		t.Errorf("service.ParsedURL: got %s, want %s", u, service.url)
	}

	got, err := service.URL("/handler")
	if err != nil {
		t.Errorf("service.URL: %v", err)
	}
	if want := service.url.String() + "/handler"; got != want {
		t.Errorf("service.URL: got %s, want %s", got, want)
	}

	got, err = service.Host()
	if err != nil {
		t.Errorf("service.Host: %v", err)
	}
	if want := service.url.Host + ":443"; got != want {
		t.Errorf("service.URL: got %s, want %s", got, want)
	}
}

func TestDeployArgs(t *testing.T) {
	service := NewService("my-serivce", "my-project")
	service.Image = "gcr.io/my-project/my-service"
	service.Env = EnvVars{
		"NAME1": "value1",
		"NAME2": "value2",
	}

	cmd := service.deployCmd()
	for i, v := range service.Env {
		if !contains(cmd.Args, fmt.Sprintf("%s=%s", i, v)) {
			t.Errorf("Environment variable (%s) missing from deploy command", i)
		}
	}
}

// contains searches for a string value in a string slice.
func contains(haystack []string, needle string) bool {
	for _, i := range haystack {
		if i == needle {
			return true
		}
	}
	return false
}
