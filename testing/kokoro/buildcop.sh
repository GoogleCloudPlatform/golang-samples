#!/bin/bash

# Copyright 2020 Google LLC
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

# This is manually copied to
# gs://cloud-devrel-kokoro-resources/trampoline/buildcop.sh.

# See https://github.com/googleapis/repo-automation-bots/tree/master/packages/buildcop.

type gcloud > /dev/null 2>&1 || { echo >&2 "gcloud is required! Not sending logs to the Build Cop Bot." && exit 1 }

gcloud auth activate-service-account --key-file $KOKORO_GFILE_DIR/kokoro-trampoline.service-account.json

REPO=$(echo $KOKORO_GITHUB_COMMIT_URL | cut -d/ -f 4,5)

if [ -z ${INSTALLATION_ID+x} ]; then
    if [ $REPO = *"GoogleCloudPlatform"* ]; then
        INSTALLATION_ID=5943459
    elif [ $REPO = *"googleapis"* ]; then
        INSTALLATION_ID=6370238
    else
        echo >&2 "INSTALLATION_ID unset. If your repo is part of GoogleCloudPlatform or googleapis and you see this error,"
        echo >&2 "file an issue at https://github.com/googleapis/repo-automation-bots/issues."
        echo >&2 "Otherwise, set INSTALLATION_ID with the numeric installation ID before calling buildcop.sh"
        echo >&2 "See https://github.com/apps/build-cop-bot/".
        exit 1
    fi
fi

# Loop over all sponge_log.xml files.
shopt -s globstar
for log in **/sponge_log.xml; do
    XML=$(base64 -w 0 $log)

    # See https://github.com/apps/build-cop-bot/installations/5943459.
    MESSAGE=$(cat <<EOF
    {
        "Name": "buildcop",
        "Type" : "function",
        "Location": "us-central1",
        "installation": {"id": "$INSTALLATION_ID"},
        "repo": "$REPO",
        "buildID": "$KOKORO_GIT_COMMIT",
        "buildURL": "[Build Status](https://source.cloud.google.com/results/invocations/$KOKORO_BUILD_ID), [Sponge](http://sponge2/$KOKORO_BUILD_ID)",
        "xunitXML": "$XML"
    }
EOF
    )

    gcloud pubsub topics publish passthrough --project=repo-automation-bots --message="$MESSAGE"
done