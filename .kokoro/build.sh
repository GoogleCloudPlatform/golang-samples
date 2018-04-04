#!/bin/bash

set -e
set -x

rm -rf $GOPATH/src/github.com/go-sql-driver/mysql

cd ${KOKORO_ARTIFACTS_DIR}/github/golang-samples

GO_IMPORTS=$(go list -f '{{join .Imports "\n"}}{{"\n"}}{{join .TestImports "\n"}}' ./... | sort | uniq | grep -v golang-samples)
go get -u -v -d $GO_IMPORTS

# pin go-sql-driver/mysql to v1.3 (which supports Go 1.6)
if go version | grep go1\.6\.; then
  pushd $GOPATH/src/github.com/go-sql-driver/mysql;
  git checkout v1.3;
  popd;
fi

go install -v $GO_IMPORTS
go get -u -v github.com/rakyll/gotest

# Basic build to test Kokoro configuration.
go build