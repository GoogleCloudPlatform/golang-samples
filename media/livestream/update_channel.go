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

// [START livestream_update_channel]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateChannel updates an existing channel with a different input.
func updateChannel(w io.Writer, projectID, location, channelID, inputID string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	// channelID := "my-channel"
	// inputID := "my-updated-input"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.UpdateChannelRequest{
		Channel: &livestreampb.Channel{
			Name: fmt.Sprintf("projects/%s/locations/%s/channels/%s", projectID, location, channelID),
			InputAttachments: []*livestreampb.InputAttachment{
				{
					Key:   "updated-input",
					Input: fmt.Sprintf("projects/%s/locations/%s/inputs/%s", projectID, location, inputID),
				},
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"input_attachments",
			},
		},
	}
	// Updates the input.
	op, err := client.UpdateChannel(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateChannel: %w", err)
	}
	response, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %w", err)
	}

	fmt.Fprintf(w, "Updated channel: %v", response.Name)
	return nil
}

// [END livestream_update_channel]
