module github.com/GoogleCloudPlatform/golang-samples

go 1.19

require (
	cloud.google.com/go/batch v0.4.0
	cloud.google.com/go/bigquery v1.43.0
	cloud.google.com/go/compute v1.12.1
	cloud.google.com/go/datastore v1.9.0
	cloud.google.com/go/errorreporting v0.2.0
	cloud.google.com/go/logging v1.5.0
	cloud.google.com/go/storage v1.28.0
	cloud.google.com/go/vision v1.2.0
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/h2non/filetype v1.1.3
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783
	google.golang.org/api v0.103.0
	google.golang.org/genproto v0.0.0-20221109142239-94d6d90a7d66
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go v0.105.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.1 // indirect
	cloud.google.com/go/iam v0.7.0 // indirect
	cloud.google.com/go/longrunning v0.3.0 // indirect
	cloud.google.com/go/vision/v2 v2.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/net v0.0.0-20221014081412-f15817d10f9b // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.50.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
