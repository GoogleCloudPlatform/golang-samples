// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRiskNumerical(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 10, time.Second, func(_ *testutil.R) {
		buf := new(bytes.Buffer)
		riskNumerical(buf, client, projectID, "bigquery-public-data", "dlp-test-topic", "dlp-test-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		wants := []string{"Created job", "Value range", "Value at"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				t.Fatalf("riskNumerical got %s, want substring %q", got, want)
			}
		}
	})
}

func TestRiskCategorical(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 10, time.Second, func(_ *testutil.R) {
		buf := new(bytes.Buffer)
		riskCategorical(buf, client, projectID, "bigquery-public-data", "dlp-test-topic", "dlp-test-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		wants := []string{"Created job", "Histogram bucket", "Most common value occurs"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				t.Fatalf("riskCategorical got %s, want substring %q", got, want)
			}
		}
	})
}

func TestRiskKAnonymity(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 10, time.Second, func(_ *testutil.R) {
		buf := new(bytes.Buffer)
		riskKAnonymity(buf, client, projectID, "bigquery-public-data", "dlp-test-topic", "dlp-test-sub", "nhtsa_traffic_fatalities", "accident_2015", "state_number", "county")
		wants := []string{"Created job", "Histogram bucket", "Size range"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				t.Fatalf("riskKAnonymity got %s, want substring %q", got, want)
			}
		}
	})
}

func TestRiskLDiversity(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 10, time.Second, func(_ *testutil.R) {
		fmt.Println("Running")
		buf := new(bytes.Buffer)
		riskLDiversity(buf, client, projectID, "bigquery-public-data", "dlp-test-topic", "dlp-test-sub", "nhtsa_traffic_fatalities", "accident_2015", "city", "state_number", "county")
		wants := []string{"Created job", "Histogram bucket", "Size range"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				t.Fatalf("riskLDiversity got %s, want substring %q", got, want)
			}
		}
	})
}

func TestRiskKMap(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 10, time.Second, func(_ *testutil.R) {
		buf := new(bytes.Buffer)
		riskKMap(buf, client, projectID, "bigquery-public-data", "dlp-test-topic", "dlp-test-sub", "san_francisco", "bikeshare_trips", "US", "zip_code")
		wants := []string{"Created job", "Histogram bucket", "Anonymity range"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				t.Fatalf("riskKMap got %s, want substring %q", got, want)
			}
		}
	})
}
