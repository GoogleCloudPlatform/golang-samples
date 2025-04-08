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

// [START monitoring_list_resources]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/api/iterator"
)

// listMonitoredResources lists all the resources available to be monitored.
func listMonitoredResources(w io.Writer, projectID string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

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
			return fmt.Errorf("Could not list time series: %w", err)
		}
		fmt.Fprintf(w, "%v\n", resp)
	}
	fmt.Fprintln(w, "Done")
	return nil
}

// [END monitoring_list_resources]
