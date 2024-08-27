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

package taillogs

// [START logging_tail_log_entries]
import (
	"context"
	"fmt"
	"io"

	logging "cloud.google.com/go/logging/apiv2"
	"cloud.google.com/go/logging/apiv2/loggingpb"
)

// tailLogs creates a channel to stream log entries that were recently ingested for a project
func tailLogs(projectID string) error {
	// projectID := "your_project_id"

	ctx := context.Background()
	client, err := logging.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient error: %w", err)
	}
	defer client.Close()

	stream, err := client.TailLogEntries(ctx)
	if err != nil {
		return fmt.Errorf("TailLogEntries error: %w", err)
	}
	defer stream.CloseSend()

	req := &loggingpb.TailLogEntriesRequest{
		ResourceNames: []string{
			"projects/" + projectID,
		},
	}
	if err := stream.Send(req); err != nil {
		return fmt.Errorf("stream.Send error: %w", err)
	}

	// read and print two or more streamed log entries
	for counter := 0; counter < 2; {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("stream.Recv error: %w", err)
		}
		fmt.Printf("received:\n%v\n", resp)
		if resp.Entries != nil {
			counter += len(resp.Entries)
		}
	}
	return nil
}

// [END logging_tail_log_entries]
