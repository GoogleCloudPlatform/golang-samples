// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package manager lets you manage Cloud IoT Core devices and registries.
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

// [START get_client]

// GetClient returns a client based on the environment variable GOOGLE_APPLICATION_CREDENTIALS
func GetClient() (*cloudiot.Service, error) {
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

// [END get_client]

// Registry Management

// [START iot_create_registry]

// CreateRegistry creates a IoT Core device registry associated with a PubSub topic
func CreateRegistry(w io.Writer, projectID string, region string, registryID string, topicName string) (*cloudiot.DeviceRegistry, error) {
	client, err := GetClient()
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

	return response, err
}

// [END iot_create_registry]

// [START iot_delete_registry]

// DeleteRegistry deletes a device registry if it is empty.
func DeleteRegistry(w io.Writer, projectID string, region string, registryID string) (*cloudiot.Empty, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Delete(name).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Deleted registry: %s\n", registryID)

	return response, err
}

// [END iot_delete_registry]

// [START iot_get_registry]

// GetRegistry gets information about a device registry given a registryID.
func GetRegistry(w io.Writer, projectID string, region string, registryID string) (*cloudiot.DeviceRegistry, error) {
	client, err := GetClient()
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

	return response, err
}

// [END iot_get_registry]

// [START iot_list_registries]

// ListRegistries gets the names of device registries given a project / region.
func ListRegistries(w io.Writer, projectID string, region string) ([]*cloudiot.DeviceRegistry, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	response, err := client.Projects.Locations.Registries.List(parent).Do()
	if err != nil {
		return nil, err
	}

	if len(response.DeviceRegistries) > 0 {
		fmt.Fprintf(w, "%d registries:\n", len(response.DeviceRegistries))
		for _, registry := range response.DeviceRegistries {
			fmt.Fprintf(w, "\t%s\n", registry.Name)
		}
	} else {
		fmt.Fprintln(w, "No registries found")
	}

	return response.DeviceRegistries, err
}

// [END iot_list_registries]

// [START iot_get_iam_policy]

// GetRegistryIAM gets the IAM policy for a device registry.
func GetRegistryIAM(w io.Writer, projectID string, region string, registryID string) (*cloudiot.Policy, error) {
	client, err := GetClient()
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

	return response, err
}

// [END iot_get_iam_policy]

// [START iot_set_iam_policy]

// SetRegistryIAM sets the IAM policy for a device registry
func SetRegistryIAM(w io.Writer, projectID string, region string, registryID string, member string, role string) (*cloudiot.Policy, error) {
	client, err := GetClient()
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

	return response, err
}

// [END iot_set_iam_policy]

// Device Management

// [START iot_create_es_device]

// CreateES creates a device in a registry with ES256 credentials.
func CreateES(w io.Writer, projectID string, region string, registryID string, deviceID string, keyPath string) (*cloudiot.Device, error) {
	client, err := GetClient()
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

	return response, err
}

// [END iot_create_es_device]

// [START iot_create_rsa_device]

// CreateRSA creates a device in a registry given RSA X.509 credentials.
func CreateRSA(w io.Writer, projectID string, region string, registryID string, deviceID string, keyPath string) (*cloudiot.Device, error) {
	client, err := GetClient()
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

	return response, err
}

// [END iot_create_rsa_device]

// [START iot_create_unauth_device]

// CreateUnauth creates a device in a registry without credentials.
func CreateUnauth(w io.Writer, projectID string, region string, registryID string, deviceID string) (*cloudiot.Device, error) {
	client, err := GetClient()
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

	return response, err
}

// [END iot_create_unauth_device]

// [START iot_delete_device]

// DeleteDevice deletes a device from a registry.
func DeleteDevice(w io.Writer, projectID string, region string, registryID string, deviceID string) (*cloudiot.Empty, error) {
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

	return response, err
}

// [END iot_delete_device]

// [START iot_get_device]

// GetDevice retrieves a specific device and prints its details.
func GetDevice(w io.Writer, projectID string, region string, registryID string, device string) (*cloudiot.Device, error) {
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

	return response, err
}

// [END iot_get_device]

// [START iot_get_device_configs]

// GetDeviceConfigs retrieves and lists device configurations.
func GetDeviceConfigs(w io.Writer, projectID string, region string, registryID string, device string) ([]*cloudiot.DeviceConfig, error) {
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

	return response.DeviceConfigs, err
}

// [END iot_get_device_configs]

// [START iot_get_device_state]

// GetDeviceStates retrieves and lists device states.
func GetDeviceStates(w io.Writer, projectID string, region string, registryID string, device string) ([]*cloudiot.DeviceState, error) {
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

	return response.DeviceStates, err
}

// [END iot_get_device_state]

// [START iot_list_devices]

// ListDevices gets the identifiers of devices for a specific registry.
func ListDevices(w io.Writer, projectID string, region string, registryID string) ([]*cloudiot.Device, error) {
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

	return response.Devices, err
}

// [END iot_list_devices]

// [START iot_patch_es]

// PatchDeviceES patches a device to use ES256 credentials.
func PatchDeviceES(w io.Writer, projectID string, region string, registryID string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	return response, err
}

// [END iot_patch_es]

// [START iot_patch_rsa]

// PatchDeviceRSA patches a device to use RSA256 X.509 credentials.
func PatchDeviceRSA(w io.Writer, projectID string, region string, registryID string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	return response, err
}

// [END iot_patch_rsa]

// [START iot_set_device_config]

// SetConfig sends a configuration change to a device.
func SetConfig(w io.Writer, projectID string, region string, registryID string, deviceID string, configData string, format string) (*cloudiot.DeviceConfig, error) {
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

	return response, err
}

// [END iot_set_device_config]

// [START iot_send_command]

// SendCommand sends a command to a device listening for commands.
func SendCommand(w io.Writer, projectID string, region string, registryID string, deviceID string, sendData string) (*cloudiot.SendCommandToDeviceResponse, error) {
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

	return response, err
}

// [END iot_send_command]

// START BETA FEATURES

// [START iot_create_gateway]

// createGateway creates a new IoT Core gateway with a given id, public key, and auth method.
// gatewayAuthMethod can be one of: ASSOCIATION_ONLY, DEVICE_CREDENTIALS_ONLY, ASSOCIATION_AND_DEVICE_AUTH_TOKEN
func createGateway(w io.Writer, projectID string, region string, registryID string, gatewayID string, gatewayAuthMethod string, publicKeyPath string) (*cloudiot.Device, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	keyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	gateway := cloudiot.Device{
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
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &gateway).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Successfully created gateway:", gatewayID)

	return response, err
}

// [END iot_create_gateway]

// [START iot_list_gateways]

// listGateways lists all the gateways in a specific registry
func listGateways(w io.Writer, projectID string, region string, registryID string) ([]*cloudiot.Device, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.List(parent).GatewayListOptionsGatewayType("GATEWAY").Do()

	if err != nil {
		return nil, err
	}

	if len(response.Devices) == 0 {
		fmt.Fprintln(w, "No gateways found")
	} else {
		fmt.Fprintln(w, len(response.Devices), "devices:")
		for _, gateway := range response.Devices {
			fmt.Fprintf(w, "\t%s\n", gateway.Id)
		}
	}

	return response.Devices, err
}

// [END iot_list_gateways]

// [START iot_bind_device_to_gateway]

// bindDeviceToGateway creates an association between an existing device and gateway
func bindDeviceToGateway(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string) (*cloudiot.BindDeviceToGatewayResponse, error) {
	client, err := GetClient()

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	bindRequest := &cloudiot.BindDeviceToGatewayRequest{
		DeviceId:  deviceID,
		GatewayId: gatewayID,
	}

	response, err := client.Projects.Locations.Registries.BindDeviceToGateway(parent, bindRequest).Do()

	if err != nil {
		return response, err
	}

	if response.HTTPStatusCode == 200 {
		fmt.Fprintf(w, "Bound %s to %s", deviceID, gatewayID)
	}

	return response, err
}

// [END iot_bind_device_to_gateway]

// [START unbind_device_from_gateway]

// unbindDeviceFromGateway unbinds a bound device from a gateway
func unbindDeviceFromGateway(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string) (*cloudiot.UnbindDeviceFromGatewayResponse, error) {
	client, err := GetClient()

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	unbindRequest := &cloudiot.UnbindDeviceFromGatewayRequest{
		DeviceId:  deviceID,
		GatewayId: gatewayID,
	}

	response, err := client.Projects.Locations.Registries.UnbindDeviceFromGateway(parent, unbindRequest).Do()

	if err != nil {
		return response, err
	}

	if response.HTTPStatusCode == 200 {
		fmt.Fprintf(w, "Unbound %s from %s", deviceID, gatewayID)
	}

	return response, err
}

// [END unbind_device_from_gateway]

// [START list_devices_for_gateway]
func listDevicesForGateway(w io.Writer, projectID string, region string, registryID, gatewayID string) ([]*cloudiot.Device, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registryID)
	response, err := client.Projects.Locations.Registries.Devices.List(parent).GatewayListOptionsAssociationsGatewayId(gatewayID).Do()

	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Devices for %s:\n", gatewayID)
	if len(response.Devices) == 0 {
		fmt.Fprintln(w, "\tNo devices found")
	} else {
		for _, gateway := range response.Devices {
			fmt.Fprintf(w, "\t%s\n", gateway.Id)
		}
	}

	return response.Devices, err
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
		{"createRegistry", CreateRegistry, []string{"cloud-region", "registry-id", "pubsub-topic"}},
		{"deleteRegistry", DeleteRegistry, []string{"cloud-region", "registry-id"}},
		{"getRegistry", GetRegistry, []string{"cloud-region", "registry-id"}},
		{"listRegistries", ListRegistries, []string{"cloud-region"}},
		{"getRegistryIAM", GetRegistryIAM, []string{"cloud-region", "registry-id"}},
		{"setRegistryIAM", SetRegistryIAM, []string{"cloud-region", "registry-id", "member", "role"}},
	}

	deviceManagementCommands := []command{
		{"createES", CreateES, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"createRSA", CreateRSA, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"createUnauth", CreateUnauth, []string{"cloud-region", "registry-id", "device-id"}},
		{"deleteDevice", DeleteDevice, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDevice", GetDevice, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDeviceConfigs", GetDeviceConfigs, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDeviceStates", GetDeviceStates, []string{"cloud-region", "registry-id", "device-id"}},
		{"listDevices", ListDevices, []string{"cloud-region", "registry-id"}},
		{"patchDevice", PatchDeviceES, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"patchDeviceRSA", PatchDeviceRSA, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"setConfig", SetConfig, []string{"cloud-region", "registry-id", "device-id", "config-data"}},
		{"sendCommand", SendCommand, []string{"cloud-region", "registry-id", "device-id", "send-data"}},
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
