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

// Command manager lets you manage Cloud IoT Core devices and registries.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	// [START imports]
	"context"
	b64 "encoding/base64"

	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
	// [END imports]
)

// [START iot_get_client]

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

// [START iot_delete_registry]

// deleteRegistry deletes a device registry if it is empty.
func deleteRegistry(w io.Writer, projectID string, region string, registryID string) (*cloudiot.Empty, error) {
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

	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Delete(name).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Deleted registry: %s\n", registryID)

	return response, nil
}

// [END iot_delete_registry]

// [START iot_get_registry]

// getRegistry gets information about a device registry given a registryID.
func getRegistry(w io.Writer, projectID string, region string, registryID string) (*cloudiot.DeviceRegistry, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Get(parent).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Got registry:")
	fmt.Fprintf(w, "\tID: %s\n", response.Id)
	fmt.Fprintf(w, "\tHTTP: %s\n", response.HttpConfig.HttpEnabledState)
	fmt.Fprintf(w, "\tMQTT: %s\n", response.MqttConfig.MqttEnabledState)
	fmt.Fprintf(w, "\tName: %s\n", response.Name)

	return response, nil
}

// [END iot_get_registry]

// [START iot_list_registries]

// listRegistries gets the names of device registries given a project / region.
func listRegistries(w io.Writer, projectID string, region string) ([]*cloudiot.DeviceRegistry, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	response, err := client.Projects.Locations.Registries.List(parent).Do()
	if err != nil {
		return nil, err
	}

	if len(response.DeviceRegistries) == 0 {
		fmt.Fprintln(w, "No registries found")
		return response.DeviceRegistries, nil
	}

	fmt.Fprintf(w, "%d registries:\n", len(response.DeviceRegistries))
	for _, registry := range response.DeviceRegistries {
		fmt.Fprintf(w, "\t%s\n", registry.Name)
	}

	return response.DeviceRegistries, nil
}

// [END iot_list_registries]

// [START iot_get_iam_policy]

// getRegistryIAM gets the IAM policy for a device registry.
func getRegistryIAM(w io.Writer, projectID string, region string, registryID string) (*cloudiot.Policy, error) {
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

	var req cloudiot.GetIamPolicyRequest

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.GetIamPolicy(path, &req).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Policy:")
	for _, binding := range response.Bindings {
		fmt.Fprintf(w, "Role: %s\n", binding.Role)
		for _, member := range binding.Members {
			fmt.Fprintf(w, "\tMember: %s\n", member)
		}
	}

	return response, nil
}

// [END iot_get_iam_policy]

// [START iot_set_iam_policy]

// setRegistryIAM sets the IAM policy for a device registry
func setRegistryIAM(w io.Writer, projectID string, region string, registryID string, member string, role string) (*cloudiot.Policy, error) {
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

	req := cloudiot.SetIamPolicyRequest{
		Policy: &cloudiot.Policy{
			Bindings: []*cloudiot.Binding{
				{
					Members: []string{member},
					Role:    role,
				},
			},
		},
	}
	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.SetIamPolicy(path, &req).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully set IAM policy for registry: %s\n", registryID)

	return response, nil
}

// [END iot_set_iam_policy]

// Device Management

// [START iot_create_es_device]

// createES creates a device in a registry with ES256 credentials.
func createES(w io.Writer, projectID string, region string, registryID string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	device := cloudiot.Device{
		Id: deviceID,
		Credentials: []*cloudiot.DeviceCredential{
			{
				PublicKey: &cloudiot.PublicKeyCredential{
					Format: "ES256_PEM",
					Key:    string(keyBytes),
				},
			},
		},
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully created ES256 device: %s\n", deviceID)

	return response, nil
}

// [END iot_create_es_device]

// [START iot_create_rsa_device]

// createRSA creates a device in a registry given RSA X.509 credentials.
func createRSA(w io.Writer, projectID string, region string, registryID string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	device := cloudiot.Device{
		Id: deviceID,
		Credentials: []*cloudiot.DeviceCredential{
			{
				PublicKey: &cloudiot.PublicKeyCredential{
					Format: "RSA_X509_PEM",
					Key:    string(keyBytes),
				},
			},
		},
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully created RSA256 X.509 device: %s", deviceID)

	return response, nil
}

// [END iot_create_rsa_device]

// [START iot_create_unauth_device]

// createUnauth creates a device in a registry without credentials.
func createUnauth(w io.Writer, projectID string, region string, registryID string, deviceID string) (*cloudiot.Device, error) {
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

	device := cloudiot.Device{
		Id: deviceID,
	}
	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully created device without credentials: %s\n", deviceID)

	return response, nil
}

// [END iot_create_unauth_device]

// [START iot_create_device]

// createDevice creates a device in a registry with one of the following public key formats:
// RSA_PEM, RSA_X509_PEM, ES256_PEM, ES256_X509_PEM, UNAUTH.
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

	// If no credentials are passed in, create an unauth device.
	if publicKeyFormat == "UNAUTH" {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.Delete(path).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Deleted device: %s\n", deviceID)

	return response, nil
}

// [END iot_delete_device]

// [START iot_get_device]

// getDevice retrieves a specific device and prints its details.
func getDevice(w io.Writer, projectID string, region string, registryID string, device string) (*cloudiot.Device, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, device)
	response, err := client.Projects.Locations.Registries.Devices.Get(path).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "\tId: %s\n", response.Id)
	for _, credential := range response.Credentials {
		fmt.Fprintf(w, "\t\tCredential Expire: %s\n", credential.ExpirationTime)
		fmt.Fprintf(w, "\t\tCredential Type: %s\n", credential.PublicKey.Format)
		fmt.Fprintln(w, "\t\t--------")
	}
	fmt.Fprintf(w, "\tLast Config Ack: %s\n", response.LastConfigAckTime)
	fmt.Fprintf(w, "\tLast Config Send: %s\n", response.LastConfigSendTime)
	fmt.Fprintf(w, "\tLast Event Time: %s\n", response.LastEventTime)
	fmt.Fprintf(w, "\tLast Heartbeat Time: %s\n", response.LastHeartbeatTime)
	fmt.Fprintf(w, "\tLast State Time: %s\n", response.LastStateTime)
	fmt.Fprintf(w, "\tNumId: %d\n", response.NumId)

	return response, nil
}

// [END iot_get_device]

// [START iot_get_device_configs]

// getDeviceConfigs retrieves and lists device configurations.
func getDeviceConfigs(w io.Writer, projectID string, region string, registryID string, device string) ([]*cloudiot.DeviceConfig, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, device)
	response, err := client.Projects.Locations.Registries.Devices.ConfigVersions.List(path).Do()
	if err != nil {
		return nil, err
	}

	for _, config := range response.DeviceConfigs {
		fmt.Fprintf(w, "%d : %s\n", config.Version, config.BinaryData)
	}

	return response.DeviceConfigs, nil
}

// [END iot_get_device_configs]

// [START iot_get_device_state]

// getDeviceStates retrieves and lists device states.
func getDeviceStates(w io.Writer, projectID string, region string, registryID string, device string) ([]*cloudiot.DeviceState, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, device)
	response, err := client.Projects.Locations.Registries.Devices.States.List(path).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Successfully retrieved device states!")

	for _, state := range response.DeviceStates {
		fmt.Fprintf(w, "%s : %s\n", state.UpdateTime, state.BinaryData)
	}

	return response.DeviceStates, nil
}

// [END iot_get_device_state]

// [START iot_list_devices]

// listDevices gets the identifiers of devices for a specific registry.
func listDevices(w io.Writer, projectID string, region string, registryID string) ([]*cloudiot.Device, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.List(parent).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Devices:")
	for _, device := range response.Devices {
		fmt.Fprintf(w, "\t%s\n", device.Id)
	}

	return response.Devices, nil
}

// [END iot_list_devices]

// [START iot_patch_es]

// patchDeviceES patches a device to use ES256 credentials.
func patchDeviceES(w io.Writer, projectID string, region string, registryID string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	device := cloudiot.Device{
		Id: deviceID,
		Credentials: []*cloudiot.DeviceCredential{
			{
				PublicKey: &cloudiot.PublicKeyCredential{
					Format: "ES256_PEM",
					Key:    string(keyBytes),
				},
			},
		},
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.
		Patch(parent, &device).UpdateMask("credentials").Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Successfully patched device with ES256 credentials")

	return response, nil
}

// [END iot_patch_es]

// [START iot_patch_rsa]

// patchDeviceRSA patches a device to use RSA256 X.509 credentials.
func patchDeviceRSA(w io.Writer, projectID string, region string, registryID string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	device := cloudiot.Device{
		Id: deviceID,
		Credentials: []*cloudiot.DeviceCredential{
			{
				PublicKey: &cloudiot.PublicKeyCredential{
					Format: "RSA_X509_PEM",
					Key:    string(keyBytes),
				},
			},
		},
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.
		Patch(parent, &device).UpdateMask("credentials").Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Successfully patched device with RSA256 X.509 credentials")

	return response, nil
}

// [END iot_patch_rsa]

// [START iot_set_device_config]

// setConfig sends a configuration change to a device.
func setConfig(w io.Writer, projectID string, region string, registryID string, deviceID string, configData string) (*cloudiot.DeviceConfig, error) {
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

	req := cloudiot.ModifyCloudToDeviceConfigRequest{
		BinaryData: b64.StdEncoding.EncodeToString([]byte(configData)),
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.ModifyCloudToDeviceConfig(path, &req).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Config set!\nVersion now: %d\n", response.Version)

	return response, nil
}

// [END iot_set_device_config]

// [START iot_send_command]

// sendCommand sends a command to a device listening for commands.
func sendCommand(w io.Writer, projectID string, region string, registryID string, deviceID string, sendData string) (*cloudiot.SendCommandToDeviceResponse, error) {
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

	req := cloudiot.SendCommandToDeviceRequest{
		BinaryData: b64.StdEncoding.EncodeToString([]byte(sendData)),
	}

	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, deviceID)

	response, err := client.Projects.Locations.Registries.Devices.SendCommandToDevice(name, &req).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Sent command to device")

	return response, nil
}

// [END iot_send_command]

// START BETA FEATURES

// [START iot_create_gateway]

// createGateway creates a new IoT Core gateway with a given id, public key, and auth method.
// gatewayAuthMethod can be one of: ASSOCIATION_ONLY, DEVICE_AUTH_TOKEN_ONLY, ASSOCIATION_AND_DEVICE_AUTH_TOKEN.
// https://cloud.google.com/iot/docs/reference/cloudiot/rest/v1/projects.locations.registries.devices#gatewayauthmethod
func createGateway(w io.Writer, projectID string, region string, registryID string, gatewayID string, gatewayAuthMethod string, publicKeyPath string) (*cloudiot.Device, error) {
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

// [START iot_list_gateways]

// listGateways lists all the gateways in a specific registry.
func listGateways(w io.Writer, projectID string, region string, registryID string) ([]*cloudiot.Device, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.List(parent).GatewayListOptionsGatewayType("GATEWAY").Do()

	if err != nil {
		return nil, fmt.Errorf("ListGateways: %v", err)
	}

	if len(response.Devices) == 0 {
		fmt.Fprintln(w, "No gateways found")
		return response.Devices, nil
	}

	fmt.Fprintln(w, len(response.Devices), "devices:")
	for _, gateway := range response.Devices {
		fmt.Fprintf(w, "\t%s\n", gateway.Id)
	}

	return response.Devices, nil
}

// [END iot_list_gateways]

// [START iot_bind_device_to_gateway]

// bindDeviceToGateway creates an association between an existing device and gateway.
func bindDeviceToGateway(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string) (*cloudiot.BindDeviceToGatewayResponse, error) {
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

// [START list_devices_for_gateway]

// listDevicesForGateway lists the devices that are bound to a gateway.
func listDevicesForGateway(w io.Writer, projectID string, region string, registryID, gatewayID string) ([]*cloudiot.Device, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.List(parent).GatewayListOptionsAssociationsGatewayId(gatewayID).Do()

	if err != nil {
		return nil, fmt.Errorf("ListDevicesForGateway: %v", err)
	}

	if len(response.Devices) == 0 {
		fmt.Fprintln(w, "\tNo devices found")
		return response.Devices, nil
	}

	fmt.Fprintf(w, "Devices for %s:\n", gatewayID)
	for _, gateway := range response.Devices {
		fmt.Fprintf(w, "\t%s\n", gateway.Id)
	}

	return response.Devices, nil
}

// [END list_devices_for_gateway]

// END BETA FEATURES

type command struct {
	name string
	fn   interface{}
	args []string
}

func (c command) usage() string {
	var buf bytes.Buffer
	buf.WriteString(c.name)
	buf.WriteString(" ")
	for _, arg := range c.args {
		buf.WriteString("<")
		buf.WriteString(arg)
		buf.WriteString("> ")
	}
	buf.UnreadByte()
	return buf.String()
}

func main() {
	registryManagementCommands := []command{
		{"createRegistry", createRegistry, []string{"cloud-region", "registry-id", "pubsub-topic"}},
		{"deleteRegistry", deleteRegistry, []string{"cloud-region", "registry-id"}},
		{"getRegistry", getRegistry, []string{"cloud-region", "registry-id"}},
		{"listRegistries", listRegistries, []string{"cloud-region"}},
		{"getRegistryIAM", getRegistryIAM, []string{"cloud-region", "registry-id"}},
		{"setRegistryIAM", setRegistryIAM, []string{"cloud-region", "registry-id", "member", "role"}},
	}

	deviceManagementCommands := []command{
		{"createES", createES, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"createRSA", createRSA, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"createUnauth", createUnauth, []string{"cloud-region", "registry-id", "device-id"}},
		{"createDevice", createDevice, []string{"cloud-region", "registry-id", "device-id", "public-key-format", "keyfile-path"}},
		{"deleteDevice", deleteDevice, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDevice", getDevice, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDeviceConfigs", getDeviceConfigs, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDeviceStates", getDeviceStates, []string{"cloud-region", "registry-id", "device-id"}},
		{"listDevices", listDevices, []string{"cloud-region", "registry-id"}},
		{"patchDevice", patchDeviceES, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"patchDeviceRSA", patchDeviceRSA, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"setConfig", setConfig, []string{"cloud-region", "registry-id", "device-id", "config-data"}},
		{"sendCommand", sendCommand, []string{"cloud-region", "registry-id", "device-id", "send-data"}},
	}

	// Beta Features: Gateway management commands
	gatewayManagementCommands := []command{
		{"createGateway", createGateway, []string{"cloud-region", "registry-id", "gateway-id", "auth-method", "public-key-path"}},
		{"listGateways", listGateways, []string{"cloud-region", "registry-id"}},
		{"bindDeviceToGateway", bindDeviceToGateway, []string{"cloud-region", "registry-id", "gateway-id", "device-id"}},
		{"unbindDeviceFromGateway", unbindDeviceFromGateway, []string{"cloud-region", "registry-id", "gateway-id", "device-id"}},
		{"listDevicesForGateway", listDevicesForGateway, []string{"cloud-region", "registry-id", "gateway-id"}},
	}

	var commands []command
	commands = append(commands, registryManagementCommands...)
	commands = append(commands, deviceManagementCommands...)
	commands = append(commands, gatewayManagementCommands...)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\tRegistry Management\n")
		fmt.Fprintf(os.Stderr, "\t-----\n")
		for _, cmd := range registryManagementCommands {
			fmt.Fprintf(os.Stderr, "\t%s %s\n", filepath.Base(os.Args[0]), cmd.usage())
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "\tDevice Management\n")
		fmt.Fprintf(os.Stderr, "\t-----\n")
		for _, cmd := range deviceManagementCommands {
			fmt.Fprintf(os.Stderr, "\t%s %s\n", filepath.Base(os.Args[0]), cmd.usage())
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "\tBeta Gateway Management\n")
		fmt.Fprintf(os.Stderr, "\t-----\n")
		for _, cmd := range gatewayManagementCommands {
			fmt.Fprintf(os.Stderr, "\t%s %s\n", filepath.Base(os.Args[0]), cmd.usage())
		}
	}
	flag.Parse()

	// Retrieve project ID from console.
	projectID := os.Getenv("GCLOUD_PROJECT")
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if projectID == "" {
		fmt.Fprintln(os.Stderr, "Set the GCLOUD_PROJECT or GOOGLE_CLOUD_PROJECT environment variable.")
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	commandName := flag.Args()[0]
	commandArgs := flag.Args()[1:]

	for _, cmd := range commands {
		if cmd.name == commandName {
			if len(commandArgs[1:]) != len(cmd.args)-1 {
				fmt.Fprintf(os.Stderr, "Wrong number of arguments. Usage:\n\t%s\n", cmd.usage())
				os.Exit(1)
			}
			var fnArgs []reflect.Value

			fnArgs = append(fnArgs, reflect.ValueOf(os.Stdout))
			fnArgs = append(fnArgs, reflect.ValueOf(projectID))
			for _, arg := range commandArgs {
				fnArgs = append(fnArgs, reflect.ValueOf(arg))
			}
			retValues := reflect.ValueOf(cmd.fn).Call(fnArgs)
			err := retValues[len(retValues)-1]
			if !err.IsNil() {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	// Unknown command
	flag.Usage()
	os.Exit(1)
}
