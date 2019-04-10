#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

protoc --proto_path ${DIR}/../ api.proto --go_out=plugins=grpc:${DIR}/../internal/pb
