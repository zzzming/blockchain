#!/bin/bash

#
# Run the CI flow and build the binary
# Prerequisite -
# 1. Go runtime
#

# absolute directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

BASE_PKG_DIR="github.com/zzzming/mbt/src/"
ALL_PKGS=""

cd $DIR/../src
gofmt -w -s -d .
# test lint, vet, and build as basic build steps in CI
echo run golint
revive -config ../revive.toml  -formatter friendly ./...
echo run go vet
go vet ./...

echo run go build
mkdir -p ${DIR}/../bin
rm -f ${DIR}/../bin/mbt
GIT_COMMIT=$(git rev-list -1 HEAD)
go build -o ${DIR}/../bin/mbt -ldflags "-X main.gitCommit=$GIT_COMMIT"

cd $DIR/../src
go test ./...
go test -v -coverpkg=./... -coverprofile=coverage.out -json ./... > report.json
# coverage visualization
# go tool cover -html=coverage.out
