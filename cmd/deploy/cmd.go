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

package deploy

import (
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

// NewCommand creates the deploy sub-command object.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy [COMMAND] [ARG...]",
		Short: "Deploy Apache APISIX with being connected to API7 Cloud.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := persistence.CheckConfigurationAndInitCloudClient(); err != nil {
				output.Errorf(err.Error())
			}
		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Profile, "profile", "", "The profile of the configuration to use")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Name, "name", "apisix", "The identifier of this deployment, it would be the container name (on Docker), the helm release (on Kubernetes) and it's useless if APISIX is deployed on bare metal")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.APISIXConfigFile, "apisix-config", "", "Specify the custom APISIX configuration file")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.APISIXInstanceID, "apisix-id", "", "Specify the custom APISIX instance ID")

	cmd.AddCommand(newDockerCommand())
	cmd.AddCommand(newBareCommand())
	cmd.AddCommand(newKubernetesCommand())

	return cmd
}
