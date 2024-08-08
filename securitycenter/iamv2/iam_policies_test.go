// Copyright 2024 Google LLC
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

package iamv2

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var sourceName = ""

func setupEntities() error {
	sourceName = os.Getenv("GCLOUD_SOURCE_NAME")
	if sourceName == "" {
		return fmt.Errorf("GCLOUD_SOURCE_NAME not set")
	}
	return nil
}

func setup(t *testing.T) string {
	if sourceName == "" {
		t.Skip("GCLOUD_SOURCE_NAME not set")
	}
	return sourceName
}

func TestMain(m *testing.M) {
	if err := setupEntities(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize IAM test environment: %v\n", err)
		return
	}
	code := m.Run()
	os.Exit(code)
}

func TestGetSourceIamPolicy(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := getSourceIamPolicy(buf, sourceName)

		if err != nil {
			r.Errorf("GetSourceIamPolicy(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if want := "Policy"; !strings.Contains(got, want) {
			r.Errorf("GetSourceIamPolicy(%s) got: %s want %s", sourceName, got, want)
		}
	})
}

func TestSetSourceIamPolicy(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		user := "testuser@example.com"
		err := setSourceIamPolicy(buf, sourceName, user)

		if err != nil {
			r.Errorf("SetSourceIamPolicy(%s, %s) had error: %v", sourceName, user, err)
			return
		}

		got := buf.String()
		if want := fmt.Sprintf("Principal: user:%s Role: roles/securitycenter.findingsEditor", user); !strings.Contains(got, want) {
			r.Errorf("SetSourceIamPolicy(%s, %s) got: %s want %s", sourceName, user, got, want)
		}
	})
}

func TestTestIam(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := testIam(buf, sourceName)

		if err != nil {
			r.Errorf("TestIam(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if want := "Permision to create/update findings? true"; !strings.Contains(got, want) {
			r.Errorf("TestIam(%s) got: %s want %s", sourceName, got, want)
		}
		if want := "Permision to update state? true"; !strings.Contains(got, want) {
			r.Errorf("TestIam(%s) got: %s want %s", sourceName, got, want)
		}
	})
}
