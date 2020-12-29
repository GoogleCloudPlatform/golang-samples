#! /bin/sh

# Run the tests with default values for all
# environment variable
go test -timeout $TIMEOUT -v "${1:-./...}" | tee sponge_log.log

# Run the tests again with custom environment
# variables that force a Unix connection
DB_HOST=; go test -timeout $TIMEOUT -v "${1:-./...}" | tee sponge_log.log