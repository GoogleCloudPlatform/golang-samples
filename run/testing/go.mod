module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	cloud.google.com/go/logging v1.0.0
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20201216233243-555da975282a
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20201216233243-555da975282a
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20201216233243-555da975282a
	google.golang.org/api v0.36.0
	google.golang.org/grpc v1.34.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
