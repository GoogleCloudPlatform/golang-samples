module github.com/GoogleCloudPlatform/golang-samples

go 1.11

require (
	cloud.google.com/go v0.87.0
	cloud.google.com/go/bigquery v1.19.0
	cloud.google.com/go/bigtable v1.4.0
	cloud.google.com/go/datastore v1.2.0
	cloud.google.com/go/firestore v1.3.0
	cloud.google.com/go/logging v1.0.0
	cloud.google.com/go/pubsub v1.9.1
	cloud.google.com/go/pubsublite v0.8.0
	cloud.google.com/go/spanner v1.21.0
	cloud.google.com/go/storage v1.16.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/aws/aws-sdk-go v1.38.69
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/fluent/fluent-logger-golang v1.6.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.2.0
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/h2non/filetype v1.1.1
	github.com/linkedin/goavro/v2 v2.10.0
	github.com/mailgun/mailgun-go/v3 v3.6.4
	github.com/mailjet/mailjet-apiv3-go v0.0.0-20190724151621-55e56f74078c
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/sendgrid/smtpapi-go v0.6.0 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	go.opencensus.io v0.23.0
	golang.org/x/exp v0.0.0-20210625193404-fa9d1d177d71
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/text v0.3.6
	google.golang.org/api v0.50.0
	google.golang.org/appengine v1.6.7
	google.golang.org/genproto v0.0.0-20210713002101-d411969a0d9a
	google.golang.org/grpc v1.39.0
	google.golang.org/grpc/examples v0.0.0-20200707005602-4258d12073b4
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/sendgrid/sendgrid-go.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
