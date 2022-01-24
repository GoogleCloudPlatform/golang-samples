module github.com/GoogleCloudPlatform/golang-samples/bigquery

go 1.11

require (
	cloud.google.com/go/bigquery v1.25.0
	cloud.google.com/go/storage v1.18.2
	github.com/GoogleCloudPlatform/golang-samples v0.0.0
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/googleapis/gax-go/v2 v2.1.1
	github.com/linkedin/goavro/v2 v2.10.1
	google.golang.org/api v0.63.0
	google.golang.org/genproto v0.0.0-20211221195035-429b39de9b1c
	google.golang.org/grpc v1.40.1
	google.golang.org/protobuf v1.27.1
)

replace github.com/GoogleCloudPlatform/golang-samples => ../
