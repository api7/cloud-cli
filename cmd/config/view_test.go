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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/persistence"
)

func TestConfigView(t *testing.T) {
	testCases := []struct {
		name      string
		config    *persistence.CloudConfiguration
		args      []string
		mockCloud func(api *cloud.MockAPI)
		outputs   []string
	}{
		{
			name: "one profile with default output format",
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
			args: []string{"view"},
			mockCloud: func(api *cloud.MockAPI) {
				api.EXPECT().GetDefaultOrganization().Return(&sdk.Organization{
					ID:   123,
					Name: "API7.AI",
				}, nil)
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID:   456,
					Name: "default",
				}, nil)
			},
			outputs: []string{`
+--------------+--------------+---------+------------+----------------------+
| PROFILE NAME | ORGANIZATION | CLUSTER | IS DEFAULT |  API7 CLOUD ADDRESS  |
+--------------+--------------+---------+------------+----------------------+
| prod         | API7.AI      | default | true       | https://prod.api7.ai |
+--------------+--------------+---------+------------+----------------------+
				`},
		},
		{
			name: "two profiles with one bad profile",
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
			args: []string{"view"},
			mockCloud: func(api *cloud.MockAPI) {
				api.EXPECT().GetDefaultOrganization().Return(nil, errors.New("organization not found"))
				api.EXPECT().GetDefaultCluster().Return(nil, errors.New("cluster not found"))
				api.EXPECT().GetDefaultOrganization().Return(&sdk.Organization{
					ID:   321,
					Name: "APACHE",
				}, nil)
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID:   654,
					Name: "default",
				}, nil)
			},
			outputs: []string{`

+--------------+--------------+---------+------------+----------------------+
| PROFILE NAME | ORGANIZATION | CLUSTER | IS DEFAULT |  API7 CLOUD ADDRESS  |
+--------------+--------------+---------+------------+----------------------+
| prod         | -            | -       | true       | https://prod.api7.ai |
| dev          | APACHE       | default | false      | https://dev.api7.ai  |
+--------------+--------------+---------+------------+----------------------+
			`,
				"WARNING: Failed to get default cluster for profile prod: cluster not found",
				"WARNING: Failed to get default organization for profile prod: organization not found",
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := persistence.SaveConfiguration(tc.config)
			assert.NoError(t, err, "prepare fake cloud configuration")

			ctrl := gomock.NewController(t)
			api := cloud.NewMockAPI(ctrl)
			cloud.NewClient = func(_ string, _ string) (cloud.API, error) {
				return api, nil
			}
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				if tc.mockCloud != nil {
					tc.mockCloud(api)
				}
				cmd := NewCommand()
				cmd.SetArgs(tc.args)
				err := cmd.Execute()
				assert.NoError(t, err, "check if the command executed successfully")
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")

			output, _ := cmd.CombinedOutput()

			for _, o := range tc.outputs {
				assert.Contains(t, string(output), strings.TrimSpace(o), "check output")
			}

		})
	}
}
