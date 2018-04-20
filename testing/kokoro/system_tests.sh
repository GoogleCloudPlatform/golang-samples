#!/bin/bash

set -e

export GOOGLE_APPLICATION_CREDENTIALS=$KOKORO_KEYSTORE_DIR/71386_golang-samples-kokoro-service-account
export GOLANG_SAMPLES_KMS_KEYRING=ring1
export GOLANG_SAMPLES_KMS_CRYPTOKEY=key1

curl https://storage.googleapis.com/gimme-proj/linux_amd64/gimmeproj > /bin/gimmeproj && chmod +x /bin/gimmeproj;
gimmeproj version;
export GOLANG_SAMPLES_PROJECT_ID=$(gimmeproj -project golang-samples-tests lease 20m);
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

# Do the easy stuff first. Fail fast!
if [ $GOLANG_SAMPLES_GO_VET ]; then
  diff -u <(echo -n) <(gofmt -d -s .)
  go vet ./...
fi

# Check use of Go 1.7 context package
! grep -R '"context"$' * || { echo "Use golang.org/x/net/context"; false; }

# Download imports.
GO_IMPORTS=$(go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' ./... | sort | uniq | grep -v golang-samples)
time go get -u -v -d $GO_IMPORTS

# Pin go-sql-driver/mysql to v1.3 (which supports Go 1.6)
if go version | grep go1\.6\.; then
  pushd $GOPATH/src/github.com/go-sql-driver/mysql;
  git checkout v1.3;
  popd;
fi

go install -v $GO_IMPORTS

date

# Run all of the tests
go test -timeout 20m -v ./...
