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
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

func newSwitchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "switch",
		Short:   "Switch the default profile used by Cloud CLI",
		Example: `cloud-cli config switch <profile>`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				output.Errorf("please specify the profile name")
				return
			}
			profileName := args[0]

			config, err := persistence.LoadConfiguration()
			if err != nil {
				output.Errorf(err.Error())
			}

			_, err = config.GetProfile(profileName)
			if err != nil {
				output.Errorf(err.Error())
				return
			}

			config.DefaultProfile = profileName
			if err := persistence.SaveConfiguration(config); err != nil {
				output.Errorf(err.Error())
				return
			}

			output.Infof("switched to profile: %s", profileName)
		},
	}

	return cmd
}
