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
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeleteChannel(t *testing.T) {
	tc := testutil.SystemTest(t)

	testutil.Retry(t, 5, 10*time.Second, func(r *testutil.R) {
		c, err := createChannel(tc.ProjectID)
		if err != nil {
			r.Errorf("Error creating test channel: %v", err)
			return
		}

		buf := &bytes.Buffer{}
		if err := deleteChannel(buf, c.GetName()); err != nil {
			r.Errorf("deleteChannel: %v", err)
			return
		}
		want := "Deleted channel"
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("deleteChannel got %q, want to contain %q", got, want)
			return
		}
	})
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
		return nil, fmt.Errorf("CreateNotificationChannel: %w", err)
	}

	return channel, nil
}
