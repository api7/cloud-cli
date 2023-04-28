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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
	sdk "github.com/api7/cloud-go-sdk"
	"strconv"
)

var (
	_resourceListHandler = map[string]func() interface{}{
		"cluster": func() interface{} {
			user, err := cloud.Client().Me()
			if err != nil {
				output.Errorf(err.Error())
			}
			limit := options.Global.Resource.List.Limit
			skip := options.Global.Resource.List.Skip
			clusters, err := cloud.DefaultClient.ListClusters(user.OrgIDs[0], limit, skip)
			if err != nil {
				output.Errorf("Failed to list clusters: %s", err.Error())
			}
			if clusters == nil {
				return nil
			}
			return clusters
		},
		"service": func() interface{} {
			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get default cluster: %s", err.Error())
			}
			limit := options.Global.Resource.List.Limit
			skip := options.Global.Resource.List.Skip
			services, err := cloud.DefaultClient.ListServices(cluster.ID, limit, skip)
			if err != nil {
				output.Errorf("Failed to list services: %s", err.Error())
			}
			if services == nil {
				return nil
			}
			return services
		},
		"ssl": func() interface{} {
			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get default cluster: %s", err.Error())
			}
			limit := options.Global.Resource.List.Limit
			skip := options.Global.Resource.List.Skip

			ssl, err := cloud.DefaultClient.ListSSL(cluster.ID, limit, skip)
			if err != nil {
				output.Errorf("Failed to list ssl: %s", err.Error())
			}
			if ssl == nil {
				return nil
			}
			return ssl
		},
		"consumer": func() interface{} {
			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get default cluster: %s", err.Error())
			}
			limit := options.Global.Resource.List.Limit
			skip := options.Global.Resource.List.Skip

			ssl, err := cloud.DefaultClient.ListConsumers(cluster.ID, limit, skip)
			if err != nil {
				output.Errorf("Failed to list ssl: %s", err.Error())
			}
			if ssl == nil {
				return nil
			}
			return ssl
		},
		"route": func() interface{} {
			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get default cluster: %s", err.Error())
			}
			limit := options.Global.Resource.List.Limit
			skip := options.Global.Resource.List.Skip

			serviceID := options.Global.Resource.List.ServiceID
			uint64ServiceID, err := strconv.ParseUint(serviceID, 10, 64)
			if err != nil {
				output.Errorf("Failed to parse service-id: %s", err.Error())
			}
			if uint64ServiceID == 0 {
				output.Errorf("service-id is required")
			}

			routes, err := cloud.DefaultClient.ListRoutes(cluster.ID, sdk.ID(uint64ServiceID), limit, skip)
			if err != nil {
				output.Errorf("Failed to list routes: %s", err.Error())
			}
			if routes == nil {
				return nil
			}
			return routes
		},
	}
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Show the resource list by the Cloud CLI",
		Example: `cloud-cli resource list [RESOURCE] [ARGS...]`,
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
			kind := options.Global.Resource.List.Kind
			handler, ok := _resourceListHandler[kind]
			if !ok {
				output.Errorf("This kind of resource is not supported")
			} else {
				resource := handler()
				if resource != nil {
					text, _ := json.MarshalIndent(resource, "", "\t")
					fmt.Println(string(text))
				}
			}
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Resource.List.Kind, "kind", "cluster", "Specify the resource kind")
	cmd.PersistentFlags().IntVar(&options.Global.Resource.List.Limit, "limit", 10, "Specify the amount of data to be listed")
	cmd.PersistentFlags().IntVar(&options.Global.Resource.List.Skip, "skip", 0, "Specifies how much data to skip ahead")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.List.ServiceID, "service-id", "", "Specify the service id")
	return cmd
}
