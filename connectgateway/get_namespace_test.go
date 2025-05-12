// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package connectgateway

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	gke "google.golang.org/api/container/v1"
)

var (
	zone        = "us-central1-a"
	region      = "us-central1"
	clusterName = fmt.Sprintf("cluster-%s", uuid.New().String()[:10])
)

func TestGetNamespace(t *testing.T) {
	ctx := context.Background()
	tc := testutil.EndToEndTest(t)
	// Setup cluster.
	if err := createCluster(ctx, tc.ProjectID, zone, clusterName); err != nil {
		t.Fatalf("failed to create cluster: %v", err)
	}
	defer deleteCluster(ctx, tc.ProjectID, zone, clusterName)

	membershipName := fmt.Sprintf("projects/%s/locations/%s/memberships/%s", tc.ProjectID, region, clusterName)
	var buf bytes.Buffer
	err := getNamespace(&buf, membershipName, region)
	if err != nil {
		t.Fatalf("getNamespace failed: %v", err)
	}

	got := buf.String()
	if want := "Name:\"default\""; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func createCluster(ctx context.Context, projectID, location, clusterName string) error {
	svc, err := gke.NewService(ctx)
	if err != nil {
		log.Fatalf("Could not initialize gke client: %v", err)
	}
	clusterLocation := fmt.Sprintf("projects/%s/locations/%s", projectID, location)

	req := &gke.CreateClusterRequest{
		Parent: clusterLocation,
		Cluster: &gke.Cluster{
			Name:             clusterName,
			InitialNodeCount: 1,
			Fleet: &gke.Fleet{
				Project: projectID,
			},
		},
	}

	fmt.Printf("Creating cluster %s in %s...\n", clusterName, clusterLocation)
	resp, err := svc.Projects.Zones.Clusters.Create(projectID, location, req).Do()
	if err != nil {
		return fmt.Errorf("failed to create cluster: %v", err.Error())
	}

	return pollOperation(svc, projectID, resp.Name)
}

func pollOperation(svc *gke.Service, projectId, opID string) error {
	fmt.Printf("Polling operation: %s\n", opID)
	for {

		op, err := svc.Projects.Zones.Operations.Get(projectId, zone, opID).Do()
		if err != nil {
			return fmt.Errorf("failed to get operation %s: %v", opID, err)
		}
		fmt.Printf("Operation status: %v\n", op)

		if op.Status == "RUNNING" {
			fmt.Println("Waiting 30 seconds before polling again...")
			time.Sleep(30 * time.Second)
			continue
		}

		if op.Status == "DONE" {
			fmt.Println("Operation completed successfully.")
			return nil
		}

		return fmt.Errorf("operation failed with status %v", op.Status)
	}
}

func deleteCluster(ctx context.Context, projectID, location, clusterName string) error {
	svc, err := gke.NewService(ctx)
	if err != nil {
		log.Fatalf("Could not initialize gke client: %v", err)
	}
	clusterFullName := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, location, clusterName)

	fmt.Printf("Deleting cluster %s...\n", clusterFullName)
	_, err = svc.Projects.Zones.Clusters.Delete(projectID, zone, clusterName).Do()
	if err != nil {
		return fmt.Errorf("failed to delete cluster %v: %v", clusterName, err)
	}
	return nil
}
