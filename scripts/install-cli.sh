#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

go install ${DIR}/../cmd/cli/ && . <(c2se-cli completion)
