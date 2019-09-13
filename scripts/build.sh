#!/bin/bash

set -e

PROJECT=c2ae

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

if [[ -z "$OUTDIR" ]]; then
OUTDIR="${DIR}/../bin/"
fi

echo "$PROJECT build script (c) Teserakt AG 2018-2019. All rights reserved."
echo ""

GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_TAG=$(git describe --exact-match HEAD 2>/dev/null || true)
NOW=$(date "+%Y%m%d")

if [[ -z "${GOOS}" ]]; then
GOOS=`uname -s | tr '[:upper:]' '[:lower:]'`
fi
if [[ -z "${GOARCH}" ]]; then
GOARCH=amd64
fi

RACEDETECTOR=$(if [[ "$RACE" -ne "" ]]; then echo "-race"; else echo ""; fi)
CGO=$(if [[ "$RACE" -eq "" ]]; then echo "0"; else echo "1"; fi)
CMDS=($(find ${DIR}/../cmd/ -mindepth 1 -maxdepth 1  -type d -exec basename {} \;))
for cmd in ${CMDS[@]}; do
    printf "Building ${PROJECT}-${cmd}:\n\tversion ${NOW}-${GIT_COMMIT}\n\tOS ${GOOS}\n\tarch: ${GOARCH}\n"

    printf "=> ${PROJECT}-${cmd}...\n"
    CGO_ENABLED=$CGO GOOS=${GOOS} GOARCH=${GOARCH} go build $RACEDETECTOR -o ${OUTDIR}/${PROJECT}-${cmd} -ldflags "-X main.gitTag=${GIT_TAG} -X main.gitCommit=${GIT_COMMIT} -X main.buildDate=${NOW}" ${DIR}/../cmd/${cmd}
done
