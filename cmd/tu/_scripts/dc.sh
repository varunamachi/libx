#!/bin/bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"

docker-compose -f "$scriptDir/fake-data.dc.yml" -p "fake-data" "$@"