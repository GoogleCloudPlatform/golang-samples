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

package main

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/oauth2/google"
	talent "google.golang.org/api/jobs/v3"
)

func checkServiceAvailable(t *testing.T, projectID string) {
	client, err := google.DefaultClient(context.Background(), talent.CloudPlatformScope)
	if err != nil {
		t.Skipf("DefaultClient: %v", err)
	}

	service, err := talent.New(client)
	if err != nil {
		t.Skipf("createCTSService: service account likely in different project: %v", err)
	}
	if _, err := service.Projects.Companies.List("projects/" + projectID).Do(); err != nil {
		t.Skip("List: service account likely in different project")
	}
}

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	checkServiceAvailable(t, tc.ProjectID)

	os.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

	want := "Request ID"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
}
