#!/usr/bin/env bash

#!/bin/bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
root=$(readlink -f "$scriptDir/..")

envFile="${scriptDir}/common.env"
if [[ -f  "${envFile}" ]]; then
    set -o allexport
    # shellcheck disable=SC1090
    source "${envFile}"
    set +o allexport
fi

echo "Building...."
#!/bin/sh

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
root=$(readlink -f "$scriptDir/../..")


cmdDir="${scriptDir}"
if [ ! -d "$cmdDir" ] ; then
    echo "Command directory $cmdDir does not exist"
fi
cd "$cmdDir" || exit 1


git --version  >/dev/null 2>&1
GIT_IS_AVAILABLE=$?
if [ $GIT_IS_AVAILABLE -eq 0 ] &&  [ -z "$GIT_TAG" ]; then 
    GIT_TAG=$(git describe --tag || echo 'latest')
    GIT_HASH=$(git rev-parse --verify HEAD)
    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    BUILD_TIME=$(date -Isec)
    BUILD_HOST=$(hostname)
    BUILD_USER=$(whoami)
fi

if [ -z "$GIT_TAG" ];    then echo "GIT_TAG not set "; exit 1; fi
if [ -z "$GIT_HASH" ];   then echo "GIT_HASH not set "; exit 1; fi
if [ -z "$GIT_BRANCH" ]; then echo "GIT_BRANCH not set "; exit 1; fi
if [ -z "$BUILD_TIME" ]; then echo "BUILD_TIME not set "; exit 1; fi
if [ -z "$BUILD_HOST" ]; then echo "BUILD_HOST not set "; exit 1; fi
if [ -z "$BUILD_USER" ]; then echo "BUILD_USER not set "; exit 1; fi


# GOMOD=${GOMOD:-"go.mod"}
# echo "Using Go mod file: ${GOMOD}"

depDir="${root}/_local/bin"
if [ ! -d "${depDir}" ]; then 
    mkdir -p "${depDir}" || exit 1
fi

export CGO_ENABLED=0
go build \
    -installsuffix 'static' \
    -ldflags "-w -s \
        -X github.com/varunamachi/libx/tools/exman/main.GitTag=${GIT_TAG}
        -X github.com/varunamachi/libx/tools/exman/main.GitHash=${GIT_HASH}
        -X github.com/varunamachi/libx/tools/exman/main.GitBranch=${GIT_BRANCH}
        -X github.com/varunamachi/libx/tools/exman/main.BuildTime=${BUILD_TIME}
        -X github.com/varunamachi/libx/tools/exman/main.BuildHost=${BUILD_HOST}
        -X github.com/varunamachi/libx/tools/exman/main.BuildUser=${BUILD_USER}    
    "\
    -o "${depDir}/exman" || exit 2


echo "Running...."
echo
"$root/_local/bin/exman" "$@"