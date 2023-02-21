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

package resource

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

func newOrgInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "org-info",
		Short:   "Show the resource details by the Cloud CLI",
		Example: `cloud-cli resource org-info`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			config, err := persistence.LoadConfiguration()
			if err != nil {
				output.Errorf(err.Error())
			}

			for _, profile := range config.Profiles {

				if api, err := cloud.NewClient(profile.Address, profile.User.AccessToken); err != nil {
					output.Warnf("Failed to create API7 Cloud client for profile %s: %s", profile.Name, err.Error())
				} else {
					if org, err := api.GetDefaultOrganization(); err != nil {
						output.Warnf("Failed to get default organization for profile %s: %s", profile.Name, err.Error())
					} else {
						out, err := json.Marshal(org)
						if err != nil {
							output.Warnf("Failed to parse organization info %s", err.Error())
							return
						}
						output.Infof(string(out))
					}
				}
			}
		},
	}

	return cmd
}
