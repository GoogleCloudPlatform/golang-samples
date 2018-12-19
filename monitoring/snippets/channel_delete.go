// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
