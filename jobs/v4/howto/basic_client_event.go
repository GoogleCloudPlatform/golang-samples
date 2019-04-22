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

package howto

// [START job_search_create_client_event]
import (
	"context"
	"fmt"
	"io"
	"time"

	talent "cloud.google.com/go/talent/apiv4beta1"
	"github.com/golang/protobuf/ptypes"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
)

// createClientEvent creates a client event.
func createClientEvent(w io.Writer, projectID string, requestID string, eventID string, relatedJobNames []string) (*talentpb.ClientEvent, error) {
	ctx := context.Background()

	// Create an eventService client.
	c, err := talent.NewEventClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("talent.NewEventClient: %v", err)
	}

	createTime, _ := ptypes.TimestampProto(time.Now())
	clientEventToCreate := &talentpb.ClientEvent{
		RequestId:  requestID,
		EventId:    eventID,
		CreateTime: createTime,
		Event: &talentpb.ClientEvent_JobEvent{
			JobEvent: &talentpb.JobEvent{
				Type: talentpb.JobEvent_VIEW,
				Jobs: relatedJobNames,
			},
		},
	}

	// Construct a createJob request.
	req := &talentpb.CreateClientEventRequest{
		Parent:      fmt.Sprintf("projects/%s", projectID),
		ClientEvent: clientEventToCreate,
	}

	resp, err := c.CreateClientEvent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateClientEvent: %v", err)
	}

	fmt.Fprintf(w, "Client event created: %v\n", resp.GetEvent())

	return resp, nil
}

// [END job_search_create_client_event]
