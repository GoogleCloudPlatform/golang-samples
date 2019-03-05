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
