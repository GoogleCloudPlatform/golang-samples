#!/bin/bash

# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Command gimmeproj provides access to a pool of projects.
#
# The metadata about the project pool is stored in Cloud Datastore in a meta-project.
# Projects are leased for a certain duration, and automatically returned to the pool when the lease expires.
# Projects should be returned before the lease expires.

set -ex

if [ -z $KOKORO_BUILD_ARTIFACTS_SUBDIR ]; then
  echo "This should only be run from Kokoro."
  exit 1
fi

gcloud -q components update
gcloud -q components install app-engine-go
gcloud -q components install beta # Needed for Cloud Run E2E tests until --pack goes to GA
gcloud -q components install alpha # Needed for Cloud Run E2E tests until --use-http2 goes GA

# Set config.
gcloud config set disable_prompts True
gcloud config set project $GOLANG_SAMPLES_PROJECT_ID
gcloud config set app/promote_by_default false
gcloud config set core/account $GOLANG_SAMPLES_PROJECT_ID@$GOLANG_SAMPLES_PROJECT_ID.iam.gserviceaccount.com
# gcloud auth activate-service-account --key-file "$GOOGLE_APPLICATION_CREDENTIALS"
gcloud auth activate-service-account --key-file "$GOOGLE_APPLICATION_CREDENTIALS" --project="golang-samples-tests"

# Diagnostic information.
gcloud info
gcloud config list
