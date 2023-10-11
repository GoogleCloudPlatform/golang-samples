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

// Package jobs contains example snippets using the DLP jobs API.
package jobs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateJob(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	// createBucketForCreatJob will create a bucket and upload a txt file
	bucketName, fileName, err := createBucketForCreatJob(t, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	gcsPath := "gs://" + bucketName + "/" + fileName
	infoTypeNames := []string{"EMAIL_ADDRESS", "PERSON_NAME", "LOCATION", "PHONE_NUMBER"}

	if err := createJob(&buf, tc.ProjectID, gcsPath, infoTypeNames); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Created a Dlp Job "; !strings.Contains(got, want) {
		t.Errorf("TestInspectWithCustomRegex got %q, want %q", got, want)
	}

	defer deleteAssetsOfCreateJobTest(t, tc.ProjectID, bucketName, fileName)
}
