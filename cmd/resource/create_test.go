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
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/persistence"
	sdk "github.com/api7/cloud-go-sdk"
)

func TestServiceCreate(t *testing.T) {

	testCases := []struct {
		name       string
		config     *persistence.CloudConfiguration
		args       []string
		mockCloud  func(api *cloud.MockAPI)
		outputs    string
		testConfig string
	}{
		{
			name: "create service",
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
			args:       []string{"create", "--from-file", path.Join(os.TempDir(), "config.json"), "--kind", "service"},
			testConfig: path.Join(os.TempDir(), "config.json"),
			mockCloud: func(api *cloud.MockAPI) {
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 123,
				}, nil)
				api.EXPECT().CreateService(sdk.ID(123), gomock.Any()).Return(&sdk.Application{
					ID:        sdk.ID(123),
					ClusterID: sdk.ID(123),
					ApplicationSpec: sdk.ApplicationSpec{
						Description: "456",
					},
				}, nil)
			},
			outputs: "{\n\t\"name\": \"\",\n\t\"description\": \"456\",\n\t\"path_prefix\": \"\",\n\t\"hosts\": null,\n\t\"upstreams\": null,\n\t\"active\": 0,\n\t\"id\": \"123\",\n\t\"cluster_id\": \"123\",\n\t\"status\": 0,\n\t\"created_at\": \"0001-01-01T00:00:00Z\",\n\t\"updated_at\": \"0001-01-01T00:00:00Z\",\n\t\"available_cert_ids\": null,\n\t\"canary_release_id\": null,\n\t\"canary_upstream_version_list\": null\n}",
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
			testFile, err := os.Create(tc.testConfig)
			assert.Nil(t, err, "create test file")
			_, err = testFile.Write([]byte(tc.outputs))
			assert.Nil(t, err, "write test file")
			assert.Nil(t, testFile.Close(), "close test file")

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
			assert.Contains(t, string(output), tc.outputs, "check output")

			os.Remove(tc.testConfig)
		})
	}
}
