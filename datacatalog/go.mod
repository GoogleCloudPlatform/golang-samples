module github.com/GoogleCloudPlatform/golang-samples/datacatalog

go 1.11

require (
	cloud.google.com/go/bigquery v1.27.0
	cloud.google.com/go/datacatalog v1.1.0
	github.com/GoogleCloudPlatform/golang-samples 2082aefea4f3
	github.com/GoogleCloudPlatform/golang-samples/bigquery v0.0.0
	google.golang.org/api v0.67.0
	google.golang.org/genproto 7721543eae58
)

replace github.com/GoogleCloudPlatform/golang-samples/bigquery => ../bigquery
