// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains examples of using the monitoring API.
package snippets

// [START monitoring_list_resources]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"google.golang.org/api/iterator"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// listMonitoredResources lists all the resources available to be monitored.
func listMonitoredResources(w io.Writer, projectID string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}

	req := &monitoringpb.ListMonitoredResourceDescriptorsRequest{
		Name: "projects/" + projectID,
	}
	iter := c.ListMonitoredResourceDescriptors(ctx, req)

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Could not list time series: %v", err)
		}
		fmt.Fprintf(w, "%v\n", resp)
	}
	fmt.Fprintln(w, "Done")
	return nil
}

// [END monitoring_list_resources]
