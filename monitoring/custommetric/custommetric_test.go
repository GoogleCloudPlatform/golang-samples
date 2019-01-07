// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package custommetric

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func TestMain(m *testing.M) {
	// These functions are noisy.
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}

func TestCustomMetric(t *testing.T) {
	hc := testutil.SystemTest(t)
	if err := createCustomMetric(hc.ProjectID, metricType); err != nil {
		t.Fatal(err)
	}

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		_, err := getCustomMetric(hc.ProjectID, metricType)
		if err != nil {
			r.Errorf("%v", err)
		}
	})

	time.Sleep(2 * time.Second)

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := writeTimeSeriesValue(hc.ProjectID, metricType); err != nil {
			t.Error(err)
		}
	})

	time.Sleep(2 * time.Second)

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := readTimeSeriesValue(hc.ProjectID, metricType); err != nil {
			r.Errorf("%v", err)
		}
	})

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := deleteMetric(hc.ProjectID, metricType); err != nil {
			t.Error(err)
		}
	})
}

// getCustomMetric reads the custom metric created.
func getCustomMetric(projectID, metricType string) (*metric.MetricDescriptor, error) {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return nil, err
	}
	req := &monitoringpb.GetMetricDescriptorRequest{
		Name: fmt.Sprintf("projects/%s/metricDescriptors/%s", projectID, metricType),
	}
	resp, err := c.GetMetricDescriptor(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not get custom metric: %v", err)
	}

	return resp, nil
}
