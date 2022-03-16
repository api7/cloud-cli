#!/usr/bin/env bash

ori_version=$(apisix version 2>/dev/null)
if [[ $? -ne 0 ]]; then
  yum install -y {{ .APISIXRepoURL }}
  yum install -y apisix-{{ .Version }}
else
  ori_version=$(echo $ori_version |grep -o "{{ .Version }}")
  if [[ $org_version == "" ]]; then
    echo "another version of Apache APISIX have already installed."
    exit 1
  fi
fi

cp -rf {{ .TLSDir }} /usr/local/apisix/certs

apisix start -c {{ .ConfigFile }}
