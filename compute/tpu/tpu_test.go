// Copyright 2024 Google LLC
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
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTPU(t *testing.T) {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	var buf bytes.Buffer
	tc := testutil.SystemTest(t)
	location := "europe-west4-a"
	resourceName := fmt.Sprintf("projects/%s/locations/%s/nodes/", tc.ProjectID, location)

	t.Run("Create TPU", func(t *testing.T) {
		buf.Reset()
		nodeName := "test-" + fmt.Sprint(seededRand.Int())
		nodeFullName := resourceName + nodeName

		err := createTPUNode(&buf, tc.ProjectID, location, nodeName)
		if err != nil {
			t.Error(err)
		}
		defer deleteTPUNode(&buf, tc.ProjectID, location, nodeName)

		expectedResult := fmt.Sprintf("Node created: %s", nodeFullName)
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("createTpuNode got %q, want %q", got, expectedResult)
		}
	})

	t.Run("Get TPU", func(t *testing.T) {
		buf.Reset()
		nodeName := "test-" + fmt.Sprint(seededRand.Int())
		nodeFullName := resourceName + nodeName

		err := createTPUNode(&buf, tc.ProjectID, location, nodeName)
		if err != nil {
			t.Errorf("failed to create node: %v", err)
		}
		defer deleteTPUNode(&buf, tc.ProjectID, location, nodeName)

		err = getTPUNode(&buf, tc.ProjectID, location, nodeName)
		if err != nil {
			t.Errorf("failed to get node: %v", err)
		}

		expectedResult := fmt.Sprintf("Got node: %s", nodeFullName)
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("getTpuNode got %q, want %q", got, expectedResult)
		}
	})

	t.Run("Delete TPU", func(t *testing.T) {
		buf.Reset()
		nodeName := "test-" + fmt.Sprint(seededRand.Int())
		nodeFullName := resourceName + nodeName

		err := createTPUNode(&buf, tc.ProjectID, location, nodeName)
		if err != nil {
			t.Errorf("failed to create node: %v", err)
		}
		err = deleteTPUNode(&buf, tc.ProjectID, location, nodeName)
		if err != nil {
			t.Errorf("failed to delete node: %v", err)
		}

		expectedResult := fmt.Sprintf("Deleted node: %s", nodeFullName)
		if got := buf.String(); !strings.Contains(got, expectedResult) {
			t.Errorf("deleteTpuNode got %q, want %q", got, expectedResult)
		}
	})
}
