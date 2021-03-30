module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210325165548-a6135c8a474c
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20210325165548-a6135c8a474c
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20210325165548-a6135c8a474c
	golang.org/x/net v0.0.0-20210330142815-c8897c278d10
	google.golang.org/grpc v1.36.1
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
