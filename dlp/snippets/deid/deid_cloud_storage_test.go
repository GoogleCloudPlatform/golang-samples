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
package deid

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeidentifyCloudStorage(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	gcsURI := fmt.Sprint("gs://" + bucketForDeidCloudStorageForInput + "/" + filePathToGCSForDeidTest)
	outputBucket := fmt.Sprint("gs://" + bucketForDeidCloudStorageForOutput)

	fullDeidentifyTemplateID := fmt.Sprint("projects/" + tc.ProjectID + "/deidentifyTemplates/" + deidentifyTemplateID)
	fullDeidentifyStructuredTemplateID := fmt.Sprint("projects/" + tc.ProjectID + "/deidentifyTemplates/" + deidentifyStructuredTemplateID)
	fullRedactImageTemplate := fmt.Sprint("projects/" + tc.ProjectID + "/deidentifyTemplates/" + redactImageTemplate)

	if err := deidentifyCloudStorage(&buf, tc.ProjectID, gcsURI, tableID, dataSetID, outputBucket, fullDeidentifyTemplateID, fullDeidentifyStructuredTemplateID, fullRedactImageTemplate); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Job created successfully:"; !strings.Contains(got, want) {
		t.Errorf("deidentifyCloudStorage got %q, want %q", got, want)
	}
}
