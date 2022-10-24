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
	"io"
	"math/rand"
	//"strings"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

func TestCreateJobWithTemplate(t *testing.T) {
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	region := "us-central1"
	jobName := fmt.Sprintf("test-job-go-template-%v-%v", time.Now().Format("2006-12-25"), r.Int())
	templateName := fmt.Sprintf("test-template-go-batch-%v-%v", time.Now().Format("2006-12-25"), r.Int())
	buf := &bytes.Buffer{}

	if err := createTemplate(buf, tc.ProjectID, templateName); err != nil {
		t.Errorf("Failed to create instance template: createTemplate got err: %v", err)
	}
	// TODO: defer deletion

	buf.Reset()

	if err := createScriptJobWithTemplate(buf, tc.ProjectID, region, jobName, templateName); err != nil {
		t.Errorf("createScriptJobWithTemplate got err: %v", err)
	}
}

// Copied from compute_template_create snippet
// createTemplate creates a new instance template with the provided name and a specific instance configuration.
func createTemplate(w io.Writer, projectID, templateName string) error {
	// projectID := "your_project_id"
	// templateName := "your_template_name"

	ctx := context.Background()
	instanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstanceTemplatesRESTClient: %v", err)
	}
	defer instanceTemplatesClient.Close()

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
			},
		},
	}

	op, err := instanceTemplatesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance template: %v", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %v", err)
	}

	fmt.Fprintf(w, "Instance template created\n")

	return nil
}
