#!/bin/bash

set -ex

if [ $TRAVIS != "true" ]; then
  echo "This should only be run from travis."
  exit 1
fi

# Get the SDK tar and untar it.
TARFILE=google-cloud-sdk.tar.gz
wget https://dl.google.com/dl/cloudsdk/release/$TARFILE
tar xzf $TARFILE
rm $TARFILE

# Install the SDK
./google-cloud-sdk/install.sh \
  --usage-reporting false \
  --path-update false \
  --command-completion false

gcloud components update

# Set config.
gcloud config set disable_prompts True
gcloud config set project $GOLANG_SAMPLES_PROJECT_ID
gcloud config set app/promote_by_default false
gcloud auth activate-service-account --key-file "$GOOGLE_APPLICATION_CREDENTIALS"

# Diagnostic information.
gcloud info
