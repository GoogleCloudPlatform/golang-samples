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

set -e

set -x

go version
date

# Re-organize files
export GOPATH=$PWD/gopath
target=$GOPATH/src/github.com/GoogleCloudPlatform
mkdir -p $target
mv github/golang-samples $target
cd $target/golang-samples

export GO111MODULE=on # Always use modules.
export GOPROXY=https://proxy.golang.org

# Don't print environment variables in case there are secrets.
# If you need a secret, use a keystore_resource in common.cfg.
set +x

export GOLANG_SAMPLES_KMS_KEYRING=ring1
export GOLANG_SAMPLES_KMS_CRYPTOKEY=key1

export GOLANG_SAMPLES_IOT_PUB=$KOKORO_GFILE_DIR/rsa_cert.pem
export GOLANG_SAMPLES_IOT_PRIV=$KOKORO_GFILE_DIR/rsa_private.pem

export GCLOUD_ORGANIZATION=1081635000895

export GOLANG_SAMPLES_SPANNER=projects/golang-samples-tests/instances/golang-samples-tests
export GOLANG_SAMPLES_BIGTABLE_PROJECT=golang-samples-tests
export GOLANG_SAMPLES_BIGTABLE_INSTANCE=testing-instance

set -x

TIMEOUT=45m

# Set application credentials before using gimmeproj so it has access.
# This is changed to a project-specific credential after a project is leased.
export GOOGLE_APPLICATION_CREDENTIALS=$KOKORO_KEYSTORE_DIR/71386_kokoro-golang-samples-tests
gimmeproj version;
export GOLANG_SAMPLES_PROJECT_ID=$(gimmeproj -project golang-samples-tests lease $TIMEOUT);
if [ -z "$GOLANG_SAMPLES_PROJECT_ID" ]; then
  echo "Lease failed."
  exit 1
fi
echo "Running tests in project $GOLANG_SAMPLES_PROJECT_ID";

# Always return the project and clean the cache so Kokoro doesn't try to copy
# it when exiting.
trap "go clean -modcache; gimmeproj -project golang-samples-tests done $GOLANG_SAMPLES_PROJECT_ID" EXIT

set +x

# Set application credentials to the project-specific account. Some APIs do not
# allow the service account project and GOOGLE_CLOUD_PROJECT to be different.
export GOOGLE_APPLICATION_CREDENTIALS=$KOKORO_KEYSTORE_DIR/71386_kokoro-$GOLANG_SAMPLES_PROJECT_ID
export GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL=kokoro-$GOLANG_SAMPLES_PROJECT_ID@$GOLANG_SAMPLES_PROJECT_ID.iam.gserviceaccount.com

set -x

pwd
date

if [[ $KOKORO_BUILD_ARTIFACTS_SUBDIR = *"system-tests"* && -n $GOLANG_SAMPLES_GO_VET ]]; then
  echo "This test run will run end-to-end tests.";
  export GOLANG_SAMPLES_E2E_TEST=1
  export PATH="$PATH:/tmp/google-cloud-sdk/bin";
  ./testing/kokoro/configure_gcloud.bash;
fi

RUN_ALL_TESTS="0"
# If this is a nightly test (not a PR), run all tests.
if [ -z ${KOKORO_GITHUB_PULL_REQUEST_NUMBER:-} ]; then
  RUN_ALL_TESTS="1"
fi

# CHANGED_DIRS is the list of significant top-level directories that changed.
# CHANGED_DIRS will be empty when run on master.
# Also see trampoline.sh - system_tests.sh is only run when there are
# significant changes.
CHANGED_DIRS=$(git --no-pager diff --name-only HEAD..master | egrep -v '(\.md$|^\.github)' | grep "/" | cut -d/ -f1 | sort -u)
# If test configuration is changed, run all tests.
if [[ $CHANGED_DIRS =~ "testing" || $CHANGED_DIRS =~ "internal" ]]; then
  RUN_ALL_TESTS="1"
fi

# Filter out directories that don't exist (the current PR deleted them).
TARGET_DIRS=""
for d in "$CHANGED_DIRS"; do
  if [ -d "$d" ]; then
    TARGET_DIRS="$TARGET_DIRS $d"
  fi
done

if [[ $RUN_ALL_TESTS = "1" ]]; then
  TARGET="./..."
  echo "Running all tests"
elif [[ -z "${TARGET_DIRS// }" ]]; then
  TARGET=""
  echo "Only running root tests"
else
  TARGET=$(printf "./%s/... " $TARGET_DIRS)
  echo "Running tests in modified directories: $TARGET"
fi

# Do the easy stuff before running tests. Fail fast!
if [ $GOLANG_SAMPLES_GO_VET ]; then
  diff -u <(echo -n) <(gofmt -d -s .)
  go vet $TARGET
fi

date

OUTFILE=gotest.out
2>&1 go test -timeout $TIMEOUT -v . $TARGET | tee $OUTFILE

cat $OUTFILE | /go/bin/go-junit-report -set-exit-code > sponge_log.xml
