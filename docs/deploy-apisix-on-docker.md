<!--
# Copyright 2023 API7.ai, Inc
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
-->

Deploy APISIX on Docker
=======================

In this section, you'll learn how to deploy APISIX on Docker through Cloud CLI.

> Note, before you go ahead, make sure you read the section
> [How to Configure Cloud CLI](./configuring-cloud-cli.md)

Cloud CLI will create a Docker container for APISIX, and mount some
information to the container:

* The Cloud Lua Module

The Cloud Lua Module contains codes to communicate with API7 Cloud (such as
heartbeat, status reporting, etc.), it'll be downloaded every time you run the command.

> Currently, the Cloud Lua Module will be downloaded from [api7/cloud-scripts](https://github.com/api7/cloud-scripts).

* The TLS Bundle

TLS Bundle (Certificate, Private Key, CA Bundle) will be downloaded from API7
Cloud, only instances with a valid client certificate can connect to API7 Cloud.

> See the
> [DP Certificate API](https://docs.az-staging.api7.cloud/swagger/#/controlplanes_operation/getCertificates)
> to learn the details.

* The APISIX Configuration Template

Cloud CLI will also download the APISIX Configuration Template, which contains
the essential parts that APISIX needs to run.

> See
> [config-default.yaml](https://github.com/apache/apisix/blob/master/conf/config-default.yaml)
> to learn the completed APISIX Configuration.
> See [APISIX Configuration Template API](https://docs.az-staging.api7.cloud/swagger/#/controlplanes_operation/getControlPlaneStartupConfig)
> for the details.

Run Command
-----------

```shell
cloud-cli deploy docker \
  --apisix-image apache/apisix:2.15.0-centos \
  --name my-apisix

Congratulations! Your APISIX instance was deployed successfully
Container ID: 1b2e54380cdc
APISIX ID: 4189c82c-fdf1-40f2-87e2-9a7bb6ad5ed7
```

In this command, we:

1. name the container to `my-apisix`;
2. use the APISIX image `apache/apisix:2.15.0-centos`.

If you see the similar output about the instance ID and container ID, then your
APISIX instance is deployed successfully. You can redirect to API7 Cloud console
to check the status of your APISIX instance.

> You can also run the `docker ps` command to check the status of the container.

Besides, the container will expose the port `9080` and `9443` to the host, so
you can access your APISIX instance through `127.0.0.1:9080` (HTTP) or
`127.0.0.1:9443` (HTTPS). Care must be taken here that you **cannot run** another
APISIX instance on the same machine due to the port conflict.

> Note: we always run the container in the background.

### Cloud Lua Module Mirror

During the deployment, Cloud CLI has to download the [Cloud Lua Module](https://api7.cloud/docs/overview/how-apisix-connects-to-api7-cloud#the-api7-cloud-lua-module)
, users in China may suffer from the slow network. In such a case, try to export the below environment.

```shell
export API7_CLOUD_LUA_MODULE_URL=https://api7-cloud-1301662268.cos.ap-nanjing.myqcloud.com/latest/assets/cloud_module_beta.tar.gz
```

### Persistent Local Configuration Cache

Apache APISIX will save the configuration to the local file (`/usr/local/apisix/conf/apisix.data`), however, this
file will disappear if the container is removed. To avoid this, you can mount a host file to the container with the `--local-cache-bind-path` option.

```shell
cloud-cli deploy docker \
--apisix-image apache/apisix:2.15.0-centos \
--name my-apisix \
--local-cache-bind-path /path/to/apisix.data
```

Now the local configuration cache will be saved to the host file `/path/to/apisix.data`. 

Stop Instance
-------------

If you want to stop the container, just run the command below:

```shell
cloud-cli stop docker --name my-apisix
```

This command will stop the container but won't remove it (unless you already add
the `--rm` flag when you deploy it), so if you want to remove the container, just
run the following command:

```shell
cloud-cli stop docker --name my-apisix --rm
```

Command Option Reference
------------------------

You can run `cloud-cli deploy docker --help` to learn the command line option meanings.
