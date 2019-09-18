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

package template

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTemplateSamples(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)
	fullID := "projects/" + tc.ProjectID + "/inspectTemplates/golang-samples-test-template"
	// Delete template before trying to create it since the test uses the same name every time.
	if err := listInspectTemplates(buf, tc.ProjectID); err != nil {
		t.Errorf("listInspectTemplates: %v", err)
	}

	if got := buf.String(); strings.Contains(got, fullID) {
		buf.Reset()
		if err := deleteInspectTemplate(buf, fullID); err != nil {
			t.Errorf("deleteInspectTemplate: %v", err)
		}
		if got, want := buf.String(), "Successfully deleted inspect template"; !strings.Contains(got, want) {
			t.Errorf("deleteInspectTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
		}
	}

	buf.Reset()
	if err := createInspectTemplate(buf, tc.ProjectID, "golang-samples-test-template", "Test Template", "Template for testing", nil); err != nil {
		t.Errorf("createInspectTemplate: %v", err)
	}
	if got, want := buf.String(), "Successfully created inspect template"; !strings.Contains(got, want) {
		t.Errorf("createInspectTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
	}

	buf.Reset()
	if err := listInspectTemplates(buf, tc.ProjectID); err != nil {
		t.Errorf("listInspectTemplates: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fullID) {
		t.Errorf("listInspectTemplates got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, fullID)
	}

	buf.Reset()
	if err := deleteInspectTemplate(buf, fullID); err != nil {
		t.Errorf("deleteInspectTemplate: %v", err)
	}
	if got, want := buf.String(), "Successfully deleted inspect template"; !strings.Contains(got, want) {
		t.Errorf("deleteInspectTemplate got\n----\n%v\n----\nWant to contain:\n----\n%v\n----", got, want)
	}
}
