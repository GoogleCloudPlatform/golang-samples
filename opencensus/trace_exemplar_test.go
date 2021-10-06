// Copyright 2020 Google LLC
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

package opencensus

import (
	"context"
	"testing"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func writeTimeSeriesData(projectID string) error {
	ctx := context.Background()
	dataPoint := createDataPointWithExemplar()
	// Creates a client.
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return err
	}
	// Writes time series data.
	err = client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: monitoring.MetricProjectPath(projectID),
		TimeSeries: []*monitoringpb.TimeSeries{{
			Metric:     &metricpb.Metric{Type: "custom.googleapis.com/distribution"},
			MetricKind: metricpb.MetricDescriptor_CUMULATIVE,
			ValueType:  metricpb.MetricDescriptor_DISTRIBUTION,
			Resource: &monitoredrespb.MonitoredResource{
				Type:   "generic_node",
				Labels: map[string]string{"location": "us-east1-a", "namespace": "space", "node_id": "1"},
			},
			Points: []*monitoringpb.Point{dataPoint},
		}},
	})
	return err
}

func TestWriteTimeSeriesData(t *testing.T) {
	tc := testutil.SystemTest(t)
	if err := writeTimeSeriesData(tc.ProjectID); err != nil {
		t.Fatalf("writeTimeSeriesData: %v", err)
	}
}
