module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210226184207-920b50d04dd4
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20210226184207-920b50d04dd4
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20210226184207-920b50d04dd4
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	google.golang.org/grpc v1.36.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
