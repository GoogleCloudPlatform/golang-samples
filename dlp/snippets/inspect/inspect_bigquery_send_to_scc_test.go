// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package inspect

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	dataSetID = "dlp_test_dataset"
	tableID   = "dlp_inspect_test_table_table_id"
)

func TestInspectBigQuerySendToScc(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectBigQuerySendToScc(&buf, tc.ProjectID, dataSetID, tableID); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Job created successfully:"; !strings.Contains(got, want) {
		t.Errorf("InspectBigQuerySendToScc got %q, want %q", got, want)
	}

	jobName := strings.SplitAfter(got, "Job created successfully: ")

	log.Printf("Job Name : %v", jobName)

	deleteJob(tc.ProjectID, jobName[1])
}
