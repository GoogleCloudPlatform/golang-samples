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
package opencensus

// [START monitoring_opencensus_configure_trace_exemplar]
import (
	"time"

	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	distributionpb "google.golang.org/genproto/googleapis/api/distribution"
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
