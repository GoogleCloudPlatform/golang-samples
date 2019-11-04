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

# [START getting_started_gce_startup_script]
set -ex

# Talk to the metadata server to get the project id
export GOOGLE_CLOUD_PROJECT=$(curl -s "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google")
echo "Project ID: ${GOOGLE_CLOUD_PROJECT}"

# Install logging monitor. The monitor will automatically pickup logs sent to syslog.
curl -s "https://storage.googleapis.com/signals-agents/logging/google-fluentd-install.sh" | bash
service google-fluentd restart &

# Install Go and git with apt-get.
apt-get install -yq ca-certificates git software-properties-common
add-apt-repository -y ppa:longsleep/golang-backports
apt-get install -yq golang-go

go version

# Get the code.
mkdir /code
cd /code
git clone https://github.com/GoogleCloudPlatform/golang-samples.git

cd golang-samples && git checkout gce && cd - # DO NOT SUBMIT. TODO: remove

cd golang-samples/getting-started/gce

# Build the binary and make sure it's executable.
GO111MODULE=on GOCACHE=on go build -o /usr/bin/gce
chmod +x /usr/bin/gce

# Create the systemd service file.
cat <<EOF > /etc/systemd/system/my-gce-app.service
[Unit]
Description=Example Go app on GCE.

[Service]
Type=simple
ExecStart=/usr/bin/gce
Environment=PORT=80

[Install]
WantedBy=multi-user.target
EOF

chmod 644 /etc/systemd/system/my-gce-app.service

# Start the service.
service my-gce-app start
# [END getting_started_gce_startup_script]
