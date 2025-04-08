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

cd $TESTING_ROOT

# Use SECRET_MANAGER_PROJECT if set, fallback to "golang-samples-tests".
PROJECT_ID="${SECRET_MANAGER_PROJECT:-golang-samples-tests}"

gcloud secrets versions add "golang-samples-test-env" \
       --project="${PROJECT_ID}" \
       --data-file="kokoro/test-env.sh"
