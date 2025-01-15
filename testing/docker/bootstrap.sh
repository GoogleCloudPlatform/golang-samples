#!/bin/sh

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

apt-get update && apt-get install -y imagemagick graphviz
rm -rf /var/lib/apt/lists/*

# Install tools used in testing workflows
GOPROXY=direct go install github.com/GoogleCloudPlatform/golang-samples/testing/gimmeproj@latest
go install golang.org/x/tools/cmd/goimports@latest
go install gotest.tools/gotestsum@latest

(
# Get the SDK tar and untar it.
cd /tmp

TARFILE=google-cloud-sdk.tar.gz
wget https://dl.google.com/dl/cloudsdk/release/$TARFILE
tar xzf $TARFILE
rm $TARFILE

# Install the SDK
./google-cloud-sdk/install.sh \
    --usage-reporting false \
    --path-update false \
    --command-completion false

./google-cloud-sdk/bin/gcloud -q components update
./google-cloud-sdk/bin/gcloud -q components install app-engine-go
)
