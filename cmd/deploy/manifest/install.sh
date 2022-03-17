#!/usr/bin/env bash

version="{{ .Version }}"

installed_version=$(apisix version 2>/dev/null)
if [[ $? -ne 0 ]]; then
  yum install -y {{ .APISIXRepoURL }}
  yum install -y apisix-$version
fi

cp -rf {{ .TLSDir }} /usr/local/apisix/conf/certs

apisix start -c {{ .ConfigFile }}
