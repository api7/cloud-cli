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

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/persistence"
)

func TestPersistentPreRunFunc(t *testing.T) {

	tests := []struct {
		name            string
		token           string
		env             string
		successExpected bool
		outputExpected  string
	}{
		{
			name:            "no credential",
			successExpected: false,
			outputExpected:  "Failed to load credential",
		},
		{
			name:            "credential in file",
			token:           "token-in-file",
			successExpected: true,
			outputExpected:  "token-in-file",
		},
		{
			name:            "credential in environment",
			env:             fmt.Sprintf("%s=token-in-environment", consts.Api7CloudAccessTokenEnv),
			successExpected: true,
			outputExpected:  "token-in-environment",
		},
		{
			name:            "environment have higher priority than file",
			token:           "token-in-file",
			env:             fmt.Sprintf("%s=token-in-environment", consts.Api7CloudAccessTokenEnv),
			successExpected: true,
			outputExpected:  "token-in-environment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Because `os.Exit(-1)` will be triggered in the failure case, so here the test is executed using a subprocess
			//The method come from: https://talks.golang.org/2014/testing.slide#23
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				options.Global.DryRun = true
				options.Global.Verbose = true
				cmd := NewCommand()
				cmd.SetArgs([]string{"docker", "--apisix-image", "apache/apisix:2.11.0-centos", "--docker-run-arg", "--detach"})
				err := cmd.Execute()
				assert.NoError(t, err, "check if the command executed successfully")
				return
			}

			if tt.token != "" {
				defer func() {
					err := os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".api7cloud/credentials"))
					assert.NoError(t, err, "remove credential file")
				}()
				err := persistence.SaveCredential(&persistence.Credential{User: persistence.User{AccessToken: tt.token}})
				assert.NoError(t, err, "save credential")
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%s", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1", tt.env)

			output, err := cmd.CombinedOutput()

			if tt.successExpected {
				assert.NoError(t, err, "checking configure command execution successful")
				if tt.outputExpected != "" {
					assert.Contains(t, string(output), tt.outputExpected)
				}
			} else {
				assert.Error(t, err, "checking configure command execution failed")
				assert.Contains(t, string(output), tt.outputExpected)
			}
		})
	}
}
