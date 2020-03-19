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
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dataprocpb "google.golang.org/genproto/googleapis/cloud/dataproc/v1"
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

func cleanBucket(ctx context.Context, t *testing.T, client *storage.Client, projectID, bucket string) {
	b := client.Bucket(bucket)
	_, err := b.Attrs(ctx)
	if err == nil {
		it := b.Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("Bucket.Objects(%q): %v", bucket, err)
			}
			if attrs.EventBasedHold || attrs.TemporaryHold {
				if _, err := b.Object(attrs.Name).Update(ctx, storage.ObjectAttrsToUpdate{
					TemporaryHold:  false,
					EventBasedHold: false,
				}); err != nil {
					t.Fatalf("Bucket(%q).Object(%q).Update: %v", bucket, attrs.Name, err)
				}
			}
			if err := b.Object(attrs.Name).Delete(ctx); err != nil {
				t.Fatalf("Bucket(%q).Object(%q).Delete: %v", bucket, attrs.Name, err)
			}
		}
		if err := b.Delete(ctx); err != nil {
			t.Fatalf("Bucket.Delete(%q): %v", bucket, err)
		}
	}
	if err := b.Create(ctx, projectID, nil); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucket, err)
	}
}

func setup(t *testing.T, projectID string) {
	ctx := context.Background()
	flag.Parse()

	clusterName = "go-qs-test-" + projectID
	bktName = "go-dataproc-qs-test-" + projectID
	jobFilePath = fmt.Sprintf("gs://%s/%s", bktName, jobFName)

	sc, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("storage.NewClient: %v", err)
	}

	cleanBucket(ctx, t, sc, projectID, bktName)
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

	deleteClusters(ctx, projectID) // Ignore any errors.
}

func teardown(t *testing.T, projectID string) {
	ctx := context.Background()

	sc, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("storage.NewClient: %v", err)
	}

	if err := sc.Bucket(bktName).Object(jobFName).Delete(ctx); err != nil {
		t.Errorf("Error deleting object: %v", err)
	}

	if err := sc.Bucket(bktName).Delete(ctx); err != nil {
		t.Errorf("Error deleting bucket: %v", err)
	}

	if err := deleteClusters(ctx, projectID); err != nil {
		t.Errorf("deleteClusters: %v", err)
	}
}

func deleteClusters(ctx context.Context, projectID string) error {
	endpoint := fmt.Sprintf("%s-dataproc.googleapis.com:443", region)
	client, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("dataproc.NewClusterControllerClient: %v", err)
	}

	lReq := &dataprocpb.ListClustersRequest{ProjectId: projectID, Region: region}
	it := client.ListClusters(ctx, lReq)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListClusters.Next: %v", err)
		}
		if resp.ClusterName == clusterName {
			dReq := &dataprocpb.DeleteClusterRequest{ProjectId: projectID, Region: region, ClusterName: clusterName}
			op, err := client.DeleteCluster(ctx, dReq)
			if err != nil {
				return fmt.Errorf("DeleteCluster: %v", err)
			}

			if err := op.Wait(ctx); err != nil {
				return fmt.Errorf("DeleteCluster.Wait: %v", err)
			}
		}
	}
	return nil
}

func TestQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)
	setup(t, tc.ProjectID)
	defer teardown(t, tc.ProjectID)

	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	stdOut, stdErr, err := m.Run(nil, 10*time.Minute,
		"--project_id", tc.ProjectID,
		"--region", region,
		"--cluster_name", clusterName,
		"--job_file_path", jobFilePath,
	)
	if err != nil {
		t.Errorf("stdout: %v", string(stdOut))
		t.Errorf("stderr: %v", string(stdErr))
		t.Errorf("execution failed: %v", err)
	}

	got := string(stdOut)
	wants := []string{
		"Cluster created successfully",
		"Submitted job",
		"finished with state DONE:",
		"successfully deleted",
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("got %q, want to contain %q", got, want)
		}
	}
}
