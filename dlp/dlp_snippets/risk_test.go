// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var riskTopicName = "dlp-risk-test-topic"
var riskSubscriptionName = "dlp-risk-test-sub"

func TestRiskNumerical(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 20, 2*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		riskNumerical(buf, client, projectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		wants := []string{"Created job", "Value range", "Value at"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				r.Errorf("riskNumerical got %s, want substring %q", got, want)
			}
		}
	})
}

func TestRiskCategorical(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 20, 2*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		riskCategorical(buf, client, projectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "nhtsa_traffic_fatalities", "accident_2015", "state_number")
		wants := []string{"Created job", "Histogram bucket", "Most common value occurs"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				r.Errorf("riskCategorical got %s, want substring %q", got, want)
			}
		}
	})
}

func TestRiskKAnonymity(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 20, 2*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		riskKAnonymity(buf, client, projectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "nhtsa_traffic_fatalities", "accident_2015", "state_number", "county")
		wants := []string{"Created job", "Histogram bucket", "Size range"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				r.Errorf("riskKAnonymity got %s, want substring %q", got, want)
			}
		}
	})
}

func TestRiskLDiversity(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 20, 2*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		riskLDiversity(buf, client, projectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "nhtsa_traffic_fatalities", "accident_2015", "city", "state_number", "county")
		wants := []string{"Created job", "Histogram bucket", "Size range"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				r.Errorf("riskLDiversity got %s, want substring %q", got, want)
			}
		}
	})
}

func TestRiskKMap(t *testing.T) {
	testutil.SystemTest(t)
	testutil.Retry(t, 20, 2*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		riskKMap(buf, client, projectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "san_francisco", "bikeshare_trips", "US", "zip_code")
		wants := []string{"Created job", "Histogram bucket", "Anonymity range"}
		got := buf.String()
		for _, want := range wants {
			if !strings.Contains(got, want) {
				r.Errorf("riskKMap got %s, want substring %q", got, want)
			}
		}
	})
}
