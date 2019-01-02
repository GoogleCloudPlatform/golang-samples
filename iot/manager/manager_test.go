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
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

var projectID string
var topicID string
var topicName string // topicName is the full path to the topic (e.g. project/{project}/topics/{topic})
var registryID string

const region = "us-central1"
const pubKeyRSA = "./resources/rsa_cert.pem"

var client *pubsub.Client

// returns a v1 UUID for a resource: e.g. topic, registry, gateway, device
func createIDForTest(resource string) string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Could not generate uuid: %v", err)
	}
	id := fmt.Sprintf("golang-test-%s-%s", resource, uuid.String())

	return id
}

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

	topicID = createIDForTest("topic")

	topic, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		log.Fatalf("Could not create topic: %v", err)
	}
	fmt.Printf("Topic created: %v\n", topic)
	topicName = topic.String()

	// Generate UUID v1 for registry used for tests
	registryID = createIDForTest("registry")

	_, err = createRegistry(os.Stdout, projectID, region, registryID, topicName)
	if err != nil {
		log.Fatalf("Could not create registry: %v\n", err)
	}
}

func shutdown() {
	ctx := context.Background()

	t := client.Topic(topicID)
	if err := t.Delete(ctx); err != nil {
		log.Fatalf("Could not delete topic: %v\n", err)
	}
	fmt.Printf("Deleted topic: %v\n", t)

	if _, err := deleteRegistry(os.Stdout, projectID, region, registryID); err != nil {
		log.Fatalf("Could not delete registry: %v\n", err)
	}
}

func TestCreateRegistry(t *testing.T) {
	testRegistryID := createIDForTest("registry")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		registry, err := createRegistry(buf, projectID, region, testRegistryID, topicName)
		if err != nil {
			r.Errorf("Could not create registry: %v\n", err)

		} else {
			if registry.Id != testRegistryID {
				r.Errorf("Created registry, but registryID is wrong. Got %q, want %q", registry.Id, testRegistryID)
			}

			got := buf.String()
			want := fmt.Sprintf("Created registry:\n\tID: %s", testRegistryID)
			if !strings.Contains(got, want) {
				r.Errorf("CreateRegistry got %s, want substring %q", got, want)
			}

			deleteRegistry(buf, projectID, region, testRegistryID)
		}
	})
}

func TestGetRegistry(t *testing.T) {
	testRegistryID := createIDForTest("registry")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := createRegistry(buf, projectID, region, testRegistryID, topicName)
		if err != nil {
			r.Errorf("Could not create registry: %v\n", err)
		}

		_, err = getRegistry(buf, projectID, region, testRegistryID)
		if err != nil {
			r.Errorf("Could not get registry: %v\n", err)
		}

		got := buf.String()
		want := "Got registry:\n\tID: " + testRegistryID
		if !strings.Contains(got, want) {
			r.Errorf("GetRegistry got %s, want substring %q", got, want)
		}

		deleteRegistry(buf, projectID, region, testRegistryID)
	})
}

func TestListRegistries(t *testing.T) {
	testRegistryID := createIDForTest("registry")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := createRegistry(buf, projectID, region, testRegistryID, topicName)
		if err != nil {
			r.Errorf("Could not create registry 1: %v\n", err)
		}

		_, err = listRegistries(buf, projectID, region)
		if err != nil {
			r.Errorf("Could not list registries: %v\n", err)
		}

		got := buf.String()
		want := testRegistryID + "\n"
		if !strings.Contains(got, want) {
			r.Errorf("listRegistries got:\n %s, want substring: \n %q", got, want)
		}

		deleteRegistry(buf, projectID, region, testRegistryID)
	})
}

func TestDeleteRegistry(t *testing.T) {
	testRegistryID := createIDForTest("registry")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := createRegistry(buf, projectID, region, testRegistryID, topicName)
		if err != nil {
			r.Errorf("Could not create registry: %v\n", err)
		}

		_, err = deleteRegistry(buf, projectID, region, testRegistryID)
		if err != nil {
			r.Errorf("Could not delete registry: %v\n", err)
		}
		got := buf.String()
		want := "Deleted registry: " + testRegistryID

		if !strings.Contains(got, want) {
			r.Errorf("deleteRegistry got %s, want substring %q", got, want)
		}
	})
}

func TestSendCommand(t *testing.T) {
	deviceID := createIDForTest("device")
	buf := new(bytes.Buffer)

	if _, err := createUnauth(buf, projectID, region, registryID, deviceID); err != nil {
		t.Fatalf("Could not create device: %v", err)
	}

	commandToSend := "test"

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := sendCommand(buf, projectID, region, registryID, deviceID, commandToSend)

		// Currently, there is no Go client to receive commands so instead test for the "not subscribed" message
		if err == nil {
			r.Errorf("Should not be able to send command")
		}

		if !strings.Contains(err.Error(), "is not subscribed to the commands topic") {
			r.Errorf("Should create an error that device is not subscribed: %v", err)
		}
	})

	deleteDevice(buf, projectID, region, registryID, deviceID)
}

func TestCreateGateway(t *testing.T) {
	gatewayID := createIDForTest("gateway")

	testutil.Retry(t, 1, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := createGateway(buf, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA)

		if err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
		}

		got := buf.String()
		want := "Successfully created gateway: " + gatewayID

		if !strings.Contains(got, want) {
			r.Errorf("CreateGateway got %s, want substring %q", got, want)
		}

		deleteDevice(buf, projectID, region, registryID, gatewayID)
	})
}

func TestListGateways(t *testing.T) {
	// list zero gateways for initial registry
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := listGateways(buf, projectID, region, registryID)
		if err != nil {
			r.Errorf("Could not list gateways: %v\v", err)
		}

		got := buf.String()
		want := "No gateways found\n"

		if !strings.Contains(got, want) {
			r.Errorf("ListGateways got %s, want substring %q", got, want)
		}
	})

	// create and list gateway
	gatewayID := createIDForTest("gateway")
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		_, err := createGateway(buf, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
		}

		_, err = listGateways(buf, projectID, region, registryID)

		got := buf.String()
		want := gatewayID + "\n"

		if !strings.Contains(got, want) {
			r.Errorf("ListGateways got %s, want substring %q", got, want)
		}

		deleteDevice(buf, projectID, region, registryID, gatewayID)
	})
}

func TestBindDeviceToGateway(t *testing.T) {
	gatewayID := createIDForTest("gateway")
	deviceID := createIDForTest("device")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := createGateway(buf, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
		}

		_, err = createRSA(buf, projectID, region, registryID, deviceID, pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create device: %v\n", err)
		}

		_, err = bindDeviceToGateway(buf, projectID, region, registryID, gatewayID, deviceID)
		if err != nil {
			r.Errorf("Could not bind device to gateway: %v\n", err)
		}

		got := buf.String()
		want := fmt.Sprintf("Bound %s to %s", deviceID, gatewayID)

		if !strings.Contains(got, want) {
			r.Errorf("BindDeviceToGateway got %s, want substring %q", got, want)
		}

		unbindDeviceFromGateway(buf, projectID, region, registryID, gatewayID, deviceID)
		deleteDevice(buf, projectID, region, registryID, deviceID)
		deleteDevice(buf, projectID, region, registryID, gatewayID)
	})
}

func TestUnbindDeviceFromGateway(t *testing.T) {
	gatewayID := createIDForTest("gateway")
	deviceID := createIDForTest("device")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := createGateway(buf, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
		}

		_, err = createRSA(buf, projectID, region, registryID, deviceID, pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create device: %v\n", err)
		}

		_, err = bindDeviceToGateway(buf, projectID, region, registryID, gatewayID, deviceID)
		if err != nil {
			r.Errorf("Could not bind device to gateway: %v\n", err)
		}

		_, err = unbindDeviceFromGateway(buf, projectID, region, registryID, gatewayID, deviceID)
		if err != nil {
			r.Errorf("Could not unbind device to gateway: %v\n", err)
		}

		got := buf.String()
		want := fmt.Sprintf("Unbound %s from %s", deviceID, gatewayID)

		if !strings.Contains(got, want) {
			r.Errorf("CreateGateway got %s, want substring %q", got, want)
		}

		deleteDevice(buf, projectID, region, registryID, deviceID)
		deleteDevice(buf, projectID, region, registryID, gatewayID)
	})
}
func TestListDevicesForGateway(t *testing.T) {
	gatewayID := createIDForTest("gateway")
	deviceID := createIDForTest("device")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		_, err := createGateway(buf, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
		}

		_, err = createRSA(buf, projectID, region, registryID, deviceID, pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create device: %v\n", err)
		}

		bindDeviceToGateway(buf, projectID, region, registryID, gatewayID, deviceID)

		_, err = listDevicesForGateway(buf, projectID, region, registryID, gatewayID)

		got := buf.String()
		want := fmt.Sprintf("Devices for %s", gatewayID)
		if !strings.Contains(got, want) {
			r.Errorf("ListDeviesForGateway got %s, want substring %q", got, want)
		}

		want = fmt.Sprintf("\t%s\n", deviceID)
		if !strings.Contains(got, want) {
			r.Errorf("ListDeviesForGateway got %s, want substring %q", got, want)
		}

		unbindDeviceFromGateway(buf, projectID, region, registryID, gatewayID, deviceID)
		deleteDevice(buf, projectID, region, registryID, deviceID)
		deleteDevice(buf, projectID, region, registryID, gatewayID)
	})
}
