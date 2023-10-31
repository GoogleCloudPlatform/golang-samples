# Makefile for running typical developer workflow actions.

.-PHONY: build test lint

DIR=${DIR:-.}

build:
	go -C ${DIR} build .

test:
	go -C ${DIR} test .

lint:
	cd ${DIR}
	go mod tidy
	go fmt .
	go vet .
