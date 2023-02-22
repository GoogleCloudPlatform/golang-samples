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

package create

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
	log.Printf("projectNumber:%s", projectNumber)
	os.Exit(m.Run())
}

func TestCreateSavedQuery(t *testing.T) {
	t.Skip("Skipped while investigating https://github.com/GoogleCloudPlatform/golang-samples/issues/2811")
	buf := new(bytes.Buffer)
	err := createSavedQuery(buf, projectID, savedQueryID)
	if err != nil {
		t.Errorf("createSavedQuery: %v", err)
	}
	fullQueryName := fmt.Sprintf("projects/%s/savedQueries/%s", projectNumber, savedQueryID)
	delete_err := client.DeleteSavedQuery(ctx, &assetpb.DeleteSavedQueryRequest{
		Name: fullQueryName,
	})
	if delete_err != nil {
		t.Errorf("DeleteSavedQuery: %v", delete_err)
	}
	got := buf.String()
	if want := fullQueryName; !strings.Contains(got, want) {
		t.Errorf("createSavedQuery got%q, want%q", got, want)
	}
}
