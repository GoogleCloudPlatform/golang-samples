module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210319194002-1e35cc23a955
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20210319194002-1e35cc23a955
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20210319194002-1e35cc23a955
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4
	google.golang.org/grpc v1.36.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
