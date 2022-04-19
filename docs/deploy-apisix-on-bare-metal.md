<!--
#  Copyright 2022 API7.ai, Inc under one or more contributor license
#  agreements.  See the NOTICE file distributed with this work for
#  additional information regarding copyright ownership.
#  The ASF licenses this file to You under the Apache License, Version 2.0
#  (the "License"); you may not use this file except in compliance with
#  the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
-->

Deploy APISIX on Bare Metal
=======================

In this section, you'll learn how to deploy APISIX on Bare Metal (**only CentOS 7
is supported**) through Cloud CLI.

> Note, before you go ahead, and please make sure you read the section
> [How to Configure Cloud CLI](./configuring-cloud-cli.md)

Cloud CLI will help you to install Apache APISIX, Cloud Lua Module, and generate
configuration (to communicate with Control Plane).

* Apache APISIX

Cloud CLI will install Apache APISIX via the RPM package, see
[Installation via RPM Repository (CentOS 7)](https://apisix.apache.org/docs/apisix/how-to-build#installation-via-rpm-repository-centos-7)
for more details.

* The Cloud Lua Module

The Cloud Lua Module contains codes to communicate with API7 Cloud (such as
heartbeat, status reporting, etc.), it'll be downloaded every time you run the command.

> Currently, the Cloud Lua Module will be downloaded from [api7/cloud-scripts](https://github.com/api7/cloud-scripts).

* The TLS Bundle

TLS Bundle (Certificate, Private Key, CA Bundle) will be downloaded from API7
Cloud, only instances with a valid client certificate can connect to API7 Cloud.

Cloud CLI will copy TLS Bundle into installation directory of Apache APISIX to
avoid permission issue.

> See the
> [DP Certificate API](https://docs.az-staging.api7.cloud/swagger/#/controlplanes_operation/getCertificates)
> to learn the details.

* The APISIX Configuration Template

Cloud CLI will also download the APISIX configuration template, which contains
the essential parts that APISIX needs to run.

> See
> [config-default.yaml](https://github.com/apache/apisix/blob/master/conf/config-default.yaml)
> to learn the completed APISIX Configuration.
> See [APISIX Configuration Template API](https://docs.az-staging.api7.cloud/swagger/#/controlplanes_operation/getControlPlaneStartupConfig)
> for the details.

Run Command
-----------

```shell
cloud-cli deploy bare --apisix-version 2.11.0

Congratulations! Your APISIX instance was deployed successfully
APISIX ID: 4189c82c-fdf1-40f2-87e2-9a7bb6ad5ed7
```

In this command, we:

1. install Apache APISIX;
2. load Cloud Lua Module and start up Apache APISIX instance;

If you see a similar output about the message
then your APISIX instance is deployed successfully. You can
redirect to API7 Cloud console to check the status of your APISIX instance.

> You can also run the `ps -ef | grep apisix` command to check the status of the
> Apache APISIX service.

Besides, Apache APISIX service will listen the ports `9080` for HTTP traffic and
`9443` for HTTPS. Care must be taken here that you may suffer from the "port is
already in use" issue if these ports were occupied.

Stop Instance
-------------

If you want to stop the APISIX instance, just run the command below:

```shell
cloud-cli stop bare
```

Command Option Reference
------------------------

You can run `cloud-cli deploy bare --help` to learn the command line option meanings.
