#!/usr/bin/env bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
root=$(readlink -f "$scriptDir/../../..")

envFile="${scriptDir}/common.env"
if [[ -f  "${envFile}" ]]; then
    set -o allexport
    # shellcheck disable=SC1090
    source "${envFile}"
    set +o allexport
fi

if [[ "$#" -lt 1 ]]; then
    echo "instance name is required"
    exit 1
fi

cmdDir="${scriptDir}"
if [ ! -d "$cmdDir" ] ; then
    echo "Command directory $cmdDir does not exist"
fi
cd "$cmdDir" || exit 1


depDir="${root}/_local/bin"
if [ ! -d "${depDir}" ]; then 
    mkdir -p "${depDir}" || exit 1
fi

export CGO_ENABLED=0
go build -o "${depDir}/exmpl" || exit 2

echo
exec "$root/tools/exman/brun.sh" exec --name "exmpl-$1" "${depDir}/exmpl"  "$@"