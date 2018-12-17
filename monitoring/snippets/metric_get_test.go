// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

const metricType = "custom.googleapis.com/golang-samples-tests/get"

func TestGetMetricDescriptor(t *testing.T) {
	tc := testutil.SystemTest(t)

	m, err := createMetric(tc.ProjectID)
	if err != nil {
		t.Fatalf("createMetric: %v", err)
	}
	defer deleteMetric(m.GetName())

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := getMetricDescriptor(buf, tc.ProjectID, metricType); err != nil {
			r.Errorf("getMetricDescriptor: %v", err)
		}
		want := "Name:"
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("getMetricDescriptor got %q, want to contain %q", got, want)
		}
	})
}

func createMetric(projectID string) (*metric.MetricDescriptor, error) {
	ctx := context.Background()

	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewMetricClient: %v", err)
	}
	req := &monitoringpb.CreateMetricDescriptorRequest{
		Name: "projects/" + projectID,
		MetricDescriptor: &metric.MetricDescriptor{
			Type:       metricType,
			MetricKind: metric.MetricDescriptor_GAUGE,
			ValueType:  metric.MetricDescriptor_INT64,
		},
	}

	m, err := c.CreateMetricDescriptor(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateMetricDescriptor: %v", err)
	}

	return m, nil
}

func deleteMetric(name string) error {
	ctx := context.Background()

	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return fmt.Errorf("NewMetricClient: %v", err)
	}

	req := &monitoringpb.DeleteMetricDescriptorRequest{
		Name: name,
	}

	if err := c.DeleteMetricDescriptor(ctx, req); err != nil {
		return fmt.Errorf("DeleteMetricDescriptor: %v", err)
	}
	return nil
}
