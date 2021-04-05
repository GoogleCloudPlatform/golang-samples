module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20201216233243-555da975282a
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20201216233243-555da975282a
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20201216233243-555da975282a
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4
	google.golang.org/grpc v1.36.1
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
