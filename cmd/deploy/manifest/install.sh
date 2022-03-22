#!/usr/bin/env bash

version="{{ .Version }}"
instance_id="{{ .InstanceID }}"

installed_version=$(apisix version 2>/dev/null)
if [[ $? -ne 0 ]]; then
  yum install -y {{ .APISIXRepoURL }}
  yum install -y apisix-$version
fi

# copy certs to apisix directory to avoid permission issue
cp -rf {{ .TLSDir }} /usr/local/apisix/conf/ssl

# get the APISIX instance id when instance id is not set
if [[ -z ${instance_id} ]]; then
  instance_id="$(cat /usr/local/apisix/conf/apisix.uid)"
fi

apisix start -c {{ .ConfigFile }}

if [[ $? -eq 0 ]]; then
  echo "Your APISIX Instance was deployed successfully!"
  echo "Instance ID: ${instance_id}"
fi
