//  Copyright 2024 Google LLC
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func deleteInstance(ctx context.Context, projectId, zone, instanceName string) error {
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	req := &computepb.DeleteInstanceRequest{
		Project:  projectId,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

// getInstance fetches the instance details by name
func getInstance(ctx context.Context, projectID, zone, instanceName string) (*computepb.Instance, error) {
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	reqInstance := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, reqInstance)
	if err != nil {
		return nil, fmt.Errorf("unable to get instance: %w", err)
	}

	fmt.Printf("Instance: %s\n", instance.GetName())

	return instance, nil
}

func TestIsSpotVM(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "spot-instance-name" + fmt.Sprint(seededRand.Int())
	buf := &bytes.Buffer{}

	// Initiate instance
	err := createSpotInstance(buf, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Fatalf("Failed to create spot instance: %v", err)
	}

	// cleanup
	defer func() {
		if err := deleteInstance(ctx, tc.ProjectID, zone, instanceName); err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}

	}()

	buf.Reset()

	// Check if the instance exists
	instance, err := getInstance(ctx, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("Failed to get instance: %v", err)
	}

	if instance == nil {
		t.Errorf("Instance %q does not exist", instanceName)
	}

	buf.Reset()

	// Verify Spot VM status
	isSpot, err := isSpotVM(buf, tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Fatalf("isSpotVM got err: %v", err)
	}

	want := fmt.Sprintf("Instance %s is spot", instanceName)
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("isSpotVM got %q, want %q to be included", got, want)
	}

	if !isSpot {
		t.Errorf("expected instance to be a Spot VM, but it was not")
	}
}
