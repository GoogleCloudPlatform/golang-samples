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
	"log"
	"fmt"
	"io"
 
	logging "cloud.google.com/go/logging/apiv2"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
)

func tailLogs(projectID string) {
	ctx := context.Background()
	client, err := logging.NewClient(ctx)
	if err != nil {
		log.Fatalf("NewClient error: %v", err)
	}
	stream, err := client.TailLogEntries(ctx)
	if err != nil {
		log.Fatalf("TailLogEntries error: %v", err)
	}
	req := &loggingpb.TailLogEntriesRequest{
		ResourceNames: []string{
			"projects/" + projectID,
		},
	}
	if err := stream.Send(req); err != nil {
		log.Fatalf("stream.Send error: %v", err)
	}

	defer stream.CloseSend()
	
	// read and print two streamed log entries
	for counter := 0; counter < 2; counter++ {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream.Recv error: %v", err)
		}
		fmt.Printf("received:\n%v\n", resp)
	}
}
// [END logging_tail_log_entries]