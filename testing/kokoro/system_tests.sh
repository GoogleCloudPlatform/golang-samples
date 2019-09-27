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

export GO111MODULE=on # Always use modules.
export GOPROXY=https://proxy.golang.org

export GOLANG_SAMPLES_KMS_KEYRING=ring1
export GOLANG_SAMPLES_KMS_CRYPTOKEY=key1
export GOLANG_SAMPLES_IOT_PUB=$KOKORO_GFILE_DIR/rsa_cert.pem
export GOLANG_SAMPLES_IOT_PRIV=$KOKORO_GFILE_DIR/rsa_private.pem
export GCLOUD_ORGANIZATION=1081635000895

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
trap "gimmeproj -project golang-samples-tests done $GOLANG_SAMPLES_PROJECT_ID" EXIT

# Set application credentials to the project-specific account. Some APIs do not
# allow the service account project and GOOGLE_CLOUD_PROJECT to be different.
export GOOGLE_APPLICATION_CREDENTIALS=$KOKORO_KEYSTORE_DIR/71386_kokoro-$GOLANG_SAMPLES_PROJECT_ID

set -x

export GOLANG_SAMPLES_SPANNER=projects/golang-samples-tests/instances/golang-samples-tests
export GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL=kokoro-$GOLANG_SAMPLES_PROJECT_ID@$GOLANG_SAMPLES_PROJECT_ID.iam.gserviceaccount.com
export GOLANG_SAMPLES_BIGTABLE_PROJECT=golang-samples-tests
export GOLANG_SAMPLES_BIGTABLE_INSTANCE=testing-instance

go version
date

# Re-organize files
export GOPATH=$PWD/gopath
target=$GOPATH/src/github.com/GoogleCloudPlatform
mkdir -p $target
mv github/golang-samples $target
cd $target/golang-samples

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

# CHANGED_DIRS is the list of top-level directories that changed. CHANGED_DIRS will be empty when run on master.
CHANGED_DIRS=$(git --no-pager diff --name-only HEAD..master | grep "/" | cut -d/ -f1 | sort | uniq || true)
# If test configuration is changed, run all tests.
if [[ $CHANGED_DIRS =~ "testing" || $CHANGED_DIRS =~ "internal" ]]; then
  RUN_ALL_TESTS="1"
fi

if [[ $RUN_ALL_TESTS = "1" ]]; then
  TARGET="./..."
  echo "Running all tests"
else
  TARGET=$(printf "./%s/... " $CHANGED_DIRS)
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

# Clear the cache so Kokoro doesn't try to copy it.
# Must happen before calling go-junit-report since it can cause a non-zero exit
# code, stopping execution.
go clean -modcache

cat $OUTFILE | /go/bin/go-junit-report -set-exit-code > sponge_log.xml
