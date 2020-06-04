module github.com/GoogleCloudPlatform/golang-samples

go 1.11

require (
	cloud.google.com/go v0.57.0
	cloud.google.com/go/bigquery v1.8.0
	cloud.google.com/go/bigtable v1.3.0
	cloud.google.com/go/datastore v1.1.0
	cloud.google.com/go/firestore v1.2.0
	cloud.google.com/go/logging v1.0.0
	cloud.google.com/go/pubsub v1.4.0
	cloud.google.com/go/spanner v1.6.0
	cloud.google.com/go/storage v1.9.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.1
	github.com/aws/aws-sdk-go v1.31.10
	github.com/bmatcuk/doublestar v1.3.1
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/fluent/fluent-logger-golang v1.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.4.1
	github.com/google/uuid v1.1.1
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.2
	github.com/h2non/filetype v1.0.12
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/lib/pq v1.6.0
	github.com/linkedin/goavro/v2 v2.9.7
	github.com/mailgun/mailgun-go/v3 v3.6.4
	github.com/mailjet/mailjet-apiv3-go v0.0.0-20190724151621-55e56f74078c
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/sendgrid/smtpapi-go v0.6.0 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	go.opencensus.io v0.22.3
	golang.org/x/exp v0.0.0-20200513190911-00229845015e
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20200604042327-9b20fe4cabe8
	google.golang.org/api v0.26.0
	google.golang.org/appengine v1.6.6
	google.golang.org/genproto v0.0.0-20200604104852-0b0486081ffb
	google.golang.org/grpc v1.29.1
	gopkg.in/sendgrid/sendgrid-go.v2 v2.0.0
	gopkg.in/yaml.v2 v2.3.0
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
