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

// Registry Management

// [START iot_create_registry]

// createRegistry creates a device registry.
func createRegistry(w io.Writer, projectID string, region string, registryID string, topicName string) (*cloudiot.DeviceRegistry, error) {
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

	fmt.Fprintf(w, "Created registry:\n")
	fmt.Fprintf(w, "\tID: %s\n", response.Id)
	fmt.Fprintf(w, "\tHTTP: %s\n", response.HttpConfig.HttpEnabledState)
	fmt.Fprintf(w, "\tMQTT: %s\n", response.MqttConfig.MqttEnabledState)
	fmt.Fprintf(w, "\tName: %s\n", response.Name)

	return response, err
}

// [END iot_create_registry]

// [START iot_delete_registry]

// deleteRegistry deletes a device registry
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

	fmt.Fprintf(w, "Deleted registry\n")

	return response, err
}

// [END iot_delete_registry]

// [START iot_get_registry]

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
	// [END iot_get_registry]

	fmt.Fprintf(w, "Got registry:\n")
	fmt.Fprintf(w, "\tID: %s\n", response.Id)
	fmt.Fprintf(w, "\tHTTP: %s\n", response.HttpConfig.HttpEnabledState)
	fmt.Fprintf(w, "\tMQTT: %s\n", response.MqttConfig.MqttEnabledState)
	fmt.Fprintf(w, "\tName: %s\n", response.Name)

	return response, err
}

// [START iot_get_iam_policy]

// getRegistryIam gets the IAM policy for a device registry.
func getRegistryIam(w io.Writer, projectID string, region string, registryID string) (*cloudiot.Policy, error) {
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

	fmt.Fprintf(w, "Policy:\n")
	for _, binding := range response.Bindings {
		fmt.Fprintf(w, "Role: %s\n", binding.Role)
		for _, member := range binding.Members {
			fmt.Fprintf(w, "\tMember: %s\n", member)
		}
	}

	return response, err
}

// [END iot_get_iam_policy]

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

	fmt.Fprintf(w, "Registries:\n")
	for _, registry := range response.DeviceRegistries {
		fmt.Fprintf(w, "\t%s\n", registry.Name)
	}

	return response.DeviceRegistries, err
}

// [END iot_list_registries]

// [START iot_set_iam_policy]

// setRegistryIam sets the IAM policy for a device registry
func setRegistryIam(w io.Writer, projectID string, region string, registryID string, member string, role string) (*cloudiot.Policy, error) {
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

	return response, err
}

// [END iot_set_iam_policy]

// Device Management

// [START iot_create_es_device]

// createEs creates a device in a registry with ES credentials
func createEs(w io.Writer, projectID string, region string, registry string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registry)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully created ESA device\n")

	return response, err
}

// [END iot_create_es_device]

// [START iot_create_rsa_device]

// createRsa creates a device in a registry with RS credentials
func createRsa(w io.Writer, projectID string, region string, registry string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registry)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully created RSA device\n")

	return response, err
}

// [END iot_create_rsa_device]

// [START iot_create_unauth_device]

// createUnauth creates a device in a registry without credentials
func createUnauth(w io.Writer, projectID string, region string, registry string, deviceID string) (*cloudiot.Device, error) {
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
	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registry)
	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully created device without credentials\n")

	return response, err
}

// [END iot_create_unauth_device]

// [START iot_delete_device]

// deleteDevice deletes a device from a registry
func deleteDevice(w io.Writer, projectID string, region string, registry string, deviceID string) (*cloudiot.Empty, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.Delete(path).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Deleted device: %s\n", deviceID)

	return response, err
}

// [END iot_delete_device]

// [START iot_get_device]

// getDevice retrieves a specific device and prints its details
func getDevice(w io.Writer, projectID string, region string, registry string, device string) (*cloudiot.Device, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, device)
	response, err := client.Projects.Locations.Registries.Devices.Get(path).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "\tId: %s\n", response.Id)
	for _, credential := range response.Credentials {
		fmt.Fprintf(w, "\t\tCredential Expire: %s\n", credential.ExpirationTime)
		fmt.Fprintf(w, "\t\tCredential Type: %s\n", credential.PublicKey.Format)
		fmt.Fprintf(w, "\t\t--------\n")
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

// getDeviceConfigs retrieves and lists device configurations
func getDeviceConfigs(w io.Writer, projectID string, region string, registry string, device string) ([]*cloudiot.DeviceConfig, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, device)
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

// getDeviceStates retrieves and lists device states
func getDeviceStates(w io.Writer, projectID string, region string, registry string, device string) ([]*cloudiot.DeviceState, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, device)
	response, err := client.Projects.Locations.Registries.Devices.States.List(path).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully retrieved device states!\n")

	for _, state := range response.DeviceStates {
		fmt.Fprintf(w, "%s : %s\n", state.UpdateTime, state.BinaryData)
	}

	return response.DeviceStates, err
}

// [END iot_get_device_state]

// [START iot_list_devices]

// listDevices gets the identifiers of devices given a registry name
func listDevices(w io.Writer, projectID string, region string, registry string) ([]*cloudiot.Device, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", projectID, region, registry)
	response, err := client.Projects.Locations.Registries.Devices.List(parent).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Devices:\n")
	for _, device := range response.Devices {
		fmt.Fprintf(w, "\t%s\n", device.Id)
	}

	return response.Devices, err
}

// [END iot_list_devices]

// [START iot_patch_es]

// patchDeviceEs patches a device to use ES credentials
func patchDeviceEs(w io.Writer, projectID string, region string, registry string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.
		Patch(parent, &device).UpdateMask("credentials").Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully patched device with ES credentials\n")

	return response, err
}

// [END iot_patch_es]

// [START iot_patch_rsa]

// patchDeviceRsa patches a device to use RSA credentials
func patchDeviceRsa(w io.Writer, projectID string, region string, registry string, deviceID string, keyPath string) (*cloudiot.Device, error) {
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

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.
		Patch(parent, &device).UpdateMask("credentials").Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Successfully patched device\n")

	return response, err
}

// [END iot_patch_rsa]

// [START iot_set_device_config]

// setConfig sends a configuration change to a device.
func setConfig(w io.Writer, projectID string, region string, registry string, deviceID string, configData string, format string) (*cloudiot.DeviceConfig, error) {
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

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.ModifyCloudToDeviceConfig(path, &req).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Config set!\nVersion now: %d\n", response.Version)

	return response, err
}

// [END iot_set_device_config]

// [START iot_send_command]

// sendCommand sends a command to a device listening for commands
func sendCommand(w io.Writer, projectID string, region string, registry string, deviceID string, sendData string) (*cloudiot.SendCommandToDeviceResponse, error) {
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

	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, deviceID)

	response, err := client.Projects.Locations.Registries.Devices.SendCommandToDevice(name, &req).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Sent command to device\n")

	return response, err
}

// [END iot_send_command]

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
		{"getRegistryIam", getRegistryIam, []string{"cloud-region", "registry-id"}},
		{"setRegistryIam", setRegistryIam, []string{"cloud-region", "registry-id", "member", "role"}},
	}

	deviceManagementCommands := []command{
		{"createEs", createEs, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"createRsa", createRsa, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"createUnauth", createUnauth, []string{"cloud-region", "registry-id", "device-id"}},
		{"deleteDevice", deleteDevice, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDevice", getDevice, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDeviceConfigs", getDeviceConfigs, []string{"cloud-region", "registry-id", "device-id"}},
		{"getDeviceStates", getDeviceStates, []string{"cloud-region", "registry-id", "device-id"}},
		{"listDevices", listDevices, []string{"cloud-region", "registry-id"}},
		{"patchDeviceEs", patchDeviceEs, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"patchDeviceRsa", patchDeviceRsa, []string{"cloud-region", "registry-id", "device-id", "keyfile-path"}},
		{"setConfig", setConfig, []string{"cloud-region", "registry-id", "device-id", "config-data"}},
		{"sendCommand", sendCommand, []string{"cloud-region", "registry-id", "device-id", "send-data"}},
	}

	var commands []command
	commands = append(commands, registryManagementCommands...)
	commands = append(commands, deviceManagementCommands...)

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
	}
	flag.Parse()

	// Retrieve project ID from console
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("GCLOUD_PROJECT")
	}
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
