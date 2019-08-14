#!/bin/bash

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Retrieve path to grpc-gateway modules folder and grep its latest version path only
GRPC_GATEWAY_SRC_PATH=$(find $GOPATH/pkg/mod/github.com/grpc-ecosystem/ -maxdepth 1 -type d -path *grpc-gateway* | sort -r | head -1)

protoc -I ${DIR}/../ -I $GRPC_GATEWAY_SRC_PATH/third_party/googleapis/ -I $GRPC_GATEWAY_SRC_PATH/ \
    --go_out=plugins=grpc:${DIR}/../internal/pb \
    --grpc-gateway_out=logtostderr=true:${DIR}/../internal/pb \
    --swagger_out=logtostderr=true:${DIR}/../doc/ \
    ${DIR}/../api.proto
