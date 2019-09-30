#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "c2ae Docker build script (c) Teserakt AG 2018-2019. All Right Reserved"
echo ""

E4_VERSION="${CI_COMMIT_REF_NAME//\//_}"
E4_GIT_COMMIT="${CI_COMMIT_SHORT_SHA}"

if [[ -z "$E4_VERSION" ]]; then
    E4_VERSION="devel"
fi

if [[ -z "$E4_GIT_COMMIT" ]]; then
    E4_GIT_COMMIT=$(git rev-list -1 HEAD)
fi

echo "Building version $E4_VERSION, commit $E4_GIT_COMMIT\n"

if [ -z $(ldd ${DIR}/../bin/c2ae-api | grep "not a dynamic executable") ]; then
    echo "c2ae-api is not a static binary, please rebuild it with CGO_ENABLED=0"
    exit 1
fi

printf "=> c2ae-api"
docker build \
    --target c2ae-api \
    --build-arg binary_path=./bin/c2ae-api \
    --tag "c2ae-api:$E4_VERSION" \
    --tag "c2ae-api:$E4_GIT_COMMIT" \
    -f "${DIR}/../docker/c2ae/Dockerfile" \
    "${DIR}/../"


if [ -z $(ldd ${DIR}/../bin/c2ae-cli | grep "not a dynamic executable") ]; then
    echo "c2ae-cli is not a static binary, please rebuild it with CGO_ENABLED=0"
    exit 1
fi

printf "=> c2ae-cli"
docker build \
    --target c2ae-cli \
    --build-arg binary_path=./bin/c2ae-cli \
    --tag "c2ae-cli:$E4_VERSION" \
    --tag "c2ae-cli:$E4_GIT_COMMIT" \
    -f "${DIR}/../docker/c2ae/Dockerfile" \
    "${DIR}/../"
