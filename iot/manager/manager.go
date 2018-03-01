// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command manager lets you manage Cloud IoT Core devices and registries.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	// [START imports]
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
	// [END imports]
)

// Registry Management

// createRegistry creates a device registry.
func createRegistry(projectID string, region string, registryID string, topicName string) (*cloudiot.DeviceRegistry, error) {
	// [START iot_create_registry]
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

	fmt.Println("Created registry:")
	fmt.Println("\tID: ", response.Id)
	fmt.Println("\tHTTP: ", response.HttpConfig.HttpEnabledState)
	fmt.Println("\tMQTT: ", response.MqttConfig.MqttEnabledState)
	fmt.Println("\tName: ", response.Name)
	// [END iot_create_registry]

	return response, err
}

// deleteRegistry deletes a device registry
func deleteRegistry(projectID string, region string, registryID string) (*cloudiot.Empty, error) {
	// [START iot_delete_registry]
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

	fmt.Println("Deleted registry")
	// [END iot_delete_registry]

	return response, err
}

func getRegistry(projectID string, region string, registryID string) (*cloudiot.DeviceRegistry, error) {
	// [START iot_get_registry]
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
	// [END iot_get_iam]

	return response, err
}

// getRegistryIam gets the IAM policy for a device registry.
func getRegistryIam(projectID string, region string, registryID string) (*cloudiot.Policy, error) {
	// [START iot_get_iam_policy]
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

	fmt.Println("Policy:")
	for _, binding := range response.Bindings {
		fmt.Fprintf(os.Stdout, "Role: %s\n", binding.Role)
		for _, member := range binding.Members {
			fmt.Fprintf(os.Stdout, "\tMember: %s\n", member)
		}
	}
	// [END iot_get_iam_policy]

	return response, err
}

// listRegistries gets the names of device registries given a project / region.
func listRegistries(projectID string, region string) ([]*cloudiot.DeviceRegistry, error) {
	// [START iot_list_registries]
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

	fmt.Println("Registries:")
	for _, registry := range response.DeviceRegistries {
		fmt.Println("\t", registry.Name)
	}
	// [END iot_list_registries]

	return response.DeviceRegistries, err
}

// setRegistryIam sets the IAM policy for a device registry
func setRegistryIam(projectID string, region string, registryID string, member string, role string) (*cloudiot.Policy, error) {
	// [START iot_set_iam_policy]
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

	fmt.Println("Set policy!")
	// [END iot_set_iam_policy]

	return response, err
}

// Device Management

// createEs creates a device in a registry with ES credentials
func createEs(projectID string, region string, registry string, deviceID string, keyPath string) (*cloudiot.Device, error) {
	// [START iot_create_es_device]
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

	fmt.Println("Successfully created device.")
	// [END iot_create_es_device]

	return response, err
}

// createRsa creates a device in a registry with RS credentials
func createRsa(projectID string, region string, registry string, deviceID string, keyPath string) (*cloudiot.Device, error) {
	// [START iot_create_rsa_device]
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

	fmt.Println("Successfully created device.")
	// [END iot_create_rsa_device]

	return response, err
}

// createUnauth creates a device in a registry without credentials
func createUnauth(projectID string, region string, registry string, deviceID string) (*cloudiot.Device, error) {
	// [START iot_create_unauth_device]
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

	fmt.Println("Successfully created device.")
	// [END iot_create_unauth_device]

	return response, err
}

// deleteDevice deletes a device from a registry
func deleteDevice(projectID string, region string, registry string, deviceID string) (*cloudiot.Empty, error) {
	// [START iot_delete_device]
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

	fmt.Println("Deleted device!")
	// [END iot_delete_device]

	return response, err
}

// getDevice retrieves a specific device and prints its details
func getDevice(projectID string, region string, registry string, device string) (*cloudiot.Device, error) {
	// [START iot_get_device]
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

	fmt.Println("\tId: ", response.Id)
	for _, credential := range response.Credentials {
		fmt.Println("\t\tCredential Expire: ", credential.ExpirationTime)
		fmt.Println("\t\tCredential Type: ", credential.PublicKey.Format)
		fmt.Println("\t\t--------")
	}
	fmt.Println("\tLast Config Ack: ", response.LastConfigAckTime)
	fmt.Println("\tLast Config Send: ", response.LastConfigSendTime)
	fmt.Println("\tLast Event Time: ", response.LastEventTime)
	fmt.Println("\tLast Heartbeat Time: ", response.LastHeartbeatTime)
	fmt.Println("\tLast State Time: ", response.LastStateTime)
	fmt.Println("\tNumId: ", response.NumId)

	return response, err
	// [END iot_get_device]
}

// getDeviceConfigs retrieves and lists device configurations
func getDeviceConfigs(projectID string, region string, registry string, device string) ([]*cloudiot.DeviceConfig, error) {
	// [START iot_get_device_configs]
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
		fmt.Println(config.Version, " : ", config.BinaryData)
	}
	// [END iot_get_device_configs]

	return response.DeviceConfigs, err
}

// getDeviceStates retrieves and lists device states
func getDeviceStates(projectID string, region string, registry string, device string) ([]*cloudiot.DeviceState, error) {
	// [START iot_get_device_state]
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

	fmt.Println("Successfully retrieved device states!")

	for _, state := range response.DeviceStates {
		fmt.Println(state.UpdateTime, " : ", state.BinaryData)
	}
	// [END iot_get_device_state]

	return response.DeviceStates, err
}

// listDevices gets the identifiers of devices given a registry name
func listDevices(projectID string, region string, registry string) ([]*cloudiot.Device, error) {
	// [START iot_list_devices]
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

	fmt.Println("Devices:")
	for _, device := range response.Devices {
		fmt.Println("\t", device.Id)
	}
	// [END iot_list_devices]

	return response.Devices, err
}

// patchDeviceEs patches a device to use ES credentials
func patchDeviceEs(projectID string, region string, registry string, deviceID string, keyPath string) (*cloudiot.Device, error) {
	// [START iot_patch_es]
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

	fmt.Println("Successfully patched device.")
	// [END iot_patch_es]

	return response, err
}

// patchDeviceRsa patches a device to use RSA credentials
func patchDeviceRsa(projectID string, region string, registry string, deviceID string, keyPath string) (*cloudiot.Device, error) {
	// [START iot_patch_rsa]
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

	fmt.Println("Successfully patched device.")
	// [END iot_patch_rsa]

	return response, err
}

// setConfig sends a configuration change to a device.
func setConfig(projectID string, region string, registry string, deviceID string, configData string) (*cloudiot.DeviceConfig, error) {
	// [START iot_set_device_config]
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
		BinaryData: configData,
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registry, deviceID)
	response, err := client.Projects.Locations.Registries.Devices.ModifyCloudToDeviceConfig(path, &req).Do()
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stdout, "Config set!\nVersion now: %d", response.Version)
	// [END iot_set_device_config]

	return response, err
}

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
			fnArgs = append(fnArgs, reflect.ValueOf(projectID))
			for _, arg := range commandArgs {
				fnArgs = append(fnArgs, reflect.ValueOf(arg))
			}
			retValues := reflect.ValueOf(cmd.fn).Call(fnArgs)
			err := retValues[len(retValues)-1].Interface().(error)
			if err != nil {
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
