// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(m *testing.M) {
	// These functions are noisy.
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}

func TestListResources(t *testing.T) {
	hc := testutil.SystemTest(t)
	if err := listMonitoredResourceDescriptors(hc.ProjectID); err != nil {
		log.Fatal(err)
	}
	if err := listMetricDescriptors(hc.ProjectID); err != nil {
		log.Fatal(err)
	}
	if err := listTimeSeries(hc.ProjectID); err != nil {
		log.Fatal(err)
	}
}
