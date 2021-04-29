module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210426214444-1bff782b3539
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20210426214444-1bff782b3539
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20210426214444-1bff782b3539
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781
	google.golang.org/grpc v1.37.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
