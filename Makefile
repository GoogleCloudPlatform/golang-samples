# Makefile for running typical developer workflow actions.
# To run actions in a subdirectory of the repo:
#   make lint build dir=translate/snippets

# Default values for make variables
dir ?= $(shell pwd)

INTERFACE_ACTIONS="build test lint"
.ONESHELL:  # run make recipies in the same shell, to ease subdirectory usage.
.-PHONY: build test lint check-env list-actions

export GOLANG_SAMPLES_E2E_TEST ?= true
export GOLANG_SAMPLES_PROJECT_ID ?=${GOOGLE_SAMPLES_PROJECT}

# Required for cloud profiler tests to run
export GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn

build:
	go -C ${dir} build .

test: check-env
	# TODO: remove when we've re-built our testing containers to include this
	go install gotest.tools/gotestsum@latest
	cd ${dir}
	gotestsum --rerun-fails=3 --packages="./..." --junitfile sponge_log.xml -f standard-verbose -- --timeout 60m

lint:
	cd ${dir}
	go mod tidy
	go fmt .
	go vet .

check-env:
ifndef GOOGLE_SAMPLES_PROJECT
	$(error GOOGLE_SAMPLES_PROJECT environment variable is required to perform this action)
endif

list-actions:
	@ echo ${INTERFACE_ACTIONS}

