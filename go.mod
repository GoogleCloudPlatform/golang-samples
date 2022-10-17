module github.com/GoogleCloudPlatform/golang-samples

go 1.15

require (
	cloud.google.com/go/bigquery v1.42.0
	cloud.google.com/go/datastore v1.2.0
	cloud.google.com/go/errorreporting v0.1.0
	cloud.google.com/go/logging v1.0.0
	cloud.google.com/go/storage v1.24.0
	cloud.google.com/go/vision v1.2.0
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/googleapis/gax-go/v2 v2.6.0 // indirect
	github.com/h2non/filetype v1.1.1
	golang.org/x/net v0.0.0-20221017152216-f25eb7ecb193 // indirect
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/api v0.99.0
	google.golang.org/genproto v0.0.0-20221014213838-99cd37c6964a
	google.golang.org/grpc v1.50.1 // indirect
	google.golang.org/protobuf v1.28.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
