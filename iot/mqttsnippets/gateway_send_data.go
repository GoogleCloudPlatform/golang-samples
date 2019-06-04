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

// [START iot_send_data_from_bound_device]

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// sendDataFromBoundDevice starts a gateway client that sends data on behalf of a bound device.
func sendDataFromBoundDevice(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string, privateKeyPath string, algorithm string, numMessages int, payload string) error {
	const (
		mqttBrokerURL      = "tls://mqtt.googleapis.com:8883"
		protocolVersion    = 4  // corresponds to MQTT 3.1.1
		minimumBackoffTime = 1  // initial backoff time in seconds
		maximumBackoffTime = 32 // maximum backoff time in seconds
	)

	var backoffTime = minimumBackoffTime
	var shouldBackoff = false

	// onConnect defines the on connect handler which resets backoff variables.
	var onConnect mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Fprintf(w, "Client connected: %t\n", client.IsConnected())

		shouldBackoff = false
		backoffTime = minimumBackoffTime
	}

	// onMessage defines the message handler for the mqtt client.
	var onMessage mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Fprintf(w, "Topic: %s\n", msg.Topic())
		fmt.Fprintf(w, "Message: %s\n", msg.Payload())
	}

	// onDisconnect defines the connection lost handler for the mqtt client.
	var onDisconnect mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Fprintln(w, "Client disconnected")
		shouldBackoff = true
	}

	jwt, _ := createJWT(projectID, privateKeyPath, algorithm, 60)
	clientID := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", projectID, region, registryID, gatewayID)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBrokerURL)
	opts.SetClientID(clientID)
	opts.SetUsername("unused")
	opts.SetPassword(jwt)
	opts.SetProtocolVersion(protocolVersion)
	opts.SetOnConnectHandler(onConnect)
	opts.SetDefaultPublishHandler(onMessage)
	opts.SetConnectionLostHandler(onDisconnect)

	// Create and connect a client using the above options.
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Fprintln(w, "Failed to connect client")
		return token.Error()
	}

	if err := attachDevice(deviceID, client, ""); err != nil {
		fmt.Fprintf(w, "Failed to attach device %s\n", err)
		return err
	}

	// Sleep for 5 seconds to allow attachDevice message to propagate.
	time.Sleep(5 * time.Second)

	gatewayStateTopic := fmt.Sprintf("/devices/%s/state", gatewayID)
	deviceStateTopic := fmt.Sprintf("/devices/%s/state", deviceID)

	gatewayInitPayload := fmt.Sprintf("Starting gateway at time: %d", time.Now().Unix())
	if token := client.Publish(gatewayStateTopic, 1, false, gatewayInitPayload); token.Wait() && token.Error() != nil {
		fmt.Fprintln(w, "Failed to publish initial gateway payload")
		return token.Error()
	}

	for i := 1; i <= numMessages; i++ {
		if shouldBackoff {
			if backoffTime > maximumBackoffTime {
				fmt.Fprintln(w, "Exceeded max backoff time.")
				return errors.New("exceeded maximum backoff time, exiting")
			}

			waitTime := backoffTime + rand.Intn(1000)/1000.0
			time.Sleep(time.Duration(waitTime) * time.Second)

			backoffTime *= 2
			client = mqtt.NewClient(opts)
		}

		deviceStatePayload := fmt.Sprintf("%s, #%d", payload, i)

		fmt.Fprintf(w, "Publishing message: %s to %s\n", deviceStatePayload, deviceStateTopic)

		if token := client.Publish(deviceStateTopic, 1, false, payload); token.Wait() && token.Error() != nil {
			fmt.Fprintln(w, "Failed to publish payload to device state topic")
			return token.Error()
		}

		// Sleep for a bit between messages to simulate real world device state publishing.
		time.Sleep(5 * time.Second)
	}

	detachDevice(deviceID, client, "")

	client.Disconnect(20)
	return nil
}

// [END iot_send_data_from_bound_device]
