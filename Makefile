# Makefile for running typical developer workflow actions.
# To run actions in a subdirectory of the repo:
#   make lint build dir=translate/snippets

INTERFACE_ACTIONS="build test lint"

.-PHONY: build test lint check-env list-actions

#TODO: name this variable something more meaningful
DIR=${dir:-.}

GOLANG_SAMPLES_E2E_TEST=true
GOLANG_SAMPLES_PROJECT_ID="${GOOGLE_SAMPLE_PROJECT}"

build:
	go -C ${DIR} build .

test: check-env
	go -C ${DIR} test .

lint:
	cd ${DIR}
	go mod tidy
	go fmt .
	go vet .

check-env:
ifndef GOOGLE_SAMPLE_PROJECT
	$(error GOOGLE_SAMPLE_PROJECT environment variable is required to perform this action)
endif

list-actions:
	@ echo ${INTERFACE_ACTIONS}

