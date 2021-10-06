module github.com/GoogleCloudPlatform/golang-samples/trace

go 1.13

replace github.com/GoogleCloudPlatform/golang-samples => ../

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	go.opencensus.io v0.23.0
)

require (
	cloud.google.com/go v0.94.0 // indirect
	cloud.google.com/go/container v0.1.0 // indirect
	cloud.google.com/go/monitoring v0.1.0 // indirect
	cloud.google.com/go/trace v0.1.0 // indirect
	github.com/aws/aws-sdk-go v1.38.69 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
