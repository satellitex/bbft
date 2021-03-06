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

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux/bbft main.go

.PHONY: test
test:
	go test -cover -v $(shell glide novendor)

.PHONY: test-ci
test-ci:
	go test -v $(shell glide novendor)

.PHONY: build-sender
build-sender:
	go build -o ./bin/sender ./demo/sender.go

.PHONY: push-docker-demo
push-docker-demo: build-linux
	docker build . -t satellitex/bbft-server
	docker push satellitex/bbft-server
