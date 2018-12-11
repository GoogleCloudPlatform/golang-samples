// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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

var projectID string
var topicID string

var client *pubsub.Client

func TestMain(m *testing.M) {
	setup(m)
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	shutdown()
	os.Exit(s)
}

func setup(m *testing.M) {
	ctx := context.Background()
	tc, ok := testutil.ContextMain(m)

	// Retrive
	if ok {
		projectID = tc.ProjectID
	} else {
		fmt.Fprintln(os.Stderr, "Project is not set up properly for system tests. Make sure GOLANG_SAMPLES_PROJECT_ID is set")
		os.Exit(1)
	}

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
