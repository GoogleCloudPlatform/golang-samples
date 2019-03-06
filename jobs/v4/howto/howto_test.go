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
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
	money "google.golang.org/genproto/googleapis/type/money"
)

var testCompany *talentpb.Company
var testJob *talentpb.Job

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Fatal("Error getting test context")
	}

	companyToCreate := &talentpb.Company{
		ExternalId:  fmt.Sprintf("company-%s", uuid.Must(uuid.NewV4()).String()),
		DisplayName: "Google Sample",
	}
	var err error
	testCompany, err = createCompany(ioutil.Discard, tc.ProjectID, companyToCreate)
	if err != nil {
		log.Fatalf("createCompany: %v", err)
	}

	jobToCreate := &talentpb.Job{
		RequisitionId: fmt.Sprintf("job-%s", uuid.Must(uuid.NewV4()).String()),
		Title:         "Software Engineer",
		CompanyName:   testCompany.Name,
		ApplicationInfo: &talentpb.Job_ApplicationInfo{
			Uris: []string{"https://googlesample.com/career"},
		},
		Description:     "Design, devolop, test, deploy, maintain and improve software.",
		LanguageCode:    "en-US",
		PromotionValue:  2,
		EmploymentTypes: []talentpb.EmploymentType{1},
		Addresses:       []string{"Mountain View, CA"},
		CustomAttributes: map[string]*talentpb.CustomAttribute{
			"someFieldString": {
				Filterable:   true,
				StringValues: []string{"someStrVal"},
			},
			"someFieldLong": {
				Filterable: true,
				LongValues: []int64{900},
			},
			"anotherFieldLong": {
				Filterable: true,
				LongValues: []int64{900},
			},
		},
		CompensationInfo: &talentpb.CompensationInfo{
			Entries: []*talentpb.CompensationInfo_CompensationEntry{
				{
					Type: 1,
					Unit: 1,
					CompensationAmount: &talentpb.CompensationInfo_CompensationEntry_Amount{
						&money.Money{
							CurrencyCode: "USD",
							Units:        12,
						},
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

	if err := deleteCompany(ioutil.Discard, testCompany.GetName()); err != nil {
		log.Fatalf("deleteCompany: %v", err)
	}

	os.Exit(result)
}
