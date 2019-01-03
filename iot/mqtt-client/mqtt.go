package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

const privateKeyPath = "../manager/resources/ec_private.pem"
const algorithm = "ES256"

const clientID = "projects/hongalex-iot-samples/locations/us-central1/registries/testRegistry/devices/testDevice"

// [START iot_mqtt_jwt]

// createJWT creates a Cloud IoT Core JWT for the given project id, signed with the given
// private key and an algorithm (RS256 or ES256)
func createJWT(projectID string, privateKeyPath string, algorithm string) (string, error) {
	claims := jwt.StandardClaims{
		Audience:  projectID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Unix() + 20*60,
	}

	keyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(algorithm), claims)

	if algorithm == "RS256" {
		privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
		return token.SignedString(privKey)
	} else if algorithm == "ES256" {
		privKey, _ := jwt.ParseECPrivateKeyFromPEM(keyBytes)
		return token.SignedString(privKey)
	} else {
		return "", errors.New("Cannot find JWT algorithm")
	}

}

// [END iot_mqtt_jwt]

func main() {
	res, _ := createJWT(projectID, privateKeyPath, algorithm)
	fmt.Println(res)

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tls://mqtt.googleapis.com:443")
	opts.SetClientID(clientID)
	opts.SetUsername("unused")
	opts.SetPassword(res)
	opts.SetProtocolVersion(4)
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe("/devices/testDevice/commands/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	time.Sleep(60 * time.Second)

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
