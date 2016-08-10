// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
	"strings"
	"testing"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	storage "cloud.google.com/go/storage"
	dataproc "google.golang.org/api/dataproc/v1"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateCluster(t *testing.T) {
	testutil.SystemTest(t)
	dc := newDataprocClient(t)

	cluster := clusterConfig{project: *projectID, region: *region, name: *clusterName, zone: *zoneID}

	err := createCluster(dc, cluster)
	if err != nil {
		t.Fatalf("createCluster - got %v, want nil err", err)
	}
}

func TestWaitForCluster(t *testing.T) {
	testutil.SystemTest(t)
	dc := newDataprocClient(t)

	cluster := clusterConfig{project: *projectID, region: *region, name: *clusterName, zone: *zoneID}

	_, err := waitForCluster(dc, cluster)
	if err != nil {
		t.Fatalf("waitForCluster - got %v, want nil err", err)
	}
}

func TestGetClusterDetails(t *testing.T) {
	testutil.SystemTest(t)
	dc := newDataprocClient(t)

	configuration := clusterConfig{project: *projectID, region: *region, name: *clusterName, zone: *zoneID}
	cluster, err := getClusterDetails(dc, configuration)
	if err != nil {
		t.Fatalf("getClusterDetails - got %v, want nil err", err)
	}

	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}$")
	if r.MatchString(cluster.uuid) != true {
		t.Fatalf("updateClusterMetadata - cluster uuid")
	}
	if strings.Contains(cluster.bucket, "dataproc-") != true {
		t.Fatalf("updateClusterMetadata - invalid bucket metadata for cluster")
	}
}

func TestJobFunctionality(t *testing.T) {
	testutil.SystemTest(t)
	dc := newDataprocClient(t)
	sc := newStorageClient(t)

	configuration := clusterConfig{project: *projectID, region: *region, name: *clusterName, zone: *zoneID}
	cluster, err := getClusterDetails(dc, configuration)

	jobId, err := submitJob(dc, sc, *pysparkFile, *bucketName, cluster)
	if err != nil {
		t.Fatalf("submitJob - got %v, want nil err", err)
	}
	_, err = waitForJob(dc, jobId, cluster)
	if err != nil {
		t.Fatalf("waitForJob - got %v, want nil err", err)
	}

	output, err := getJobOutput(sc, jobId, cluster)
	if err != nil {
		t.Fatalf("getJobOutput - got %v, want nil err", err)
	}
	if strings.Contains(string(output), "['Hello,', 'dog', 'elephant', 'panther', 'world!']") != true {
		t.Fatalf("getJobOutput - unexpected job output")
	}

}

func TestDeleteCluster(t *testing.T) {
	testutil.SystemTest(t)
	dc := newDataprocClient(t)

	configuration := clusterConfig{project: *projectID, region: *region, name: *clusterName, zone: *zoneID}
	cluster, err := getClusterDetails(dc, configuration)

	_, err = deleteCluster(dc, cluster)
	if err != nil {
		t.Fatalf("deleteCluster - got %v, want nil err", err)
	}
}

func newDataprocClient(t *testing.T) *dataproc.Service {
	ctx := context.Background()
	hc, err := google.DefaultClient(ctx, dataproc.CloudPlatformScope)
	if err != nil {
		t.Fatalf("DefaultClient: %v", err)
	}
	client, err := dataproc.New(hc)
	if err != nil {
		t.Fatalf("dataproc.New: %v", err)
	}
	return client
}

func newStorageClient(t *testing.T) *storage.Client {
	ctx := context.Background()
	sc, err := storage.NewClient(ctx)

	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}

	return sc
}
