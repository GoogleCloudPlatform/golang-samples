// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Prompt: You are a Go programmer that knows Google Cloud. Write a test for importDataImageClassification that is similar to TestCreateDataset.

package snippets

import (
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestImportDataImageClassification(t *testing.T) {
	tc := testutil.SystemTest(t)
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := importDataImageClassification(tc.ProjectID, tc.DatasetID, tc.Location); err != nil {
			r.Errorf("importDataImageClassification got err: %v", err)
		}
	})
}
