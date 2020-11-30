module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20201125203603-a12e2fd198c0
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20201125203603-a12e2fd198c0
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20201125203603-a12e2fd198c0
	google.golang.org/grpc v1.33.2
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
