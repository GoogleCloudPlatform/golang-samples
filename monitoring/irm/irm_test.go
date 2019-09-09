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

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestIncident(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	incident, err := createIncident(buf, tc.ProjectID)
	if err != nil {
		t.Fatalf("createIncident: %v", err)
	}
	if got, want := buf.String(), "Created incident"; !strings.Contains(got, want) {
		t.Errorf("createIncident got\n----\n%s\n----\nWant to contain\n----\n%s\n", got, want)
	}

	buf.Reset()
	if err := changeStage(buf, incident.Name); err != nil {
		t.Errorf("changeStage(%q): %v", incident.Name, err)
	}
	if got, want := buf.String(), "Changed stage"; !strings.Contains(got, want) {
		t.Errorf("changeStage(%q) got\n----\n%s\n----\nWant to contain\n----\n%s\n", incident.Name, got, want)
	}

	buf.Reset()
	if err := changeSeverity(buf, incident.Name); err != nil {
		t.Errorf("changeSeverity: %v", err)
	}
	if got, want := buf.String(), "Changed severity"; !strings.Contains(got, want) {
		t.Errorf("changeSeverity got\n----\n%s\n----\nWant to contain\n----\n%s\n", got, want)
	}

	buf.Reset()
	if err := annotateIncident(buf, incident.Name); err != nil {
		t.Errorf("annotateIncident(%q): %v", incident.Name, err)
	}
	if got, want := buf.String(), "Created annotation"; !strings.Contains(got, want) {
		t.Errorf("annotateIncident(%q) got\n----\n%s\n----\nWant to contain\n----\n%s\n", incident.Name, got, want)
	}
}

func TestCreateSignal(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	if err := createSignal(buf, tc.ProjectID); err != nil {
		t.Fatalf("createSignal: %v", err)
	}
	want := "Created signal"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("createSignal got\n----\n%s\n----\nWant to contain\n----\n%s\n", got, want)
	}
}
