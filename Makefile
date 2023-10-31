# Makefile for running typical developer workflow actions.
# To run actions in a subdirectory of the repo:
#   make lint build dir=translate/snippets

.-PHONY: build test lint

#TODO: name this variable something more meaningful
DIR=${dir:-.}

GOLANG_SAMPLES_E2E_TEST=true
GOLANG_SAMPLES_PROJECT_ID="${GOOGLE_PROJECT_ID}"

build:
	go -C ${DIR} build .

test:
	go -C ${DIR} test .

lint:
	cd ${DIR}
	go mod tidy
	go fmt .
	go vet .
