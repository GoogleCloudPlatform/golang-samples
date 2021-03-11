module github.com/GoogleCloudPlatform/golang-samples

go 1.11

require (
	cloud.google.com/go v0.79.0
	cloud.google.com/go/bigquery v1.16.0
	cloud.google.com/go/bigtable v1.8.0
	cloud.google.com/go/datastore v1.5.0
	cloud.google.com/go/firestore v1.5.0
	cloud.google.com/go/logging v1.3.0
	cloud.google.com/go/pubsub v1.10.1
	cloud.google.com/go/pubsublite v0.7.0
	cloud.google.com/go/spanner v1.15.0
	cloud.google.com/go/storage v1.14.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/aws/aws-sdk-go v1.37.28
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eclipse/paho.mqtt.golang v1.3.2
	github.com/fluent/fluent-logger-golang v1.5.0
	github.com/go-chi/chi v4.1.2+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.3 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.2.0
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/h2non/filetype v1.1.1
	github.com/linkedin/goavro/v2 v2.10.0
	github.com/mailgun/mailgun-go/v3 v3.6.4
	github.com/mailjet/mailjet-apiv3-go v0.0.0-20201009050126-c24bc15a9394
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/sendgrid/smtpapi-go v0.6.6 // indirect
	github.com/tinylib/msgp v1.1.5 // indirect
	go.opencensus.io v0.23.0
	golang.org/x/exp v0.0.0-20210220032938-85be41e4509f
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/oauth2 v0.0.0-20210220000619-9bb904979d93
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
	golang.org/x/text v0.3.5
	google.golang.org/api v0.41.0
	google.golang.org/appengine v1.6.7
	google.golang.org/genproto v0.0.0-20210310155132-4ce2db91004e
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/examples v0.0.0-20210310172623-a45f13b16073
	google.golang.org/protobuf v1.25.0
	gopkg.in/sendgrid/sendgrid-go.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
