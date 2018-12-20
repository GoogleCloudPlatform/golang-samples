// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains examples of using the monitoring API.
package snippets

// [START monitoring_get_resource]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// getMonitoredResource gets the descriptor for the given resourceType and
// prints information about it. resource should be of the form
// "projects/[PROJECT_ID]/monitoredResourceDescriptors/[RESOURCE_TYPE]".
func getMonitoredResource(w io.Writer, resource string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return fmt.Errorf("NewMetricClient: %v", err)
	}
	req := &monitoringpb.GetMonitoredResourceDescriptorRequest{
		Name: fmt.Sprintf(resource),
	}
	resp, err := c.GetMonitoredResourceDescriptor(ctx, req)
	if err != nil {
		return fmt.Errorf("could not get custom metric: %v", err)
	}

	fmt.Fprintf(w, "Name: %v\n", resp.GetName())
	fmt.Fprintf(w, "Description: %v\n", resp.GetDescription())
	fmt.Fprintf(w, "Type: %v\n", resp.GetType())
	fmt.Fprintf(w, "Labels:\n")
	for _, l := range resp.GetLabels() {
		fmt.Fprintf(w, "\t%s (%s) - %s", l.GetKey(), l.GetValueType(), l.GetDescription())
	}
	return nil
}

// [END monitoring_get_resource]
