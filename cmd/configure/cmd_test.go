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

package configure

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/config"
)

func TestConfigureCommand(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		successExpected bool
		outputExpected  string
	}{
		{
			name:            "invalid token",
			token:           "invalid token",
			successExpected: false,
			outputExpected:  "invalid access token",
		},
		{
			name: "expired token",
			// expire at 2000-01-01T00:00:00Z
			token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0NjY4NDgwMCwianRpIjoidXNlcl9pZCJ9.bc90jznuC_WjbKMPUl2Sf-MKs3xifG1GG6pG3JQABNls3aVCLd7rqQDIs4yywLoKE80jDYl4pLtIDfXnPb-YTLTuy5xJdP4lYDYWCO7M91ECtW4PzfN4noTM6IkPlwJAixjtcRIeN6MU6CidOjvkeeoHKgdDF7cOVxlgksxrlFTTcj76KwuR-d9TzDe0z21tB7Qx21lXDBx5gPXlr1P8h7M0A_6mqs16cGQQQqfsetVPModaVVH8H4yQG8Sbt-MGdj4MYwQNqQYjK3ezf041I5KTYZDxuId0_IVZliNNvZA0FJw-06yiRVW-knw6M23wZzlkBLeZoqVal-vbRJx9pg",
			successExpected: false,
			outputExpected:  "access token expired",
		},
		{
			name: "success",
			// expire at 2100-01-01T00:00:00Z
			token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMDI0NDQ4MDAsImp0aSI6InVzZXJfaWQifQ.cPY9qlkzlya6PLLCmTfr611nzcaETO2vrBATr45QkpmDHPFctv_zmxgkBiWlvJNMHcgCua7vwXgO-uPfdqDFsDJryI3lDj3w3CHhq85ZGUU9HFBOXVX9NKdBw3eDn4WyHVTDfSNpLNrLSr1xBuvBQs0jTYSUHk2RHHeSfViOrcq91EKfEzFXX8lOikXKbHVs0bYHryrjJeCheW_Z5shIimfgMqLZIIA8F7INPpAeCppkicUkStBixiCO0YGRZAdQcmI3QTBttIwd-mnBc8SQqwMfwFc9DCpwvdcdyZ6B8tRwpZuPJM1u8k2XuH17wUfeCLgkaHgczcAsWQm3T5Ldew",
			successExpected: true,
			outputExpected:  fmt.Sprintf("your access token will expire at %s", time.Unix(4102444800, 0).Format(time.RFC3339)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Because `os.Exit(-1)` will be triggered in the failure case, so here the test is executed using a subprocess
			//The method come from: https://talks.golang.org/2014/testing.slide#23
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				cmd := NewCommand()
				err := cmd.Execute()
				assert.NoError(t, err)
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%s", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")
			cmd.Stdin = strings.NewReader(tt.token + "\n")

			output, err := cmd.CombinedOutput()

			if tt.successExpected {
				assert.NoError(t, err, "checking configure command execution successful")
				if tt.outputExpected != "" {
					assert.Contains(t, string(output), tt.outputExpected)
				}

				cfg, err := config.Load()
				assert.NoError(t, err, "checking load config error")
				assert.Equal(t, tt.token, cfg.User.AccessToken, "checking token")
			} else {
				assert.Error(t, err, "checking configure command execution failed")
				assert.Contains(t, string(output), tt.outputExpected)
			}
		})
	}
}
