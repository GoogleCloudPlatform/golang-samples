// Copyright 2018 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

var topicID string
var projectID string

var client *pubsub.Client

func TestMain(m *testing.M) {
	setup()
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	shutdown()
	os.Exit(s)
}

func setup() {
	ctx := context.Background()
	// Retrieve project ID from console
	projectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}
	os.Setenv("GOOGLE_CLOUD_PROJECT", projectID)

	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	client = pubsubClient

	pubsubUUID, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Could not generate uuid: %v", err)
	}
	topicID = "golang-iot-topic-" + pubsubUUID.String()

	t, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		log.Fatalf("Could not create topic: %v", err)
	}
	fmt.Printf("Topic created: %v\n", t)
}

func shutdown() {
	ctx := context.Background()

	t := client.Topic(topicID)
	if err := t.Delete(ctx); err != nil {
		log.Fatalf("Could not delete topic: %v", err)
	}
	fmt.Printf("Deleted topic: %v\n", t)
}

func TestSendCommand(t *testing.T) {
	testutil.SystemTest(t)

	// Generate UUID v1 for test registry and device
	registryUUID, _ := uuid.NewRandom()
	deviceUUID, _ := uuid.NewRandom()

	region := "us-central1"
	registryID := "golang-test-registry-" + registryUUID.String()
	deviceID := "golang-test-device-" + deviceUUID.String()

	topic := client.Topic(topicID)

	var buf bytes.Buffer

	if _, err := createRegistry(&buf, projectID, region, registryID, topic.String()); err != nil {
		log.Fatalf("Could not create registry: %v", err)
	}

	if _, err := createUnauth(&buf, projectID, region, registryID, deviceID); err != nil {
		log.Fatalf("Could not create device: %v", err)
	}

	commandToSend := "test"

	_, err := sendCommand(&buf, projectID, region, registryID, deviceID, commandToSend)

	// Currently, there is no Go client to receive commands so instead test for the "not subscribed" message
	if err == nil {
		t.Error("Should not be able to send command")
	}

	if !strings.Contains(err.Error(), "is not subscribed to the commands topic") {
		t.Error("Should create an error that device is not subscribed", err)
	}

	deleteDevice(&buf, projectID, region, registryID, deviceID)
	deleteRegistry(&buf, projectID, region, registryID)
}
