module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210618111007-7e892a0cddf1
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20210618111007-7e892a0cddf1
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20210618111007-7e892a0cddf1
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420
	google.golang.org/grpc v1.38.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
