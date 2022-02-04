module github.com/GoogleCloudPlatform/golang-samples/datacatalog

go 1.11

require (
	cloud.google.com/go/bigquery v1.26.0
	cloud.google.com/go/datacatalog v1.0.0
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20220204002944-f20d8abe1519
	github.com/GoogleCloudPlatform/golang-samples/bigquery v0.0.0
	google.golang.org/api v0.66.0
	google.golang.org/genproto v0.0.0-20220201184016-50beb8ab5c44
)

replace github.com/GoogleCloudPlatform/golang-samples/bigquery => ../bigquery
