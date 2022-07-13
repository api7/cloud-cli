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

package debug

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

func newShowConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-config [RESOURCE] [ARG...]",
		Short: "Show translated Apache APISIX configurations related to the specify API7 Cloud resource.",
		Example: `
cloud-cli debug show-config api \
	--id 0e3851a5f4a7`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				output.Errorf("Please specify an API7 Cloud resource. Resource can be application, api, consumer and certificate.")
			}

			id := options.Global.Debug.ShowConfig.ID
			if id == "" {
				output.Errorf("Empty resource ID, please specify --id option")
			}

			defaultCP, err := cloud.DefaultClient.GetDefaultControlPlane()
			if err != nil {
				output.Errorf(err.Error())
			}

			data, err := cloud.DefaultClient.DebugShowConfig(defaultCP.ID, args[0], id)
			if err != nil {
				output.Errorf("Failed to show config: %s", err.Error())
			}
			fmt.Println(data)
		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Debug.ShowConfig.ID, "id", "", "Specify the API7 Cloud resource ID")

	return cmd
}
