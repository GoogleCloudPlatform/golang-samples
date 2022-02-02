module github.com/GoogleCloudPlatform/golang-samples/iot

go 1.11

require (
	cloud.google.com/go/kms v1.1.0 // indirect
	cloud.google.com/go/pubsub v1.3.1
	github.com/GoogleCloudPlatform/golang-samples v0.0.0
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	google.golang.org/api v0.65.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../
