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

func teardown(t *testing.T, tc testutil.Context, clusterName, region string) {
	t.Helper()
	ctx := context.Background()

	clusterClient, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(fmt.Sprintf("%s-dataproc.googleapis.com:443", region)))

	req := &dataprocpb.DeleteClusterRequest{
		ProjectId:   tc.ProjectID,
		Region:      region,
		ClusterName: clusterName,
	}
	op, err := clusterClient.DeleteCluster(ctx, req)

	op.Wait(ctx)
	if err != nil {
		t.Errorf("Error deleting cluster %q: %v", clusterName, err)
	}
}

func TestCreateCluster(t *testing.T) {
	tc := testutil.SystemTest(t)

	clusterName := fmt.Sprintf("go-cc-test-%s", tc.ProjectID)
	region := "us-central1"

	defer teardown(t, tc, clusterName, region)

	buf := new(bytes.Buffer)

	if err := createCluster(buf, tc.ProjectID, region, clusterName); err != nil {
		t.Fatalf("createCluster got err: %v", err)
	}

	got := buf.String()
	if want := fmt.Sprintf("successfully: %s", clusterName); !strings.Contains(got, want) {
		t.Fatalf("CreateCluster: got %s, want %s", got, want)
	}
}
