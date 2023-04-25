// Copyright 2022 API7.ai, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Show the services list by the Cloud CLI",
		Example: `cloud-cli service list`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := options.Global.Service.List.Validate(); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			kind := options.Global.Service.List.Kind
			limit := options.Global.Service.List.Limit
			skip := options.Global.Service.List.Skip

			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf(err.Error())
			}
			if kind == "service" {
				servicesList, err := cloud.DefaultClient.ListServices(cluster.ID, limit, skip)
				if err != nil {
					output.Errorf("Failed to list services: %s", err.Error())
				}
				json, _ := json.MarshalIndent(servicesList, "", "\t")
				fmt.Println(string(json))
				return
			}
			output.Errorf("This kind of resource is not supported")
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Service.List.Kind, "kind", "service", "Specify the resource kind")
	cmd.PersistentFlags().IntVar(&options.Global.Service.List.Limit, "limit", 10, "Specify the amount of data to be listed")
	cmd.PersistentFlags().IntVar(&options.Global.Service.List.Skip, "skip", 0, "Specifies how much data to skip ahead")

	return cmd
}
