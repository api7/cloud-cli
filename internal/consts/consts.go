//  Licensed to the Apache Software Foundation (ASF) under one or more
//  contributor license agreements.  See the NOTICE file distributed with
//  this work for additional information regarding copyright ownership.
//  The ASF licenses this file to You under the Apache License, Version 2.0
//  (the "License"); you may not use this file except in compliance with
//  the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package consts

const (
	// Api7CloudAddrEnv is the environment variable used to specify the API7 Cloud address,
	// e.g. https://console.api7.cloud
	Api7CloudAddrEnv = "API7_CLOUD_ADDR"
	// Api7CloudAccessTokenEnv is environment variable used to specify the access token for API7 Cloud.
	Api7CloudAccessTokenEnv = "API7_CLOUD_ACCESS_TOKEN"
	// Api7CloudLuaModuleURL is the environment variable used to specify the API7 Cloud Lua module address.
	// e.g. https://github.com/api7/cloud-scripts/raw/main/assets/cloud_module_beta.tar.gz.
	// Note this variable should be deprecated once we can download the module from API7 Cloud.
	Api7CloudLuaModuleURL = "API7_CLOUD_LUA_MODULE_URL"
)

const (
	// CloudLuaModuleDir specifies the directory where the cloud module is stored.
	CloudLuaModuleDir = "cloud_module_dir"
)
