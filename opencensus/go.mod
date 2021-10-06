module github.com/GoogleCloudPlatform/golang-samples/opencensus

go 1.13

require (
	cloud.google.com/go/container v0.1.0 // indirect
	cloud.google.com/go/monitoring v0.1.0
	cloud.google.com/go/spanner v1.25.0
	cloud.google.com/go/trace v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20210915194233-767ebcbf0013
	github.com/golang/protobuf v1.5.2
	go.opencensus.io v0.23.0
	golang.org/x/exp v0.0.0-20210625193404-fa9d1d177d71
	google.golang.org/genproto v0.0.0-20210924002016-3dee208752a0
)

replace github.com/GoogleCloudPlatform/golang-samples => ../
