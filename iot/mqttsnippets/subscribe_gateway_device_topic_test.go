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

package mqttsnippets

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

// setConfig sends a configuration message to a device.
func setConfig(w io.Writer, projectID string, region string, registryID string, deviceID string, configData string) (*cloudiot.DeviceConfig, error) {
	client, err := getClient()
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

func TestSubscribeGatewayToDeviceTopic(t *testing.T) {
	if pubKeyRSA == "" {
		t.Skip("GOLANG_SAMPLES_IOT_PUB not set")
	}

	projectID := testutil.SystemTest(t).ProjectID

	registryID := "golang-iot-test-registry"
	gatewayID := createID("gateway")
	deviceID := createID("device")

	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		if _, err := createGateway(ioutil.Discard, projectID, region, registryID, gatewayID, "ASSOCIATION_ONLY", pubKeyRSA); err != nil {
			r.Errorf("Could not create gateway: %v\n", err)
		}
	})
	if t.Failed() {
		return
	}

	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		if _, err := createDevice(ioutil.Discard, projectID, region, registryID, deviceID, "RSA_X509_PEM", pubKeyRSA); err != nil {
			r.Errorf("Could not create device: %v\n", err)
		}
	})
	if t.Failed() {
		return
	}

	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		if _, err := bindDeviceToGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID); err != nil {
			r.Errorf("Could not bind device to gateway: %v\n", err)
		}
	})
	if t.Failed() {
		return
	}

	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		// sample test config message.
		message := "{'threshold':'high'}"

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
	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		if _, err := unbindDeviceFromGateway(ioutil.Discard, projectID, region, registryID, gatewayID, deviceID); err != nil {
			r.Errorf("Could not unbind device: %v\n", err)
		}
	})

	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		if _, err := deleteDevice(ioutil.Discard, projectID, region, registryID, deviceID); err != nil {
			r.Errorf("Could not unbind device: %v\n", err)
		}
	})

	testutil.Retry(t, 10, 5*time.Second, func(r *testutil.R) {
		if _, err := deleteDevice(ioutil.Discard, projectID, region, registryID, gatewayID); err != nil {
			r.Errorf("Could not unbind device: %v\n", err)
		}
	})
}
