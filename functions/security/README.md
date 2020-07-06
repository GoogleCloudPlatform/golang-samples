# Function Security Samples

## System Tests

1. `gcloud alpha functions deploy HelloHTTP --region=us-central1 --runtime=go111 --trigger-http --source=../helloworld`

1. `export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service/account"`

1. `export BASE_URL=https://REGION-PROJECT.cloudfunctions.net/`

1. `go test -v ./security_system_test.go`