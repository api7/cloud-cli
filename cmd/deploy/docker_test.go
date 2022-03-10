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
	"archive/tar"
	"bytes"
	"compress/gzip"
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
			args:       []string{"docker", "--apisix-image", "apache/apisix:2.11.0-centos", "--docker-run-arg", "--detach"},
			cmdPattern: `docker run --detach --mount type=bind,source=,target=/cloud_lua_module,readonly -p 9080:9080 -p 9443:9443 apache/apisix:2.11.0-centos`,
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
				buffer := bytes.NewBuffer(nil)
				gzipWriter, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
				assert.NoError(t, err, "create gzip writer")
				tarWriter := tar.NewWriter(gzipWriter)
				body := "hello world"
				hdr := &tar.Header{
					Name: "foo.txt",
					Size: int64(len(body)),
				}
				err = tarWriter.WriteHeader(hdr)
				assert.NoError(t, err, "write tar header")
				_, err = tarWriter.Write([]byte(body))
				assert.NoError(t, err, "write tar body")
				err = tarWriter.Flush()
				assert.NoError(t, err, "flush tar body")
				err = tarWriter.Close()
				assert.NoError(t, err, "close tar writer")
				err = gzipWriter.Close()
				assert.NoError(t, err, "close gzip writer")

				api.EXPECT().GetCloudLuaModule().Return(buffer.Bytes(), nil)
				cloud.DefaultClient = api
			},
		},
		{
			name:       "test deploy docker command with apisix config",
			args:       []string{"docker", "--apisix-image", "apache/apisix:2.11.0-centos", "--docker-run-arg", "--detach", "--apisix-config", "./testdata/apisix.yaml"},
			cmdPattern: `docker run --detach --mount type=bind,source=/.+?/apisix-config-\d+.yaml,target=/usr/local/apisix/conf/config.yaml,readonly --mount type=bind,source=,target=/cloud_lua_module,readonly -p 9080:9080 -p 9443:9443 apache/apisix:2.11.0-centos`,
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

				buffer := bytes.NewBuffer(nil)
				gzipWriter, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
				assert.NoError(t, err, "create gzip writer")
				tarWriter := tar.NewWriter(gzipWriter)
				body := "hello world"
				hdr := &tar.Header{
					Name: "foo.txt",
					Size: int64(len(body)),
				}
				err = tarWriter.WriteHeader(hdr)
				assert.NoError(t, err, "write tar header")
				_, err = tarWriter.Write([]byte(body))
				assert.NoError(t, err, "write tar body")
				err = tarWriter.Flush()
				assert.NoError(t, err, "flush tar body")
				err = tarWriter.Close()
				assert.NoError(t, err, "close tar writer")
				err = gzipWriter.Close()
				assert.NoError(t, err, "close gzip writer")

				api.EXPECT().GetCloudLuaModule().Return(buffer.Bytes(), nil)
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
