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

package config

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

func newViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "Show the configuration details currently used by the Cloud CLI",
		Example: `cloud-cli config view`,
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

			// output as ascii table
			rows := [][]string{}
			for _, profile := range config.Profiles {
				var (
					clusterName = "-"
					orgName     = "-"
				)

				if api, err := cloud.NewClient(profile.Address, profile.User.AccessToken); err != nil {
					output.Warnf("Failed to create API7 Cloud client for profile %s: %s", profile.Name, err.Error())
				} else {
					if cluster, err := api.GetDefaultCluster(); err != nil {
						output.Warnf("Failed to get default cluster for profile %s: %s", profile.Name, err.Error())
					} else {
						clusterName = cluster.Name
					}

					if org, err := api.GetDefaultOrganization(); err != nil {
						output.Warnf("Failed to get default organization for profile %s: %s", profile.Name, err.Error())
					} else {
						orgName = org.Name
					}
				}

				rows = append(rows, []string{profile.Name, orgName, clusterName, strconv.FormatBool(profile.Name == config.DefaultProfile), profile.Address})
			}

			table := tablewriter.NewWriter(os.Stdout)

			table.SetHeader([]string{"Profile Name", "Organization", "Cluster", "Is Default", "API7 Cloud Address"})

			for _, r := range rows {
				table.Append(r)
			}
			table.Render()
		},
	}

	return cmd
}
