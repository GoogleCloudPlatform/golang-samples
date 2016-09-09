#!/bin/bash

set -e

if [ $SYSTEM_TESTS ]; then
  # IMPORTANT -x IS NOT SET HERE
  echo $GOOGLE_CREDENTIALS | base64 -d > /tmp/key.json
  export GOOGLE_APPLICATION_CREDENTIALS=/tmp/key.json

  curl https://storage.googleapis.com/gimme-proj/linux_amd64/gimmeproj > /bin/gimmeproj && chmod +x /bin/gimmeproj;
  gimmeproj version;
  export GOLANG_SAMPLES_PROJECT_ID=$(gimmeproj -project golang-samples-tests lease 12m);
  if [ -z "$GOLANG_SAMPLES_PROJECT_ID" ]; then
    echo "Lease failed."
    exit 1
  fi
  echo "Running tests in project $GOLANG_SAMPLES_PROJECT_ID";
  trap "gimmeproj -project golang-samples-tests done $GOLANG_SAMPLES_PROJECT_ID" EXIT
else
  echo "Not running system tests.";
fi

set -x

export GOPATH=/gopath

# Do the easy stuff first. Fail fast!
diff -u <(echo -n) <(gofmt -d -s .)
go vet ./...

# Check use of Go 1.7 context package
grep -R '"context"$' * && { echo "Use golang.org/x/net/context"; false; } || true

# Update imports from the cached image.
go get -u -v $(go list -f '{{join .Imports "\n"}}' ./... | sort | uniq | grep -v golang-samples)

# Run all of the tests
go test -v ./...
