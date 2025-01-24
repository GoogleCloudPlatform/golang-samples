// Copyright 2023 Google LLC
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

package topics

// [START pubsub_publisher_retry_settings]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub/v2"
	vkit "cloud.google.com/go/pubsub/v2/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
)

func publishWithRetrySettings(w io.Writer, projectID, topicID, msg string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// msg := "Hello World"
	ctx := context.Background()

	config := &pubsub.ClientConfig{
		PublisherCallOptions: &vkit.TopicAdminCallOptions{
			Publish: []gax.CallOption{
				gax.WithRetry(func() gax.Retryer {
					return gax.OnCodes([]codes.Code{
						codes.Aborted,
						codes.Canceled,
						codes.Internal,
						codes.ResourceExhausted,
						codes.Unknown,
						codes.Unavailable,
						codes.DeadlineExceeded,
					}, gax.Backoff{
						Initial:    250 * time.Millisecond, // default 100 milliseconds
						Max:        60 * time.Second,       // default 60 seconds
						Multiplier: 1.45,                   // default 1.3
					})
				}),
			},
		},
	}

	client, err := pubsub.NewClientWithConfig(ctx, projectID, config)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	// Make sure to reuse this publisher across publishes.
	p := client.Publisher(topicID)
	result := p.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return nil
}

// [END pubsub_publisher_retry_settings]
