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

func bulkInsertInstance(w io.Writer, projectID, zone string, template *computepb.InstanceTemplate, count int64, namePattern string, minCount *wrapperspb.Int64Value, labels map[string]string) ([]*computepb.Instance, error) {
	// projectID := "your_project_id"
	// zone := "us-central1-a"
	// template := &computepb.InstanceTemplate{...}  // Fill in template details
	// count := int64(5)
	// namePattern := "instance-####"
	// minCount := wrapperspb.Int64Value{Value: 3}  // Optional
	// labels := map[string]string{"key": "value"}  // Optional

	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer client.Close()

	if minCount == nil {
		minCount = &wrapperspb.Int64Value{Value: count}
	}

	bulkInsertResource := &computepb.BulkInsertInstanceResource{
		SourceInstanceTemplate: proto.String(template.GetSelfLink()),
		Count:                  proto.Int64(count),
		MinCount:               proto.Int64(minCount.GetValue()),
		NamePattern:            &namePattern,
	}

	if labels == nil {
		labels = make(map[string]string)
	}

	labels["bulk_batch"] = uuid.New().String()
	instanceProp := &computepb.InstanceProperties{
		Labels: labels,
	}
	bulkInsertResource.InstanceProperties = instanceProp

	bulkInsertRequest := &computepb.BulkInsertInstanceRequest{
		Project:                            projectID,
		Zone:                               zone,
		BulkInsertInstanceResourceResource: bulkInsertResource,
	}

	op, err := client.BulkInsert(ctx, bulkInsertRequest)
	if err != nil {
		return nil, fmt.Errorf("BulkInsert: %w", err)
	}

	// Wait for the operation to complete
	if err := op.Wait(ctx); err != nil {
		return nil, fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Bulk instance creation completed\n")

	// List the created instances
	listReq := &computepb.ListInstancesRequest{
		Project: projectID,
		Zone:    zone,
		Filter:  proto.String("labels.bulk_batch = " + labels["bulk_batch"]),
	}
	it := client.List(ctx, listReq)

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

func createFiveInstances(w io.Writer, projectID, zone, templateName, namePattern string) ([]*computepb.Instance, error) {
	// projectID := "your_project_id"
	// zone := "us-central1-a"
	// templateName := "your_instance_template_name"
	// namePattern := "instance-####"

	// 	   namePattern NOTE: The string pattern used for the names of the VMs. The pattern
	//     must contain one continuous sequence of placeholder hash characters (#)
	//     with each character corresponding to one digit of the generated instance
	//     name. Example: a namePattern of inst-#### generates instance names such
	//     as inst-0001 and inst-0002. If existing instances in the same project and
	//     zone have names that match the name pattern then the generated instance
	//     numbers start after the biggest existing number. For example, if there
	//     exists an instance with name inst-0050, then instance names generated
	//     using the pattern inst-#### begin with inst-0051. The name pattern
	//     placeholder #...# can contain up to 18 characters.

	template, err := getInstanceTemplate(projectID, templateName)
	if err != nil {
		return nil, fmt.Errorf("getInstanceTemplate: %w", err)
	}

	instances, err := bulkInsertInstance(w, projectID, zone, template, 5, namePattern, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("bulkInsertInstance: %w", err)
	}

	return instances, nil
}

func getInstanceTemplate(projectID, templateName string) (*computepb.InstanceTemplate, error) {
	// projectID := "your_project_id"
	// templateName := "your_instance_template_name"

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

	return client.Get(ctx, req)
}

// [END compute_instances_bulk_insert]
