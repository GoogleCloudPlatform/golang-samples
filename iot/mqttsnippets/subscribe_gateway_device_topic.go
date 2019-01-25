// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mqttsnippets

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

// [START iot_subuscribe_config_bound_device]

// subscribeGatewayToDeviceTopic creates a gateway client that subscribes to a topic of a bound device.
// Currently supported topics include: "config", "state", "commands", "errors"
func subscribeGatewayToDeviceTopic(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string, privateKeyPath string, algorithm string, clientDuration int, topic string) {

	const (
		mqttBrokerURL   = "tls://mqtt.googleapis.com:443"
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
		return
	}

	attachDevice(deviceID, client, "")

	time.Sleep(5 * time.Second)

	// Subscribe to the config topic of the current gateway and a device bound to the gateway.
	gatewayTopic := fmt.Sprintf("/devices/%s/%s", gatewayID, topic)
	if token := client.Subscribe(gatewayTopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Fprintln(w, token.Error())
		os.Exit(1)
	}

	deviceTopic := fmt.Sprintf("/devices/%s/%s", deviceID, topic)
	if token := client.Subscribe(deviceTopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Fprintln(w, token.Error())
		os.Exit(1)
	}

	time.Sleep(time.Duration(clientDuration) * time.Second)

	detachDevice(deviceID, client, "")

	if token := client.Unsubscribe(gatewayTopic, deviceTopic); token.Wait() && token.Error() != nil {
		fmt.Fprintln(w, token.Error())
		os.Exit(1)
	}

	client.Disconnect(10)
}

// [END iot_subscribe_to_bound_device]
