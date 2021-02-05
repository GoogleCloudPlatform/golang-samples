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
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// [START iot_attach_device]

// attachDevice attaches a device to a gateway.
func attachDevice(deviceID string, client mqtt.Client, jwt string) error {
	attachTopic := fmt.Sprintf("/devices/%s/attach", deviceID)
	fmt.Printf("Attaching device: %s\n", attachTopic)

	attachPayload := "{}"
	if jwt != "" {
		attachPayload = fmt.Sprintf("{ 'authorization' : %s }", jwt)
	}

	if token := client.Publish(attachTopic, 1, false, attachPayload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// [END iot_attach_device]

// [START iot_detach_device]

// detatchDevice detaches a device from a gateway.
func detachDevice(deviceID string, client mqtt.Client, jwt string) error {
	detachTopic := fmt.Sprintf("/devices/%s/detach", deviceID)
	fmt.Printf("Detaching device: %s\n", detachTopic)

	detachPayload := "{}"
	if jwt != "" {
		detachPayload = fmt.Sprintf("{ 'authorization' : %s }", jwt)
	}

	if token := client.Publish(detachTopic, 1, false, detachPayload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// [END iot_detatch_device]
