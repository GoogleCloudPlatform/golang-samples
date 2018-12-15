// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
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
