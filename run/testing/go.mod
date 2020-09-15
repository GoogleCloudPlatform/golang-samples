module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples 4f178324dbb0
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping 4f178324dbb0
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming 4f178324dbb0
	google.golang.org/grpc v1.32.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
