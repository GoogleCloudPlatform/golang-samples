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

func TestInspectDataToHybridJobTrigger(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	trigger := jobTriggerForInspectSample
	log.Print("Name:" + trigger)
	if err := inspectDataToHybridJobTrigger(&buf, tc.ProjectID, "My email is test@example.org and my name is Gary.", trigger); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "successfully inspected data using hybrid job trigger"; !strings.Contains(got, want) {
		t.Errorf("TestInspectDataToHybridJobTrigger got %q, want %q", got, want)
	}
	if want := "Findings"; !strings.Contains(got, want) {
		t.Errorf("TestInspectDataToHybridJobTrigger got %q, want %q", got, want)
	}
	if want := "Job State: ACTIVE"; !strings.Contains(got, want) {
		t.Errorf("TestInspectDataToHybridJobTrigger got %q, want %q", got, want)
	}

}
