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

# Only for the mtls_smoketest nightly job.
# TODO(cbro): remove with mtls_smoketest.cfg at some point.

./testing/kokoro/configure_gcloud.bash

# Keep these deps at HEAD so we don't need to cut a release to check for a fix.
for f in $(find . -name go.mod); do
  pushd $(dirname $f)
    go get google.golang.org/api@master
    go get cloud.google.com/go@master
  popd
done