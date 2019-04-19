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

package inspect

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInspectString(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := inspectString(buf, tc.ProjectID, "I'm Gary and my email is gary@example.com")
	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}

	got := buf.String()
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInspectTextFile(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := inspectTextFile(buf, tc.ProjectID, "testdata/test.txt")
	if err != nil {
		t.Errorf("TestInspectTextFile: %v", err)
	}

	got := buf.String()
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestInspectImageFile(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := inspectImageFile(buf, tc.ProjectID, "testdata/test.png")
	if err != nil {
		t.Errorf("TestInspectImageFile: %v", err)
	}

	got := buf.String()
	if want := "Info type: PHONE_NUMBER"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
