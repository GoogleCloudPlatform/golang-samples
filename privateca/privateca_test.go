// Copyright 2023 Google LLC
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

package snippets

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/googleapi"
)

// Global variables used in testing
var location string
var projectId string
var r *rand.Rand
var buf bytes.Buffer

// Setup for all tests
func setupTests(t *testing.T) {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	location = "us-central1"
	tc := testutil.SystemTest(t)
	projectId = tc.ProjectID
}

// Setup and teardown functions for CaPoolTests
func setupCaPool(t *testing.T) (string, func(t *testing.T)) {
	caPoolId := fmt.Sprintf("test-ca-pool-%v-%v", time.Now().Format("2006-01-02"), r.Int())

	if err := createCaPool(&buf, projectId, location, caPoolId); err != nil {
		t.Fatal("setupCaPool got err:", err)
	}

	// Return a function to teardown the test
	return caPoolId, func(t *testing.T) {
		if err := deleteCaPool(&buf, projectId, location, caPoolId); err != nil {
			var gerr *googleapi.Error
			if errors.As(err, &gerr) {
				if gerr.Code == 404 {
					t.Log("setupCaPool teardown - skipped CA Pool deletion (not found)")
				} else {
					t.Errorf("setupCaPool teardown got err: %v", err)
				}
			}
		}
	}
}

func TestCreateCaPool(t *testing.T) {
	setupTests(t)

	t.Run("createCaPool", func(t *testing.T) {
		caPoolId := fmt.Sprintf("test-ca-pool-%v-%v", time.Now().Format("2006-01-02"), r.Int())

		buf.Reset()
		if err := createCaPool(&buf, projectId, location, caPoolId); err != nil {
			t.Fatal("createCaPool got err:", err)
		}

		expectedResult := "CA Pool created"
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("createCaPool got %q, want %q", got, expectedResult)
		}

		if err := deleteCaPool(&buf, projectId, location, caPoolId); err != nil {
			t.Fatal("createCaPool teardown got err:", err)
		}
	})

	t.Run("deleteCaPool", func(t *testing.T) {
		caPoolId, teardownCaPoolTests := setupCaPool(t)
		defer teardownCaPoolTests(t)

		buf.Reset()
		if err := deleteCaPool(&buf, projectId, location, caPoolId); err != nil {
			t.Fatal("deleteCaPool got err:", err)
		}

		expectedResult := "CA Pool deleted"
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("deleteCaPool got %q, want %q", got, expectedResult)
		}
	})
}
