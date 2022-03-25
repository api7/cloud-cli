#!/usr/bin/env bash

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
cp -rf {{ .TLSDir }} ${apisix_home}/conf/ssl

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
