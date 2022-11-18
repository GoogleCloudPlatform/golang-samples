module github.com/GoogleCloudPlatform/golang-samples/spanner

go 1.17

require (
	cloud.google.com/go v0.105.0
	cloud.google.com/go/kms v1.6.0
	cloud.google.com/go/spanner v1.39.1-0.20221118045942-492382ef4872
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20220328195317-2183fb3440ed
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/googleapis/gax-go/v2 v2.6.0
	go.opencensus.io v0.23.0
	google.golang.org/api v0.102.0
	google.golang.org/genproto v0.0.0-20221027153422-115e99e71e1c
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
)

require (
	cloud.google.com/go/compute v1.12.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.1 // indirect
	cloud.google.com/go/iam v0.7.0 // indirect
	cloud.google.com/go/longrunning v0.3.0 // indirect
	cloud.google.com/go/monitoring v1.8.0 // indirect
	cloud.google.com/go/storage v1.27.0 // indirect
	cloud.google.com/go/trace v1.4.0 // indirect
	github.com/aws/aws-sdk-go v1.44.105 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe // indirect
	github.com/cncf/xds/go v0.0.0-20220520190051-1e77728a1eaa // indirect
	github.com/envoyproxy/go-control-plane v0.10.3 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.8 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/prometheus/prometheus v0.38.0 // indirect
	golang.org/x/net v0.0.0-20221014081412-f15817d10f9b // indirect
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
)

// replace cloud.google.com/go/spanner v1.39.0 => ../../google-cloud-go-changes/google-cloud-go-proto-support-v2/spanner
// replace cloud.google.com/go/spanner v1.39.0 => cloud.google.com/go/spanner proto-column-enhancement-alpha
// replace cloud.google.com/go/spanner => github.com/harshachinta/google-cloud-go/spanner v0.0.0-20221117175741-ca271287f57d

replace google.golang.org/genproto => github.com/cloudspannerecosystem/temp-resources v0.0.0-20221117065524-b1f320c13693
