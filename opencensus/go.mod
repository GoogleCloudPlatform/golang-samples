module github.com/GoogleCloudPlatform/golang-samples/opencensus

go 1.13

require (
	cloud.google.com/go/container v1.0.0 // indirect
	cloud.google.com/go/monitoring v1.0.0
	cloud.google.com/go/spanner v1.25.0
	cloud.google.com/go/trace v1.0.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20211020191508-fa7b610d56d1
	github.com/golang/protobuf v1.5.2
	go.opencensus.io v0.23.0
	golang.org/x/exp v0.0.0-20211025140241-8418b01e8c3b
	google.golang.org/genproto v0.0.0-20211027151537-807f52c398cb
)

replace github.com/GoogleCloudPlatform/golang-samples => ../
