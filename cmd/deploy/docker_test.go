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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/options"
)

func TestDockerDeployCommand(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		configFile string
		cmdPattern string
	}{
		{
			name:       "test deploy docker command",
			args:       []string{"docker", "--apisix-image", "apache/apisix:2.11.0-centos", "--docker-run-arg", "--detach"},
			cmdPattern: `docker run --detach -p 9080:9080 -p 9443:9443 apache/apisix:2.11.0-centos`,
		},
		{
			name:       "test deploy docker command with apisix config",
			args:       []string{"docker", "--apisix-image", "apache/apisix:2.11.0-centos", "--docker-run-arg", "--detach", "--apisix-config", "./testdata/apisix.yaml"},
			cmdPattern: `docker run --detach --mount type=bind,source=/.+?/apisix-config-\d+.yaml,target=/usr/local/apisix/conf/config.yaml,readonly -p 9080:9080 -p 9443:9443 apache/apisix:2.11.0-centos`,
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			//Because `os.Exit(-1)` will be triggered in the failure case, so here the test is executed using a subprocess
			//The method come from: https://talks.golang.org/2014/testing.slide#23
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				options.Global.DryRun = true
				cmd := NewCommand()
				cmd.SetArgs(tc.args)
				err := cmd.Execute()
				assert.NoError(t, err, "check if the command executed successfully")
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%s", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")

			output, err := cmd.CombinedOutput()
			assert.NoError(t, err, "check if the command executed successfully")

			assert.Regexp(t, tc.cmdPattern, string(output), "check if the composed docker command is correct")
		})
	}
}
