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

// Package risk contains example snippets using the DLP API to create risk jobs.
package risk

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	riskTopicName        = "dlp-risk-test-topic"
	riskSubscriptionName = "dlp-risk-test-sub"
)

func TestRisk(t *testing.T) {
	tc := testutil.SystemTest(t)
	tests := []struct {
		name string
		fn   func(r *testutil.R)
	}{
		{
			name: "Numerical",
			fn: func(r *testutil.R) {
				buf := new(bytes.Buffer)
				riskNumerical(buf, tc.ProjectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "nhtsa_traffic_fatalities", "accident_2015", "state_number")
				wants := []string{"Created job", "Value range", "Value at"}
				got := buf.String()
				for _, want := range wants {
					if !strings.Contains(got, want) {
						r.Errorf("riskNumerical got %s, want substring %q", got, want)
					}
				}
			},
		},
		{
			name: "Categorical",
			fn: func(r *testutil.R) {
				buf := new(bytes.Buffer)
				riskCategorical(buf, tc.ProjectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "nhtsa_traffic_fatalities", "accident_2015", "state_number")
				wants := []string{"Created job", "Histogram bucket", "Most common value occurs"}
				got := buf.String()
				for _, want := range wants {
					if !strings.Contains(got, want) {
						r.Errorf("riskCategorical got %s, want substring %q", got, want)
					}
				}
			},
		},
		{
			name: "K Anonymity",
			fn: func(r *testutil.R) {
				buf := new(bytes.Buffer)
				riskKAnonymity(buf, tc.ProjectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "nhtsa_traffic_fatalities", "accident_2015", "state_number", "county")
				wants := []string{"Created job", "Histogram bucket", "Size range"}
				got := buf.String()
				for _, want := range wants {
					if !strings.Contains(got, want) {
						r.Errorf("riskKAnonymity got %s, want substring %q", got, want)
					}
				}
			},
		},
		{
			name: "L Diversity",
			fn: func(r *testutil.R) {
				buf := new(bytes.Buffer)
				riskLDiversity(buf, tc.ProjectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "nhtsa_traffic_fatalities", "accident_2015", "city", "state_number", "county")
				wants := []string{"Created job", "Histogram bucket", "Size range"}
				got := buf.String()
				for _, want := range wants {
					if !strings.Contains(got, want) {
						r.Errorf("riskLDiversity got %s, want substring %q", got, want)
					}
				}
			},
		},
		{
			name: "K Map",
			fn: func(r *testutil.R) {
				buf := new(bytes.Buffer)
				riskKMap(buf, tc.ProjectID, "bigquery-public-data", riskTopicName, riskSubscriptionName, "san_francisco", "bikeshare_trips", "US", "zip_code")
				wants := []string{"Created job", "Histogram bucket", "Anonymity range"}
				got := buf.String()
				for _, want := range wants {
					if !strings.Contains(got, want) {
						r.Errorf("riskKMap got %s, want substring %q", got, want)
					}
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testutil.Retry(t, 20, 2*time.Second, test.fn)
		})
	}
}
