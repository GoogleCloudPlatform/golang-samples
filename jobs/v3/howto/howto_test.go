// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
		log.Fatal("Error getting test context")
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
