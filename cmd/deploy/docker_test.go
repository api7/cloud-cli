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

package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/api7/cloud-cli/internal/types"
)

func TestDockerDeployCommand(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		configFile string
		cmdPattern string
		mockCloud  func(t *testing.T)
	}{
		{
			name:       "test deploy docker command",
			args:       []string{"docker", "--apisix-image", "apache/apisix:2.15.0-centos"},
			cmdPattern: `docker run --detach --mount type=bind,source=/.+?/\.api7cloud,target=/cloud_lua_module,readonly --mount type=bind,source=/.+?/\.api7cloud/tls,target=/cloud/tls,readonly --mount type=bind,source=/.+?/\.api7cloud/apisix\.uid,target=/usr/local/apisix/conf/apisix.uid,readonly -p 9080:9080 -p 9443:9443 --name apisix --hostname apisix apache/apisix:2.15.0-centos`,
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "12345",
					},
					OrganizationID: "org1",
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				api.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)
				api.EXPECT().GetStartupConfig("12345", cloud.APISIX).Return(_apisixStartupConfigTpl, nil)

				cloud.DefaultClient = api
			},
		},
		{
			name:       "test deploy docker command with custom http and https host ports",
			args:       []string{"docker", "--apisix-image", "apache/apisix:2.15.0-centos", "--http-host-port", "8080", "--https-host-port", "443"},
			cmdPattern: `docker run --detach --mount type=bind,source=/.+?/\.api7cloud,target=/cloud_lua_module,readonly --mount type=bind,source=/.+?/\.api7cloud/tls,target=/cloud/tls,readonly --mount type=bind,source=/.+?/\.api7cloud/apisix\.uid,target=/usr/local/apisix/conf/apisix.uid,readonly -p 8080:9080 -p 443:9443 --name apisix --hostname apisix apache/apisix:2.15.0-centos`,
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "12345",
					},
					OrganizationID: "org1",
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				api.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)
				api.EXPECT().GetStartupConfig("12345", cloud.APISIX).Return(_apisixStartupConfigTpl, nil)

				cloud.DefaultClient = api
			},
		},
		{
			name:       "test deploy docker command with apisix config",
			args:       []string{"docker", "--apisix-image", "apache/apisix:2.15.0-centos", "--apisix-config", "./testdata/apisix.yaml"},
			cmdPattern: `docker run --detach --mount type=bind,source=/.+?/apisix-config-\d+.yaml,target=/usr/local/apisix/conf/config.yaml,readonly --mount type=bind,source=/.+?/\.api7cloud,target=/cloud_lua_module,readonly --mount type=bind,source=/.+?/\.api7cloud/tls,target=/cloud/tls,readonly --mount type=bind,source=/.+?/\.api7cloud/apisix\.uid,target=/usr/local/apisix/conf/apisix.uid,readonly -p 9080:9080 -p 9443:9443 --name apisix --hostname apisix apache/apisix:2.15.0-centos`,
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "12345",
					},
					OrganizationID: "org1",
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				api.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)
				api.EXPECT().GetStartupConfig("12345", cloud.APISIX).Return(_apisixStartupConfigTpl, nil)
				cloud.DefaultClient = api
			},
		},
		{
			name:       "test deploy docker command with complicated docker run arg",
			args:       []string{"docker", "--apisix-image", "apache/apisix:2.15.0-centos", "--docker-run-arg", "\"--mount=type=bind,source=/etc/hosts,target=/etc/hosts,readonly\""},
			cmdPattern: `docker run --mount type=bind,source=/etc/hosts,target=/etc/hosts,readonly --detach --mount type=bind,source=/.+?/\.api7cloud,target=/cloud_lua_module,readonly --mount type=bind,source=/.+?/\.api7cloud/tls,target=/cloud/tls,readonly --mount type=bind,source=/.+?/\.api7cloud/apisix\.uid,target=/usr/local/apisix/conf/apisix.uid,readonly -p 9080:9080 -p 9443:9443 --name apisix --hostname apisix apache/apisix:2.15.0-centos`,
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "12345",
					},
					OrganizationID: "org1",
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				api.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)
				api.EXPECT().GetStartupConfig("12345", cloud.APISIX).Return(_apisixStartupConfigTpl, nil)

				cloud.DefaultClient = api
			},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			persistence.HomeDir = filepath.Join(os.TempDir(), ".api7cloud")
			certFilename := filepath.Join(persistence.HomeDir, "tls", "tls.crt")
			certKeyFilename := filepath.Join(persistence.HomeDir, "tls", "tls.key")
			certCAFilename := filepath.Join(persistence.HomeDir, "tls", "ca.crt")
			defer func() {
				os.Remove(certFilename)
				os.Remove(certKeyFilename)
				os.Remove(certCAFilename)
			}()
			//Because `os.Exit(-1)` will be triggered in the failure case, so here the test is executed using a subprocess
			//The method come from: https://talks.golang.org/2014/testing.slide#23
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				options.Global.DryRun = true
				tc.mockCloud(t)
				cmd := NewCommand()
				cmd.SetArgs(tc.args)
				err := cmd.Execute()
				assert.NoError(t, err, "check if the command executed successfully")
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1", fmt.Sprintf("%s=test-token", consts.Api7CloudAccessTokenEnv))

			output, err := cmd.CombinedOutput()
			assert.NoError(t, err, "check if the command executed successfully")

			assert.Regexp(t, tc.cmdPattern, string(output), "check if the composed docker command is correct")
		})
	}
}
