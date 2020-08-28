module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20200828000024-04937a2b81cf
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20200828000024-04937a2b81cf
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20200828000024-04937a2b81cf
	google.golang.org/grpc v1.31.1
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
