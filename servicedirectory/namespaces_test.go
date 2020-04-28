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

package servicedirectory

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateNamespace(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := createNamespace(buf, tc.ProjectID)

	got := buf.String()

	if err != nil {
		t.Errorf("CreateNamespace: %v", err)
		return
	}
	if want := "namespaces/golang-test-namespace"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
		return
	}
}

func TestDeleteNamespace(t *testing.T) {
	tc := testutil.SystemTest(t)
	err := deleteNamespace(tc.ProjectID)

	if err != nil {
		t.Errorf("DeleteNamespace: %v", err)
		return
	}
}
