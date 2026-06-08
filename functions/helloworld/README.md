# Hello World Samples

## System Tests

### Storage

1. `export BUCKET_NAME=...`

1. `gcloud storage buckets create gs://$BUCKET_NAME`

1. `gcloud functions deploy HelloGCS --runtime=go113 --entry-point=HelloGCS --trigger-resource=$BUCKET_NAME --trigger-event=providers/cloud.storage/eventTypes/object.change`

1. `go test -v ./hello_cloud_storage_system_test.go`

### HTTP

1. `gcloud functions deploy HelloHTTP --region=us-central1 --runtime=go113 --trigger-http`

1. `export BASE_URL=https://REGION-PROJECT.cloudfunctions.net/`

1. `go test -v ./hello_http_system_test.go`

### Pub/Sub

1. `export FUNCTIONS_TOPIC=example-topic`

1. `gcloud functions deploy HelloPubSub --runtime=go113 --trigger-topic=$FUNCTIONS_TOPIC`

1. `go test -v ./hello_pubsub_system_test.go`
