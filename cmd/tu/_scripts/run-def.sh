#!/bin/bash

scriptDir="$(cd "$(dirname "$0")" || exit ; pwd -P)"
cmdDir=$(readlink -f "$scriptDir/..")

export PG_URL="postgres://postgres:onetwothree@localhost:5432/fake-data?sslmode=disable"

cd "$cmdDir" || exit 1

go run main.go "$@"