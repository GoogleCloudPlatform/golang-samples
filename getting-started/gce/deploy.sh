#! /bin/sh

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

if [ -z "$DEPLOY_DIR" ]; then
  echo "Must set \$DEPLOY_DIR. For example: DEPLOY_DIR=gs://my-bucket"
  exit 1
fi

if [ -z "$DEPLOY_FILENAME" ]; then
  echo "Must set \$DEPLOY_FILENAME. For example: DEPLOY_FILENAME=app.tar.gz"
  exit 1
fi

gcloud builds submit --substitutions=_DEPLOY_DIR="$DEPLOY_DIR",_DEPLOY_FILENAME="$DEPLOY_FILENAME"

ZONE=us-central1-f

gcloud compute instances create my-app-instance \
    --image-family=debian-10 \
    --image-project=debian-cloud \
    --machine-type=g1-small \
    --scopes userinfo-email,cloud-platform \
    --metadata-from-file startup-script=startup-script.sh \
    --metadata app-location="${DEPLOY_DIR%/}/${DEPLOY_FILENAME}" \
    --zone $ZONE \
    --tags http-server

gcloud compute firewall-rules create default-allow-http-80 \
    --allow tcp:80 \
    --source-ranges 0.0.0.0/0 \
    --target-tags http-server \
    --description "Allow port 80 access to http-server"
