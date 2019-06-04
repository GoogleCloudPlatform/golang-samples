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

// [START iot_subscribe_gateway_bound_device]

import (
	"fmt"
	"io"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// subscribeGatewayToDeviceTopic creates a gateway client that subscribes to a topic of a bound device.
// Currently supported topics include: "config", "state", "commands", "errors"
func subscribeGatewayToDeviceTopic(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string, privateKeyPath string, algorithm string, clientDuration int, topic string) error {

	const (
		mqttBrokerURL   = "tls://mqtt.googleapis.com:8883"
		protocolVersion = 4 // corresponds to MQTT 3.1.1
	)

	// onConnect defines the on connect handler which resets backoff variables.
	var onConnect mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Fprintf(w, "Client connected: %t\n", client.IsConnected())
	}

	// onMessage defines the message handler for the mqtt client.
	var onMessage mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Fprintf(w, "Topic: %s\n", msg.Topic())
		fmt.Fprintf(w, "Message: %s\n", msg.Payload())
	}

	// onDisconnect defines the connection lost handler for the mqtt client.
	var onDisconnect mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Println("Client disconnected")
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
		fmt.Fprintf(w, "AttachDevice error: %v\n", err)
		return err
	}

	// Sleep for 5 seconds to allow attachDevice message to propagate.
	time.Sleep(5 * time.Second)

	// Subscribe to the config topic of the current gateway and a device bound to the gateway.
	gatewayTopic := fmt.Sprintf("/devices/%s/%s", gatewayID, topic)
	if token := client.Subscribe(gatewayTopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Fprintln(w, token.Error())
		return token.Error()
	}

	deviceTopic := fmt.Sprintf("/devices/%s/%s", deviceID, topic)
	if token := client.Subscribe(deviceTopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Fprintln(w, token.Error())
		return token.Error()
	}

	time.Sleep(time.Duration(clientDuration) * time.Second)

	if err := detachDevice(deviceID, client, ""); err != nil {
		fmt.Fprintf(w, "DetachDevice error: %v\n", err)
		return err
	}

	if token := client.Unsubscribe(gatewayTopic, deviceTopic); token.Wait() && token.Error() != nil {
		fmt.Fprintln(w, token.Error())
		return token.Error()
	}

	client.Disconnect(10)
	return nil
}

// [END iot_subscribe_gateway_bound_device]
