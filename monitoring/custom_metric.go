// Copyright 2015 Google, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This example shows how to create a custom metric and write TimeSeries value
// to it. It writes a GAUGE measurement, which is a measure of value at a
// specific point in time. This means the startTime and endTime of the interval
// are the same. To make it easier to see the output, a random value is written.
// When reading the TimeSeries back, a window of the last 5 minutes is used.
// See README.md for instructions on how to run.

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"math/rand"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/monitoring/v3"
)

// createCustomMetric creates a custom metric specified by metricName
func createCustomMetric(s *monitoring.Service, projectResource string, metricType string, metricName string) error {
	ld := monitoring.LabelDescriptor{Key: "environment", ValueType: "STRING", Description: "An arbitrary measurement"}
	md := monitoring.MetricDescriptor{
		Name:        metricName,
		Type:        metricType,
		Labels:      []*monitoring.LabelDescriptor{&ld},
		MetricKind:  "GAUGE",
		ValueType:   "INT64",
		Unit:        "items",
		Description: "An arbitrary measurement",
		DisplayName: "Custom Metric",
	}
	resp, err := s.Projects.MetricDescriptors.Create(projectResource, &md).Do()
	if err != nil {
		return fmt.Errorf("Could not create custom metric: %v", err)
	}

	log.Printf("createCustomMetric: %s\n", formatResource(resp))
	return nil
}

// getCustomMetric reads the custom metric created
func getCustomMetric(s *monitoring.Service, projectResource string, metricType string, metricName string) (*monitoring.ListMetricDescriptorsResponse, error) {
	resp, err := s.Projects.MetricDescriptors.List(projectResource).
		Filter(fmt.Sprintf("metric.type=\"%s\"", metricType)).Do()
	if err != nil {
		return nil, fmt.Errorf("Could not get custom metric: %v", err)
	}

	log.Printf("getCustomMetric: %s\n", formatResource(resp))
	return resp, nil
}

// writeTimeSeriesValue writes a value for the custom metric created
func writeTimeSeriesValue(s *monitoring.Service, projectResource string, metricType string, metricName string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
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
		MetricKind: "GAUGE",
		ValueType:  "INT64",
		Points: []*monitoring.Point{
			&monitoring.Point{
				Interval: &monitoring.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoring.TypedValue{
					Int64Value: rand.Int63n(10),
				},
			},
		},
	}

	createTimeseriesRequest := monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{&timeseries}}

	log.Printf("writeTimeseriesRequest: %s\n", formatResource(createTimeseriesRequest))
	_, err := s.Projects.TimeSeries.
		Create(projectResource, &createTimeseriesRequest).Do()
	if err != nil {
		return fmt.Errorf("Could not write time series value, %v ", err)
	} else {
		return nil
	}
}

// readTimeSeriesValue reads the TimeSeries for the value specified by metricName in a time window from
// the last 5 minutes
func readTimeSeriesValue(s *monitoring.Service, projectResource string, metricType string, metricName string) error {
	startTime := time.Now().UTC().Add(time.Minute * -5)
	endTime := time.Now().UTC()
	resp, err := s.Projects.TimeSeries.List(projectResource).
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

func main() {
	client, err := google.DefaultClient(
		oauth2.NoContext,
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/monitoring",
		"https://www.googleapis.com/auth/monitoring.read",
		"https://www.googleapis.com/auth/monitoring.write",
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 2 {
		fmt.Println("Usage: auth.go <project_id>")
		return
	}
	projectResource := "projects/" + os.Args[1]
	// Create the Google Cloud Monitoring Service
	svc, err := monitoring.New(client)
	if err != nil {
		log.Fatal(err)
	}

	metricType := "custom.googleapis.com/custom_measurement"
	metricName := projectResource + "/metricDescriptors/" + metricType

	// create the metric
	err = createCustomMetric(svc, projectResource, metricType, metricName)
	if err != nil {
		log.Fatal(err)
	}

	var resp *monitoring.ListMetricDescriptorsResponse
	// wait until the new metric can be read back
	for resp == nil || resp.MetricDescriptors == nil {
		resp, err = getCustomMetric(svc, projectResource, metricType, metricName)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	// write a TimeSeries value for that metric
	err = writeTimeSeriesValue(svc, projectResource, metricType, metricName)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	// read the TimeSeries for the last 5 minutes for that metric
	readTimeSeriesValue(svc, projectResource, metricType, metricName)
}
