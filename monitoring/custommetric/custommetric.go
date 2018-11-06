// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command custommetric creates a custom metric and writes TimeSeries value
// to it. It writes a GAUGE measurement, which is a measure of value at a
// specific point in time. This means the startTime and endTime of the interval
// are the same. To make it easier to see the output, a random value is written.
// When reading the TimeSeries back, a window of the last 5 minutes is used.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/monitoring/apiv3"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/api/label"
	"google.golang.org/genproto/googleapis/api/metric"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredres "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

const metricType = "custom.googleapis.com/custom_measurement"

func projectResource(projectID string) string {
	return "projects/" + projectID
}

// [START monitoring_create_metric]

// createCustomMetric creates a custom metric specified by the metric type.
func createCustomMetric(projectID, metricType string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	md := &metric.MetricDescriptor{
		Name: "Custom Metric",
		Type: metricType,
		Labels: []*label.LabelDescriptor{{
			Key:         "environment",
			ValueType:   label.LabelDescriptor_STRING,
			Description: "An arbitrary measurement",
		}},
		MetricKind:  metric.MetricDescriptor_GAUGE,
		ValueType:   metric.MetricDescriptor_INT64,
		Unit:        "s",
		Description: "An arbitrary measurement",
		DisplayName: "Custom Metric",
	}
	req := &monitoringpb.CreateMetricDescriptorRequest{
		Name:             "projects/" + projectID,
		MetricDescriptor: md,
	}
	resp, err := c.CreateMetricDescriptor(ctx, req)
	if err != nil {
		return fmt.Errorf("could not create custom metric: %v", err)
	}

	log.Printf("createCustomMetric: %s\n", formatResource(resp))
	return nil
}

// [END monitoring_create_metric]

// [START monitoring_list_descriptors]

// getCustomMetric reads the custom metric created.
func getCustomMetric(projectID, metricType string) (*metric.MetricDescriptor, error) {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, err
	}
	req := &monitoringpb.GetMetricDescriptorRequest{
		Name: fmt.Sprintf("projects/%s/metricDescriptors/%s", projectID, metricType),
	}
	resp, err := c.GetMetricDescriptor(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not get custom metric: %v", err)
	}

	log.Printf("getCustomMetric: %s\n", formatResource(resp))
	return resp, nil
}

// [END monitoring_list_descriptors]

// [START monitoring_delete_metric]

// deleteMetric deletes the given metric.
func deleteMetric(projectID, metricType string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	metricResource := "projects/" + projectID + "/metricDescriptors/" + metricType
	req := &monitoringpb.DeleteMetricDescriptorRequest{
		Name: metricResource,
	}
	err = c.DeleteMetricDescriptor(ctx, req)
	if err != nil {
		return fmt.Errorf("could not delete metric: %v", err)
	}
	log.Printf("Deleted metric: %q\n", metricType)
	return nil
}

// [END monitoring_delete_metric]

// [START monitoring_write_timeseries]

// writeTimeSeriesValue writes a value for the custom metric created
func writeTimeSeriesValue(projectID, metricType string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	now := &timestamp.Timestamp{
		Seconds: time.Now().Unix(),
	}
	req := &monitoringpb.CreateTimeSeriesRequest{
		Name: projectResource(projectID),
		TimeSeries: []*monitoringpb.TimeSeries{{
			Metric: &metricpb.Metric{
				Type: metricType,
				Labels: map[string]string{
					"environment": "STAGING",
				},
			},
			Resource: &monitoredres.MonitoredResource{
				Type: "gce_instance",
				Labels: map[string]string{
					"instance_id": "test-instance",
					"zone":        "us-central1-f",
				},
			},
			Points: []*monitoringpb.Point{{
				Interval: &monitoringpb.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_Int64Value{
						Int64Value: rand.Int63n(10),
					},
				},
			}},
		}},
	}
	log.Printf("writeTimeseriesRequest: %s\n", formatResource(req))

	err = c.CreateTimeSeries(ctx, req)
	if err != nil {
		return fmt.Errorf("could not write time series value, %v ", err)
	}
	return nil
}

// [END monitoring_write_timeseries]

// [START monitoring_read_timeseries_simple]

// readTimeSeriesValue reads the TimeSeries for the value specified by metric type in a time window from the last 5 minutes.
func readTimeSeriesValue(projectID, metricType string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	startTime := time.Now().UTC().Add(time.Minute * -5).Unix()
	endTime := time.Now().UTC().Unix()

	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   projectResource(projectID),
		Filter: fmt.Sprintf("metric.type=\"%s\"", metricType),
		Interval: &monitoringpb.TimeInterval{
			StartTime: &timestamp.Timestamp{Seconds: startTime},
			EndTime:   &timestamp.Timestamp{Seconds: endTime},
		},
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

// [END monitoring_read_timeseries_simple]

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println("Usage: custommetric <project_id>")
		return
	}

	projectID := os.Args[1]

	// Create the metric.
	if err := createCustomMetric(projectID, metricType); err != nil {
		log.Fatal(err)
	}

	// Wait until the new metric can be read back.
	for {
		resp, err := getCustomMetric(projectID, metricType)
		if err != nil {
			log.Fatal(err)
		}
		if resp != nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// Write a TimeSeries value for that metric
	if err := writeTimeSeriesValue(projectID, metricType); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	// Read the TimeSeries for the last 5 minutes for that metric.
	if err := readTimeSeriesValue(projectID, metricType); err != nil {
		log.Fatal(err)
	}

	if err := deleteMetric(projectID, metricType); err != nil {
		log.Fatal(err)
	}
}

// formatResource marshals a response object as JSON.
func formatResource(resource interface{}) []byte {
	b, err := json.MarshalIndent(resource, "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}
