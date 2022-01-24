module github.com/GoogleCloudPlatform/golang-samples/datacatalog

go 1.11

require github.com/GoogleCloudPlatform/golang-samples v0.0.0

require (
	cloud.google.com/go/bigquery v1.25.0
	cloud.google.com/go/datacatalog v1.0.0
	github.com/GoogleCloudPlatform/golang-samples/bigquery v0.0.0
	google.golang.org/api v0.63.0
	google.golang.org/genproto v0.0.0-20211221195035-429b39de9b1c
)

replace github.com/GoogleCloudPlatform/golang-samples => ../

replace github.com/GoogleCloudPlatform/golang-samples/bigquery => ../bigquery
