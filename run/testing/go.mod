module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20201113163810-8bf1fa39bf2b
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20201113163810-8bf1fa39bf2b
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20201113163810-8bf1fa39bf2b
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	google.golang.org/grpc v1.33.2
)

replace github.com/GoogleCloudPlatform/golang-samples => ../..

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping => ../grpc-ping

replace github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping-streaming => ../grpc-ping-streaming
