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

	"golang.org/x/net/context"

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
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	s, err := createService(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := listMonitoredResourceDescriptors(s, hc.ProjectID); err != nil {
		log.Fatal(err)
	}
	if err := listMetricDescriptors(s, hc.ProjectID); err != nil {
		log.Fatal(err)
	}
	if err := listTimeSeries(s, hc.ProjectID); err != nil {
		log.Fatal(err)
	}
}
