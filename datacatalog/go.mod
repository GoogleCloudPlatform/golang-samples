module github.com/GoogleCloudPlatform/golang-samples/datacatalog

go 1.15

require (
	cloud.google.com/go/bigquery v1.39.1-0.20220908212230-60e120cef30c
	cloud.google.com/go/datacatalog v1.3.0
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-20220623032819-b45447226621
	github.com/GoogleCloudPlatform/golang-samples/bigquery v0.0.0
	google.golang.org/api v0.94.0
	google.golang.org/genproto v0.0.0-20220902135211-223410557253
)

replace github.com/GoogleCloudPlatform/golang-samples/bigquery => ../bigquery
