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

package deploy

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/consts"
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
			accessToken := os.Getenv(consts.Api7CloudAccessTokenEnv)
			if accessToken == "" {
				credential, err := persistence.LoadCredential()
				if err != nil {
					output.Errorf("Failed to load credential: %s,\nPlease run 'cloud-cli configure' first, access token can be created from https://console.api7.cloud", err)
				}
				accessToken = credential.User.AccessToken
			}

			output.Verbosef("Loaded access token: %s", accessToken)

			if err := cloud.InitDefaultClient(accessToken); err != nil {
				output.Errorf("Failed to init api7 cloud client: %s", err)
				return
			}
			if err := persistence.PrepareCertificate(); err != nil {
				output.Errorf("Failed to prepare certificate: %s", err)
				return
			}
		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Deploy.APISIXConfigFile, "apisix-config", "", "Specify the custom APISIX configuration file")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.APISIXInstanceID, "apisix-id", "", "Specify the custom APISIX instance ID")

	cmd.AddCommand(newDockerCommand())
	cmd.AddCommand(newBareCommand())

	return cmd
}
