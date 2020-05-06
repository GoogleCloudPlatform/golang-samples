#!/bin/bash

set -ex

if [ -z $KOKORO_BUILD_ARTIFACTS_SUBDIR ]; then
  echo "This should only be run from Kokoro."
  exit 1
fi

gcloud -q components update
gcloud -q components install app-engine-go

# Set config.
gcloud config set disable_prompts True
gcloud config set project $GOLANG_SAMPLES_PROJECT_ID
gcloud config set app/promote_by_default false
gcloud auth activate-service-account --key-file "$GOOGLE_APPLICATION_CREDENTIALS"

# Diagnostic information.
gcloud info
