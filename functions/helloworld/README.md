# Hello World Samples

## System Tests

### Storage

1. `export BUCKET_NAME=...`

1. `gsutil mb gs://$BUCKET_NAME`

1. `gcloud alpha functions deploy HelloGCS --runtime=go111 --entry-point=HelloGCS --trigger-resource=$BUCKET_NAME --trigger-event=providers/cloud.storage/eventTypes/object.change`

1. `go test -v ./hello_cloud_storage_system_test.go`

### HTTP

1. `gcloud alpha functions deploy HelloHTTP --region=us-central1 --runtime=go111 --trigger-http`

1. `export BASE_URL=https://REGION-PROJECT.cloudfunctions.net/`

1. `go test -v ./hello_http_system_test.go`

### Pub/Sub

1. `export FUNCTIONS_TOPIC=example-topic`

1. `gcloud alpha functions deploy HelloPubSub --runtime=go111 --trigger-topic=$FUNCTIONS_TOPIC`

1. `go test -v ./hello_pubsub_system_test.go`
