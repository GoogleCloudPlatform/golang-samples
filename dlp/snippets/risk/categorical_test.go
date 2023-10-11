// Copyright 2023 Google LLC
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

package risk

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
)

func TestCategorical(t *testing.T) {
	tc := testutil.SystemTest(t)
	client, err := pubsub.NewClient(context.Background(), tc.ProjectID)
	if err != nil {
		t.Fatalf("pubsub.NewClient: %v", err)
	}

	buf := new(bytes.Buffer)
	u := uuid.Must(uuid.NewV4()).String()[:8]

	err = riskCategorical(buf, tc.ProjectID, "bigquery-public-data", riskTopicName+u, riskSubscriptionName+u, "nhtsa_traffic_fatalities", "accident_2015", "state_number")
	defer cleanupPubsub(t, client, riskTopicName+u, riskSubscriptionName+u)
	if err != nil {
		t.Errorf("riskCategorical got err: %v", err)
		return
	}
	if got, want := buf.String(), "Created job"; !strings.Contains(got, want) {
		t.Errorf("riskCategorical got %s, want substring %q", got, want)
	}

}
