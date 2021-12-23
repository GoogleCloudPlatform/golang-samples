// Copyright 2021 Google LLC
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

// [START compute_template_create_from_instance]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

// createTemplateFromInstance creates a new instance template based on an existing instance.
// This new template specifies a different boot disk.
func createTemplateFromInstance(w io.Writer, projectID, instance, templateName string) error {
	// projectID := "your_project_id"
	// instance := "projects/project/zones/zone/instances/instance"
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
			Name:           proto.String(templateName),
			SourceInstance: proto.String(instance),
			SourceInstanceParams: &computepb.SourceInstanceParams{
				DiskConfigs: []*computepb.DiskInstantiationConfig{
					{
						// Device name must match the name of a disk attached to the instance
						// your template is based on.
						DeviceName: proto.String("disk-1"),
						// Replace the original boot disk image used in your instance with a Rocky Linux image.
						InstantiateFrom: proto.String(computepb.DiskInstantiationConfig_CUSTOM_IMAGE.String()),
						CustomImage:     proto.String("projects/rocky-linux-cloud/global/images/family/rocky-linux-8"),
						// Override the auto_delete setting.
						AutoDelete: proto.Bool(true),
					},
				},
			},
		},
	}

	op, err := instanceTemplatesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance template: %v", err)
	}

	globalOperationsClient, err := compute.NewGlobalOperationsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGlobalOperationsRESTClient: %v", err)
	}
	defer globalOperationsClient.Close()

	for {
		waitReq := &computepb.WaitGlobalOperationRequest{
			Operation: op.Proto().GetName(),
			Project:   projectID,
		}
		zoneOp, err := globalOperationsClient.Wait(ctx, waitReq)
		if err != nil {
			return fmt.Errorf("unable to wait for the operation: %v", err)
		}

		if *zoneOp.Status.Enum() == computepb.Operation_DONE {
			fmt.Fprintf(w, "Instance template created\n")
			break
		}
	}

	return nil
}

// [END compute_template_create_from_instance]
