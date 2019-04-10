#!/bin/bash

PROJECT=c2se

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "$PROJECT build script (c) Teserakt AG 2018-2019. All rights reserved."
echo ""

GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_TAG=$(git describe --exact-match HEAD 2>/dev/null)
NOW=$(date "+%Y%m%d")

GOOS=`uname -s | tr '[:upper:]' '[:lower:]'` 
GOARCH=amd64

CMDS=($(find ${DIR}/../cmd/ -mindepth 1 -maxdepth 1  -type d -exec basename {} \;))
for cmd in ${CMDS[@]}; do

    printf "building ${PROJECT}-${cmd}:\n\tversion ${NOW}-${GIT_COMMIT}\n\tOS ${GOOS}\n\tarch: ${GOARCH}\n"

    printf "=> ${PROJECT}-${cmd}...\n"
    GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${DIR}/../bin/${PROJECT}-${cmd} -ldflags "-X main.gitTag=${GIT_TAG} -X main.gitCommit=${GIT_COMMIT} -X main.buildDate=${NOW}" ${DIR}/../cmd/${cmd}
done
