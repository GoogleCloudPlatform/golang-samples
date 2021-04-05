module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210405151107-4e093192e115
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20210405151107-4e093192e115
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20210405151107-4e093192e115
	golang.org/x/net v0.0.0-20210331212208-0fccb6fa2b5c
	google.golang.org/grpc v1.36.1
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
