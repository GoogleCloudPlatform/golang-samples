module github.com/GoogleCloudPlatform/golang-samples/opencensus

go 1.13

require (
	cloud.google.com/go/container v1.0.0 // indirect
	cloud.google.com/go/monitoring v1.1.0
	cloud.google.com/go/spanner v1.28.0
	cloud.google.com/go/trace v1.0.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20211221192803-31bb00e8dd74
	github.com/golang/protobuf v1.5.2
	go.opencensus.io v0.23.0
	golang.org/x/exp v0.0.0-20220104160115-025e73f80486
	google.golang.org/genproto v0.0.0-20211223182754-3ac035c7e7cb
)

replace github.com/GoogleCloudPlatform/golang-samples => ../
