module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210416175205-e85b572b9ebb
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20210416175205-e85b572b9ebb
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20210416175205-e85b572b9ebb
	golang.org/x/net v0.0.0-20210415231046-e915ea6b2b7d
	google.golang.org/grpc v1.37.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
