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

package snippets

// [START monitoring_create_metric]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"google.golang.org/genproto/googleapis/api/label"
	"google.golang.org/genproto/googleapis/api/metric"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// createCustomMetric creates a custom metric specified by the metric type.
func createCustomMetric(w io.Writer, projectID, metricType string) (*metricpb.MetricDescriptor, error) {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, err
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
	m, err := c.CreateMetricDescriptor(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not create custom metric: %v", err)
	}

	fmt.Fprintf(w, "Created %s\n", m.GetName())
	return m, nil
}

// [END monitoring_create_metric]
