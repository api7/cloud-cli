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

func TestKubernetesDeployCommand(t *testing.T) {
	defaultMockCloud := func(t *testing.T) {
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
		api.EXPECT().GetStartupConfig(sdk.ID(12345), cloud.HELM).Return(_helmStartupConfigTpl, nil)

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
				`kubectl create configmap cloud-module --from-file cloud.ljbc=.*?cloud.ljbc --from-file cloud-agent.ljbc=.*?agent.ljbc --from-file cloud-metrics.ljbc=.*?metrics.ljbc --from-file cloud-utils.ljbc=.*?utils.ljbc --from-file cloud-file.ljbc=.*?file.ljbc --from-file apisix-local-storage.ljbc=.*?local_storage.ljbc --from-file apisix-core-config-etcd.ljbc=.*?config_etcd.ljbc --from-file apisix-cli-etcd.ljbc=.*?etcd.ljbc --from-file apisix-cli-local-storage.ljbc=.*?local_storage.ljbc --namespace apisix`,
				`helm repo add apisix https://charts.apiseven.com`,
				`helm repo update`,
				`helm install apisix apisix/apisix --namespace apisix --values .*?.yaml`,
				`The Helm release name is: apisix`,
				`kubectl get deployment -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[0\].metadata.name\}"`,
				`kubectl get pods -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[\*\].metadata.name\}"`,
				`kubectl exec  -n apisix -- cat /usr/local/apisix/conf/apisix.uid`,
				`kubectl wait --for condition=Ready --timeout 60s pod/ -n apisix`,
				`kubectl get service -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[0\].metadata.name\}"`,
			},
			mockCloud: defaultMockCloud,
		},
		{
			name: "deploy on kubernetes with customize options",
			args: []string{"kubernetes", "--name", "apisix-test", "--namespace", "my-apisix",
				"--helm-cli-path", "/tmp/helm", "--kubectl-cli-path", "/tmp/kubectl",
				"--helm-install-arg", "--output=table", "--helm-install-arg", "--wait"},
			cmdPatterns: []string{
				`/tmp/kubectl create ns my-apisix`,
				`/tmp/kubectl create secret generic cloud-ssl --from-file tls.crt=.*?tls.crt --from-file tls.key=.*?tls.key --from-file ca.crt=.*?ca.crt --namespace my-apisix`,
				`/tmp/kubectl create configmap cloud-module --from-file cloud.ljbc=.*?cloud.ljbc --from-file cloud-agent.ljbc=.*?agent.ljbc --from-file cloud-metrics.ljbc=.*?metrics.ljbc --from-file cloud-utils.ljbc=.*?utils.ljbc --from-file cloud-file.ljbc=.*?file.ljbc --from-file apisix-local-storage.ljbc=.*?local_storage.ljbc --from-file apisix-core-config-etcd.ljbc=.*?config_etcd.ljbc --from-file apisix-cli-etcd.ljbc=.*?etcd.ljbc --from-file apisix-cli-local-storage.ljbc=.*?local_storage.ljbc --namespace my-apisix`,
				`/tmp/helm repo add apisix https://charts.apiseven.com`,
				`/tmp/helm repo update`,
				`/tmp/helm install apisix-test apisix/apisix --namespace my-apisix --output table --wait --values .*?.yaml`,
				`Congratulations! Your APISIX cluster was deployed successfully on Kubernetes.`,
				`The Helm release name is: apisix-test`,
				`/tmp/kubectl get deployment -n my-apisix -l app.kubernetes.io/instance=apisix-test -o jsonpath="\{.items\[0\].metadata.name\}"`,
				`/tmp/kubectl get pods -n my-apisix -l app.kubernetes.io/instance=apisix-test -o jsonpath="\{.items\[\*\].metadata.name\}"`,
				`/tmp/kubectl exec  -n my-apisix -- cat /usr/local/apisix/conf/apisix.uid`,
				`/tmp/kubectl wait --for condition=Ready --timeout 60s pod/ -n my-apisix`,
				`/tmp/kubectl get service -n my-apisix -l app.kubernetes.io/instance=apisix-test -o jsonpath="\{.items\[0\].metadata.name\}"`,
			},
			mockCloud: defaultMockCloud,
		},
		{
			name: "deploy on kubernetes with customize helm install values",
			args: []string{"kubernetes", "--helm-install-arg", "--values=./testdata/apisix_chart_values.yaml"},
			cmdPatterns: []string{
				`kubectl create ns apisix`,
				`kubectl create secret generic cloud-ssl --from-file tls.crt=.*?tls.crt --from-file tls.key=.*?tls.key --from-file ca.crt=.*?ca.crt --namespace apisix`,
				`kubectl create configmap cloud-module --from-file cloud.ljbc=.*?cloud.ljbc --from-file cloud-agent.ljbc=.*?agent.ljbc --from-file cloud-metrics.ljbc=.*?metrics.ljbc --from-file cloud-utils.ljbc=.*?utils.ljbc --from-file cloud-file.ljbc=.*?file.ljbc --from-file apisix-local-storage.ljbc=.*?local_storage.ljbc --from-file apisix-core-config-etcd.ljbc=.*?config_etcd.ljbc --from-file apisix-cli-etcd.ljbc=.*?etcd.ljbc --from-file apisix-cli-local-storage.ljbc=.*?local_storage.ljbc --namespace apisix`,
				`helm repo add apisix https://charts.apiseven.com`,
				`helm repo update`,
				`helm install apisix apisix/apisix --namespace apisix --values .*?.yaml`,
				`Congratulations! Your APISIX cluster was deployed successfully on Kubernetes.`,
				`The Helm release name is: apisix`,
				`kubectl get deployment -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[0\].metadata.name\}"`,
				`kubectl get pods -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[\*\].metadata.name\}"`,
				`kubectl exec  -n apisix -- cat /usr/local/apisix/conf/apisix.uid`,
				`kubectl wait --for condition=Ready --timeout 60s pod/ -n apisix`,
				`kubectl get service -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[0\].metadata.name\}"`,
			},
			mockCloud: defaultMockCloud,
		},
		{
			name: "deploy on kubernetes with customize helm install --set options",
			args: []string{"kubernetes", "--helm-install-arg", "--set=apisix.ingress.enabled=false"},
			cmdPatterns: []string{
				`kubectl create ns apisix`,
				`kubectl create secret generic cloud-ssl --from-file tls.crt=.*?tls.crt --from-file tls.key=.*?tls.key --from-file ca.crt=.*?ca.crt --namespace apisix`,
				`kubectl create configmap cloud-module --from-file cloud.ljbc=.*?cloud.ljbc --from-file cloud-agent.ljbc=.*?agent.ljbc --from-file cloud-metrics.ljbc=.*?metrics.ljbc --from-file cloud-utils.ljbc=.*?utils.ljbc --from-file cloud-file.ljbc=.*?file.ljbc --from-file apisix-local-storage.ljbc=.*?local_storage.ljbc --from-file apisix-core-config-etcd.ljbc=.*?config_etcd.ljbc --from-file apisix-cli-etcd.ljbc=.*?etcd.ljbc --from-file apisix-cli-local-storage.ljbc=.*?local_storage.ljbc --namespace apisix`,
				`helm repo add apisix https://charts.apiseven.com`,
				`helm repo update`,
				`helm install apisix apisix/apisix --namespace apisix --set apisix.ingress.enabled=false`,
				`Congratulations! Your APISIX cluster was deployed successfully on Kubernetes.`,
				`The Helm release name is: apisix`,
				`kubectl get deployment -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[0\].metadata.name\}"`,
				`kubectl get pods -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[\*\].metadata.name\}"`,
				`kubectl exec  -n apisix -- cat /usr/local/apisix/conf/apisix.uid`,
				`kubectl wait --for condition=Ready --timeout 60s pod/ -n apisix`,
				`kubectl get service -n apisix -l app.kubernetes.io/instance=apisix -o jsonpath="\{.items\[0\].metadata.name\}"`,
			},
			mockCloud: defaultMockCloud,
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			persistence.HomeDir = filepath.Join(os.TempDir(), ".api7cloud")
			defer func() {
				os.RemoveAll(persistence.HomeDir)
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

			for _, pattern := range tc.cmdPatterns {
				assert.Regexp(t, pattern, string(output), "check if the kubectl and helm command is correct")
			}
		})
	}
}
