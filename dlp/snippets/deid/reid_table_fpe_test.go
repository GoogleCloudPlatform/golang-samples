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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestReidTableDataWithFPE(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)

	keyRingName, err := createKeyRing(t, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	kmsKeyName, wrappedAesKey, keyVersion, err := createKey(t, tc.ProjectID, keyRingName)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyKey(t, tc.ProjectID, keyVersion)

	if err := reidTableDataWithFPE(buf, tc.ProjectID, kmsKeyName, wrappedAesKey); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if want := "Table after re-identification "; !strings.Contains(got, want) {
		t.Errorf("TestReidTableDataWithFPE got %q, want %q", got, want)
	}

}
