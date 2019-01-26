// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mqttsnippets

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

var privateKeyRSA = os.Getenv("GOLANG_SAMPLES_IOT_PRIV")
var pubKeyRSA = os.Getenv("GOLANG_SAMPLES_IOT_PUB")
var region = "us-central1"

// returns a v1 UUID for a resource: e.g. topic, registry, gateway, device
func createIDForTest(resource string) string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Could not generate uuid: %v", err)
	}
	id := fmt.Sprintf("golang-test-%s-%s", resource, uuid.String())

	return id
}

// getClient returns a client based on the environment variable GOOGLE_APPLICATION_CREDENTIALS
func getClient() (*cloudiot.Service, error) {
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	client, err := cloudiot.New(httpClient)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// [END iot_get_client]

// Registry Management

// [START iot_create_registry]

// createRegistry creates a IoT Core device registry associated with a PubSub topic
func createRegistry(w io.Writer, projectID string, region string, registryID string, topicName string) (*cloudiot.DeviceRegistry, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	registry := cloudiot.DeviceRegistry{
		Id: registryID,
		EventNotificationConfigs: []*cloudiot.EventNotificationConfig{
			{
				PubsubTopicName: topicName,
			},
		},
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	response, err := client.Projects.Locations.Registries.Create(parent, &registry).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Created registry:")
	fmt.Fprintf(w, "\tID: %s\n", response.Id)
	fmt.Fprintf(w, "\tHTTP: %s\n", response.HttpConfig.HttpEnabledState)
	fmt.Fprintf(w, "\tMQTT: %s\n", response.MqttConfig.MqttEnabledState)
	fmt.Fprintf(w, "\tName: %s\n", response.Name)

	return response, nil
}

// [END iot_create_registry]

// [START iot_create_device]

// createDevice creates a device in a registry with one of the following public key formats
// RSA_PEM, RSA_X509_PEM, ES256_PEM, ES256_X509_PEM
func createDevice(w io.Writer, projectID string, region string, registryID string, deviceID string, publicKeyFormat string, keyPath string) (*cloudiot.Device, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	var device cloudiot.Device

	// if no credentials are passed in, create an unauth device
	if publicKeyFormat == "" {
		device = cloudiot.Device{
			Id: deviceID,
		}
	} else {
		device = cloudiot.Device{
			Id: deviceID,
			Credentials: []*cloudiot.DeviceCredential{
				{
					PublicKey: &cloudiot.PublicKeyCredential{
						Format: publicKeyFormat,
						Key:    string(keyBytes),
					},
				},
			},
		}
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully created a device with %s public key: %s", publicKeyFormat, deviceID)

	return response, nil
}

// [END iot_create_device]

// [START iot_delete_device]

// deleteDevice deletes a device from a registry.
func deleteDevice(w io.Writer, projectID string, region string, registryID string, deviceID string) (*cloudiot.Empty, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.Delete(path).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Deleted device: %s\n", deviceID)

	return response, nil
}

// [END iot_delete_device]

// [START iot_create_gateway]

// createGateway creates a new IoT Core gateway with a given id, public key, and auth method.
// gatewayAuthMethod can be one of: ASSOCIATION_ONLY, DEVICE_AUTH_TOKEN_ONLY, ASSOCIATION_AND_DEVICE_AUTH_TOKEN.
// https://cloud.google.com/iot/docs/reference/cloudiot/rest/v1/projects.locations.registries.devices#gatewayauthmethod
func createGateway(w io.Writer, projectID string, region string, registryID string, gatewayID string, gatewayAuthMethod string, publicKeyPath string) (*cloudiot.Device, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	keyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	gateway := &cloudiot.Device{
		Id: gatewayID,
		Credentials: []*cloudiot.DeviceCredential{
			{
				PublicKey: &cloudiot.PublicKeyCredential{
					Format: "RSA_X509_PEM",
					Key:    string(keyBytes),
				},
			},
		},
		GatewayConfig: &cloudiot.GatewayConfig{
			GatewayType:       "GATEWAY",
			GatewayAuthMethod: gatewayAuthMethod,
		},
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, gateway).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Successfully created gateway:", gatewayID)

	return response, nil
}

// [END iot_create_gateway]

// [START iot_bind_device_to_gateway]

// bindDeviceToGateway creates an association between an existing device and gateway.
func bindDeviceToGateway(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string) (*cloudiot.BindDeviceToGatewayResponse, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	bindRequest := &cloudiot.BindDeviceToGatewayRequest{
		DeviceId:  deviceID,
		GatewayId: gatewayID,
	}

	response, err := client.Projects.Locations.Registries.BindDeviceToGateway(parent, bindRequest).Do()

	if err != nil {
		return nil, fmt.Errorf("BindDeviceToGateway: %v", err)
	}

	if response.HTTPStatusCode/100 != 2 {
		return nil, fmt.Errorf("BindDeviceToGateway: HTTP status code not 2xx\n %v", response)
	}

	fmt.Fprintf(w, "Bound %s to %s", deviceID, gatewayID)

	return response, nil
}

// [END iot_bind_device_to_gateway]
// [START unbind_device_from_gateway]

// unbindDeviceFromGateway unbinds a bound device from a gateway.
func unbindDeviceFromGateway(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string) (*cloudiot.UnbindDeviceFromGatewayResponse, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	unbindRequest := &cloudiot.UnbindDeviceFromGatewayRequest{
		DeviceId:  deviceID,
		GatewayId: gatewayID,
	}

	response, err := client.Projects.Locations.Registries.UnbindDeviceFromGateway(parent, unbindRequest).Do()

	if err != nil {
		return nil, fmt.Errorf("UnbindDeviceFromGateway error: %v", err)
	}

	if response.HTTPStatusCode/100 != 2 {
		return nil, fmt.Errorf("UnbindDeviceFromGateway: HTTP status code not 2xx\n %v", response)
	}

	fmt.Fprintf(w, "Unbound %s from %s", deviceID, gatewayID)

	return response, nil
}

// [END unbind_device_from_gateway]

func TestSendDataFromBoundDevice(t *testing.T) {
	projectID := testutil.SystemTest(t).ProjectID

	registryID := "golang-iot-test-registry"
	gatewayID := createIDForTest("gateway")
	deviceID := createIDForTest("device")

	testutil.Retry(t, 1, 10*time.Second, func(r *testutil.R) {
		if _, err := createGateway(ioutil.Discard, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA); err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
			return
		}

		if _, err := createDevice(ioutil.Discard, projectID, region, registryID, deviceID, "RSA_X509_PEM", pubKeyRSA); err != nil {
			r.Errorf("Could not create device: %v\n", err)
			return
		}

		if _, err := bindDeviceToGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID); err != nil {
			r.Errorf("Could not bind device to gateway: %v\n", err)
			return
		}

		buf := new(bytes.Buffer)
		sendDataFromBoundDevice(buf, projectID, region, registryID, gatewayID, deviceID, privateKeyRSA, "RS256", 2, "test")

		want := fmt.Sprintf("Publishing message: %s", "test")
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("SendDataFromBoundDevice got %s, want substring %q", got, want)
		}
	})

	// cleanup
	unbindDeviceFromGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID)
	deleteDevice(ioutil.Discard, projectID, region, registryID, deviceID)
	deleteDevice(ioutil.Discard, projectID, region, registryID, gatewayID)
}
