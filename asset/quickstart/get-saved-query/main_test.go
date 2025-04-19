// Copyright 2022 Google LLC
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

package get

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

var (
	projectID     string
	savedQueryID  string
	projectNumber string

	ctx    context.Context
	client *asset.Client
)

func TestMain(m *testing.M) {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	savedQueryID = fmt.Sprintf("query-%s", strconv.FormatInt(time.Now().UnixNano(), 10))
	ctx = context.Background()
	var err error
	client, err = asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}

	cloudResourceManagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		log.Fatalf("cloudresourcemanager.NewService: %v", err)
	}

	project, err := cloudResourceManagerClient.Projects.Get(projectID).Do()
	if err != nil {
		log.Fatalf("cloudResourceManagerClient.Projects.Get.Do: %v", err)
	}
	projectNumber = strconv.FormatInt(project.ProjectNumber, 10)
	parent := fmt.Sprintf("projects/%s", projectID)
	log.Printf("projectNumber:%s", projectNumber)

	req := &assetpb.CreateSavedQueryRequest{
		Parent:       parent,
		SavedQueryId: savedQueryID,
		SavedQuery: &assetpb.SavedQuery{
			Content: &assetpb.SavedQuery_QueryContent{
				QueryContent: &assetpb.SavedQuery_QueryContent_IamPolicyAnalysisQuery{
					IamPolicyAnalysisQuery: &assetpb.IamPolicyAnalysisQuery{
						Scope: parent,
						AccessSelector: &assetpb.IamPolicyAnalysisQuery_AccessSelector{
							Permissions: []string{"iam.serviceAccount.actAs"},
						},
					},
				},
			},
		}}
	// TODO(#2811): uncomment when tests are fixed.
	_ = req
	// _, err = client.CreateSavedQuery(ctx, req)
	// if err != nil {
	// 	log.Fatalf("client.CreateSavedQuery: %v", err)
	// }

	os.Exit(m.Run())
}

func TestGetSavedQuery(t *testing.T) {
	buf := new(bytes.Buffer)
	get_err := getSavedQuery(buf, projectID, savedQueryID)

	fullQueryName := fmt.Sprintf("projects/%s/savedQueries/%s", projectNumber, savedQueryID)
	delete_error := client.DeleteSavedQuery(ctx, &assetpb.DeleteSavedQueryRequest{
		Name: fullQueryName,
	})

	if get_err != nil {
		t.Errorf("getSavedQuery: %v", get_err)
	}

	if delete_error != nil {
		t.Errorf("DeleteSavedQuery: %v", delete_error)
	}

	got := buf.String()
	if want := fullQueryName; !strings.Contains(got, want) {
		t.Errorf("getSavedQuery got%q, want%q", got, want)
	}
}
