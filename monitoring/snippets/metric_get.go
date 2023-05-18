// Copyright 2019 Google LLC
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

// Package snippets contains examples of using the monitoring API.
package snippets

// [START monitoring_get_descriptor]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
)

// getMetricDescriptor gets the descriptor for the given metricType and prints
// information about it. metricType is the type of the metric, for example
// compute.googleapis.com/firewall/dropped_packets_count.
func getMetricDescriptor(w io.Writer, projectID, metricType string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return fmt.Errorf("NewMetricClient: %w", err)
	}
	defer c.Close()
	req := &monitoringpb.GetMetricDescriptorRequest{
		Name: fmt.Sprintf("projects/%s/metricDescriptors/%s", projectID, metricType),
	}
	resp, err := c.GetMetricDescriptor(ctx, req)
	if err != nil {
		return fmt.Errorf("could not get custom metric: %w", err)
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
