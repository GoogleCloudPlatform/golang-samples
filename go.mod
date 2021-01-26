module github.com/GoogleCloudPlatform/golang-samples

go 1.11

require (
	cloud.google.com/go v0.74.0
	cloud.google.com/go/bigquery v1.14.0
	cloud.google.com/go/bigtable v1.4.0
	cloud.google.com/go/datastore v1.2.0
	cloud.google.com/go/firestore v1.3.0
	cloud.google.com/go/logging v1.0.0
	cloud.google.com/go/pubsub v1.6.1
	cloud.google.com/go/spanner v1.10.0
	cloud.google.com/go/storage v1.11.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.4
	github.com/aws/aws-sdk-go v1.36.20
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eclipse/paho.mqtt.golang v1.3.1
	github.com/fluent/fluent-logger-golang v1.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.4
	github.com/google/uuid v1.1.4
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/h2non/filetype v1.1.0
	github.com/lib/pq v1.9.0
	github.com/linkedin/goavro/v2 v2.10.0
	github.com/mailgun/mailgun-go/v3 v3.6.4
	github.com/mailjet/mailjet-apiv3-go v0.0.0-20190724151621-55e56f74078c
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/sendgrid/smtpapi-go v0.6.0 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	go.opencensus.io v0.22.5
	golang.org/x/exp v0.0.0-20201229011636-eab1b5eb1a03
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5
	golang.org/x/text v0.3.4
	google.golang.org/api v0.36.0
	google.golang.org/appengine v1.6.7
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d
	google.golang.org/grpc v1.34.0
	google.golang.org/grpc/examples v0.0.0-20200707005602-4258d12073b4
	google.golang.org/protobuf v1.25.0
	gopkg.in/sendgrid/sendgrid-go.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
