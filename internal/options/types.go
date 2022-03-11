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
package options

var (
	// Global contains all options.
	Global Options
)

// Options contains all options.
type Options struct {
	// Verbose controls if the output should be elaborate.
	Verbose bool
	// DryRun controls if all the actions should be simulated instead of executed.
	DryRun bool
	// Deploy contains the options for the deploy command.
	Deploy DeployOptions
}

// DeployOptions contains options for the deploy command.
type DeployOptions struct {
	// Name is an identifier of this deployment.
	// It'll be container name if deploy on Docker;
	// It'll be the Helm release name if deploy on Kubernetes;
	// It'll be noop if deploy on Bare metal.
	Name string
	// APISIXInstanceID specifies the ID of the APISIX instance to deploy.
	// When this field is empty, the instance ID will be generated automatically.
	APISIXInstanceID string `validate:"min=1 max=128"`
	// APISIXConfigFile is the path to the APISIX configuration file.
	APISIXConfigFile string
	// Docker contains the options for the deploy docker command.
	Docker DockerDeployOptions
}

// DockerDeployOptions contains options for the deploy docker command.
type DockerDeployOptions struct {
	// APISIXImage is the name of the APISIX image to deploy.
	APISIXImage string `validate:"image"`
	// DockerRunArgs contains a series of arguments to pass to the docker run command.
	DockerRunArgs []string
	// DockerCLIPath is the filepath of the docker command.
	DockerCLIPath string
}
