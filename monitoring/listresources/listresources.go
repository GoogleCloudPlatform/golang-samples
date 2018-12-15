// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command listresources lists the Google Cloud Monitoring v3 Environment against an authenticated user.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

const metric = "compute.googleapis.com/instance/cpu/usage_time"

func projectResource(projectID string) string {
	return "projects/" + projectID
}

// listMonitoredResourceDescriptor lists all the resources available to be monitored in the API.
func listMonitoredResourceDescriptors(projectID string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}

	req := &monitoringpb.ListMonitoredResourceDescriptorsRequest{
		Name: projectResource(projectID),
	}
	iter := c.ListMonitoredResourceDescriptors(ctx, req)

	var list []*monitoredrespb.MonitoredResourceDescriptor
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Could not list time series: %v", err)
		}
		list = append(list, resp)
	}

	log.Printf("listMonitoredResourceDescriptors: %s\n", formatResource(list))
	return nil
}

// listMetricDescriptors lists the metrics specified by the metric constant.
func listMetricDescriptors(projectID string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}

	req := &monitoringpb.ListMetricDescriptorsRequest{
		Name:   projectResource(projectID),
		Filter: fmt.Sprintf("metric.type=%q", metric),
	}
	iter := c.ListMetricDescriptors(ctx, req)

	var list []*metricpb.MetricDescriptor
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("could not list metric descriptors: %v", err)
		}
		list = append(list, resp)

	}
	log.Printf("listMetricDescriptors %s\n", formatResource(list))
	return nil
}

// listTimesSeries lists all the timeseries created for metric created in a 5
// minute interval an hour ago
func listTimeSeries(projectID string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}

	startTime := time.Now().UTC().Add(time.Minute * -5).Unix()
	endTime := time.Now().UTC().Unix()

	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   projectResource(projectID),
		Filter: fmt.Sprintf("metric.type=\"%s\"", metric),
		Interval: &monitoringpb.TimeInterval{
			StartTime: &timestamp.Timestamp{Seconds: startTime},
			EndTime:   &timestamp.Timestamp{Seconds: endTime},
		},
		PageSize: 5,
	}
	iter := c.ListTimeSeries(ctx, req)

	var series []*monitoringpb.TimeSeries
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read time series value, %v ", err)
		}
		series = append(series, resp)
	}

	log.Printf("readTimeseriesValue: %s\n", formatResource(series))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: listresources <project_id>")
		return
	}

	projectID := os.Args[1]

	if err := listMonitoredResourceDescriptors(projectID); err != nil {
		log.Fatal(err)
	}
	if err := listMetricDescriptors(projectID); err != nil {
		log.Fatal(err)
	}
	if err := listTimeSeries(projectID); err != nil {
		log.Fatal(err)
	}
}

// formatResource marshals a response objects as JSON.
func formatResource(resource interface{}) []byte {
	b, err := json.MarshalIndent(resource, "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}
