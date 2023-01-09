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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
)

func TestDebugShowConfig(t *testing.T) {
	testCases := []struct {
		name         string
		args         []string
		errorMessage string
		output       string
		mockCloud    func(t *testing.T)
	}{
		{
			name:         "invalid cluster",
			args:         []string{"show-config", "api", "--id", "123"},
			errorMessage: "ERROR: mock error\n",
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultCluster().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = api
			},
		},
		{
			name:         "app not found",
			args:         []string{"show-config", "application", "--id", "123"},
			errorMessage: "ERROR: Failed to show config: not found\n",
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 12345,
					ClusterSpec: sdk.ClusterSpec{
						OrganizationID: 1,
					},
				}, nil)
				api.EXPECT().DebugShowConfig(sdk.ID(12345), "application", "123").Return("", errors.New("not found"))
				cloud.DefaultClient = api
			},
		},
		{
			name: "show application related APISIX objects",
			args: []string{"show-config", "application", "--id", "123"},
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 12345,
					ClusterSpec: sdk.ClusterSpec{
						OrganizationID: 1,
					},
				}, nil)
				resources := `
{
  "upstreams": [
    {
       "nodes": [
         {"host": "127.0.0.1", "port": 9080, "weight": 1}
       ]
    }
  ]
}`
				api.EXPECT().DebugShowConfig(sdk.ID(12345), "application", "123").Return(resources, nil)
				cloud.DefaultClient = api
			},
			output: `
{
  "upstreams": [
    {
       "nodes": [
         {"host": "127.0.0.1", "port": 9080, "weight": 1}
       ]
    }
  ]
}
`,
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				tc.mockCloud(t)
				cmd := NewCommand()
				cmd.SetArgs(tc.args)
				err := cmd.Execute()
				assert.NoError(t, err, "check if the command executed successfully")
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")

			// Ignore error since it won't be nil if we mock a failed command.
			output, _ := cmd.CombinedOutput()

			// Use assert.Contains since the output might contains some noises like
			// assert output in the subprocess.
			if tc.errorMessage != "" {
				assert.Contains(t, string(output), tc.errorMessage, "check error message")
			} else {
				assert.Contains(t, string(output), tc.output, "check output")
			}
		})
	}
}
