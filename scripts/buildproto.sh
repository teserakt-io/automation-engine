#!/bin/bash

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

GOOGLEAPI=$(find $GOPATH/pkg/mod/github.com/grpc-ecosystem/ -path */grpc-gateway*/third_party/googleapis -type d | sort -r | head -1)

protoc -I ${DIR}/../ -I $GOOGLEAPI \
    --go_out=plugins=grpc:${DIR}/../internal/pb \
    --grpc-gateway_out=logtostderr=true:${DIR}/../internal/pb \
    --swagger_out=logtostderr=true:${DIR}/../doc/ \
    ${DIR}/../api.proto
