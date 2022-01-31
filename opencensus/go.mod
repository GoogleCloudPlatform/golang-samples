module github.com/GoogleCloudPlatform/golang-samples/opencensus

go 1.13

require (
	cloud.google.com/go/container v1.0.0 // indirect
	cloud.google.com/go/monitoring v1.2.0
	cloud.google.com/go/spanner v1.29.0
	cloud.google.com/go/trace v1.0.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20220129100446-dc42837e161f
	github.com/golang/protobuf v1.5.2
	go.opencensus.io v0.23.0
	golang.org/x/exp v0.0.0-20220128181451-c853b6ddb95e
	google.golang.org/genproto v0.0.0-20220126215142-9970aeb2e350
)

replace github.com/GoogleCloudPlatform/golang-samples => ../
