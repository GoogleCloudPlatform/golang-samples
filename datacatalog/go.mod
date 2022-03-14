module github.com/GoogleCloudPlatform/golang-samples/datacatalog

go 1.15

require (
	cloud.google.com/go/bigquery v1.29.0
	cloud.google.com/go/datacatalog v1.3.0
	github.com/GoogleCloudPlatform/golang-samples 2bd24627dd5e
	github.com/GoogleCloudPlatform/golang-samples/bigquery v0.0.0
	google.golang.org/api v0.71.0
	google.golang.org/genproto 1973136f34c6
)

replace github.com/GoogleCloudPlatform/golang-samples/bigquery => ../bigquery
