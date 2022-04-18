module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20220415190337-c0583a59e8b6
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20220415190337-c0583a59e8b6
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20220415190337-c0583a59e8b6
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5
	google.golang.org/grpc v1.45.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../../
