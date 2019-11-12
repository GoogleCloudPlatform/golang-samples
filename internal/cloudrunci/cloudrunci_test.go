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
	service := NewService("my-serivce", "my-project")

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
