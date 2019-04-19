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
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

var testCompany *talentpb.Company
var testJob *talentpb.Job

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Fatal("Error getting test context")
	}

	externalID := fmt.Sprintf("company-%s", uuid.Must(uuid.NewV4()).String())
	displayName := "Google Sample"
	var err error
	testCompany, err = createCompany(ioutil.Discard, tc.ProjectID, externalID, displayName)
	if err != nil {
		log.Fatalf("createCompany: %v", err)
	}

	companyID := strings.SplitAfter(testCompany.Name, "companies/")[1]
	requisitionID := fmt.Sprintf("job-%s", uuid.Must(uuid.NewV4()).String())
	title := "Software Engineer"
	URI := "https://googlesample.com/career"
	description := "Design, devolop, test, deploy, maintain and improve software."
	address1 := "1600 Amphitheatre PkwyMountain View, CA 94043"
	address2 := "85 10th Ave, New York, NY 10011"
	languageCode := "en-US"

	testJob, err = createJob(ioutil.Discard, tc.ProjectID, companyID, requisitionID, title, URI, description, address1, address2, languageCode)
	if err != nil {
		log.Fatalf("createJob: %v", err)
	}
	jobID := strings.SplitAfter(testJob.Name, "jobs/")[1]

	result := m.Run()

	if err := deleteJob(ioutil.Discard, tc.ProjectID, jobID); err != nil {
		log.Fatalf("deleteJob: %v", err)
	}

	if err := deleteCompany(ioutil.Discard, tc.ProjectID, companyID); err != nil {
		log.Fatalf("deleteCompany: %v", err)
	}

	os.Exit(result)
}
