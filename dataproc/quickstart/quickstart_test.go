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

package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	dataproc "cloud.google.com/go/dataproc/apiv1"
	"cloud.google.com/go/dataproc/apiv1/dataprocpb"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var (
	clusterName string
	bktName     string
	jobFilePath string
	jobFName    = "sum.py"
	code        = `import pyspark
sc = pyspark.SparkContext()
rdd = sc.parallelize((1,2,3,4,5))
sum = rdd.reduce(lambda x, y: x + y)`
	region = "us-central1"
)

func setup(t *testing.T, projectID string) {
	ctx := context.Background()
	flag.Parse()

	uuid := uuid.New().String()

	clusterName = "go-qs-test-" + uuid
	jobFilePath = fmt.Sprintf("gs://%s/%s", bktName, jobFName)

	sc, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { sc.Close() })

	bktName = testutil.CreateTestBucket(ctx, t, sc, projectID, "go-dataproc-qs-test")
	bkt := sc.Bucket(bktName)

	obj := bkt.Object(jobFName)

	w := obj.NewWriter(ctx)

	if _, err := fmt.Fprintf(w, code); err != nil {
		if err2 := w.Close(); err != nil {
			t.Errorf("Error writing to file and closing it: %v", err2)
		}
		t.Errorf("Error writing to file: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Errorf("Error closing file: %v", err)
	}

	// Opportunistically delete colliding cluster name.  Ignore errors.
	deleteCluster(ctx, projectID, region, clusterName)
}

func teardown(t *testing.T, projectID string) {
	ctx := context.Background()

	// Post-hoc cleanup, ignore errors.
	deleteCluster(ctx, projectID, region, clusterName)
}

func deleteCluster(ctx context.Context, projectID, region, clusterName string) error {
	endpoint := fmt.Sprintf("%s-dataproc.googleapis.com:443", region)
	client, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("dataproc.NewClusterControllerClient: %w", err)
	}

	dReq := &dataprocpb.DeleteClusterRequest{ProjectId: projectID, Region: region, ClusterName: clusterName}
	op, err := client.DeleteCluster(ctx, dReq)
	if err != nil {
		return fmt.Errorf("DeleteCluster: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("DeleteCluster.Wait: %w", err)
	}
	return nil
}

func TestQuickstart(t *testing.T) {
	t.Skip("Skipped until https://github.com/GoogleCloudPlatform/golang-samples/issues/4350 is resolved.")
	tc := testutil.EndToEndTest(t)
	m := testutil.BuildMain(t)
	setup(t, tc.ProjectID)
	defer teardown(t, tc.ProjectID)

	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	testutil.Retry(t, 3, 30*time.Second, func(r *testutil.R) {

		stdOut, stdErr, err := m.Run(nil, 10*time.Minute,
			"--project_id", tc.ProjectID,
			"--region", region,
			"--cluster_name", clusterName,
			"--job_file_path", jobFilePath,
		)
		if err != nil {
			r.Errorf("stdout: %v", string(stdOut))
			r.Errorf("stderr: %v", string(stdErr))
			r.Errorf("execution failed: %v", err)
			// We may have created the cluster in the failed invocation; try deleting.
			deleteCluster(context.Background(), tc.ProjectID, region, clusterName)
			return
		}

		got := string(stdOut)
		wants := []string{
			"Cluster created successfully",
			"Job finished successfully",
			"successfully deleted",
		}
		for _, want := range wants {
			if !strings.Contains(got, want) {
				r.Errorf("got %q, want to contain %q", got, want)
			}
		}
	})
}
