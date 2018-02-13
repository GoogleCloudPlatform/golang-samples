// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command manager lets you manage Cloud IoT Core devices and registries.
package main

import (
	// [START imports]
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
	// [END imports]
)

// Registry Management
// createRegistry creates a device registry
func createRegistry(proj string, region string, registryId string, topicName string) (*cloudiot.DeviceRegistry, error) {
	// [START iot_create_registry]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", proj, region)
	var notify cloudiot.EventNotificationConfig
	notifyConfigs := []*cloudiot.EventNotificationConfig{&notify}
	notify.PubsubTopicName = topicName
	reg := &cloudiot.DeviceRegistry{
		Id: registryId,
		EventNotificationConfigs: notifyConfigs,
	}

	response, err := client.Projects.Locations.Registries.Create(parent, reg).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating registry: %v", err)
		return nil, err
	}

	fmt.Println("Created registry:")
	fmt.Println("\tId: ", response.Id)
	fmt.Println("\tHTTP: ", response.HttpConfig.HttpEnabledState)
	fmt.Println("\tMQTT: ", response.MqttConfig.MqttEnabledState)
	fmt.Println("\tName: ", response.Name)

	return response, err
	// [END iot_create_registry]
}

// deleteRegistry deletes a device registry
func deleteRegistry(proj string, region string, registryId string) (*cloudiot.Empty, error) {
	// [START iot_delete_registry]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s", proj, region, registryId)

	response, err := client.Projects.Locations.Registries.Delete(name).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing registries: %v", err)
		return nil, err
	} else {
		fmt.Println("Deleted registry")
	}

	return response, err
}

func getRegistry(proj string, region string, registryId string) (*cloudiot.DeviceRegistry, error) {
	// [START iot_get_registry]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", proj, region, registryId)
	response, err := client.Projects.Locations.Registries.Get(parent).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing registries: %v", err)
		return nil, err
	}

	return response, err
	// [END iot_get_iam]
}

// getRegistryIam gets the IAM policy for a device registry
func getRegistryIam(proj string, region string, registryId string) (*cloudiot.Policy, error) {
	// [START iot_get_iam]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s", proj, region, registryId)
  var req cloudiot.GetIamPolicyRequest
	response, err := client.Projects.Locations.Registries.GetIamPolicy(path, &req).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting registry policy: %v", err)
		return nil, err
	}

	fmt.Println("Found policy:")
	for _, binding := range response.Bindings {
		fmt.Fprintf(os.Stdout, "Role: %s\n", binding.Role)
    for _, member := range binding.Members {
      fmt.Fprintf(os.Stdout, "\tMember: %s\n", member)
    }
	}
	return response, err
	// [END iot_get_iam]
}

// listRegistries gets the names of device registries given a project / region.
func listRegistries(proj string, region string) ([]*cloudiot.DeviceRegistry, error) {
	// [START iot_list_registries]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", proj, region)
	response, err := client.Projects.Locations.Registries.List(parent).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing registries: %v", err)
		return nil, err
	}

	fmt.Println("Registries:")
	for _, registry := range response.DeviceRegistries {
		fmt.Println("\t", registry.Name)
	}

	return response.DeviceRegistries, err
	// [END iot_list_registries]
}

// setRegistryIam gets the IAM policy for a device registry
func setRegistryIam(proj string, region string, registryId string, member string, role string) (*cloudiot.Policy, error) {
	// [START iot_set_iam]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s", proj, region, registryId)
  var policy cloudiot.Policy
  var binding cloudiot.Binding
  var req cloudiot.SetIamPolicyRequest
  binding.Role = role
  binding.Members = []string {member}
  policy.Bindings = []*cloudiot.Binding {&binding}
  req.Policy = &policy

	response, err := client.Projects.Locations.Registries.SetIamPolicy(path, &req).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting registry policy: %v", err)
		return nil, err
	}

	fmt.Println("Set policy!")
	return response, err
	// [END iot_set_iam]
}

// Device Management
// createEs creates a device in a registry with ES credentials
func createEs(proj string, region string, registry string, device string, keyPath string) (*cloudiot.Device, error) {
	// [START iot_create_es]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", proj, region, registry)

	var dev cloudiot.Device
	var cred cloudiot.DeviceCredential
	var key cloudiot.PublicKeyCredential

	key.Format = "ES256_PEM"
	var keybytes, keyerr = ioutil.ReadFile(keyPath)
	key.Key = string(keybytes[:])
	if keyerr != nil {
		fmt.Fprintf(os.Stderr, "Unable to open certificate file: %v", err)
	}

	cred.PublicKey = &key
	dev.Id = device
	dev.Credentials = []*cloudiot.DeviceCredential{&cred}

	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &dev).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating device: %v", err)
		return nil, err
	} else {
		fmt.Println("Successfully created device.")
	}

	return response, err
	// [END iot_create_es]
}

// createRsa creates a device in a registry with RS credentials
func createRsa(proj string, region string, registry string, device string, keyPath string) (*cloudiot.Device, error) {
	// [START iot_create_rsa]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", proj, region, registry)

	var dev cloudiot.Device
	var cred cloudiot.DeviceCredential
	var key cloudiot.PublicKeyCredential

	key.Format = "RSA_X509_PEM"
	var keybytes, keyerr = ioutil.ReadFile(keyPath)
	key.Key = string(keybytes[:])
	if keyerr != nil {
		fmt.Fprintf(os.Stderr, "Unable to open certificate file: %v", err)
	}

	cred.PublicKey = &key
	dev.Id = device
	dev.Credentials = []*cloudiot.DeviceCredential{&cred}

	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &dev).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating device: %v", err)
		return nil, err
	} else {
		fmt.Println("Successfully created device.")
	}

	return response, err
	// [END iot_create_rsa]
}

// createUnauth creates a device in a registry without credentials
func createUnauth(proj string, region string, registry string, device string) (*cloudiot.Device, error) {
	// [START iot_create_unauth]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", proj, region, registry)
	var dev cloudiot.Device
	dev.Id = device

	response, err := client.Projects.Locations.Registries.Devices.Create(parent, &dev).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating device: %v", err)
		return nil, err
	} else {
		fmt.Println("Successfully created device.")
	}

	return response, err
	// [END iot_create_unauth]
}

// deleteDevice deletes a device from a registry
func deleteDevice(proj string, region string, registry string, device string) (*cloudiot.Empty, error) {
	// [START iot_delete_device]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", proj, region, registry, device)

	response, err := client.Projects.Locations.Registries.Devices.Delete(path).Do()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating device: %v", err)
		return nil, err
	} else {
		fmt.Println("Deleted device!")
	}

	return response, err
	// [END iot_delete_device]
}

// getDevice retrieves a specific device and prints its details
func getDevice(proj string, region string, registry string, device string) (*cloudiot.Device, error) {
	// [START iot_get_device]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", proj, region, registry, device)
	response, err := client.Projects.Locations.Registries.Devices.Get(path).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting device: %v", err)
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
func getDeviceConfigs(proj string, region string, registry string, device string) ([]*cloudiot.DeviceConfig, error) {
	// [START iot_get_device_configs]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", proj, region, registry, device)
	response, err := client.Projects.Locations.Registries.Devices.ConfigVersions.List(path).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting device configs: %v", err)
		return nil, err
	}

	for _, config := range response.DeviceConfigs {
		fmt.Println(config.Version, " : ", config.BinaryData)
	}

	return response.DeviceConfigs, err
	// [END iot_get_device_configs]
}

// getDeviceStates retrieves and lists device states
func getDeviceStates(proj string, region string, registry string, device string) ([]*cloudiot.DeviceState, error) {
	// [START iot_get_device_configs]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", proj, region, registry, device)
	response, err := client.Projects.Locations.Registries.Devices.States.List(path).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting device states: %v", err)
		return nil, err
	}

  fmt.Println("Successfully retrieved device states!")

	for _, state := range response.DeviceStates {
		fmt.Println(state.UpdateTime, " : ", state.BinaryData)
	}

	return response.DeviceStates, err
	// [END iot_get_device_configs]
}

// listDevices gets the identifiers of devices given a registry name
func listDevices(proj string, region string, registry string) ([]*cloudiot.Device, error) {
	// [START iot_list_devices]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", proj, region, registry)
	response, err := client.Projects.Locations.Registries.Devices.List(parent).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing devices: %v", err)
		return nil, err
	}

	fmt.Println("Devices:")
	for _, device := range response.Devices {
		fmt.Println("\t", device.Id)
	}

	return response.Devices, err
	// [END iot_list_devices]
}

// patchDeviceEs patches a device to use ES credentials
func patchDeviceEs(proj string, region string, registry string, device string, keyPath string) (*cloudiot.Device, error) {
	// [START iot_patch_es]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", proj, region, registry, device)

	var dev cloudiot.Device
	var cred cloudiot.DeviceCredential
	var key cloudiot.PublicKeyCredential

	key.Format = "ES256_PEM"
	var keybytes, keyerr = ioutil.ReadFile(keyPath)
	key.Key = string(keybytes[:])
	if keyerr != nil {
		fmt.Fprintf(os.Stderr, "Unable to open certificate file: %v", err)
	}

	cred.PublicKey = &key
	dev.Credentials = []*cloudiot.DeviceCredential{&cred}

	response, err := client.Projects.Locations.Registries.Devices.Patch(
		parent, &dev).UpdateMask("credentials").Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error patching device: %v", err)
		return nil, err
	}

	fmt.Println("Successfully patched device.")
	return response, err
	// [END iot_patch_es]
}

// patchDeviceRsa patches a device to use RSA credentials
func patchDeviceRsa(proj string, region string, registry string, device string, keyPath string) (*cloudiot.Device, error) {
	// [START iot_patch_rsa]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", proj, region, registry, device)

	var dev cloudiot.Device
	var cred cloudiot.DeviceCredential
	var key cloudiot.PublicKeyCredential

	key.Format = "RSA_X509_PEM"
	var keybytes, keyerr = ioutil.ReadFile(keyPath)
	key.Key = string(keybytes[:])
	if keyerr != nil {
		fmt.Fprintf(os.Stderr, "Unable to open certificate file: %v", err)
	}

	cred.PublicKey = &key
	dev.Credentials = []*cloudiot.DeviceCredential{&cred}

	response, err := client.Projects.Locations.Registries.Devices.Patch(
		parent, &dev).UpdateMask("credentials").Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error patching device: %v", err)
		return nil, err
	}

	fmt.Println("Successfully patched device.")
	return response, err
	// [END iot_patch_rsa]
}

// setConfig sends a configuration change to a device
func setConfig(proj string, region string, registry string, device string, configdata string) (*cloudiot.DeviceConfig, error) {
	// [START iot_set_config]
	// Authorize the client using Application Default Credentials.
	// See https://g.co/dv/identity/protocols/application-default-credentials
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
	client, err := cloudiot.New(c)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create Cloud IoT service: %v", err)
		return nil, err
	}


	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", proj, region, registry, device)
  var req cloudiot.ModifyCloudToDeviceConfigRequest
  req.BinaryData = configdata

	response, err := client.Projects.Locations.Registries.Devices.ModifyCloudToDeviceConfig(path, &req).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting config: %v", err)
		return nil, err
	}

  fmt.Fprintf(os.Stdout, "Config set!\nVersion now: %d", response.Version)

	return response, err
	// [END iot_set_config]
}


func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\tRegistry Management\n")
		fmt.Fprintf(os.Stderr, "\t-----\n")
		fmt.Fprintf(os.Stderr, "\t%s createRegistry <cloud-region> <registry-id> <pubsub-topic>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s deleteRegistry <cloud-region> <registry-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s getRegistry <cloud-region> <registry-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s getRegistryIam <cloud-region> <registry-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s listRegistries <cloud-region>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s setRegistryIam <cloud-region> <registry-id> <member> <role>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\tDevice Management\n")
		fmt.Fprintf(os.Stderr, "\t-----\n")
		fmt.Fprintf(os.Stderr, "\t%s createEs <cloud-region> <registry-id> <device-id> <keyfile-path>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s createRsa <cloud-region> <registry-id> <device-id> <keyfile-path>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s createUnauth <cloud-region> <registry-id> <device-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s deleteDevice <cloud-region> <registry-id> <device-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s getDevice <cloud-region> <registry-id> <device-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s getDeviceConfigs <cloud-region> <registry-id> <device-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s getDeviceStates <cloud-region> <registry-id> <device-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s listDevices <cloud-region> <registry-id>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s patchEs <cloud-region> <registry-id> <device-id> <keyfile-path>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s patchRsa <cloud-region> <registry-id> <device-id> <keyfile-path>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s setConfig <cloud-region> <registry-id> <device-id> <config-data>\n", filepath.Base(os.Args[0]))
	}
	flag.Parse()

	// Retrieve project ID from console
	proj := os.Getenv("GCLOUD_PROJECT")
	if proj == "" {
		proj = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if proj == "" {
		fmt.Fprintf(os.Stderr, "Set the GCLOUD_PROJECT or GOOGLE_CLOUD_PROJECT environment variable.")
	}

	args := flag.Args()
	if len(args) > 0 && args[0] == "createRegistry" {
		createRegistry(proj, args[1], args[2], args[3])
	} else if len(args) > 0 && args[0] == "deleteRegistry" {
		deleteRegistry(proj, args[1], args[2])
	} else if len(args) > 0 && args[0] == "getRegistry" {
		getRegistry(proj, args[1], args[2])
	} else if len(args) > 0 && args[0] == "getRegistryIam" {
		getRegistryIam(proj, args[1], args[2])
	} else if len(args) > 0 && args[0] == "listRegistries" {
		listRegistries(proj, args[1])
	} else if len(args) > 0 && args[0] == "setRegistryIam" {
		setRegistryIam(proj, args[1], args[2], args[3], args[4])
	} else if len(args) > 0 && args[0] == "createEs" {
		createEs(proj, args[1], args[2], args[3], args[4])
	} else if len(args) > 0 && args[0] == "createRsa" {
		createRsa(proj, args[1], args[2], args[3], args[4])
	} else if len(args) > 0 && args[0] == "createUnauth" {
		createUnauth(proj, args[1], args[2], args[3])
	} else if len(args) > 0 && args[0] == "deleteDevice" {
		deleteDevice(proj, args[1], args[2], args[3])
	} else if len(args) > 0 && args[0] == "getDevice" {
		getDevice(proj, args[1], args[2], args[3])
	} else if len(args) > 0 && args[0] == "getDeviceConfigs" {
		getDeviceConfigs(proj, args[1], args[2], args[3])
	} else if len(args) > 0 && args[0] == "getDeviceStates" {
		getDeviceStates(proj, args[1], args[2], args[3])
	} else if len(args) > 0 && args[0] == "listDevices" {
		listDevices(proj, args[1], args[2])
	} else if len(args) > 0 && args[0] == "patchEs" {
		patchDeviceEs(proj, args[1], args[2], args[3], args[4])
	} else if len(args) > 0 && args[0] == "patchRsa" {
		patchDeviceRsa(proj, args[1], args[2], args[3], args[4])
	} else if len(args) > 0 && args[0] == "setConfig" {
		setConfig(proj, args[1], args[2], args[3], args[4])
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
