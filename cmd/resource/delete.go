// Copyright 2023 API7.ai, Inc
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
	"strconv"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
	sdk "github.com/api7/cloud-go-sdk"
)

var (
	_resourceDeleteHandler = map[string]func(id sdk.ID, args ...any){
		"ssl": func(id sdk.ID, args ...any) {
			cluster, err := cloud.DefaultClient.GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get the default cluster: %s", err.Error())
			}
			if err := cloud.DefaultClient.DeleteSSL(cluster.ID, id); err != nil {
				output.Errorf("Failed to delete ssl: %s", err.Error())
			}
		},
		"service": func(id sdk.ID, args ...any) {
			cluster, err := cloud.DefaultClient.GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get the default cluster: %s", err.Error())
			}
			if err := cloud.DefaultClient.DeleteService(cluster.ID, id); err != nil {
				output.Errorf("Failed to delete service: %s", err.Error())
			}
		},
		"route": func(id sdk.ID, args ...any) {
			cluster, err := cloud.DefaultClient.GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get the default cluster: %s", err.Error())
			}
			if len(args) == 0 {
				output.Errorf("Please specify the service id")
			}
			serviceID, ok := args[0].(sdk.ID)
			if !ok {
				output.Errorf("Please specify the correct service id")
			}
			if serviceID == 0 {
				output.Errorf("Please specify the correct service id")
			}

			if err := cloud.DefaultClient.DeleteAPI(cluster.ID, serviceID, id); err != nil {
				output.Errorf("Failed to delete route: %s", err.Error())
			}
		},
	}
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a resource",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			kind := options.Global.Resource.Delete.Kind
			id := options.Global.Resource.Delete.ID
			serviceID := options.Global.Resource.Delete.ServiceID
			handler, ok := _resourceDeleteHandler[kind]
			if !ok {
				output.Errorf("This kind of resource is not supported")
			} else {
				uint64ID, err := strconv.ParseUint(id, 10, 64)
				if err != nil {
					output.Errorf("Failed to parse id: %s", id)
					return
				}

				uint64ServiceID, _ := strconv.ParseUint(serviceID, 10, 64)
				handler(sdk.ID(uint64ID), sdk.ID(uint64ServiceID))
			}
		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Resource.Delete.Kind, "kind", "cluster", "Specify the resource kind")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Delete.ID, "id", "", "Specify the id of resource")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Delete.ServiceID, "service-id", "0", "Specify the id of service resource, when delete API this value should be set")

	return cmd
}
