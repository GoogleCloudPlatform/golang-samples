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

// [START livestream_get_channel_event]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// getChannelEvent gets a previously-created channel event.
func getChannelEvent(w io.Writer, projectID, location, channelID, eventID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// channelID := "my-channel-id"
	// eventID := "my-channel-event"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.GetEventRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/channels/%s/events/%s", projectID, location, channelID, eventID),
	}

	response, err := client.GetEvent(ctx, req)
	if err != nil {
		return fmt.Errorf("GetEvent: %w", err)
	}

	fmt.Fprintf(w, "Channel event: %v", response.Name)
	return nil
}

// [END livestream_get_channel_event]
