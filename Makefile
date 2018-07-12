SHELL := /bin/bash

.PHONY: proto
proto:
	script/proto.sh

.PHONY: build
build:
	go build -o ./bin/bbft main.go

.PHONY: build-osx
build-osx:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin64/bbft main.go

.PHONY: test
test:
	go test -cover -v $(shell glide novendor)