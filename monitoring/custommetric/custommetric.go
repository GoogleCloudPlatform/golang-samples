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

// Package custommetric contains custom metric samples.
package custommetric

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/api/iterator"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredres "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

const metricType = "custom.googleapis.com/custom_measurement"

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
		Name: "projects/" + projectID,
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
	log.Printf("writeTimeseriesRequest: %+v\n", req)

	err = c.CreateTimeSeries(ctx, req)
	if err != nil {
		return fmt.Errorf("could not write time series value, %v ", err)
	}
	return nil
}

// [END monitoring_write_timeseries]

// [START monitoring_read_timeseries_simple]

// readTimeSeriesValue reads the TimeSeries for the value specified by metric type in a time window from the last 20 minutes.
func readTimeSeriesValue(projectID, metricType string) error {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	startTime := time.Now().UTC().Add(time.Minute * -20).Unix()
	endTime := time.Now().UTC().Unix()

	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   "projects/" + projectID,
		Filter: fmt.Sprintf("metric.type=\"%s\"", metricType),
		Interval: &monitoringpb.TimeInterval{
			StartTime: &timestamp.Timestamp{Seconds: startTime},
			EndTime:   &timestamp.Timestamp{Seconds: endTime},
		},
	}
	iter := c.ListTimeSeries(ctx, req)

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read time series value, %v ", err)
		}
		log.Printf("%+v\n", resp)
	}

	return nil
}

// [END monitoring_read_timeseries_simple]
