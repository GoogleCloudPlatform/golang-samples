module github.com/GoogleCloudPlatform/golang-samples

go 1.15

require (
	cloud.google.com/go/bigquery v1.31.0
	cloud.google.com/go/datastore v1.6.0
	cloud.google.com/go/errorreporting v0.2.0
	cloud.google.com/go/gaming v1.2.0
	cloud.google.com/go/logging v1.4.2
	cloud.google.com/go/storage v1.22.0
	cloud.google.com/go/vision v1.2.0
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/h2non/filetype v1.1.3
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	google.golang.org/api v0.74.0
	google.golang.org/genproto v0.0.0-20220414192740-2d67ff6cf2b4
	google.golang.org/protobuf v1.28.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
