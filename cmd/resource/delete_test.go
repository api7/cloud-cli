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
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/persistence"
	sdk "github.com/api7/cloud-go-sdk"
)

func TestServiceDelete(t *testing.T) {
	testCases := []struct {
		name      string
		config    *persistence.CloudConfiguration
		args      []string
		mockCloud func(api *cloud.MockAPI)
		outputs   []string
	}{
		{
			name: "delete service",
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
			args: []string{"delete", "--kind", "service", "--id", "123"},
			mockCloud: func(api *cloud.MockAPI) {
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 123,
				}, nil)
				api.EXPECT().DeleteService(sdk.ID(123), sdk.ID(123)).Return(nil)
			},
			outputs: []string{""},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := persistence.SaveConfiguration(tc.config)
			assert.NoError(t, err, "prepare fake cloud configuration")

			ctrl := gomock.NewController(t)
			api := cloud.NewMockAPI(ctrl)
			cloud.NewClient = func(_ string, _ string, _ bool) (cloud.API, error) {
				return api, nil
			}
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				if tc.mockCloud != nil {
					tc.mockCloud(api)
				}
				cloud.DefaultClient = api
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
