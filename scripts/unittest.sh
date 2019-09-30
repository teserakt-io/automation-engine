#!/bin/sh

set -e

if [ -z $(which golint) ]; then
    go get golang.org/x/lint/golint
fi

if [ -z $(which staticcheck) ]; then
    go get honnef.co/go/tools/cmd/staticcheck
fi

echo "Running golint..."
golint -set_exit_status ./...

echo "Running staticcheck..."
staticcheck ./...

echo "Running go test..."
C2AETEST_POSTGRES="${C2AETEST_POSTGRES:-}" \
    go test -failfast -race ./...
