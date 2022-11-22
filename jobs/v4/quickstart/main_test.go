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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"cloud.google.com/go/talent/apiv4beta1/talentpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

func checkServiceAvailable(t *testing.T, projectID string) {
	ctx := context.Background()
	c, err := talent.NewCompanyClient(ctx)
	if err != nil {
		t.Skipf("Client: service account likely in different project: %v", err)
	}

	req := &talentpb.ListCompaniesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
	}
	it := c.ListCompanies(ctx, req)
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Skip("List: service account likely in different project")
		}
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

	want := "Done."
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
}
