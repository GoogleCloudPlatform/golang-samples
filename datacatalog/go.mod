module github.com/GoogleCloudPlatform/golang-samples/datacatalog

go 1.15

require (
	cloud.google.com/go/bigquery v1.34.1
	cloud.google.com/go/datacatalog v1.3.0
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20220623032819-b45447226621
	github.com/GoogleCloudPlatform/golang-samples/bigquery v0.0.0
	google.golang.org/api v0.85.0
	google.golang.org/genproto v0.0.0-20220624142145-8cd45d7dbd1f
)

replace github.com/GoogleCloudPlatform/golang-samples/bigquery => ../bigquery
