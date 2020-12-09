module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20201204214603-1ac3a85e136e
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20201204214603-1ac3a85e136e
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20201204214603-1ac3a85e136e
	google.golang.org/grpc v1.34.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
