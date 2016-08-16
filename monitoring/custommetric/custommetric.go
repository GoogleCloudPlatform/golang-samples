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
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/monitoring/v3"
)

const metricType = "custom.googleapis.com/custom_measurement"

func projectResource(projectID string) string {
	return "projects/" + projectID
}

// createCustomMetric creates a custom metric specified by the metric type.
func createCustomMetric(s *monitoring.Service, projectID, metricType string) error {
	ld := monitoring.LabelDescriptor{Key: "environment", ValueType: "STRING", Description: "An arbitrary measurement"}
	md := monitoring.MetricDescriptor{
		Type:        metricType,
		Labels:      []*monitoring.LabelDescriptor{&ld},
		MetricKind:  "GAUGE",
		ValueType:   "INT64",
		Unit:        "items",
		Description: "An arbitrary measurement",
		DisplayName: "Custom Metric",
	}
	resp, err := s.Projects.MetricDescriptors.Create(projectResource(projectID), &md).Do()
	if err != nil {
		return fmt.Errorf("Could not create custom metric: %v", err)
	}

	log.Printf("createCustomMetric: %s\n", formatResource(resp))
	return nil
}

// getCustomMetric reads the custom metric created.
func getCustomMetric(s *monitoring.Service, projectID, metricType string) (*monitoring.ListMetricDescriptorsResponse, error) {
	resp, err := s.Projects.MetricDescriptors.List(projectResource(projectID)).
		Filter(fmt.Sprintf("metric.type=\"%s\"", metricType)).Do()
	if err != nil {
		return nil, fmt.Errorf("Could not get custom metric: %v", err)
	}

	log.Printf("getCustomMetric: %s\n", formatResource(resp))
	return resp, nil
}

// writeTimeSeriesValue writes a value for the custom metric created
func writeTimeSeriesValue(s *monitoring.Service, projectID, metricType string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	randVal := rand.Int63n(10)
	timeseries := monitoring.TimeSeries{
		Metric: &monitoring.Metric{
			Type: metricType,
			Labels: map[string]string{
				"environment": "STAGING",
			},
		},
		Resource: &monitoring.MonitoredResource{
			Labels: map[string]string{
				"instance_id": "test-instance",
				"zone":        "us-central1-f",
			},
			Type: "gce_instance",
		},
		Points: []*monitoring.Point{
			{
				Interval: &monitoring.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoring.TypedValue{
					Int64Value: &randVal,
				},
			},
		},
	}

	createTimeseriesRequest := monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{&timeseries},
	}

	log.Printf("writeTimeseriesRequest: %s\n", formatResource(createTimeseriesRequest))
	_, err := s.Projects.TimeSeries.Create(projectResource(projectID), &createTimeseriesRequest).Do()
	if err != nil {
		return fmt.Errorf("Could not write time series value, %v ", err)
	}
	return nil
}

// readTimeSeriesValue reads the TimeSeries for the value specified by metric type in a time window from the last 5 minutes.
func readTimeSeriesValue(s *monitoring.Service, projectID, metricType string) error {
	startTime := time.Now().UTC().Add(time.Minute * -5)
	endTime := time.Now().UTC()
	resp, err := s.Projects.TimeSeries.List(projectResource(projectID)).
		Filter(fmt.Sprintf("metric.type=\"%s\"", metricType)).
		IntervalStartTime(startTime.Format(time.RFC3339Nano)).
		IntervalEndTime(endTime.Format(time.RFC3339Nano)).
		Do()
	if err != nil {
		return fmt.Errorf("Could not read time series value, %v ", err)
	}
	log.Printf("readTimeseriesValue: %s\n", formatResource(resp))
	return nil
}

func createService(ctx context.Context) (*monitoring.Service, error) {
	hc, err := google.DefaultClient(ctx, monitoring.MonitoringScope)
	if err != nil {
		return nil, err
	}
	s, err := monitoring.New(hc)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) < 2 {
		fmt.Println("Usage: custommetric <project_id>")
		return
	}

	ctx := context.Background()
	s, err := createService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	projectID := os.Args[1]

	// Create the metric.
	if err := createCustomMetric(s, projectID, metricType); err != nil {
		log.Fatal(err)
	}

	// Wait until the new metric can be read back.
	for {
		resp, err := getCustomMetric(s, projectID, metricType)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.MetricDescriptors) != 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// Write a TimeSeries value for that metric
	if err := writeTimeSeriesValue(s, projectID, metricType); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	// Read the TimeSeries for the last 5 minutes for that metric.
	if err := readTimeSeriesValue(s, projectID, metricType); err != nil {
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
