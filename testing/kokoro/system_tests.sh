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

##
# system_tests.sh
# Runs CI checks for entire repository.
#
# Parameters
#
# [ARG 1]: Directory for the samples. Default: github/golang-samples.
# KOKORO_GFILE_DIR: Persistent filesystem location. (environment variable)
# KOKORO_KEYSTORE_DIR: Secret storage location. (environment variable)
# GOLANG_SAMPLES_GO_VET: If set, run code analysis checks. (environment variable)
##

set -ex

go version
date

cd "${1:-github/golang-samples}"
PROJECT_ROOT=$(pwd)

export GO111MODULE=on # Always use modules.
export GOPROXY=https://proxy.golang.org
TIMEOUT=60m
export GOLANG_SAMPLES_E2E_TEST=""

# Also see trampoline.sh - system_tests.sh is only run for PRs when there are
# significant changes.
# allow files to be owned by a different user than our current uid.
# Kokoro runs a double-nested container, and UIDs may not match.
git config --global --add safe.directory $(pwd)
# Allow $GIT_CHANGES to be set in the env, enabling local testing of the change detection below.
GIT_CHANGES=${GIT_CHANGES:-$(git --no-pager diff --name-only main..HEAD)}
if [[ -z $GIT_CHANGES && $KOKORO_JOB_NAME != *"system-tests"* ]]; then
  echo "No diffs detected. This is unexpected - check above for additional error messages."
  exit 2
fi
SIGNIFICANT_CHANGES=$(echo $GIT_CHANGES | grep -Ev '(\.md$|^\.github)' || true )
# CHANGED_DIRS is the list of significant top-level directories that changed,
# but weren't deleted by the current PR.
# CHANGED_DIRS will be empty when run on main.
CHANGED_DIRS=$(echo "$SIGNIFICANT_CHANGES" | tr ' ' '\n' | grep "/" | cut -d/ -f1 | sort -u | tr '\n' ' ' | xargs --no-run-if-empty ls -d 2>/dev/null || true)
GO_CHANGED_PKGS=$(echo "$SIGNIFICANT_CHANGES" | tr ' ' '\n' | grep "/" | tr '\n' ' ' | xargs --no-run-if-empty dirname | xargs --no-run-if-empty ls -d 2>/dev/null || true)

# List all modules in changed directories.
# If running on main will collect all modules in the repo, including the root module.
# shellcheck disable=SC2086
GO_CHANGED_MODULES="$(find ${GO_CHANGED_PKGS:-.} -name go.mod | xargs --no-run-if-empty dirname)"
# # If we didn't find any modules, use the root module.
# GO_CHANGED_MODULES=${GO_CHANGED_MODULES:-./go.mod}
# # Exclude the root module, if present, from the list of sub-modules.
# GO_CHANGED_SUBMODULES=${GO_CHANGED_MODULES#./go.mod}

# Override to determine if all go tests should be run.
# Does not include static analysis checks.
RUN_ALL_TESTS="0"
if [[ $KOKORO_JOB_NAME == *"system-tests"* ]]; then
  # If this is a standard nightly test, run all modules tests.
  RUN_ALL_TESTS="1"
  # If this is a nightly test for a specific submodule, run submodule tests only.
  # Submodule job name must have the format: "golang-samples/system-tests/[OPTIONAL_MODULE_NAME]/[GO_VERSION]"
  ARR=(${KOKORO_JOB_NAME//// })
  # Gets the "/" deliminated token after "system-tests".
  SUBMODULE_NAME=${ARR[4]}
  if [[ -n $SUBMODULE_NAME ]] && [[ -d "./$SUBMODULE_NAME" ]]; then
    RUN_ALL_TESTS="0"
    CHANGED_DIRS=$SUBMODULE_NAME
  fi
# If the change touches a repo-spanning file or directory of significance, run all tests.
elif echo "$SIGNIFICANT_CHANGES" | tr ' ' '\n' | grep "^go.mod$" || [[ $CHANGED_DIRS =~ "testing" || $CHANGED_DIRS =~ "internal" ]]; then
  RUN_ALL_TESTS="1"
fi

# Don't print environment variables in case there are secrets.
# If you need a secret, use a keystore_resource in common.cfg.
set +x

export GOLANG_SAMPLES_KMS_KEYRING=ring1
export GOLANG_SAMPLES_KMS_CRYPTOKEY=key1

export GOLANG_SAMPLES_IOT_PUB="$KOKORO_GFILE_DIR/rsa_cert.pem"
export GOLANG_SAMPLES_IOT_PRIV="$KOKORO_GFILE_DIR/rsa_private.pem"

export GCLOUD_ORGANIZATION=1081635000895
export SCC_PUBSUB_PROJECT="project-a-id"
export SCC_PUBSUB_TOPIC="projects/project-a-id/topics/notifications-sample-topic"
export SCC_PUBSUB_SUBSCRIPTION="notification-sample-subscription"
# gcp-sec-demo-org.joonix.net
export SCC_PROJECT_ORG_ID=688851828130
export SCC_PROJECT_ID=sharp-quest

export GOLANG_SAMPLES_SPANNER=projects/golang-samples-tests/instances/golang-samples-tests
export GOLANG_SAMPLES_SPANNER_INSTANCE_CONFIG="regional-us-west1"
export GOLANG_SAMPLES_BIGTABLE_PROJECT=golang-samples-tests
export GOLANG_SAMPLES_BIGTABLE_INSTANCE=testing-instance

export GOLANG_SAMPLES_FIRESTORE_PROJECT=golang-samples-fire-0
# This flag is added to avoid protobuf conflicts while running tests for profiler.
# TODO: Remove this after https://github.com/googleapis/google-cloud-go/issues/9207 is resolved.
export GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn

set -x

# Set application credentials before using gimmeproj so it has access.
# This is changed to a project-specific credential after a project is leased.
export GOOGLE_APPLICATION_CREDENTIALS=$KOKORO_KEYSTORE_DIR/71386_kokoro-golang-samples-tests
gimmeproj version;
GOLANG_SAMPLES_PROJECT_ID=$(gimmeproj -project golang-samples-tests lease $TIMEOUT);
export GOLANG_SAMPLES_PROJECT_ID
if [ -z "$GOLANG_SAMPLES_PROJECT_ID" ]; then
  echo "Lease failed."
  exit 1
fi
echo "Running tests in project $GOLANG_SAMPLES_PROJECT_ID";

# Always return the project and clean the cache so Kokoro doesn't try to copy
# it when exiting.
# shellcheck disable=SC2064
trap "go clean -modcache; gimmeproj -project golang-samples-tests done $GOLANG_SAMPLES_PROJECT_ID" EXIT

set +x

# Set application credentials to the project-specific account. Some APIs do not
# allow the service account project and GOOGLE_CLOUD_PROJECT to be different.
export GOOGLE_APPLICATION_CREDENTIALS=$KOKORO_KEYSTORE_DIR/71386_kokoro-$GOLANG_SAMPLES_PROJECT_ID
export GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL=kokoro-$GOLANG_SAMPLES_PROJECT_ID@$GOLANG_SAMPLES_PROJECT_ID.iam.gserviceaccount.com
export GOOGLE_API_GO_EXPERIMENTAL_ENABLE_NEW_AUTH_LIB="true"

set -x

pwd
date

export PATH="$PATH:/tmp/google-cloud-sdk/bin";
./testing/kokoro/configure_gcloud.bash;

# fetch secrets used by storagetransfer tests
set +x

export STS_AWS_SECRET=`gcloud secrets versions access latest --project cloud-devrel-kokoro-resources --secret=go-storagetransfer-aws`
export AWS_ACCESS_KEY_ID=`S="$STS_AWS_SECRET" python3 -c 'import json,sys,os;obj=json.loads(os.getenv("S"));print (obj["AccessKeyId"]);'`
export AWS_SECRET_ACCESS_KEY=`S="$STS_AWS_SECRET" python3 -c 'import json,sys,os;obj=json.loads(os.getenv("S"));print (obj["SecretAccessKey"]);'`
export STS_AZURE_SECRET=`gcloud secrets versions access latest --project cloud-devrel-kokoro-resources --secret=go-storagetransfer-azure`
export AZURE_STORAGE_ACCOUNT=`S="$STS_AZURE_SECRET" python3 -c 'import json,sys,os;obj=json.loads(os.getenv("S"));print (obj["StorageAccount"]);'`
export AZURE_CONNECTION_STRING=`S="$STS_AZURE_SECRET" python3 -c 'import json,sys,os;obj=json.loads(os.getenv("S"));print (obj["ConnectionString"]);'`
export AZURE_SAS_TOKEN=`S="$STS_AZURE_SECRET" python3 -c 'import json,sys,os;obj=json.loads(os.getenv("S"));print (obj["SAS"]);'`

set -x

if [[ $KOKORO_BUILD_ARTIFACTS_SUBDIR = *"system-tests"* ]] || [[ $CHANGED_DIRS =~ "run" ]] && [[ -n $GOLANG_SAMPLES_GO_VET ]]; then
  echo "This test run will run end-to-end tests.";

  # Download and load secrets
  ./testing/kokoro/pull-secrets.sh

  if [[ -f "./testing/kokoro/test-env.sh" ]]; then
    source ./testing/kokoro/test-env.sh
  else
    echo "Could not find environment file"
    echo "ls -lah ./testing"
    ls -lah ./testing
    echo "ls -lah ./testing/kokoro"
    ls -lah ./testing/kokoro
    exit 1
  fi

  export GOLANG_SAMPLES_E2E_TEST=1
  ./testing/kokoro/configure_cloudsql.bash;
fi


# only set with mtls_smoketest
# TODO(cbro): remove with mtls_smoketest.cfg
if [[ $GOOGLE_API_USE_MTLS = "always" ]]; then
  ./testing/kokoro/mtls_smoketest.bash
fi

date

# exit_code collects all of the exit codes of the tests, and is used to set the
# exit code at the end of the script.
exit_code=0
set +e # Don't exit on errors to make sure we run all tests.

# runTests runs the tests in the current directory. If an argument is specified,
# it is used as the argument to `go test`.
runTests() {
  if goVersionShouldSkip; then
    set +x
    echo "SKIPPING: module's minimum version is newer than the current Go version."
    set -x
    return 0
  fi

  set +x
  test_dir=$(realpath --relative-to $PROJECT_ROOT $(pwd))
  echo "Running 'go test' in '${test_dir}'..."
  set -x
  pushd $PROJECT_ROOT
  GOOGLE_SAMPLES_PROJECT=${GOLANG_SAMPLES_PROJECT_ID} make test dir=${test_dir}
  exit_code=$((exit_code + $?))
  popd
  set +x
}

# Returns 0 if the test should be skipped because the current Go
# version is too old for the current module.
goVersionShouldSkip() {
  modVersion="$(go list -m -f '{{.GoVersion}}')"
  if [ -z "$modVersion" ]; then
    # Not in a module or minimum Go version not specified, don't skip.
    return 1
  fi

  go list -f "{{context.ReleaseTags}}" ./... | grep -q -v "go$modVersion\b"
}

if [[ $RUN_ALL_TESTS = "1" ]]; then
  echo "Running all tests"
  # shellcheck disable=SC2044
  for i in $(find . -name go.mod); do
    pushd "$(dirname "$i")" > /dev/null;
      runTests
    popd > /dev/null;
  done
elif [[ -z "${CHANGED_DIRS// }" ]]; then
  echo "Only running root tests"
  runTests .
else
  runTests . # Always run root tests.
  echo "Running tests in modified directories: $GO_CHANGED_PKGS"
  for d in $GO_CHANGED_PKGS; do
    mods=$(find "$d" -name go.mod)
    # If there are no modules, just run the tests directly.
    if [[ -z "$mods" ]]; then
      pushd "$d" > /dev/null;
        runTests
      popd > /dev/null;
    # Otherwise, run the tests in all Go directories. This way, we don't have to
    # check to see if there are tests that aren't in a sub-module.
    else
      goDirectories="$(find "$d" -name "*.go" -printf "%h\n" | sort -u)"
      if [[ -n "$goDirectories" ]]; then
        for gd in $goDirectories; do
          pushd "$gd" > /dev/null;
            runTests .
          popd > /dev/null;
        done
      fi
    fi
  done
fi

exit $exit_code
