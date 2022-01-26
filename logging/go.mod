module github.com/GoogleCloudPlatform/golang-samples/logging

go 1.11

replace github.com/GoogleCloudPlatform/golang-samples => ./..

require (
	cloud.google.com/go/logging v1.4.2
	github.com/GoogleCloudPlatform/golang-samples v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0
	google.golang.org/api v0.63.0
	google.golang.org/genproto v0.0.0-20211221195035-429b39de9b1c
)
