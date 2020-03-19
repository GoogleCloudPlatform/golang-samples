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

package dataproc

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/option"
	dataprocpb "google.golang.org/genproto/googleapis/cloud/dataproc/v1"
)

func deleteCluster(projectID string, clusterName, region string) error {
	ctx := context.Background()

	clusterClient, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(fmt.Sprintf("%s-dataproc.googleapis.com:443", region)))
	if err != nil {
		return fmt.Errorf("dataproc.NewClusterControllerClient: %v", err)
	}

	req := &dataprocpb.DeleteClusterRequest{
		ProjectId:   projectID,
		Region:      region,
		ClusterName: clusterName,
	}
	op, err := clusterClient.DeleteCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteCluster: %v", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("DeleteCluster.Wait: %v", err)
	}
	return nil
}

func TestCreateCluster(t *testing.T) {
	tc := testutil.SystemTest(t)

	clusterName := fmt.Sprintf("go-cc-test-%s", tc.ProjectID)
	region := "us-central1"

	deleteCluster(tc.ProjectID, clusterName, region) // Delete the cluster if it already exists, ignoring any errors.
	defer deleteCluster(tc.ProjectID, clusterName, region)

	buf := new(bytes.Buffer)

	if err := createCluster(buf, tc.ProjectID, region, clusterName); err != nil {
		t.Fatalf("createCluster got err: %v", err)
	}

	got := buf.String()
	if want := fmt.Sprintf("successfully: %s", clusterName); !strings.Contains(got, want) {
		t.Fatalf("CreateCluster: got %s, want %s", got, want)
	}
}
