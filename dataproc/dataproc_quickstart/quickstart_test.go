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
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

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
	jobFPath    string
	jobFName    = "sum.py"
	code        = `import pyspark
sc = pyspark.SparkContext()
rdd = sc.parallelize((1,2,3,4,5))
sum = rdd.reduce(lambda x, y: x + y)`
	region = "us-central1"
)

func setup(t *testing.T, tc testutil.Context) {
	ctx := context.Background()

	clusterName = "golang-dataproc-quickstart-test-" + tc.ProjectID
	bktName = "golang-dataproc-quickstart-test-" + tc.ProjectID
	jobFPath = fmt.Sprintf("gs://%s/%s", bktName, jobFName)

	sc, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("Error creating storage client with error: %v", err)
	}

	bkt := sc.Bucket(bktName)
	if err := bkt.Create(ctx, tc.ProjectID, nil); err != nil {
		t.Errorf("Error creating bucket %q: %v", bktName, err)
	}

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
}

func teardown(t *testing.T, tc testutil.Context) {
	ctx := context.Background()

	sc, err := storage.NewClient(ctx)
	if err != nil {
		t.Errorf("Error creating storage client with error: %v", err)
	}

	if err := sc.Bucket(bktName).Object(jobFName).Delete(ctx); err != nil {
		t.Errorf("Error deleting object: %v", err)
	}

	if err := sc.Bucket(bktName).Delete(ctx); err != nil {
		t.Errorf("Error deleting bucket: %v", err)
	}

	ep := fmt.Sprintf("%s-dataproc.googleapis.com:443", region)
	client, err := dataproc.NewClusterControllerClient(ctx, option.WithEndpoint(ep))
	if err != nil {
		t.Errorf("Error creating the cluster client: %s", err)
	}

	lReq := &dataprocpb.ListClustersRequest{ProjectId: tc.ProjectID, Region: region}
	it := client.ListClusters(ctx, lReq)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("Error listing clusters: %v", err)
		}
		if resp.ClusterName == clusterName {
			dReq := &dataprocpb.DeleteClusterRequest{ProjectId: tc.ProjectID, Region: region, ClusterName: clusterName}
			op, err := client.DeleteCluster(ctx, dReq)

			op.Wait(ctx)
			if err != nil {
				t.Fatalf("Error deleting cluster %s: %s", clusterName, err)
			}
		}
	}
}

func TestQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)

	setup(t, tc)

	defer teardown(t, tc)

	buf := new(bytes.Buffer)

	if err := quickstart(buf, tc.ProjectID, region, clusterName, jobFPath); err != nil {
		t.Errorf("quickstart got err: %v", err)
	}

	got := buf.String()
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
