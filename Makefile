# Makefile for running typical developer workflow actions.
# To run actions in a subdirectory of the repo:
#   make lint build dir=translate/snippets

# Default values for make variables
dir ?= $(shell pwd)

INTERFACE_ACTIONS="build test lint"
.ONESHELL:  # run make recipies in the same shell, to ease subdirectory usage.
.-PHONY: build test lint check-env list-actions

export GOLANG_SAMPLES_E2E_TEST ?= true
export GOLANG_SAMPLES_PROJECT_ID ?=${GOOGLE_SAMPLE_PROJECT}

build:
	go -C ${dir} build .

test: check-env
	go -C ${dir} test .

lint:
	cd ${dir}
	go mod tidy
	go fmt .
	go vet .

check-env:
ifndef GOOGLE_SAMPLE_PROJECT
	$(error GOOGLE_SAMPLE_PROJECT environment variable is required to perform this action)
endif

list-actions:
	@ echo ${INTERFACE_ACTIONS}

