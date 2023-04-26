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

func TestResourceList(t *testing.T) {
	testCases := []struct {
		name      string
		config    *persistence.CloudConfiguration
		args      []string
		mockCloud func(api *cloud.MockAPI)
		outputs   []string
	}{
		{
			name: "list clusters",
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
			args: []string{"list", "--kind", "cluster"},
			mockCloud: func(api *cloud.MockAPI) {
				api.EXPECT().Me().Return(&sdk.User{
					Email:  "demo@api7.cloud",
					OrgIDs: []sdk.ID{123},
				}, nil)
				api.EXPECT().ListClusters(sdk.ID(123), gomock.Any(), gomock.Any()).Return([]*sdk.Cluster{
					{
						ID:   123,
						Name: "API7.AI",
					},
				}, nil)
			},
			outputs: []string{"[\n\t{\n\t\t\"org_id\": \"0\",\n\t\t\"region_id\": \"0\",\n\t\t\"status\": 0,\n\t\t\"domain\": \"\",\n\t\t\"settings\": {\n\t\t\t\"client_settings\": {\n\t\t\t\t\"client_real_ip\": {\n\t\t\t\t\t\"replace_from\": {},\n\t\t\t\t\t\"recursive_search\": false,\n\t\t\t\t\t\"enabled\": false\n\t\t\t\t},\n\t\t\t\t\"maximum_request_body_size\": 0\n\t\t\t},\n\t\t\t\"observability_settings\": {\n\t\t\t\t\"metrics\": {\n\t\t\t\t\t\"enabled\": false\n\t\t\t\t},\n\t\t\t\t\"show_upstream_status_in_response_header\": false,\n\t\t\t\t\"access_log_rotate\": {\n\t\t\t\t\t\"enabled\": false,\n\t\t\t\t\t\"enable_compression\": false\n\t\t\t\t}\n\t\t\t},\n\t\t\t\"api_proxy_settings\": {\n\t\t\t\t\"enable_request_buffering\": false,\n\t\t\t\t\"url_handling_options\": null\n\t\t\t}\n\t\t},\n\t\t\"config_version\": 0,\n\t\t\"id\": \"123\",\n\t\t\"name\": \"API7.AI\",\n\t\t\"created_at\": \"0001-01-01T00:00:00Z\",\n\t\t\"updated_at\": \"0001-01-01T00:00:00Z\"\n\t}\n]"},
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

			fmt.Println(string(output))
			for _, o := range tc.outputs {
				assert.Contains(t, string(output), strings.TrimSpace(o), "check output")
			}

		})
	}
}
