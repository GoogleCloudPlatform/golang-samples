// Copyright 2024 Google LLC
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

// [START livestream_delete_channel_clip]
import (
	"context"
	"fmt"
	"io"

	livestream "cloud.google.com/go/video/livestream/apiv1"
	"cloud.google.com/go/video/livestream/apiv1/livestreampb"
)

// deleteChannelClip deletes a previously-created channel clip.
func deleteChannelClip(w io.Writer, projectID, channelID, clipID string) error {
	// projectID := "my-project-id"
	// channelID := "my-channel"
	// clipID := "my-channel-clip"
	location := "us-central1"
	ctx := context.Background()
	client, err := livestream.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &livestreampb.DeleteClipRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/channels/%s/clips/%s", projectID, location, channelID, clipID),
	}

	op, err := client.DeleteClip(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteClip: %w", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Deleted channel clip")
	return nil
}

// [END livestream_delete_channel_clip]
