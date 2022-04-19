#!/bin/bash

# Copyright 2022 API7.ai, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")
PROJECT_ROOT="$SCRIPT_ROOT/.."
COPYRIGHT_FILE=$(pwd)/$SCRIPT_ROOT/copyright_file

pushd "$PROJECT_ROOT"/internal/cloud
mockgen -source=./types.go -package=cloud -self_package=github.com/api7/cloud-cli/internal/cloud -copyright_file="$COPYRIGHT_FILE"  > ./api_mock.go
popd

pushd "$PROJECT_ROOT"/internal/commands
mockgen -source=./types.go -package=commands -copyright_file="$COPYRIGHT_FILE" > ./cmd_mock.go
popd
