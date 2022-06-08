module github.com/GoogleCloudPlatform/golang-samples/datacatalog

go 1.15

require (
	cloud.google.com/go/bigquery v1.31.0
	cloud.google.com/go/datacatalog v1.3.0
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20220204002944-f20d8abe1519
	github.com/GoogleCloudPlatform/golang-samples/bigquery v0.0.0
	google.golang.org/api v0.74.0
	google.golang.org/genproto v0.0.0-20220518221133-4f43b3371335
)

replace github.com/GoogleCloudPlatform/golang-samples/bigquery => ../bigquery
