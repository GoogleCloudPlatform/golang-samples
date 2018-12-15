// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains examples of using the monitoring API.
package snippets

// [START monitoring_read_timeseries_align]
import (
	"context"
	"fmt"
	"io"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// readTimeSeriesAlign reads the last 20 minutes of the given metric and aligns
// everything on 10 minute intervals.
func readTimeSeriesAlign(w io.Writer, projectID string) error {
	ctx := context.Background()
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return fmt.Errorf("NewMetricClient: %v", err)
	}
	startTime := time.Now().UTC().Add(time.Minute * -20)
	endTime := time.Now().UTC()
	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   "projects/" + projectID,
		Filter: `metric.type="compute.googleapis.com/instance/cpu/utilization"`,
		Interval: &monitoringpb.TimeInterval{
			StartTime: &timestamp.Timestamp{
				Seconds: startTime.Unix(),
			},
			EndTime: &timestamp.Timestamp{
				Seconds: endTime.Unix(),
			},
		},
		Aggregation: &monitoringpb.Aggregation{
			PerSeriesAligner: monitoringpb.Aggregation_ALIGN_MEAN,
			AlignmentPeriod: &duration.Duration{
				Seconds: 600,
			},
		},
	}
	it := client.ListTimeSeries(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read time series value: %v", err)
		}
		fmt.Fprintln(w, resp.GetMetric().GetLabels()["instance_name"])
		fmt.Fprintf(w, "\tNow: %.4f\n", resp.GetPoints()[0].GetValue().GetDoubleValue())
		fmt.Fprintf(w, "\t10 minutes ago: %.4f\n", resp.GetPoints()[1].GetValue().GetDoubleValue())
	}
	fmt.Fprintln(w, "Done")
	return nil
}

// [END monitoring_read_timeseries_align]
