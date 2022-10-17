#!/usr/bin/env bash

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

set -e

version="{{ .Version }}"
instance_id="{{ .InstanceID }}"
apisix_home="/usr/local/apisix"

installed_version=$(apisix version 2>/dev/null) || true
if [[ -z ${installed_version} ]]; then
  yum install -y {{ .APISIXRepoURL }}
  yum install -y apisix-$version
fi

# copy certs to apisix directory to avoid permission issue
cp -prf {{ .TLSDir }} ${apisix_home}/conf/ssl

if [[ -n ${instance_id} ]]; then
  echo "${instance_id}" > ${apisix_home}/conf/apisix.uid
fi

apisix start -c {{ .ConfigFile }}
status=$?

# wait for APISIX started and generated instance id
sleep 1

# get the APISIX instance id when instance id is not set
if [[ -z ${instance_id} ]]; then
  instance_id="$(cat ${apisix_home}/conf/apisix.uid)"
fi

if [[ $status -eq 0 ]]; then
  echo "Your APISIX Instance was deployed successfully!"
  echo "Instance ID: ${instance_id}"
fi
