// Copyright 2020 Google LLC
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
package notifications

// [START scc_receive_notifications]
import (
	"bytes"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/jsonpb"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
)

func receiveMessages(w io.Writer, projectID string, subscriptionName string) error {
	// projectID := "your-project-id"
	// subsriptionName := "your-subscription-name"

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subscriptionName)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		var notificationMessage = new(securitycenterpb.NotificationMessage)
		jsonpb.Unmarshal(bytes.NewReader(msg.Data), notificationMessage)

		fmt.Fprintln(w, "Got finding: ", notificationMessage.GetFinding())
		msg.Ack()
		cancel()
	})
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}

	return nil
}

// [END scc_receive_notifications]
