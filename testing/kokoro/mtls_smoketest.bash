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
    go get google.golang.org/api@main
    go get cloud.google.com/go@main
  popd
done

# List of tests to include during mtls_smoketest.
scope=(
  automl/
  bigtable/
  cloudsql/
  container/
  container_registry/
  dataproc/
  datastore/
  dlp/
  kms/
  logging/
  pubsub/
  spanner/
  speech/
  trace/
  translate/
)

for d in */; do
  in_scope=0
  for pkg in "${scope[@]}"; do
    if [ $pkg = $d ]; then
      in_scope=1
      break
    fi
  done
  if [ $in_scope = 0 ]; then
    find "./$d" -name '*_test.go' -exec rm -r {} \;
  fi
done
