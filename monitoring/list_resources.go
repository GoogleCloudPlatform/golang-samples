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

// This example shows how to do basic authorization and list the
// Google Cloud Monitoring v3 Environment. See README.md for instructions on
// how to run.

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/monitoring/v3"
)

const metric = "compute.googleapis.com/instance/cpu/usage_time"

// listMonitoredResourceDescriptor lists all the resources available to be monitored in the API.
func listMonitoredResourceDescriptors(s *monitoring.Service, project string) error {
	resp, err := s.Projects.MonitoredResourceDescriptors.List(project).Do()
	if err != nil {
		return fmt.Errorf("Could not list time series: %v", err)
	}

	log.Printf("listMonitoredResourceDescriptors: %s\n", formatResource(resp))
	return nil
}

// listMetricDescriptors lists the metrics specified by the metric constant.
func listMetricDescriptors(s *monitoring.Service, projectResource string) error {
	resp, err := s.Projects.MetricDescriptors.List(projectResource).
		Filter(fmt.Sprintf("metric.name=%q", metric)).
		Do()
	if err != nil {
		return fmt.Errorf("Could not list metric descriptors: %v", err)
	}

	log.Printf("listMetricDescriptors %s\n", formatResource(resp))
	return nil
}

// listTimesSeries lists all the timeseries created for metric created in a 5
// minute interval an hour ago
func listTimeSeries(s *monitoring.Service, projectResource string) error {
	startTime := time.Now().UTC().Add(-time.Hour)
	endTime := startTime.Add(5 * time.Minute)

	resp, err := s.Projects.TimeSeries.List(projectResource).
		PageSize(3).
		Filter(fmt.Sprintf("metric.type=\"%s\"", metric)).
		IntervalStartTime(startTime.Format(time.RFC3339)).
		IntervalEndTime(endTime.Format(time.RFC3339)).
		Do()
	if err != nil {
		return fmt.Errorf("Could not list time series: %v", err)
	}

	log.Printf("listTimeseries %s\n", formatResource(resp))
	return nil
}

func main() {
	// walk through the basic calls of the Monitoring API
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
	s, err := monitoring.New(client)
	if err != nil {
		log.Fatal(err)
	}

	err = listMonitoredResourceDescriptors(s, projectResource)
	if err != nil {
		log.Fatal(err)
	}
	err = listMetricDescriptors(s, projectResource)
	if err != nil {
		log.Fatal(err)
	}
	err = listTimeSeries(s, projectResource)
	if err != nil {
		log.Fatal(err)
	}
}
