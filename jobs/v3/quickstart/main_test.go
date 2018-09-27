// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
