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

package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/api7/cloud-cli/internal/testutils"
)

func TestBareMetalDeployCommand(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		configFile    string
		cmdPattern    string
		installScript string
		mockCloud     func(t *testing.T)
	}{
		{
			name:       "test deploy bare metal command",
			args:       []string{"bare", "--apisix-version", "2.15.0"},
			cmdPattern: "/usr/bin/bash -C .*/scripts/install\\.sh",
			installScript: `#!/usr/bin/env bash

set -e

version="2\.11\.0"
instance_id=""
apisix_home="/usr/local/apisix"

installed_version=\$\(apisix version 2>/dev/null\) || true
if \[\[ -z \$\{installed_version\} \]\]; then
  yum install -y https://repos.apiseven.com/packages/centos/apache-apisix-repo-1\.0-1\.noarch\.rpm
  yum install -y apisix-\$version
fi

# copy certs to apisix directory to avoid permission issue
cp -rf .*/\.api7cloud/tls /usr/local/apisix/conf/ssl

if \[\[ -n \$\{instance_id\} \]\]; then
  echo "\$\{instance_id\}" > \$\{apisix_home\}/conf/apisix\.uid
fi

apisix start -c .*/apisix-config-\d+\.yaml
status=\$\?

# wait for APISIX started and generated instance id
sleep 1

# get the APISIX instance id when instance id is not set
if \[\[ -z \$\{instance_id\} \]\]; then
  instance_id="\$\(cat \$\{apisix_home\}/conf/apisix\.uid\)"
fi

if \[\[ \$status -eq 0 \]\]; then
  echo "Your APISIX Instance was deployed successfully!"
  echo "Instance ID: \$\{instance_id\}"
fi
`,
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 12345,
					ClusterSpec: sdk.ClusterSpec{
						OrganizationID: 1,
					},
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				api.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)
				api.EXPECT().GetStartupConfig(sdk.ID(12345), cloud.APISIX).Return(_apisixStartupConfigTpl, nil)

				{
					file, err := os.CreateTemp(os.TempDir(), "apisix-cli-etcd-*.lua")
					assert.Nil(t, err, "create temp file path for /usr/local/apisix/apisix/cli/etcd.lua")
					_targetApisixCliEtcdPath = file.Name()
				}
				{
					file, err := os.CreateTemp(os.TempDir(), "apisix-cli-local-storage-*.lua")
					assert.Nil(t, err, "create temp file path for /usr/local/apisix/apisix/cli/local_storage.lua")
					_targetApisixCliLocalStoragePath = file.Name()
				}

				cloud.DefaultClient = api
			},
		},
		{
			name:       "test deploy bare metal command with apisix config",
			args:       []string{"bare", "--apisix-version", "2.15.0", "--apisix-config", "./testdata/apisix.yaml"},
			cmdPattern: "/usr/bin/bash -C .*/scripts/install\\.sh",
			installScript: `#!/usr/bin/env bash

set -e

version="2\.11\.0"
instance_id=""
apisix_home="/usr/local/apisix"

installed_version=\$\(apisix version 2>/dev/null\) || true
if \[\[ -z \$\{installed_version\} \]\]; then
  yum install -y https://repos.apiseven.com/packages/centos/apache-apisix-repo-1\.0-1\.noarch\.rpm
  yum install -y apisix-\$version
fi

# copy certs to apisix directory to avoid permission issue
cp -rf .*/\.api7cloud/tls /usr/local/apisix/conf/ssl

if \[\[ -n \$\{instance_id\} \]\]; then
  echo "\$\{instance_id\}" > \$\{apisix_home\}/conf/apisix\.uid
fi

apisix start -c .*/apisix-config-\d+\.yaml
status=\$\?

# wait for APISIX started and generated instance id
sleep 1

# get the APISIX instance id when instance id is not set
if \[\[ -z \$\{instance_id\} \]\]; then
  instance_id="\$\(cat \$\{apisix_home\}/conf/apisix\.uid\)"
fi

if \[\[ \$status -eq 0 \]\]; then
  echo "Your APISIX Instance was deployed successfully!"
  echo "Instance ID: \$\{instance_id\}"
fi
`,
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 12345,
					ClusterSpec: sdk.ClusterSpec{
						OrganizationID: 1,
					},
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				api.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)
				api.EXPECT().GetStartupConfig(sdk.ID(12345), cloud.APISIX).Return(_apisixStartupConfigTpl, nil)

				{

					file, err := os.CreateTemp(os.TempDir(), "apisix-cli-etcd-*.lua")
					assert.Nil(t, err, "create temp file path for /usr/local/apisix/apisix/cli/etcd.lua")
					_targetApisixCliEtcdPath = file.Name()
				}
				{
					file, err := os.CreateTemp(os.TempDir(), "apisix-cli-local-storage-*.lua")
					assert.Nil(t, err, "create temp file path for /usr/local/apisix/apisix/cli/local_storage.lua")
					_targetApisixCliLocalStoragePath = file.Name()
				}

				cloud.DefaultClient = api
			},
		},
		{
			name:       "test deploy bare metal command with apisix id",
			args:       []string{"bare", "--apisix-version", "2.15.0", "--apisix-id", "1234-abcd"},
			cmdPattern: "/usr/bin/bash -C .*/scripts/install\\.sh",
			installScript: `#!/usr/bin/env bash

set -e

version="2\.11\.0"
instance_id="1234-abcd"
apisix_home="/usr/local/apisix"

installed_version=\$\(apisix version 2>/dev/null\) || true
if \[\[ -z \$\{installed_version\} \]\]; then
  yum install -y https://repos.apiseven.com/packages/centos/apache-apisix-repo-1\.0-1\.noarch\.rpm
  yum install -y apisix-\$version
fi

# copy certs to apisix directory to avoid permission issue
cp -rf .*/\.api7cloud/tls /usr/local/apisix/conf/ssl

if \[\[ -n \$\{instance_id\} \]\]; then
  echo "\$\{instance_id\}" > \$\{apisix_home\}/conf/apisix\.uid
fi

apisix start -c .*/apisix-config-\d+\.yaml
status=\$\?

# wait for APISIX started and generated instance id
sleep 1

# get the APISIX instance id when instance id is not set
if \[\[ -z \$\{instance_id\} \]\]; then
  instance_id="\$\(cat \$\{apisix_home\}/conf/apisix\.uid\)"
fi

if \[\[ \$status -eq 0 \]\]; then
  echo "Your APISIX Instance was deployed successfully!"
  echo "Instance ID: \$\{instance_id\}"
fi
`,
			mockCloud: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				api := cloud.NewMockAPI(ctrl)
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 12345,
					ClusterSpec: sdk.ClusterSpec{
						OrganizationID: 1,
					},
				}, nil)
				api.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				api.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)
				api.EXPECT().GetStartupConfig(sdk.ID(12345), cloud.APISIX).Return(_apisixStartupConfigTpl, nil)

				{

					file, err := os.CreateTemp(os.TempDir(), "apisix-cli-etcd-*.lua")
					assert.Nil(t, err, "create temp file path for /usr/local/apisix/apisix/cli/etcd.lua")
					_targetApisixCliEtcdPath = file.Name()
				}
				{
					file, err := os.CreateTemp(os.TempDir(), "apisix-cli-local-storage-*.lua")
					assert.Nil(t, err, "create temp file path for /usr/local/apisix/apisix/cli/local_storage.lua")
					_targetApisixCliLocalStoragePath = file.Name()
				}

				cloud.DefaultClient = api
			},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".api7cloud", "tls"))
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

			testutils.PrepareFakeConfiguration(t)
			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")

			output, err := cmd.CombinedOutput()
			assert.NoError(t, err, "check if the command executed successfully")

			assert.Regexp(t, tc.cmdPattern, string(output), "check if the composed docker command is correct")

			installFile := filepath.Join(persistence.HomeDir, "scripts/install.sh")
			file, err := os.ReadFile(installFile)
			assert.NoError(t, err, "check if dump the install script successful")
			fmt.Println(string(file))
			assert.Regexp(t, tc.installScript, string(file), "checking")
		})
	}
}
