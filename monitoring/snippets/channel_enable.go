// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets contains examples of using the monitoring API.
package snippets

// [START monitoring_alert_enable_channel]

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/wrappers"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/genproto/protobuf/field_mask"
)

// enableChannel enables the given channel. channelName should be of the form
// "projects/[PROJECT_ID]/notificationChannels/[CHANNEL_ID]".
func enableChannel(w io.Writer, channelName string) error {
	ctx := context.Background()

	client, err := monitoring.NewNotificationChannelClient(ctx)
	if err != nil {
		return err
	}

	req := &monitoringpb.UpdateNotificationChannelRequest{
		UpdateMask: &field_mask.FieldMask{Paths: []string{"enabled"}},
		NotificationChannel: &monitoringpb.NotificationChannel{
			Name:    channelName,
			Enabled: &wrappers.BoolValue{Value: true},
		},
	}

	if _, err := client.UpdateNotificationChannel(ctx, req); err != nil {
		return fmt.Errorf("EnableNotificationChannel: %v", err)
	}

	fmt.Fprintf(w, "Enabled channel %q", channelName)
	return nil
}

// [END monitoring_alert_enable_channel]
