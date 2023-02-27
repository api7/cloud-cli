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
	"fmt"
	"strconv"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "get the resource detail by the Cloud CLI",
		Example: `cloud-cli resource get [RESOURCE] [ARGS...]`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := options.Global.Resource.List.Validate(); err != nil {
				output.Errorf(err.Error())
				return
			}

			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			kind := options.Global.Resource.Get.Kind
			ID := options.Global.Resource.Get.ID
			if kind == "cluster" {
				uint64ID, _ := strconv.ParseUint(ID, 10, 64)
				cluster, err := cloud.DefaultClient.GetClusterDetail(sdk.ID(uint64ID))
				if err != nil {
					output.Errorf("Failed to get cluster detail: %s", err.Error())
				}

				clusterDetail, _ := json.MarshalIndent(cluster, "", "\t")
				fmt.Println(string(clusterDetail))
				return
			}
			output.Errorf("This kind of resource is not supported")
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Get.Kind, "kind", "cluster", "Specify the resource kind")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Get.ID, "id", "", "Specify the id of resource")
	return cmd
}
