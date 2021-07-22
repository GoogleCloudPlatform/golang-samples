module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210622144803-cc3dd230907c
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20210622144803-cc3dd230907c
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20210622144803-cc3dd230907c
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420
	google.golang.org/grpc v1.39.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
