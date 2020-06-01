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
SIGNIFICANT_CHANGES=$(git --no-pager diff --name-only master..HEAD | egrep -v '(\.md$|^\.github)' || true)
# CHANGED_DIRS is the list of significant top-level directories that changed,
# but weren't deleted by the current PR.
# CHANGED_DIRS will be empty when run on master.
CHANGED_DIRS=$(echo $SIGNIFICANT_CHANGES | tr ' ' '\n' | grep "/" | cut -d/ -f1 | sort -u | tr '\n' ' ' | xargs ls -d 2>/dev/null || true)

# List all modules in changed directories.
# If running on master will collect all modules in the repo, including the root module.
GO_CHANGED_MODULES=$(find ${CHANGED_DIRS:-.} -name go.mod)
# If we didn't find any modules, use the root module.
GO_CHANGED_MODULES=${GO_CHANGED_MODULES:-./go.mod}
# Exclude the root module, if present, from the list of sub-modules.
GO_CHANGED_SUBMODULES=${GO_CHANGED_MODULES#./go.mod}

# Override to determine if all go tests should be run.
# Does not include static analysis checks.
RUN_ALL_TESTS="0"
# If this is a nightly test (not a PR), run all tests.
if [ -z ${KOKORO_GITHUB_PULL_REQUEST_NUMBER:-} ]; then
  RUN_ALL_TESTS="1"
# If the change touches a repo-spanning file or directory of significance, run all tests.
elif echo $SIGNIFICANT_CHANGES | tr ' ' '\n' | grep "^go.mod$" || [[ $CHANGED_DIRS =~ "testing" || $CHANGED_DIRS =~ "internal" ]]; then
  RUN_ALL_TESTS="1"
fi

## Static Analysis
# Do the easy stuff before running tests or reserving a project. Fail fast!
set +x

if [ $GOLANG_SAMPLES_GO_VET ]; then
  echo "Running 'goimports compliance check'"
  set -x
  diff -u <(echo -n) <(goimports -d .)
  set +x
  for i in $GO_CHANGED_MODULES; do
    mod="$(dirname $i)"
    pushd $mod > /dev/null;
      # Fail if a dependency was added without the necessary go.mod/go.sum change
      # being part of the commit.
      echo "Running 'go.mod/go.sum sync check' in '$mod'..."
      set -x
      go mod tidy;
      git diff go.mod | tee /dev/stderr | (! read)
      [ -f go.sum ] && git diff go.sum | tee /dev/stderr | (! read)
      set +x
    popd > /dev/null;
  done

  # Always run 'go vet' from the root, which does not look at sub-modules.
  set +x
  echo "Running 'go vet' in golang-samples root..."
  set -x
  go vet ./...
  set +x

  # Run go vet inside each sub-module.
  # Recursive submodules are not supported.
  set +x
  for i in $GO_CHANGED_SUBMODULES; do
    mod="$(dirname $i)"
    pushd $mod > /dev/null;
      echo "Running 'go vet' in '$mod'..."
      set -x
      go vet ./...
      set +x
    popd > /dev/null;
  done
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

go install ./testing/sampletests

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
fi

export PATH="$PATH:/tmp/google-cloud-sdk/bin";
if [[ $KOKORO_BUILD_ARTIFACTS_SUBDIR = *"system-tests"* ]]; then
  ./testing/kokoro/configure_gcloud.bash;
fi

date

# exit_code collects all of the exit codes of the tests, and is used to set the
# exit code at the end of the script.
exit_code=0
set +e # Don't exit on errors to make sure we run all tests.

# runTests runs the tests in the current directory. If an argument is specified,
# it is used as the argument to `go test`.
runTests() {
  set +x
  echo "Running 'go test' in '$(pwd)'..."
  set -x
  2>&1 go test -timeout $TIMEOUT -v ${1:-./...} | tee sponge_log.log
  cat sponge_log.log | /go/bin/go-junit-report -set-exit-code > raw_log.xml
  exit_code=$(($exit_code + $?))
  # Add region tags tested to test case properties.
  cat raw_log.xml | sampletests > sponge_log.xml
  rm raw_log.xml # No need to keep this around.
  set +x
}

if [[ $RUN_ALL_TESTS = "1" ]]; then
  echo "Running all tests"
  for i in $(find . -name go.mod); do
    pushd "$(dirname $i)" > /dev/null;
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
    mods="$(find $d -name go.mod)"
    # If there are no modules, just run the tests directly.
    if [[ -z "$mods" ]]; then
      pushd "$d" > /dev/null;
        runTests
      popd > /dev/null;
    # Otherwise, run the tests in all Go directories. This way, we don't have to
    # check to see if there are tests that aren't in a sub-module.
    else
      goDirectories="$(find $d -name "*.go" -printf "%h\n" | sort -u)"
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

# If we're running system tests, send the test log to the Build Cop Bot.
# See https://github.com/googleapis/repo-automation-bots/tree/master/packages/buildcop.
if [[ $KOKORO_BUILD_ARTIFACTS_SUBDIR = *"system-tests"* ]]; then
  chmod +x $KOKORO_GFILE_DIR/linux_amd64/buildcop
  $KOKORO_GFILE_DIR/linux_amd64/buildcop
fi

exit $exit_code
