#!/usr/bin/env bash
#
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

SUPPORTED_OS_LIST="darwin linux"
SUPPORTED_ARCH_LIST="amd64 arm64"

error() {
  echo -e "\033[31m> $1\033[0m"
  exit 1
}

validate_os() {
  for os in $SUPPORTED_OS_LIST; do
    if [ "$os" = "$1" ]; then
      return 0
    fi
  done

  error "Unsupported OS: $1"
}

validate_arch() {
  for arch in $SUPPORTED_ARCH_LIST; do
    if [ "$arch" = "$1" ]; then
      return 0
    fi
  done

  error "Unsupported Arch: $1"
}

install_cloud_cli() {
  curl https://github.com/api7/cloud-cli/releases/download/${1}/cloud-cli-${2}-${3}-${1}.gz -sLo /tmp/cloud-cli.gz
  gzip -d -f /tmp/cloud-cli.gz
  chmod a+x /tmp/cloud-cli
  echo "Cloud CLI was installed successfully in /tmp/cloud-cli"
}

OS=`uname -s | tr A-Z a-z`
ARCH=`uname -m | tr A-Z a-z`
CLOUD_CLI_VER=`cat ../VERSION`

if [ "$ARCH" = "x86_64" ]; then
  ARCH=amd64
fi

validate_arch $ARCH
validate_os $OS
install_cloud_cli $CLOUD_CLI_VER $OS $ARCH
