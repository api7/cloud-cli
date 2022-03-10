//  Licensed to the Apache Software Foundation (ASF) under one or more
//  contributor license agreements.  See the NOTICE file distributed with
//  this work for additional information regarding copyright ownership.
//  The ASF licenses this file to You under the Apache License, Version 2.0
//  (the "License"); you may not use this file except in compliance with
//  the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

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
	"github.com/api7/cloud-cli/internal/types"
)

func TestBareMetalDeployCommand(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		configFile string
		cmdPattern string
		mockCloud  func(t *testing.T)
	}{
		{
			name: "test deploy bare metal command",
			args: []string{"bare", "--apisix-version", "2.11.0"},
			cmdPattern: fmt.Sprintf(`yum install -y https://repos.apiseven.com/packages/centos/apache-apisix-repo-1\.0-1\.noarch\.rpm
yum-config-manager --add-repo https://repos\.apiseven\.com/packages/centos/apache-apisix\.repo
yum install -y --downloadonly --downloaddir=%s/\.api7cloud/rpm/2\.11\.0 apisix-2\.11\.0
yum install -y %s/\.api7cloud/rpm/2.11.0/\*.rpm
apisix start
PASS
`, os.Getenv("HOME"), os.Getenv("HOME")),
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().Me().Return(&types.User{
					ID:        "12345",
					FirstName: "Bob",
					LastName:  "Alice",
					Email:     "test@api7.ci",
					OrgIDs:    []string{"org1"},
				}, nil)
				api.EXPECT().ListControlPlanes(gomock.Any()).Return([]*types.ControlPlaneSummary{
					{
						ControlPlane: types.ControlPlane{
							TypeMeta: types.TypeMeta{
								ID: "12345",
							},
							OrganizationID: "org1",
						},
					},
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				cloud.DefaultClient = api
			},
		},
		{
			name: "test deploy bare metal command with apisix config",
			args: []string{"bare", "--apisix-version", "2.11.0", "--apisix-config", "./testdata/apisix.yaml"},
			cmdPattern: fmt.Sprintf(`yum install -y https://repos\.apiseven\.com/packages/centos/apache-apisix-repo-1\.0-1\.noarch\.rpm
yum-config-manager --add-repo https://repos\.apiseven\.com/packages/centos/apache-apisix\.repo
yum install -y --downloadonly --downloaddir=%s/\.api7cloud/rpm/2\.11\.0 apisix-2\.11\.0
yum install -y %s/\.api7cloud/rpm/2.11.0/\*.rpm
cp -f .*/apisix-config-\w+\.yaml /usr/local/apisix/conf/config\.yaml
apisix start
PASS
`, os.Getenv("HOME"), os.Getenv("HOME")),
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().Me().Return(&types.User{
					ID:        "12345",
					FirstName: "Bob",
					LastName:  "Alice",
					Email:     "test@api7.ci",
					OrgIDs:    []string{"org1"},
				}, nil)
				api.EXPECT().ListControlPlanes(gomock.Any()).Return([]*types.ControlPlaneSummary{
					{
						ControlPlane: types.ControlPlane{
							TypeMeta: types.TypeMeta{
								ID: "12345",
							},
							OrganizationID: "org1",
						},
					},
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				cloud.DefaultClient = api
			},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				certFilename := filepath.Join(os.Getenv("HOME"), ".api7cloud", "tls", "tls.crt")
				certKeyFilename := filepath.Join(os.Getenv("HOME"), ".api7cloud", "tls", "tls.key")
				certCAFilename := filepath.Join(os.Getenv("HOME"), ".api7cloud", "tls", "ca.crt")
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
