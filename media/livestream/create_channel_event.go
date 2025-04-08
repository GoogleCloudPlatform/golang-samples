// Copyright 2022 Google LLC
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

package livestream

// [START livestream_create_channel_event]
import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/duration"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// createChannelEvent creates a channel event. An event is a sub-resource of a
// channel, which can be scheduled by the user to execute operations on a
// channel resource without having to stop the channel. This sample creates an
// ad break event.
func createChannelEvent(w io.Writer, projectID, location, channelID, eventID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// channelID := "my-channel"
	// eventID := "my-channel-event"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.CreateEventRequest{
		Parent:  fmt.Sprintf("projects/%s/locations/%s/channels/%s", projectID, location, channelID),
		EventId: eventID,
		Event: &livestreampb.Event{
			Task: &livestreampb.Event_AdBreak{
				AdBreak: &livestreampb.Event_AdBreakTask{
					Duration: &duration.Duration{
						Seconds: 30,
					},
				},
			},
			ExecuteNow: true,
		},
	}
	// Creates the channel event.
	response, err := client.CreateEvent(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateEvent: %w", err)
	}

	fmt.Fprintf(w, "Channel event: %v", response.Name)
	return nil
}

// [END livestream_create_channel_event]
