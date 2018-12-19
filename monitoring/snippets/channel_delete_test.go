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
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func TestDeleteChannel(t *testing.T) {
	tc := testutil.SystemTest(t)

	c, err := createChannel(tc.ProjectID)
	if err != nil {
		t.Fatalf("Error creating test channel: %v", err)
	}

	buf := &bytes.Buffer{}
	if err := deleteChannel(buf, c.GetName()); err != nil {
		t.Fatalf("deleteChannel: %v", err)
	}
	want := "Deleted channel"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("deleteChannel got %q, want to contain %q", got, want)
	}
}

// createChannel creates a channel.
func createChannel(projectID string) (*monitoringpb.NotificationChannel, error) {
	ctx := context.Background()

	client, err := monitoring.NewNotificationChannelClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &monitoringpb.CreateNotificationChannelRequest{
		Name: "projects/" + projectID,
		NotificationChannel: &monitoringpb.NotificationChannel{
			Type:        "email",
			DisplayName: "Email",
			Labels:      map[string]string{"email_address": "noreply@google.com"},
		},
	}

	channel, err := client.CreateNotificationChannel(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateNotificationChannel: %v", err)
	}

	return channel, nil
}
