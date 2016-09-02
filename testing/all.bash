#!/bin/bash

set -e

if [ $SYSTEM_TESTS ]; then
  # IMPORTANT -x IS NOT SET HERE
  echo $GOOGLE_CREDENTIALS | base64 -d > /tmp/key.json
  export GOOGLE_APPLICATION_CREDENTIALS=/tmp/key.json
  export GOLANG_SAMPLES_PROJECT_ID=golang-samples-tests
else
  echo "Not running system tests.";
fi

set -x

export GOPATH=/gopath

# Do the easy stuff first. Fail fast!
diff -u <(echo -n) <(gofmt -d -s .)
go vet ./...

# Update imports from the cached image.
go get -u -v $(go list -f '{{join .Imports "\n"}}' ./... | sort | uniq | grep -v golang-samples)

# Run all of the tests
go test -v ./...
