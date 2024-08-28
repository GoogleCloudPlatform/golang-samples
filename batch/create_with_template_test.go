// Copyright 2022 Google LLC
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

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/protobuf/proto"
)

func TestCreateJobWithTemplate(t *testing.T) {
	t.Parallel()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	region := "us-central1"
	jobName := fmt.Sprintf("test-job-go-template-%v-%v", time.Now().Format("2006-01-02"), r.Int())
	templateName := fmt.Sprintf("test-template-go-batch-%v-%v", time.Now().Format("2006-01-02"), r.Int())
	buf := &bytes.Buffer{}

	if err := createTemplate(tc.ProjectID, templateName); err != nil {
		t.Errorf("Failed to create instance template: createTemplate got err: %v", err)
	}

	if err := createScriptJobWithTemplate(buf, tc.ProjectID, region, jobName, templateName); err != nil {
		t.Errorf("createScriptJobWithTemplate got err: %v", err)
	}

	succeeded, err := jobSucceeded(tc.ProjectID, region, jobName)
	if err != nil {
		t.Errorf("Could not verify job completion: %v", err)
	}
	if !succeeded {
		t.Errorf("The test job has failed: %v", err)
	}

	// clean up after the test
	if err := deleteInstanceTemplate(tc.ProjectID, templateName); err != nil {
		t.Errorf("Failed to delete instance template: deleteInstanceTemplate got err: %v", err)
	}
}

// createTemplate creates a new instance template with the provided name and a specific instance configuration.
// Includes all the setup needed for Batch specifically, such as service accounts.
func createTemplate(projectID, templateName string) error {
	ctx := context.Background()
	instanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstanceTemplatesRESTClient: %w", err)
	}
	defer instanceTemplatesClient.Close()

	projectNumber, err := projectIDtoNumber(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Could not resolve project ID '%s' to project number: %w", projectID, err)
	}

	serviceAccountAddress := fmt.Sprintf("%d-compute@developer.gserviceaccount.com", projectNumber)

	req := &computepb.InsertInstanceTemplateRequest{
		Project: projectID,
		InstanceTemplateResource: &computepb.InstanceTemplate{
			Name: proto.String(templateName),
			Properties: &computepb.InstanceProperties{
				// The template describes the size and source image of the boot disk
				// to attach to the instance.
				Disks: []*computepb.AttachedDisk{
					{
						InitializeParams: &computepb.AttachedDiskInitializeParams{
							DiskSizeGb:  proto.Int64(10),
							SourceImage: proto.String("projects/debian-cloud/global/images/family/debian-11"),
						},
						AutoDelete: proto.Bool(true),
						Boot:       proto.Bool(true),
					},
				},
				MachineType: proto.String("e2-standard-4"),
				// The template connects the instance to the `default` network,
				// without specifying a subnetwork.
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Name: proto.String("global/networks/default"),
						// The template lets the instance use an external IP address.
						AccessConfigs: []*computepb.AccessConfig{
							{
								Name:        proto.String("External NAT"),
								Type:        proto.String(computepb.AccessConfig_ONE_TO_ONE_NAT.String()),
								NetworkTier: proto.String(computepb.AccessConfig_PREMIUM.String()),
							},
						},
					},
				},
				ServiceAccounts: []*computepb.ServiceAccount{{
					Email: &serviceAccountAddress,
					Scopes: []string{
						"https://www.googleapis.com/auth/devstorage.read_only",
						"https://www.googleapis.com/auth/logging.write",
						"https://www.googleapis.com/auth/monitoring.write",
						"https://www.googleapis.com/auth/servicecontrol",
						"https://www.googleapis.com/auth/service.management.readonly",
						"https://www.googleapis.com/auth/trace.append",
					},
				}},
			},
		},
	}

	op, err := instanceTemplatesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance template: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	return nil
}

func projectIDtoNumber(ctx context.Context, projectID string) (int64, error) {
	// Resolve the project ID to project number
	resourceManagerClient, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return 0, fmt.Errorf("cloudresourcemanager.NewService: %w", err)
	}
	// resourceManagerClient doesn't have a Close() method
	projectsClient := cloudresourcemanager.NewProjectsService(resourceManagerClient)
	projectData, err := projectsClient.Get(projectID).Do()
	if err != nil {
		return 0, fmt.Errorf("Could not resolve project ID '%s' to project number: %w", projectID, err)
	}
	return projectData.ProjectNumber, nil
}

func deleteInstanceTemplate(projectID, templateName string) error {
	ctx := context.Background()
	instanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstanceTemplatesRESTClient: %w", err)
	}
	defer instanceTemplatesClient.Close()

	req := &computepb.DeleteInstanceTemplateRequest{
		Project:          projectID,
		InstanceTemplate: templateName,
	}

	op, err := instanceTemplatesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance template: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}
	return nil
}
