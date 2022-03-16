#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")
PROJECT_ROOT="$SCRIPT_ROOT/.."

pushd "$PROJECT_ROOT"/internal/cloud
mockgen -source=./types.go -package=cloud > ./api_mock.go
popd

pushd "$PROJECT_ROOT"/internal/commands
mockgen -source=./types.go -package=commands > ./cmd_mock.go
popd