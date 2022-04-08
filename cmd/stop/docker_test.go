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

package stop

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
)

func TestNewStopDockerCommand(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name       string
		args       []string
		cmdPattern string
	}{
		{
			name:       "default deployment name",
			args:       []string{"docker"},
			cmdPattern: "docker stop apisix",
		},
		{
			name:       "test deploy docker command with customize deployment name",
			args:       []string{"docker", "--name", "apisix-0"},
			cmdPattern: "docker stop apisix-0",
		},
		{
			name:       "test deploy docker command with customize docker cli path",
			args:       []string{"docker", "--docker-cli-path", "/opt/docker"},
			cmdPattern: "/opt/docker stop apisix",
		},
		{
			name:       "test deploy docker command with rm",
			args:       []string{"docker", "--name", "apisix-0", "--rm"},
			cmdPattern: "docker rm -f apisix-0",
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				options.Global.DryRun = true
				options.Global.Verbose = true
				cmd := NewStopCommand()
				cmd.SetArgs(tc.args)
				err := cmd.Execute()
				assert.NoError(t, err, "check if the command executed successfully")
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1", fmt.Sprintf("%s=test-token", consts.Api7CloudAccessTokenEnv))

			output, err := cmd.CombinedOutput()
			assert.NoError(t, err, "check if the command executed successfully")
			assert.Contains(t, string(output), tc.cmdPattern, "check if the composed docker command is correct")
		})
	}
}
