#! /bin/bash

# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

if [ ! $(dirname $0) = "." ]; then
  echo "Must run $(basename $0) from the gce_deployment directory."
  exit 1
fi

if [ -z "$BOOKSHELF_DEPLOY_LOCATION" ]; then
  echo "Must set \$BOOKSHELF_DEPLOY_LOCATION. For example: BOOKSHELF_DEPLOY_LOCATION=gs://my-bucket/bookshelf-VERSION.tar"
  exit 1
fi

TMP=$(mktemp -d -t gce-deploy-XXXXXX)

# [START cross_compile]
# Cross compile the app for linux/amd64
GOOS=linux GOARCH=amd64 go build -v -o $TMP/app ../app
# [END cross_compile]

# [START tar]
# Add the app binary
tar -c -f $TMP/bundle.tar -C $TMP app

# Add static files.
tar -u -f $TMP/bundle.tar -C ../app templates
# [END tar]

# [START gcs_push]
# BOOKSHELF_DEPLOY_LOCATION is something like "gs://my-bucket/bookshelf-VERSION.tar".
gsutil cp $TMP/bundle.tar $BOOKSHELF_DEPLOY_LOCATION
# [END gcs_push]

rm -rf $TMP
