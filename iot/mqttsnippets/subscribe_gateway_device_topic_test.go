// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mqttsnippets

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

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
		BinaryData: base64.StdEncoding.EncodeToString([]byte(configData)),
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

func TestSubscribeGatewayToDeviceTopic(t *testing.T) {
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

		message := "enable low power mode"

		go func() {
			time.Sleep(10 * time.Second)
			setConfig(ioutil.Discard, projectID, region, registryID, deviceID, message)
		}()

		buf := new(bytes.Buffer)
		subscribeGatewayToDeviceTopic(buf, projectID, region, registryID, gatewayID, deviceID, privateKeyRSA, "RS256", 10, "config")

		want := fmt.Sprintf("Message: %s\n", message)
		if got := buf.String(); !strings.Contains(got, want) {
			r.Errorf("TestSubscribeGatewayToDeviceTopic got %s, want substring %q", got, want)
		}
	})

	// cleanup
	unbindDeviceFromGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID)
	deleteDevice(ioutil.Discard, projectID, region, registryID, deviceID)
	deleteDevice(ioutil.Discard, projectID, region, registryID, gatewayID)
}
