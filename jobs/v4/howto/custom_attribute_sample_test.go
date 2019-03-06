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

package howto

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateJobWithCustomAttributes(t *testing.T) {
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		tc := testutil.SystemTest(t)
		customJob, err := createJobWithCustomAttributes(ioutil.Discard, tc.ProjectID, testJob.CompanyName, testJob.Title)
		if err != nil {
			log.Fatalf("createJob: %v", err)
		}
		if err := deleteJob(ioutil.Discard, customJob.Name); err != nil {
			log.Fatalf("deleteJob: %v", err)
		}
	})
}
