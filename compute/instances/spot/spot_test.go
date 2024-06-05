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
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
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

func listInstances(w io.Writer, projectID, zone string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
	}

	it := instancesClient.List(ctx, req)
	fmt.Fprintf(w, "Instances found in zone %s:\n", zone)
	for {
		instance, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "- %s %s\n", instance.GetName(), instance.GetMachineType())
	}
	return nil
}

func TestSpotCreate(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-central2-b"
	instanceName := "spot-instance-name" + fmt.Sprint(seededRand.Int())
	buf := &bytes.Buffer{}

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

	buf.Reset() // to check list output

	err = listInstances(buf, tc.ProjectID, zone)
	if err != nil {
		t.Error(err)
		return
	}

	if got := buf.String(); !strings.Contains(got, instanceName) {
		t.Errorf("listInstances got %q, want %q to be contained", got, instanceName)
	}
}
