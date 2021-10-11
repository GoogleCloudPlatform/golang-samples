module github.com/GoogleCloudPlatform/golang-samples/opencensus

go 1.13

require (
	cloud.google.com/go/container v0.1.0 // indirect
	cloud.google.com/go/monitoring v0.97.0
	cloud.google.com/go/spanner v1.25.0
	cloud.google.com/go/trace v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/GoogleCloudPlatform/golang-samples 553d451c8611
	github.com/golang/protobuf v1.5.2
	go.opencensus.io v0.23.0
	golang.org/x/exp 95152d363a1c
	google.golang.org/genproto 270636b82663
)

replace github.com/GoogleCloudPlatform/golang-samples => ../
