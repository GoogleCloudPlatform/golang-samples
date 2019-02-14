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

// +build ignore
// Disabled until system tests are working on Kokoro.

// [START functions_pubsub_system_test]

package helloworld

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gobuffalo/uuid"
)

func TestHelloPubSubSystem(t *testing.T) {
	ctx := context.Background()

	topicName := os.Getenv("FUNCTIONS_TOPIC")
	projectID := os.Getenv("GCP_PROJECT")

	startTime := time.Now().UTC().Format(time.RFC3339)

	// Create the Pub/Sub client and topic.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	topic := client.Topic(topicName)

	// Publish a message with a random string to verify.
	// We use a random string to make sure the function is logging the correct
	// message for this test invocation.
	u := uuid.Must(uuid.NewV4())
	msg := &pubsub.Message{
		Data: []byte(u.String()),
	}
	topic.Publish(ctx, msg).Get(ctx)

	// Wait for logs to be consistent.
	time.Sleep(20 * time.Second)

	// Check logs after a delay.
	cmd := exec.Command("gcloud", "alpha", "functions", "logs", "read", "HelloPubSub", "--start-time", startTime)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec.Command: %v", err)
	}
	if got := string(out); !strings.Contains(got, u.String()) {
		t.Errorf("HelloPubSub got %q, want to contain %q", got, u.String())
	}
}

// [END functions_pubsub_system_test]
