// Copyright 2024 Google LLC
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

package snippets

// [START compute_instances_bulk_insert]

import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func createFiveInstances(w io.Writer, projectID, zone, templateName, namePattern string) ([]*computepb.Instance, error) {
	// projectID := "your_project_id"
	// zone := "us-central1-a"
	// templateName := "your_instance_template_name"
	// namePattern := "instance-####"

	// Get instance template
	ctx := context.Background()
	client, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewInstanceTemplatesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.GetInstanceTemplateRequest{
		Project:          projectID,
		InstanceTemplate: templateName,
	}

	template, err := client.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("getInstanceTemplate: %w", err)
	}

	// Initialize Instances REST Client
	clientInstances, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer clientInstances.Close()

	// Prepare Bulk Insert Request
	minCount := &wrapperspb.Int64Value{Value: 5}
	bulkInsertResource := &computepb.BulkInsertInstanceResource{
		SourceInstanceTemplate: proto.String(template.GetSelfLink()),
		Count:                  proto.Int64(5),
		MinCount:               proto.Int64(minCount.GetValue()),
		NamePattern:            &namePattern,
	}

	labels := map[string]string{
		"bulk_batch": uuid.New().String(),
	}
	instanceProp := &computepb.InstanceProperties{
		Labels: labels,
	}
	bulkInsertResource.InstanceProperties = instanceProp

	bulkInsertRequest := &computepb.BulkInsertInstanceRequest{
		Project:                            projectID,
		Zone:                               zone,
		BulkInsertInstanceResourceResource: bulkInsertResource,
	}

	// Perform Bulk Insert
	op, err := clientInstances.BulkInsert(ctx, bulkInsertRequest)
	if err != nil {
		return nil, fmt.Errorf("BulkInsert: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return nil, fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Bulk instance creation completed\n")

	// Fetch the created instances
	listReq := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
		Filter:  proto.String("labels.bulk_batch = " + labels["bulk_batch"]),
	}
	it := clientInstances.List(ctx, listReq)

	var instances []*computepb.Instance
	for {
		instance, err := it.Next()
		if err == context.Canceled || err == context.DeadlineExceeded {
			return nil, err
		}
		if err != nil {
			break
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// [END compute_instances_bulk_insert]
