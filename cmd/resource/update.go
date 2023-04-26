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
)

var (
	_resourceUpdateHandler = map[string]func(config string){
		"service": func(config string) {
			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to list services: %s", err.Error())
			}
			servicesList, err := cloud.DefaultClient.UpdateService(cluster.ID, config)
			if err != nil {
				output.Errorf("Failed to list services: %s", err.Error())
			}
			json, _ := json.MarshalIndent(servicesList, "", "\t")
			fmt.Println(string(json))
		},
	}
)

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update the service by the Cloud CLI",
		Example: `cloud-cli service update [RESOURCE] [ARGS...]`,
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
			kind := options.Global.Resource.Update.Kind
			config := options.Global.Resource.Update.Config
			handler, ok := _resourceUpdateHandler[kind]
			if !ok {
				output.Errorf("This kind of resource is not supported")
			} else {
				handler(config)
			}
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Update.Kind, "kind", "service", "Specify the service kind")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Update.Config, "config", "", "Specify the config path of service")
	return cmd
}
