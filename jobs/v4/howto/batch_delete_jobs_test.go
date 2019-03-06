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
	"testing"
  "time"

  // talent "cloud.google.com/go/talent/apiv4beta1"
	"github.com/gofrs/uuid"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)


func TestBatchDeleteJobs(t *testing.T) {
	testutil.Retry(t, 10, 1*time.Second, func(r *testutil.R) {
		tc := testutil.SystemTest(t)

    companyToCreate := &talentpb.Company{
      ExternalId:  fmt.Sprintf("company-%s", uuid.Must(uuid.NewV4()).String()),
      DisplayName: "Google Subsidiary Sample",
    }
    var err error
    tempCompany, err := createCompany(ioutil.Discard, tc.ProjectID, companyToCreate)
    if err != nil {
      log.Fatalf("createCompany: %v", err)
    }
    requisitionId := fmt.Sprintf("job-%s", uuid.Must(uuid.NewV4()).String())
    job1ToCreate := &talentpb.Job{
      RequisitionId: requisitionId,
      Title:         "Software Engineer",
      CompanyName:   tempCompany.Name,
      Description:   "Design, devolop, deploy, maintain and improve software.",
      ApplicationInfo: &talentpb.Job_ApplicationInfo{
  			Uris: []string{"https://googlesample.com/career"},
  		},
      LanguageCode: "en_US",
    }
    job2ToCreate := &talentpb.Job{
      RequisitionId: requisitionId,
      Title:         "Test Engineer",
      CompanyName:   tempCompany.Name,
      Description:   "Test software.",
      ApplicationInfo: &talentpb.Job_ApplicationInfo{
  			Uris: []string{"https://googlesample.com/career"},
  		},
      LanguageCode: "sr_Latn", 
    }

    if _, err := createJob(ioutil.Discard, tc.ProjectID, job1ToCreate); err != nil {
      log.Fatalf("createJob1: %v", err)
    }

    if _, err := createJob(ioutil.Discard, tc.ProjectID, job2ToCreate); err != nil {
      log.Fatalf("createJob2: %v", err)
    }


    // if err := listJobs(buf, tc.ProjectID, fmt.Sprintf("companyName=%q", testCompany.Name)); err != nil {
  	// 	t.Fatalf("listJobs: %v", err)
  	// }


    if err := batchDeleteJobs(ioutil.Discard, tc.ProjectID, fmt.Sprintf("companyName=%q AND requisitionId=%q", tempCompany.Name, requisitionId)); err != nil {
      log.Fatalf("batchDeleteJob: %v", err)
    }

    if err := deleteCompany(ioutil.Discard, tempCompany.Name); err != nil {
  		log.Fatalf("deleteCompany: %v", err)
  	}

	})
}
