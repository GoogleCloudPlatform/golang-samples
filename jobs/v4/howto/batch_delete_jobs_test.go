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
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
)

func TestBatchDeleteJobs(t *testing.T) {
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		tc := testutil.SystemTest(t)

		companyID := strings.SplitAfter(testCompany.Name, "companies/")[1]
		requisitionID := fmt.Sprintf("job-%s", uuid.Must(uuid.NewV4()).String())
		title := "Software Engineer"
		URI := "https://googlesample.com/career"
		description := "Design, devolop, test, deploy, maintain and improve software."
		address1 := "Mountain View, CA"
		address2 := "New York City, NY"
		languageCode1 := "en-US"
		languageCode2 := "sr_Latn"

		// Create two identical jobs with different language codes.
		if _, err := createJob(ioutil.Discard, tc.ProjectID, companyID, requisitionID, title, URI, description, address1, address2, languageCode1); err != nil {
			log.Fatalf("createJob1: %v", err)
		}
		if _, err := createJob(ioutil.Discard, tc.ProjectID, companyID, requisitionID, title, URI, description, address1, address2, languageCode2); err != nil {
			log.Fatalf("createJob2: %v", err)
		}

		if err := batchDeleteJobs(ioutil.Discard, tc.ProjectID, companyID, requisitionID); err != nil {
			log.Fatalf("batchDeleteJob: %v", err)
		}
	})
}
