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

package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/persistence"
)

func TestConfigSwitch(t *testing.T) {
	testCases := []struct {
		name           string
		config         *persistence.CloudConfiguration
		args           []string
		output         string
		profileExpectd string
	}{
		{
			name: "switch to non exist profile",
			config: &persistence.CloudConfiguration{
				DefaultProfile: "prod",
				Profiles: []persistence.Profile{
					{
						Name:    "prod",
						Address: "https://prod.api7.ai",
						User: persistence.User{
							AccessToken: "prod-token",
						},
					},
				},
			},
			args:           []string{"switch", "dev"},
			output:         `ERROR: profile dev not found`,
			profileExpectd: "prod",
		},
		{
			name: "switch to exist profile",
			config: &persistence.CloudConfiguration{
				DefaultProfile: "prod",
				Profiles: []persistence.Profile{
					{
						Name:    "prod",
						Address: "https://prod.api7.ai",
						User: persistence.User{
							AccessToken: "prod-token",
						},
					},
					{
						Name:    "dev",
						Address: "https://dev.api7.ai",
						User: persistence.User{
							AccessToken: "dev-token",
						},
					},
				},
			},
			args:           []string{"switch", "dev"},
			output:         `switched to profile: dev`,
			profileExpectd: "dev",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := persistence.SaveConfiguration(tc.config)
			assert.NoError(t, err, "prepare fake cloud configuration")

			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				cmd := NewCommand()
				cmd.SetArgs(tc.args)
				err := cmd.Execute()
				assert.NoError(t, err, "check if the command executed successfully")
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")

			output, _ := cmd.CombinedOutput()

			assert.Contains(t, string(output), strings.TrimSpace(tc.output), "check output")

			config, err := persistence.LoadConfiguration()
			assert.NoError(t, err, "load configuration after switch")
			assert.Equal(t, tc.profileExpectd, config.DefaultProfile, "check default profile after switch")
		})
	}
}
