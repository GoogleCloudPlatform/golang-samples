// Copyright 2019 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mqttsnippets

import (
	"fmt"

	"github.com/eclipse/paho.mqtt.golang"
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
