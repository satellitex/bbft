#!/bin/bash -e
SRC=$(git rev-parse --show-toplevel)
protoc --proto_path=${GOPATH}/src/github.com/google/protobuf/src --proto_path=${SRC} --go_out=plugins=grpc:${SRC} --govalidators_out=${SRC} ${SRC}/*.proto
