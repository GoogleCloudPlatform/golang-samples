# Makefile for running typical developer workflow actions.
# To run actions in a subdirectory of the repo:
#   make lint build dir=translate/snippets

.-PHONY: build test lint check-env

#TODO: name this variable something more meaningful
DIR=${dir:-.}

GOLANG_SAMPLES_E2E_TEST=true
GOLANG_SAMPLES_PROJECT_ID="${GOOGLE_PROJECT_ID}"

build: check-env
	go -C ${DIR} build .

test: check-env
	go -C ${DIR} test .

lint:
	cd ${DIR}
	go mod tidy
	go fmt .
	go vet .

check-env:
ifndef GOOGLE_PROJECT_ID
	$(error GOOGLE_PROJECT_ID environment variable is required to perform this action)
endif
