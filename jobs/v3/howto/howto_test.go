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
	"os"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	talent "google.golang.org/api/jobs/v3"
)

var testCompany *talent.Company
var testJob *talent.Job

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Print("GOLANG_SAMPLES_PROJECT_ID is unset. Skipping.")
		return
	}

	companyToCreate := &talent.Company{
		ExternalId:  fmt.Sprintf("company-%s", uuid.Must(uuid.NewV4()).String()),
		DisplayName: "Google Sample",
	}
	var err error
	testCompany, err = createCompany(ioutil.Discard, tc.ProjectID, companyToCreate)
	if err != nil {
		log.Fatalf("createCompany: %v", err)
	}

	jobToCreate := &talent.Job{
		RequisitionId: fmt.Sprintf("job-%s", uuid.Must(uuid.NewV4()).String()),
		Title:         "SWE",
		CompanyName:   testCompany.Name,
		ApplicationInfo: &talent.ApplicationInfo{
			Uris: []string{"https://googlesample.com/career"},
		},
		Description:     "Design, devolop, test, deploy, maintain and improve software.",
		LanguageCode:    "en-US",
		PromotionValue:  2,
		EmploymentTypes: []string{"FULL_TIME"},
		Addresses:       []string{"Mountain View, CA"},
		CustomAttributes: map[string]talent.CustomAttribute{
			"someFieldString": {
				Filterable:   true,
				StringValues: []string{"someStrVal"},
			},
			"someFieldLong": {
				Filterable: true,
				LongValues: []int64{900},
			},
		},
		CompensationInfo: &talent.CompensationInfo{
			Entries: []*talent.CompensationEntry{
				{
					Type: "BASE",
					Unit: "HOURLY",
					Amount: &talent.Money{
						CurrencyCode: "USD",
						Units:        12,
					},
				},
			},
		},
	}
	testJob, err = createJob(ioutil.Discard, tc.ProjectID, jobToCreate)
	if err != nil {
		log.Fatalf("createJob: %v", err)
	}

	result := m.Run()

	if err := deleteJob(ioutil.Discard, testJob.Name); err != nil {
		log.Fatalf("deleteJob: %v", err)
	}

	if err := deleteCompany(ioutil.Discard, testCompany.Name); err != nil {
		log.Fatalf("deleteCompany: %v", err)
	}

	os.Exit(result)
}
