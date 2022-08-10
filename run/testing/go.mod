module github.com/GoogleCloudPlatform/golang-samples/run/testing

go 1.15

require (
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20220204002944-f20d8abe1519
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping v0.0.0-20211008220018-553d451c8611
	github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming v0.0.0-20211008220018-553d451c8611
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9
	google.golang.org/grpc v1.47.0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../../
