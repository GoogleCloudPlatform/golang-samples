#!/bin/bash

# Copyright 2021 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eEuo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TESTING_ROOT=$( dirname "$DIR" )

# Work from the project root.
cd $TESTING_ROOT

# Prevent it from overriding files.
if [[ -f "kokoro/test-env.sh" ]]; then
    echo "testing/kokoro/test-env.sh already exists. Aborting."
    exit 1
fi

# Use SECRET_MANAGER_PROJECT if set, fallback to "golang-samples-tests".
PROJECT_ID="${SECRET_MANAGER_PROJECT:-golang-samples-tests}"

gcloud secrets versions access latest \
    --secret="golang-samples-test-env" \
    --project="${PROJECT_ID}" \
    > kokoro/test-env.sh
