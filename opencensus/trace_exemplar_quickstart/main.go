package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	distributionpb "google.golang.org/genproto/googleapis/api/distribution"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create client"))
	}

	// Sets your Google Cloud Platform project ID.
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// [START monitoring_opencensus_configure_trace_exemplar]
	// Prepares an individual data point
	end := time.Now().Unix()
	dataPoint := &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime:   &googlepb.Timestamp{Seconds: end},
			StartTime: &googlepb.Timestamp{Seconds: end - 60},
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
	// [END monitoring_opencensus_configure_trace_exemplar]

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
		log.Fatal(errors.Wrap(err, "failed to write time series data"))
	}

	// Reads that time series data.
	it := client.ListTimeSeries(ctx, &monitoringpb.ListTimeSeriesRequest{
		Name:   monitoring.MetricProjectPath(projectID),
		Filter: "resource.type=generic_node metric.type=\"custom.googleapis.com/distribution\"",
		Interval: &monitoringpb.TimeInterval{
			EndTime:   &googlepb.Timestamp{Seconds: end},
			StartTime: &googlepb.Timestamp{Seconds: end - 60},
		},
	})
	for {
		if response, err := it.Next(); err == iterator.Done {
			break
		} else if err != nil {
			log.Fatal(errors.Wrap(err, "failed to query time series data"))
		} else {
			fmt.Println(response.GetPoints()[0].GetValue().GetDistributionValue().GetExemplars())
		}
	}

	// Closes the client and flushes the data to Stackdriver.
	if err := client.Close(); err != nil {
		log.Fatal(errors.Wrap(err, "failed to close client"))
	}
}
