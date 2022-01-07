module github.com/GoogleCloudPlatform/golang-samples

go 1.11

require (
	cloud.google.com/go v0.98.0 // indirect
	cloud.google.com/go/bigquery v1.25.0
	cloud.google.com/go/datastore v1.2.0
	cloud.google.com/go/gaming v1.0.0
	cloud.google.com/go/logging v1.0.0
	cloud.google.com/go/storage v1.18.2
	cloud.google.com/go/vision v1.0.0
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/h2non/filetype v1.1.1
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/api v0.61.0
	google.golang.org/genproto v0.0.0-20211203200212-54befc351ae9
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
