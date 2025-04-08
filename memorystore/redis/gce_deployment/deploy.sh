#!/bin/sh

# [START memorystore_deploy_sh]
if [ -z "$REDISHOST" ]; then
  echo "Must set \$REDISHOST. For example: REDISHOST=127.0.0.1"
  exit 1
fi

if [ -z "$REDISPORT" ]; then
  echo "Must set \$REDISPORT. For example: REDISPORT=6379"
  exit 1
fi

if [ -z "$GCS_BUCKET_NAME" ]; then
  echo "Must set \$GCS_BUCKET_NAME. For example: GCS_BUCKET_NAME=my-bucket"
  exit 1
fi

if [ -z "$ZONE" ]; then
  ZONE=$(gcloud config get-value compute/zone -q)
  echo "$ZONE"
fi


# Cross compile the app for linux/amd64
GOOS=linux GOARCH=amd64 go build -v -o app ../main.go
# Add the app binary
tar -cvf app.tar app
# Copy to GCS bucket
gsutil cp app.tar gs://"$GCS_BUCKET_NAME"/gce/

# Create an instance
gcloud compute instances create my-instance \
    --image-family=debian-9 \
    --image-project=debian-cloud \
    --machine-type=g1-small \
    --scopes cloud-platform \
    --metadata-from-file startup-script=startup-script.sh \
    --metadata gcs-bucket="$GCS_BUCKET_NAME",redis-host="$REDISHOST",redis-port="$REDISPORT" \
    --zone "$ZONE" \
    --tags http-server

gcloud compute firewall-rules create allow-http-server-8080 \
    --allow tcp:8080 \
    --source-ranges 0.0.0.0/0 \
    --target-tags http-server \
    --description "Allow port 8080 access to http-server"
# [END memorystore_deploy_sh]
