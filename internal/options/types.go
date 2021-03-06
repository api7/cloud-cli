// Copyright 2022 API7.ai, Inc
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

package options

import (
	"errors"
)

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
	// Stop contains the options for the stop command.
	Stop StopOptions
	// Debug contains the options for the debug command.
	Debug DebugOptions
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
	// Bare contains the options for the bare metal deployment command.
	Bare BareDeployOptions
	// KubernetesDeployOptions contains options for the kubernetes or helm command.
	Kubernetes KubernetesDeployOptions
}

// DockerDeployOptions contains options for the deploy docker command.
type DockerDeployOptions struct {
	// APISIXImage is the name of the APISIX image to deploy.
	APISIXImage string `validate:"image"`
	// DockerRunArgs contains a series of arguments to pass to the docker run command.
	DockerRunArgs []string
	// DockerCLIPath is the filepath of the docker command.
	DockerCLIPath string
	// Specify the host port for HTTP
	HTTPHostPort int
	// Specify the host port for HTTPS
	HTTPSHostPort int
}

// Validate validates the docker deploy options.
func (o *DockerDeployOptions) Validate() error {
	if o.HTTPHostPort <= 0 || o.HTTPHostPort > 65535 {
		return errors.New("invalid http host port")
	}

	if o.HTTPSHostPort <= 0 || o.HTTPSHostPort > 65535 {
		return errors.New("invalid https host port")
	}

	return nil
}

// KubernetesDeployOptions contains options for the kubectl or helm command.
type KubernetesDeployOptions struct {
	// NameSpace is the name space of kubernetes
	NameSpace string
	// APISIXImage is the name of the APISIX image to deploy.
	APISIXImage string `validate:"image"`
	// APISIXImageRepo is the APISIXImage name
	APISIXImageRepo string
	// APISIXImageTag is the APISIXImage tag
	APISIXImageTag string
	// ReplicaCount is the pod replica count
	ReplicaCount uint
	// HelmInstallArgs contains a series of arguments to pass to the helm install command.
	HelmInstallArgs []string
	// KubectlCLIPath is the filepath of the kubectl command.
	KubectlCLIPath string
	// HelmCLIPath is the filepath of the helm command.
	HelmCLIPath string
}

// StopOptions contains options for the stop command.
type StopOptions struct {
	// Name is an identifier of this deployment.
	// It'll be container name if deploy on Docker;
	// It'll be the Helm release name if deploy on Kubernetes;
	// It'll be noop if deploy on Bare metal.
	Name string
	// Remove controls whether to delete containers in docker
	Remove bool
	// Docker contains the options for the stop docker command.
	Docker DockerStopOptions
	// Kubernetes contains options for the kubectl or helm command.
	Kubernetes KubernetesStopOptions
}

// DockerStopOptions contains options for the stop docker command.
type DockerStopOptions struct {
	// DockerCLIPath is the filepath of the docker command.
	DockerCLIPath string
}

// BareDeployOptions contains options for the bare metal deployment command.
type BareDeployOptions struct {
	// APISIXVersion specifies the APISIX version to deploy.
	APISIXVersion string
}

// KubernetesStopOptions contains options for the kubectl or helm command.
type KubernetesStopOptions struct {
	// NameSpace is the name space of kubernetes
	NameSpace string
	// HelmUnInstallArgs contains a series of arguments to pass to the helm uninstall command.
	HelmUnInstallArgs []string
	// KubectlCLIPath is the filepath of the kubectl command.
	KubectlCLIPath string
	// HelmCLIPath is the filepath of the helm command.
	HelmCLIPath string
}

// DebugOptions contains options for `cloud-cli debug` command.
type DebugOptions struct {
	// ShowConfig contains options for `cloud-cli debug show-config` command.
	ShowConfig DebugShowConfigOptions
}

// DebugShowConfigOptions contains options for `cloud-cli debug show-config` command.
type DebugShowConfigOptions struct {
	// ID is the API7 Cloud resource id.
	ID string
}
