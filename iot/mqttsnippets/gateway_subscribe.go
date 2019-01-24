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

// [START iot_subscribe_to_bound_device]

// subcribeToBoundDevice subscribes and listens on a topic for a bound device
func subscribeToBoundDevice(w io.Writer, projectID string, region string, registryID string, gatewayID string, deviceID string, privateKeyPath string, algorithm string) {

	const (
		mqttBrokerURL   = "tls://mqtt.googleapis.com:443"
		protocolVersion = 4 // corresponds to MQTT 3.1.1
	)

	// onConnect defines the on connect handler which resets backoff variables.
	var onConnect mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Printf("Client connected: %t\n", client.IsConnected())
	}

	// onMessage defines the message handler for the mqtt client.
	var onMessage mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Topic: %s\n", msg.Topic())
		fmt.Printf("Message: %s\n", msg.Payload())
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

	//subscribe to the topic and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	gatewayConfigTopic := fmt.Sprintf("/devices/%s/config", gatewayID)
	if token := client.Subscribe(gatewayConfigTopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	deviceConfigTopic := fmt.Sprintf("/devices/%s/config", deviceID)
	if token := client.Subscribe(deviceConfigTopic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	time.Sleep(60 * time.Second)

	detachDevice(deviceID, client, "")

	if token := client.Unsubscribe(deviceConfigTopic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Unsubscribe(deviceConfigTopic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	client.Disconnect(250)
}

// [END iot_subscribe_to_bound_device]
