module github.com/GoogleCloudPlatform/golang-samples/datacatalog

go 1.15

require (
	cloud.google.com/go/bigquery v1.31.0
	cloud.google.com/go/datacatalog v1.3.0
	github.com/GoogleCloudPlatform/golang-samples c0583a59e8b6
	github.com/GoogleCloudPlatform/golang-samples/bigquery v0.0.0
	google.golang.org/api v0.74.0
	google.golang.org/genproto 2d67ff6cf2b4
)

replace github.com/GoogleCloudPlatform/golang-samples/bigquery => ../bigquery
