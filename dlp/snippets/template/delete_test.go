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

package template

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestDeleteTemplate(t *testing.T) {
	tc := testutil.SystemTest(t)

	u := uuid.New().String()[:8]

	templateID := "go-lang-template-test-" + u

	var buf bytes.Buffer
	if err := createInspectTemplateForTest(t, tc.ProjectID, templateID, "Test Template", "Template for testing", nil); err != nil {
		t.Fatalf("createInspectTemplate: %v", err)
	}

	fullID := fmt.Sprintf("projects/%s/locations/global/inspectTemplates/%s", tc.ProjectID, templateID)

	if err := deleteInspectTemplate(&buf, fullID); err != nil {
		t.Fatalf("deleteInspectTemplate: %v", err)
	}
	if got, want := buf.String(), "Successfully deleted inspect template"; !strings.Contains(got, want) {
		t.Errorf("deleteInspectTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
	}

}
