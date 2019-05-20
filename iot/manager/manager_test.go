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

var (
	projectID  string
	topicID    string
	topicName  string // topicName is the full path to the topic (e.g. project/{project}/topics/{topic}).
	registryID string
	client     *pubsub.Client
)

const region = "us-central1"

var pubKeyRSA = os.Getenv("GOLANG_SAMPLES_IOT_PUB")

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
	_, ok := testutil.ContextMain(m)
	if !ok {
		log.Print("GOLANG_SAMPLES_PROJECT_ID is unset. Skipping.")
		return
	}
	if err := setup(m); err != nil {
		log.Fatal(err)
	}
	s := m.Run()
	shutdown()
	os.Exit(s)
}

func setup(m *testing.M) error {
	ctx := context.Background()
	tc, _ := testutil.ContextMain(m)
	projectID = tc.ProjectID

	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	client = pubsubClient

	topicID = createIDForTest("topic")

	topic, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		return err
	}
	fmt.Printf("Topic created: %v\n", topic)
	topicName = topic.String()

	// Generate UUID v1 for registry used for tests
	registryID = createIDForTest("registry")

	if _, err := createRegistry(os.Stdout, projectID, region, registryID, topicName); err != nil {
		return err
	}
	return nil
}

func shutdown() {
	ctx := context.Background()

	topic := client.Topic(topicID)
	if err := topic.Delete(ctx); err != nil {
		fmt.Printf("Could not delete topic: %v\n", err)
		return
	}
	fmt.Printf("Deleted topic: %v\n", topic)

	if _, err := deleteRegistry(os.Stdout, projectID, region, registryID); err != nil {
		fmt.Printf("Could not delete registry: %v\n", err)
	}
}

func TestCreateRegistry(t *testing.T) {
	testRegistryID := createIDForTest("registry")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		registry, err := createRegistry(buf, projectID, region, testRegistryID, topicName)
		if err != nil {
			r.Errorf("Could not create registry: %v\n", err)
			return
		}

		if registry.Id != testRegistryID {
			r.Errorf("Created registry, but registryID is wrong. Got %q, want %q", registry.Id, testRegistryID)
		}

		want := fmt.Sprintf("Created registry:\n\tID: %s", testRegistryID)
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("CreateRegistry got %s, want substring %q", got, want)
		}

		deleteRegistry(ioutil.Discard, projectID, region, testRegistryID)

	})
}

func TestGetRegistry(t *testing.T) {
	testRegistryID := createIDForTest("registry")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if _, err := createRegistry(ioutil.Discard, projectID, region, testRegistryID, topicName); err != nil {
			r.Errorf("Could not create registry: %v\n", err)
			return
		}

		if _, err := getRegistry(buf, projectID, region, testRegistryID); err != nil {
			r.Errorf("Could not get registry: %v\n", err)
		}

		want := "Got registry:\n\tID: " + testRegistryID
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("GetRegistry got %s, want substring %q", got, want)
		}

		deleteRegistry(ioutil.Discard, projectID, region, testRegistryID)
	})
}

func TestListRegistries(t *testing.T) {
	testRegistryID := createIDForTest("registry")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if _, err := createRegistry(ioutil.Discard, projectID, region, testRegistryID, topicName); err != nil {
			r.Errorf("Could not create registry 1: %v\n", err)
		}

		if _, err := listRegistries(buf, projectID, region); err != nil {
			r.Errorf("Could not list registries: %v\n", err)
		}

		want := testRegistryID + "\n"
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("listRegistries got:\n %s, want substring: \n %q", got, want)
		}

		deleteRegistry(ioutil.Discard, projectID, region, testRegistryID)
	})
}

func TestDeleteRegistry(t *testing.T) {
	testRegistryID := createIDForTest("registry")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if _, err := createRegistry(ioutil.Discard, projectID, region, testRegistryID, topicName); err != nil {
			r.Errorf("Could not create registry: %v\n", err)
		}

		if _, err := deleteRegistry(buf, projectID, region, testRegistryID); err != nil {
			r.Errorf("Could not delete registry: %v\n", err)
		}

		want := "Deleted registry: " + testRegistryID
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("deleteRegistry got %s, want substring %q", got, want)
		}
	})
}

func TestSendCommand(t *testing.T) {
	deviceID := createIDForTest("device")

	if _, err := createUnauth(ioutil.Discard, projectID, region, registryID, deviceID); err != nil {
		t.Fatalf("Could not create device: %v", err)
		return
	}

	commandToSend := "test"

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := sendCommand(buf, projectID, region, registryID, deviceID, commandToSend)

		// Currently, there is no Go client to receive commands so instead test for the "not subscribed" message
		if err == nil {
			r.Errorf("Should not be able to send command")
		}

		if !strings.Contains(err.Error(), "not connected") {
			r.Errorf("Should create an error that device is not connected: %v", err)
		}
	})

	deleteDevice(ioutil.Discard, projectID, region, registryID, deviceID)
}

func TestCreateGateway(t *testing.T) {
	if pubKeyRSA == "" {
		t.Skip("GOLANG_SAMPLES_IOT_PUB not set")
	}

	gatewayID := createIDForTest("gateway")

	testutil.Retry(t, 1, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if _, err := createGateway(buf, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA); err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
			return
		}

		want := fmt.Sprintf("Successfully created gateway: %s", gatewayID)
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("CreateGateway got %s, want substring %q", got, want)
		}

		deleteDevice(ioutil.Discard, projectID, region, registryID, gatewayID)
	})
}

func TestListGateways(t *testing.T) {
	if pubKeyRSA == "" {
		t.Skip("GOLANG_SAMPLES_IOT_PUB not set")
	}

	// list zero gateways for initial registry
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if _, err := listGateways(buf, projectID, region, registryID); err != nil {
			r.Errorf("Could not list gateways: %v\v", err)
			return
		}

		want := "No gateways found\n"
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("ListGateways got %s, want substring %q", got, want)
		}
	})

	// create and list gateway
	gatewayID := createIDForTest("gateway")
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		if _, err := createGateway(ioutil.Discard, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA); err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
			return
		}

		if _, err := listGateways(buf, projectID, region, registryID); err != nil {
			r.Errorf("ListGateways error: %v\n", err)
		}

		want := gatewayID + "\n"
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("ListGateways got %s, want substring %q", got, want)
		}

		deleteDevice(ioutil.Discard, projectID, region, registryID, gatewayID)
	})
}

func TestBindDeviceToGateway(t *testing.T) {
	if pubKeyRSA == "" {
		t.Skip("GOLANG_SAMPLES_IOT_PUB not set")
	}

	gatewayID := createIDForTest("gateway")
	deviceID := createIDForTest("device")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		if _, err := createGateway(ioutil.Discard, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA); err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
			return
		}

		if _, err := createRSA(ioutil.Discard, projectID, region, registryID, deviceID, pubKeyRSA); err != nil {
			r.Errorf("Could not create device: %v\n", err)
		}

		if _, err := bindDeviceToGateway(buf, projectID, region, registryID, gatewayID, deviceID); err != nil {
			r.Errorf("Could not bind device to gateway: %v\n", err)
		}

		want := fmt.Sprintf("Bound %s to %s", deviceID, gatewayID)
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("BindDeviceToGateway got %s, want substring %q", got, want)
		}

		unbindDeviceFromGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID)
		deleteDevice(ioutil.Discard, projectID, region, registryID, deviceID)
		deleteDevice(ioutil.Discard, projectID, region, registryID, gatewayID)
	})
}

func TestUnbindDeviceFromGateway(t *testing.T) {
	if pubKeyRSA == "" {
		t.Skip("GOLANG_SAMPLES_IOT_PUB not set")
	}

	gatewayID := createIDForTest("gateway")
	deviceID := createIDForTest("device")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		_, err := createGateway(ioutil.Discard, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
		}

		_, err = createRSA(ioutil.Discard, projectID, region, registryID, deviceID, pubKeyRSA)
		if err != nil {
			r.Errorf("Could not create device: %v\n", err)
		}

		_, err = bindDeviceToGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID)
		if err != nil {
			r.Errorf("Could not bind device to gateway: %v\n", err)
		}

		_, err = unbindDeviceFromGateway(buf, projectID, region, registryID, gatewayID, deviceID)
		if err != nil {
			r.Errorf("Could not unbind device to gateway: %v\n", err)
		}

		want := fmt.Sprintf("Unbound %s from %s", deviceID, gatewayID)
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("CreateGateway got %s, want substring %q", got, want)
		}

		deleteDevice(ioutil.Discard, projectID, region, registryID, deviceID)
		deleteDevice(ioutil.Discard, projectID, region, registryID, gatewayID)
	})
}
func TestListDevicesForGateway(t *testing.T) {
	if pubKeyRSA == "" {
		t.Skip("GOLANG_SAMPLES_IOT_PUB not set")
	}

	gatewayID := createIDForTest("gateway")
	deviceID := createIDForTest("device")

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		if _, err := createGateway(ioutil.Discard, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA); err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
			return
		}

		if _, err := createRSA(ioutil.Discard, projectID, region, registryID, deviceID, pubKeyRSA); err != nil {
			r.Errorf("Could not create device: %v\n", err)
		}

		if _, err := bindDeviceToGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID); err != nil {
			r.Errorf("Could not bind device to gateway: %v\n", err)
		}

		if _, err := listDevicesForGateway(buf, projectID, region, registryID, gatewayID); err != nil {
			r.Errorf("Could not list gateways")
		}

		got := buf.String()
		want := fmt.Sprintf("Devices for %s", gatewayID)
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("ListDeviesForGateway got %s, want substring %q", got, want)
		}

		want = fmt.Sprintf("\t%s\n", deviceID)
		if !strings.Contains(got, want) {
			r.Errorf("ListDeviesForGateway got %s, want substring %q", got, want)
		}

		// cleanup
		unbindDeviceFromGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID)
		deleteDevice(ioutil.Discard, projectID, region, registryID, deviceID)
		deleteDevice(ioutil.Discard, projectID, region, registryID, gatewayID)
	})
}
