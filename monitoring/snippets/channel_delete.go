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

// Package snippets contains examples of using the monitoring API.
package snippets

// [START monitoring_alert_delete_channel]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// deleteChannel deletes the given channel. channelName should be of the form
// "projects/[PROJECT_ID]/notificationChannels/[CHANNEL_ID]".
func deleteChannel(w io.Writer, channelName string) error {
	ctx := context.Background()

	client, err := monitoring.NewNotificationChannelClient(ctx)
	if err != nil {
		return err
	}

	req := &monitoringpb.DeleteNotificationChannelRequest{
		Name: channelName,
	}

	if err := client.DeleteNotificationChannel(ctx, req); err != nil {
		return fmt.Errorf("DeleteNotificationChannel: %v", err)
	}

	fmt.Fprintf(w, "Deleted channel %q", channelName)
	return nil
}

// [END monitoring_alert_delete_channel]
