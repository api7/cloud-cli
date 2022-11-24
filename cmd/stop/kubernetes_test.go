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

package stop

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/testutils"
)

func TestKubernetesStopCommand(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		cmdPatterns []string
	}{
		{
			name: "stop on kubernetes with default options",
			args: []string{"kubernetes"},
			cmdPatterns: []string{
				`kubectl delete configmap cloud-module --namespace apisix`,
				`kubectl delete secret cloud-ssl --namespace apisix`,
				`helm uninstall apisix --namespace apisix`,
			},
		},
		{
			name: "stop on kubernetes with customize options",
			args: []string{"kubernetes", "--name", "apisix-test", "--namespace", "apisix-test",
				"--helm-cli-path", "/tmp/helm", "--kubectl-cli-path", "/tmp/kubectl",
				"--helm-uninstall-arg", "--keep-history", "--helm-uninstall-arg", "--wait"},
			cmdPatterns: []string{
				`/tmp/kubectl delete configmap cloud-module`,
				`/tmp/kubectl delete secret cloud-ssl`,
				`/tmp/helm uninstall apisix-test --namespace apisix-test --keep-history --wait`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(t.Name(), func(t *testing.T) {
			//Because `os.Exit(-1)` will be triggered in the failure case, so here the test is executed using a subprocess
			//The method come from: https://talks.golang.org/2014/testing.slide#23
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				options.Global.DryRun = true
				cmd := NewStopCommand()
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
