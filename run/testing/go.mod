module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20211008220018-553d451c8611
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20211008220018-553d451c8611
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20211008220018-553d451c8611
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	google.golang.org/grpc v1.40.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
