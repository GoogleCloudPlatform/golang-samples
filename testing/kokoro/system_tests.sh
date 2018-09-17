#!/bin/bash

set -e

export GOOGLE_APPLICATION_CREDENTIALS=$KOKORO_KEYSTORE_DIR/71386_golang-samples-kokoro-service-account
export GOLANG_SAMPLES_KMS_KEYRING=ring1
export GOLANG_SAMPLES_KMS_CRYPTOKEY=key1

TIMEOUT=25m

curl https://storage.googleapis.com/gimme-proj/linux_amd64/gimmeproj > /bin/gimmeproj && chmod +x /bin/gimmeproj;
gimmeproj version;
export GOLANG_SAMPLES_PROJECT_ID=$(gimmeproj -project golang-samples-tests lease $TIMEOUT);
if [ -z "$GOLANG_SAMPLES_PROJECT_ID" ]; then
  echo "Lease failed."
  exit 1
fi
echo "Running tests in project $GOLANG_SAMPLES_PROJECT_ID";
trap "gimmeproj -project golang-samples-tests done $GOLANG_SAMPLES_PROJECT_ID" EXIT

set -x

export GOLANG_SAMPLES_SPANNER=projects/golang-samples-tests/instances/golang-samples-tests

go version
date

if [[ -d /cache ]]; then
  time mv /cache/* .
  echo 'Uncached'
fi

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

# Download imports.
GO_IMPORTS=$(go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' $TARGET | \
  sort | uniq | \
  grep -v golang-samples)
time go get -u -v -d $GO_IMPORTS

# Always download internal dependencies.
go get ./internal/...

go get github.com/jstemmer/go-junit-report
go install golang.org/x/tools/imports;
go install -v $GO_IMPORTS

# Do the easy stuff before running tests. Fail fast!
if [ $GOLANG_SAMPLES_GO_VET ]; then
  diff -u <(echo -n) <(gofmt -d -s .)
  go vet $TARGET
fi

date

OUTFILE=gotest.out
2>&1 go test -timeout $TIMEOUT -v . $TARGET | tee $OUTFILE
cat $OUTFILE | $GOPATH/bin/go-junit-report -set-exit-code > sponge_log.xml
