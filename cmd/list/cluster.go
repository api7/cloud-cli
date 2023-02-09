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
package list

import (
	"encoding/json"
	"fmt"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/spf13/cobra"
)

func newClustersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clusters",
		Short: "List the clusters obtained from the API7 Cloud.",
		Example: `
cloud-cli list clusters --count 10 --skip 0 `,
		Run: func(cmd *cobra.Command, args []string) {
			user, err := cloud.Client().Me()
			if err != nil {
				output.Errorf(err.Error())
			}
			clustersList, err := cloud.DefaultClient.ListClusters(user.OrgIDs[0], options.Global.List.Clusters.Count, options.Global.List.Clusters.Skip)
			if err != nil {
				output.Errorf("Failed to list clusters: %s", err.Error())
			}
			json, _ := json.MarshalIndent(clustersList, "", "\t")
			fmt.Println(string(json))
		},
	}

	cmd.PersistentFlags().IntVar(&options.Global.List.Clusters.Count, "count", 10, "Specify the amount of data to be listed")
	cmd.PersistentFlags().IntVar(&options.Global.List.Clusters.Skip, "skip", 1, "Specifies how much data to skip ahead")

	return cmd
}
