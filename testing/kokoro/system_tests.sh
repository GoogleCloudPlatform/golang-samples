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

export GO111MODULE=on # Always use modules.
export GOPROXY=https://proxy.golang.org
TIMEOUT=60m

# Also see trampoline.sh - system_tests.sh is only run for PRs when there are
# significant changes.
SIGNIFICANT_CHANGES=$(git --no-pager diff --name-only master..HEAD | grep -Ev '(\.md$|^\.github)' || true)
# CHANGED_DIRS is the list of significant top-level directories that changed,
# but weren't deleted by the current PR.
# CHANGED_DIRS will be empty when run on master.
CHANGED_DIRS=$(echo "$SIGNIFICANT_CHANGES" | tr ' ' '\n' | grep "/" | cut -d/ -f1 | sort -u | tr '\n' ' ' | xargs ls -d 2>/dev/null || true)

# List all modules in changed directories.
# If running on master will collect all modules in the repo, including the root module.
# shellcheck disable=SC2086
GO_CHANGED_MODULES="$(find ${CHANGED_DIRS:-.} -name go.mod)"
# If we didn't find any modules, use the root module.
GO_CHANGED_MODULES=${GO_CHANGED_MODULES:-./go.mod}
# Exclude the root module, if present, from the list of sub-modules.
GO_CHANGED_SUBMODULES=${GO_CHANGED_MODULES#./go.mod}

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

export GOLANG_SAMPLES_SPANNER=projects/golang-samples-tests/instances/golang-samples-tests
export GOLANG_SAMPLES_BIGTABLE_PROJECT=golang-samples-tests
export GOLANG_SAMPLES_BIGTABLE_INSTANCE=testing-instance

export GOLANG_SAMPLES_FIRESTORE_PROJECT=golang-samples-fire-0

set -x

pushd testing/sampletests
go install .
popd

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

set -x

pwd
date

export PATH="$PATH:/tmp/google-cloud-sdk/bin";
if [[ $KOKORO_BUILD_ARTIFACTS_SUBDIR = *"system-tests"* ]] || [[ $CHANGED_DIRS =~ "run" ]]; then
  ./testing/kokoro/configure_gcloud.bash;
fi



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
  echo "Running 'go test' in '$(pwd)'..."
  set -x
  2>&1 go test -timeout $TIMEOUT -v "${1:-./...}" | tee sponge_log.log
  /go/bin/go-junit-report -set-exit-code < sponge_log.log > raw_log.xml
  exit_code=$((exit_code + $?))
  # Add region tags tested to test case properties.
  sampletests < raw_log.xml > sponge_log.xml
  rm raw_log.xml # No need to keep this around.
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
  echo "Running tests in modified directories: $CHANGED_DIRS"
  for d in $CHANGED_DIRS; do
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

# If we're running system tests, send the test log to Flaky Bot.
# See https://github.com/googleapis/repo-automation-bots/tree/master/packages/flakybot.
if [[ $KOKORO_BUILD_ARTIFACTS_SUBDIR = *"system-tests"* ]]; then
  chmod +x "$KOKORO_GFILE_DIR"/linux_amd64/flakybot
  "$KOKORO_GFILE_DIR"/linux_amd64/flakybot
fi

exit $exit_code
