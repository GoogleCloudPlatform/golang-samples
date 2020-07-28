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

// Sample opencensus_spanner_quickstart contains a sample application that
// uses Google Spanner Go client, and reports metrics
// and traces for the outgoing requests.
package trace_exemplar_quickstart

// [START monitoring_opencensus_configure_trace_exemplar]
import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	distributionpb "google.golang.org/genproto/googleapis/api/distribution"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func createDataPointWithExemplar() *monitoringpb.Point {
	end := time.Now().Unix()
	dataPoint := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			StartTime: &googlepb.Timestamp{Seconds: end - 60},
			EndTime:   &googlepb.Timestamp{Seconds: end},
		},
		Value: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DistributionValue{
			DistributionValue: &distributionpb.Distribution{
				Count: 14,
				BucketOptions: &distributionpb.Distribution_BucketOptions{Options: &distributionpb.Distribution_BucketOptions_LinearBuckets{
					LinearBuckets: &distributionpb.Distribution_BucketOptions_Linear{NumFiniteBuckets: 2, Width: 3, Offset: 0},
				}},
				BucketCounts: []int64{5, 6, 3},
				Exemplars: []*distributionpb.Distribution_Exemplar{
					&distributionpb.Distribution_Exemplar{Value: 1, Timestamp: &googlepb.Timestamp{Seconds: end - 30}},
					&distributionpb.Distribution_Exemplar{Value: 4, Timestamp: &googlepb.Timestamp{Seconds: end - 30}},
				},
			},
		}},
	}
	return dataPoint
}

// [END monitoring_opencensus_configure_trace_exemplar]

func writeTimeSeriesData(w io.Writer, projectID string) error {
	ctx := context.Background()
	dataPoint := createDataPointWithExemplar()
	// Creates a client.
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create client"))
	}
	// Writes time series data.
	if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
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
	}); err != nil {
		return err
	} else {
		fmt.Fprintln(w, "Done")
	}
	return nil
}
