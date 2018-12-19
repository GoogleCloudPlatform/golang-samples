// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains examples of using the monitoring API.
package snippets

// [START monitoring_get_descriptor]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// getMetricDescriptor gets the descriptor for the given metricType and prints
// information about it. metricType is the type of the metric, for example
// compute.googleapis.com/firewall/dropped_packets_count.
func getMetricDescriptor(w io.Writer, projectID, metricType string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return fmt.Errorf("NewMetricClient: %v", err)
	}
	req := &monitoringpb.GetMetricDescriptorRequest{
		Name: fmt.Sprintf("projects/%s/metricDescriptors/%s", projectID, metricType),
	}
	resp, err := c.GetMetricDescriptor(ctx, req)
	if err != nil {
		return fmt.Errorf("could not get custom metric: %v", err)
	}

	fmt.Fprintf(w, "Name: %v\n", resp.GetName())
	fmt.Fprintf(w, "Description: %v\n", resp.GetDescription())
	fmt.Fprintf(w, "Type: %v\n", resp.GetType())
	fmt.Fprintf(w, "Metric Kind: %v\n", resp.GetMetricKind())
	fmt.Fprintf(w, "Value Type: %v\n", resp.GetValueType())
	fmt.Fprintf(w, "Unit: %v\n", resp.GetUnit())
	fmt.Fprintf(w, "Labels:\n")
	for _, l := range resp.GetLabels() {
		fmt.Fprintf(w, "\t%s (%s) - %s", l.GetKey(), l.GetValueType(), l.GetDescription())
	}
	return nil
}

// [END monitoring_get_descriptor]
