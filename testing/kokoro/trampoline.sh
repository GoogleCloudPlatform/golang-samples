#!/bin/bash

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

set -x

date

cd github/golang-samples || exit 1

# SIGNIFICANT_CHANGES="$(git --no-pager diff --name-only main..HEAD | grep -Ev '(\.md$|^\.github)' || true)"

# # If this is a PR with only insignificant changes, don't run any tests.
# if [[ -n ${KOKORO_GITHUB_PULL_REQUEST_NUMBER:-} ]] && [[ -z "$SIGNIFICANT_CHANGES" ]]; then
#   echo "No big changes. Not running any tests."
#   exit 0
# fi

cd - || exit 1

function cleanup() {
    # This clutters logs with chmod errors when used in a docker context.
    # TODO(muncus): determine if we ever need this, or make it conditional.
    # chmod +x "${KOKORO_GFILE_DIR}"/trampoline_cleanup.sh
    # "${KOKORO_GFILE_DIR}"/trampoline_cleanup.sh
    echo "cleanup";
}
trap cleanup EXIT

$(dirname $0)/populate-secrets.sh # Secret Manager secrets.
gcloud auth configure-docker us-docker.pkg.dev #set up docker to use gcloud for this url.
python3 "${KOKORO_GFILE_DIR}/trampoline_v1.py"
