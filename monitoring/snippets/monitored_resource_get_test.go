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

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
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
		return "", fmt.Errorf("Could not list time series: %v", err)
	}
	return resp.GetName(), nil
}
