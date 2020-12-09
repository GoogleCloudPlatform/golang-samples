// Copyright 2020 Google LLC
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

package dataproc

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInstantiateInlineWorkflowTemplate(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	region := "us-central1"

	buf := new(bytes.Buffer)

	if err := instantiateInlineWorkflowTemplate(buf, tc.ProjectID, region); err != nil {
		t.Fatalf("instantiateInlineWorkflowTemplate got err: %v", err)
	}

	got := buf.String()
	if want := fmt.Sprintf("successfully"); !strings.Contains(got, want) {
		t.Fatalf("instantiateInlineWorkflowTemplate got %q, want %q", got, want)
	}
}
