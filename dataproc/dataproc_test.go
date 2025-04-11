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
	"time"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"cloud.google.com/go/dataproc/apiv1/dataprocpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

func deleteCluster(projectID string, clusterName, region string) error {
	ctx := context.Background()

	clusterClient, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(fmt.Sprintf("%s-dataproc.googleapis.com:443", region)))
	if err != nil {
		return fmt.Errorf("dataproc.NewClusterControllerClient: %w", err)
	}

	req := &dataprocpb.DeleteClusterRequest{
		ProjectId:   projectID,
		Region:      region,
		ClusterName: clusterName,
	}
	op, err := clusterClient.DeleteCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteCluster: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("DeleteCluster.Wait: %w", err)
	}
	return nil
}

func TestDataproc(t *testing.T) {
	t.Skip("skipping until https://github.com/GoogleCloudPlatform/golang-samples/issues/4350 is resolved.")
	tc := testutil.SystemTest(t)

	clusterName := fmt.Sprintf("go-dp-test-%s", uuid.New().String())
	region := "us-central1"

	deleteCluster(tc.ProjectID, clusterName, region) // Delete the cluster if it already exists, ignoring any errors.
	defer deleteCluster(tc.ProjectID, clusterName, region)

	buf := new(bytes.Buffer)

	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := createCluster(buf, tc.ProjectID, region, clusterName); err != nil {
			r.Errorf("createCluster got err: %v", err)
			return
		}

		got := buf.String()
		if want := fmt.Sprintf("successfully: %s", clusterName); !strings.Contains(got, want) {
			r.Errorf("CreateCluster: got %s, want %s", got, want)
			return
		}

		if err := submitJob(buf, tc.ProjectID, region, clusterName); err != nil {
			r.Errorf("submitJob got err: %v", err)
			return
		}

		got = buf.String()
		if want := fmt.Sprint("Job finished successfully"); !strings.Contains(got, want) {
			r.Errorf("submitJob: got %s, want %s", got, want)
			return
		}
	})
}
