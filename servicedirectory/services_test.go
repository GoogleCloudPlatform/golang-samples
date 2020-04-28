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
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)


func TestService(t * testing.T) {
	// <test setup code>
	tc := testutil.SystemTest(t)
	createBuf := new(bytes.Buffer)
	err := createNamespace(createBuf, tc.ProjectID)
	if err != nil {
		fmt.Printf("Failed to create namespace in test setup. Skipping: %v", err)
		os.Exit(0)
	}
	t.Run("create", func(t *testing.T) {
		tc := testutil.SystemTest(t)
		buf := new(bytes.Buffer)
		err := createService(buf, tc.ProjectID)

		if err != nil {
			t.Errorf("CreateService: %v", err)
		}

		got := buf.String()
		fmt.Println(got)
		if want := "services/golang-test-service"; !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("resolve", func(t *testing.T) {
		tc := testutil.SystemTest(t)
		buf := new(bytes.Buffer)
		err := resolveService(buf, tc.ProjectID)

		if err != nil {
			t.Errorf("CreateService: %v", err)
		}

		got := buf.String()
		fmt.Println(got)
		if want := "services/golang-test-service"; !strings.Contains(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("delete", func(t *testing.T) {
		tc := testutil.SystemTest(t)
		err := deleteService(tc.ProjectID)

		if err != nil {
			t.Errorf("DeleteService: %v", err)
		}
	})

	// <test tear-down>
	deleteErr := deleteNamespace(tc.ProjectID)

	if deleteErr != nil {
		fmt.Printf("Failed to delete namespace in test tear down: %v.", deleteErr)
	}
}

