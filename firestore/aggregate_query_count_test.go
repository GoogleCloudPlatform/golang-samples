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

package firestore

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCreateCountQuery(t *testing.T) {
	// Note: This test assumes a pre-populated Firestore collection in the
	// below-referenced project.
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	var bytes bytes.Buffer
	err := createCountQuery(&bytes, projectID)
	if err != nil {
		t.Fatal(err)
	}
	got := bytes.String()
	if !strings.Contains(got, "Number of") {
		t.Error("firestore: COUNT sample did not provide correct output")
	}
}
