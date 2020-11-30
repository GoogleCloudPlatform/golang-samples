// Copyright 2020 Google LLC
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
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/option"
	dataprocpb "google.golang.org/genproto/googleapis/cloud/dataproc/v1"
)

func setupForSubmitJobTest(projectID, region, clusterName string) error {
	ctx := context.Background()

	clusterClient, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(fmt.Sprintf("%s-dataproc.googleapis.com:443", region)))
	if err != nil {
		return fmt.Errorf("dataproc.NewClusterControllerClient: %v", err)
	}

	req := &dataprocpb.CreateClusterRequest{
		ProjectId: projectID,
		Region:    region,
		Cluster: &dataprocpb.Cluster{
			ProjectId:   projectID,
			ClusterName: clusterName,
			Config: &dataprocpb.ClusterConfig{
				MasterConfig: &dataprocpb.InstanceGroupConfig{
					NumInstances:   1,
					MachineTypeUri: "n1-standard-1",
				},
				WorkerConfig: &dataprocpb.InstanceGroupConfig{
					NumInstances:   2,
					MachineTypeUri: "n1-standard-1",
				},
			},
		},
	}

	// Delete the cluster if it already exists, ignoring any errors.
	deleteClusterForSubmitJobTest(projectID, region, clusterName)

	// Create the cluster.
	op, err := clusterClient.CreateCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateCluster: %v", err)
	}

	_, err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("CreateCluster.Wait: %v", err)
	}
	return nil
}

func deleteClusterForSubmitJobTest(projectID, region, clusterName string) error {
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

func TestSubmitJob(t *testing.T) {
	tc := testutil.SystemTest(t)

	clusterName := fmt.Sprintf("go-sj-test-%s", tc.ProjectID)
	region := "us-central1"

	setupForSubmitJobTest(tc.ProjectID, region, clusterName)
	defer deleteClusterForSubmitJobTest(tc.ProjectID, region, clusterName)

	buf := new(bytes.Buffer)

	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {
		if err := submitJob(buf, tc.ProjectID, region, clusterName); err != nil {
			r.Errorf("submitJob got err: %v", err)
			return
		}

		got := buf.String()
		if want := fmt.Sprint("Job finished successfully"); !strings.Contains(got, want) {
			r.Errorf("submitJob: got %s, want %s", got, want)
			return
		}
	})
}
