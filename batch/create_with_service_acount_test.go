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

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateJobWithSA(t *testing.T) {
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	jobName := fmt.Sprintf("test-job-go-%v-%v", time.Now().Format("2006-01-02"), r.Int())

	ctx := context.Background()
	projectNumber, err := projectIDtoNumber(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("Could not convert projectID to projectNumber: %v", err)
	}
	serviceAccountAddress := fmt.Sprintf("%d-compute@developer.gserviceaccount.com", projectNumber)
	region := "us-central1"

	buf := &bytes.Buffer{}

	err = createJobWithSA(buf, tc.ProjectID, region, jobName, serviceAccountAddress)

	if err != nil {
		t.Errorf("createBatchUsingServiceAccount got err: %v", err)
	}

	succeeded, err := jobSucceeded(tc.ProjectID, region, jobName)
	if err != nil {
		t.Errorf("Could not verify job completion: %v", err)
	}
	if !succeeded {
		t.Errorf("The test job has failed: %v", err)
	}

	if err := deleteJob(buf, tc.ProjectID, region, jobName); err != nil {
		t.Errorf("deleteJob got err: %v", err)
	}
}
