package gateway

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	container "cloud.google.com/go/container/apiv1"
	"cloud.google.com/go/container/apiv1/containerpb"
	"github.com/google/uuid"
)

var (
	zone        = "us-central1-a"
	region      = "us-central1"
	clusterName = fmt.Sprintf("cluster-%s", uuid.New().String()[:10])
)

func pollOperation(ctx context.Context, client *container.ClusterManagerClient, opName string) error {
	fmt.Printf("Polling operation: %s\n", opName)
	for {
		op, err := client.GetOperation(ctx, &containerpb.GetOperationRequest{Name: opName})
		if err != nil {
			return fmt.Errorf("failed to get operation %s: %v", opName, err)
		}
		fmt.Printf("Operation status: %v\n", op)

		if op.Status == containerpb.Operation_RUNNING {
			fmt.Println("Waiting 30 seconds before polling again...")
			time.Sleep(30 * time.Second)
			continue
		}

		if op.Status == containerpb.Operation_DONE {
			fmt.Println("Operation completed successfully.")
			return nil
		}

		return fmt.Errorf("operation failed with status %v", op.Status)
	}
}

func createCluster(projectID, location, clusterName string) error {
	ctx := context.Background()
	client, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create cluster manager client: %v", err)
	}
	defer client.Close()

	clusterLocation := fmt.Sprintf("projects/%s/locations/%s", projectID, location)
	clusterDef := &containerpb.Cluster{
		Name:             clusterName,
		InitialNodeCount: 1,
		Fleet: &containerpb.Fleet{
			Project: projectID,
		},
	}

	req := &containerpb.CreateClusterRequest{
		Parent:  clusterLocation,
		Cluster: clusterDef,
	}

	fmt.Printf("Creating cluster %s in %s...\n", clusterName, clusterLocation)
	fmt.Printf("cl %+v", clusterDef)
	resp, err := client.CreateCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create cluster: %v", err)
	}

	opIdentifier := fmt.Sprintf("%s/operations/%s", clusterLocation, resp.Name)
	return pollOperation(ctx, client, opIdentifier)
}

func deleteCluster(projectID, location, clusterName string) error {
	ctx := context.Background()
	client, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create cluster manager client: %v", err)
	}
	defer client.Close()

	clusterFullName := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, location, clusterName)

	fmt.Printf("Deleting cluster %s...\n", clusterFullName)
	_, err = client.DeleteCluster(ctx, &containerpb.DeleteClusterRequest{Name: clusterFullName})
	return err
}

func TestGetNamespace(t *testing.T) {
	projectid := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	// Setup cluster.
	if err := createCluster(projectid, zone, clusterName); err != nil {
		t.Fatalf("failed to create cluster: %v", err)
	}
	defer deleteCluster(projectid, zone, clusterName)

	membershipName := fmt.Sprintf("projects/%s/locations/%s/memberships/%s", projectid, region, clusterName)
	results, err := getNamespace(membershipName, region)
	if err != nil {
		t.Fatalf("getNamespace failed: %v", err)
	}

	if results == nil {
		t.Fatalf("getNamespace returned nil results")
	}

	if results.ObjectMeta.Name != "default" {
		t.Errorf("expected namespace name 'default', got '%s'", results.ObjectMeta.Name)
	}
}
