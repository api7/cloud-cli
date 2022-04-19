<!--
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
-->

Installation
============

In this section, you'll learn how to install Cloud CLI on your local machine.

Build From Source
-----------------

You can also build the Cloud CLI from source.

> Make sure you installed the [Go](https://go.dev/) environment (version >= `1.16`).

```shell
git clone https://github.com/api7/cloud-cli.git
cd /path/to/cloud-cli
make build
```

The executable will be saved in the `./bin` directory.

Download by Go Install
----------------------

> Will suffer from permission problem before this project opening source.
> Remove this note after the project opening source.

Alternatively, you can download the Cloud CLI by using the `go install` command.

```shell
# Install at tree head:
go install github.com/api7/cloud-cli@main
```

See [Versions](https://go.dev/ref/mod#versions) and
[Pseudo-versions](https://go.dev/ref/mod#pseudo-versions) for how to format the
version suffixes.
