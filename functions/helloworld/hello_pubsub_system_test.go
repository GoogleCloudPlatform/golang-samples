// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build ignore
// Disabled until system tests are working on Kokoro.

// TODO: Use testutil.SystemTest and os.Setenv for GOOGLE_CLOUD_PROJECT.

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
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

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
