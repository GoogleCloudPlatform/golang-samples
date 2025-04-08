// Copyright 2020 Google LLC
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

package servicedirectory

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestServiceDirectory(t *testing.T) {
	tc := testutil.SystemTest(t)
	t.Run("createNamespace", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := createNamespace(buf, tc.ProjectID); err != nil {
			t.Fatalf("CreateNamespace: %v", err)
		}
		got := buf.String()
		if want := "namespaces/golang-test-namespace"; !strings.Contains(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("createService", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := createService(buf, tc.ProjectID); err != nil {
			t.Fatalf("CreateService: %v", err)
		}

		got := buf.String()
		if want := "services/golang-test-service"; !strings.Contains(got, want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("resolveService", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := resolveService(buf, tc.ProjectID); err != nil {
			t.Errorf("CreateService: %v", err)
		}

		got := buf.String()
		if want := "services/golang-test-service"; !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("createEndpoint", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := createEndpoint(buf, tc.ProjectID); err != nil {
			t.Errorf("CreateEndpoint %v", err)
		}

		got := buf.String()
		if want := "endpoints/golang-test-endpoint"; !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("deleteEndpoint", func(t *testing.T) {
		if err := deleteEndpoint(tc.ProjectID); err != nil {
			t.Errorf("DeleteEndpoint: %v", err)
		}
	})
	t.Run("deleteService", func(t *testing.T) {
		if err := deleteService(tc.ProjectID); err != nil {
			t.Errorf("DeleteService: %v", err)
		}
	})
	t.Run("deleteNamespace", func(t *testing.T) {
		if err := deleteNamespace(tc.ProjectID); err != nil {
			t.Errorf("DeleteNamespace: %v", err)
		}
	})
}
