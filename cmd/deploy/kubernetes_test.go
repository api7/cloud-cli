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
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/api7/cloud-cli/internal/types"
)

func TestKubernetesDeployCommand(t *testing.T) {
	defaultMockCloud := func(t *testing.T) {
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
	}
	testCases := []struct {
		name        string
		args        []string
		configFile  string
		cmdPatterns []string
		mockCloud   func(t *testing.T)
	}{
		{
			name: "deploy on kubernetes with default options",
			args: []string{"kubernetes"},
			cmdPatterns: []string{
				`kubectl create ns apisix`,
				`kubectl create secret generic cloud-ssl --from-file tls.crt=.*?tls.crt --from-file tls.key=.*?tls.key --from-file ca.crt=.*?ca.crt --namespace apisix`,
				`kubectl create configmap cloud-module --from-file cloud.ljbc=.*?cloud.ljbc --from-file cloud-agent.ljbc=.*?agent.ljbc --from-file cloud-metrics.ljbc=.*?metrics.ljbc --from-file cloud-utils.ljbc=.*?utils.ljbc --namespace apisix`,
				`helm repo add apisix https://charts.apiseven.com`,
				`helm repo update`,
				`helm install apisix apisix/apisix --namespace apisix --values .*?.yaml`,
			},
			mockCloud: defaultMockCloud,
		},
		{
			name: "deploy on kubernetes with customize options",
			args: []string{"kubernetes", "--name", "apisix-test", "--namespace", "my-apisix",
				"--helm-install-arg", "--output=table", "--helm-install-arg", "--wait"},
			cmdPatterns: []string{
				`kubectl create ns my-apisix`,
				`kubectl create secret generic cloud-ssl --from-file tls.crt=.*?tls.crt --from-file tls.key=.*?tls.key --from-file ca.crt=.*?ca.crt --namespace my-apisix`,
				`kubectl create configmap cloud-module --from-file cloud.ljbc=.*?cloud.ljbc --from-file cloud-agent.ljbc=.*?agent.ljbc --from-file cloud-metrics.ljbc=.*?metrics.ljbc --from-file cloud-utils.ljbc=.*?utils.ljbc --namespace my-apisix`,
				`helm repo add apisix https://charts.apiseven.com`,
				`helm repo update`,
				`helm install apisix-test apisix/apisix --namespace my-apisix --output table --wait --values .*?.yaml`,
			},
			mockCloud: defaultMockCloud,
		},
		{
			name: "deploy on kubernetes with customize helm install values",
			args: []string{"kubernetes", "--helm-install-arg", "--values=./testdata/apisix_chart_values.yaml"},
			cmdPatterns: []string{
				`kubectl create ns apisix`,
				`kubectl create secret generic cloud-ssl --from-file tls.crt=.*?tls.crt --from-file tls.key=.*?tls.key --from-file ca.crt=.*?ca.crt --namespace apisix`,
				`kubectl create configmap cloud-module --from-file cloud.ljbc=.*?cloud.ljbc --from-file cloud-agent.ljbc=.*?agent.ljbc --from-file cloud-metrics.ljbc=.*?metrics.ljbc --from-file cloud-utils.ljbc=.*?utils.ljbc --namespace apisix`,
				`helm repo add apisix https://charts.apiseven.com`,
				`helm repo update`,
				`helm install apisix apisix/apisix --namespace apisix --values .*?.yaml`,
			},
			mockCloud: defaultMockCloud,
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

			for _, pattern := range tc.cmdPatterns {
				assert.Regexp(t, pattern, string(output), "check if the composed docker command is correct")
			}
		})
	}
}
