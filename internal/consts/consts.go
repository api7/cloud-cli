// Copyright 2023 API7.ai, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consts

import "time"

const (
	// Api7CloudLuaModuleURL is the environment variable used to specify the API7 Cloud Lua module address.
	// e.g. https://github.com/api7/cloud-scripts/raw/main/assets/cloud_module_beta.tar.gz.
	// Note this variable should be deprecated once we can download the module from API7 Cloud.
	Api7CloudLuaModuleURL = "API7_CLOUD_LUA_MODULE_URL"
	// Api7CloudProfile is the environment variable used to specify the API7 Cloud profile.
	Api7CloudProfile = "API7_CLOUD_PROFILE"
)

const (
	// DefaultDeploymentName is the default name for the cloud-cli deploy operation.
	DefaultDeploymentName = "apisix"
)

const (
	// DefaultConfigMapName is the default name for the configMap when deploy on kubernetes
	DefaultConfigMapName = "cloud-module"
	// DefaultSecretName is the default name for the secret when deploy on kubernetes
	DefaultSecretName = "cloud-ssl"
)

const (
	// DefaultKubectlTimeout is the default timeout for execute kubectl command.
	DefaultKubectlTimeout = time.Minute
	// DefaultHelmTimeout is the default timeout for execute helm command.
	DefaultHelmTimeout = time.Minute * 5
)
