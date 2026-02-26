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

# Install logging monitor. The monitor will automatically pickup logs sent to syslog.
curl "https://storage.googleapis.com/signals-agents/logging/google-fluentd-install.sh" --output google-fluentd-install.sh
checksum=$(sha256sum google-fluentd-install.sh | awk '{print $1;}')
if [ "$checksum" != "ec78e9067f45f6653a6749cf922dbc9d79f80027d098c90da02f71532b5cc967" ]; then
    echo "Checksum does not match"
    exit 1
fi
chmod +x google-fluentd-install.sh && ./google-fluentd-install.sh
service google-fluentd restart &

APP_LOCATION=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/app-location" -H "Metadata-Flavor: Google")
gcloud storage cp "$APP_LOCATION" app.tar.gz
tar -xzf app.tar.gz

# Start the service included in app.tar.gz.
service my-app start
# [END getting_started_gce_startup_script]
