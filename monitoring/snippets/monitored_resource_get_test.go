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

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

func TestGetMonitoredResource(t *testing.T) {
	tc := testutil.SystemTest(t)

	m, err := randomMonitoredResource(tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	if err := getMonitoredResource(buf, m); err != nil {
		t.Fatalf("getMonitoredResource: %v", err)
	}
	want := "Name"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("getMonitoredResource got %q, want to contain %q", got, want)
	}
}

// randomMonitoredResource returns the name of a random resource in the project.
func randomMonitoredResource(projectID string) (string, error) {
	ctx := context.Background()
	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return "", err
	}

	req := &monitoringpb.ListMonitoredResourceDescriptorsRequest{
		Name: "projects/" + projectID,
	}
	iter := c.ListMonitoredResourceDescriptors(ctx, req)

	resp, err := iter.Next()
	if err == iterator.Done {
		return "", fmt.Errorf("no resources")
	}
	if err != nil {
		return "", fmt.Errorf("Could not list time series: %w", err)
	}
	return resp.GetName(), nil
}
