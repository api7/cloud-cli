#!/usr/bin/env bash

version="{{ .Version }}"

installed_version=$(apisix version 2>/dev/null)
if [[ $? -ne 0 ]]; then
  yum install -y {{ .APISIXRepoURL }}
  yum install -y apisix-$version
fi

# copy certs to apisix directory to avoid permission issue
cp -rf {{ .TLSDir }} /usr/local/apisix/conf/ssl

apisix start -c {{ .ConfigFile }}
