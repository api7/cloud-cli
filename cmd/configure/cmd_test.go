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

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/api7/cloud-cli/internal/types"
)

func TestConfigureCommand(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		successExpected bool
		outputExpected  []string
		mockFn          func(api *cloud.MockAPI)
	}{
		{
			name:            "invalid token",
			token:           "invalid token",
			successExpected: false,
			outputExpected:  []string{"invalid access token"},
		},
		{
			name: "expired token",
			// expire at 2000-01-01T00:00:00Z
			token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0NjY4NDgwMCwianRpIjoidXNlcl9pZCJ9.bc90jznuC_WjbKMPUl2Sf-MKs3xifG1GG6pG3JQABNls3aVCLd7rqQDIs4yywLoKE80jDYl4pLtIDfXnPb-YTLTuy5xJdP4lYDYWCO7M91ECtW4PzfN4noTM6IkPlwJAixjtcRIeN6MU6CidOjvkeeoHKgdDF7cOVxlgksxrlFTTcj76KwuR-d9TzDe0z21tB7Qx21lXDBx5gPXlr1P8h7M0A_6mqs16cGQQQqfsetVPModaVVH8H4yQG8Sbt-MGdj4MYwQNqQYjK3ezf041I5KTYZDxuId0_IVZliNNvZA0FJw-06yiRVW-knw6M23wZzlkBLeZoqVal-vbRJx9pg",
			successExpected: false,
			outputExpected:  []string{"access token expired"},
		},
		{
			name: "token validate fail in backend",
			// expire at 2100-01-01T00:00:00Z
			token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMDI0NDQ4MDAsImp0aSI6InVzZXJfaWQifQ.cPY9qlkzlya6PLLCmTfr611nzcaETO2vrBATr45QkpmDHPFctv_zmxgkBiWlvJNMHcgCua7vwXgO-uPfdqDFsDJryI3lDj3w3CHhq85ZGUU9HFBOXVX9NKdBw3eDn4WyHVTDfSNpLNrLSr1xBuvBQs0jTYSUHk2RHHeSfViOrcq91EKfEzFXX8lOikXKbHVs0bYHryrjJeCheW_Z5shIimfgMqLZIIA8F7INPpAeCppkicUkStBixiCO0YGRZAdQcmI3QTBttIwd-mnBc8SQqwMfwFc9DCpwvdcdyZ6B8tRwpZuPJM1u8k2XuH17wUfeCLgkaHgczcAsWQm3T5Ldew",
			successExpected: false,
			outputExpected: []string{
				fmt.Sprintf("your access token will expire at %s", time.Unix(4102444800, 0).Format(time.RFC3339)),
				fmt.Sprintf("mock error"),
			},
			mockFn: func(api *cloud.MockAPI) {
				api.EXPECT().Me().Return(nil, errors.New("mock error"))
			},
		},
		{
			name: "success",
			// expire at 2100-01-01T00:00:00Z
			token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQxMDI0NDQ4MDAsImp0aSI6InVzZXJfaWQifQ.cPY9qlkzlya6PLLCmTfr611nzcaETO2vrBATr45QkpmDHPFctv_zmxgkBiWlvJNMHcgCua7vwXgO-uPfdqDFsDJryI3lDj3w3CHhq85ZGUU9HFBOXVX9NKdBw3eDn4WyHVTDfSNpLNrLSr1xBuvBQs0jTYSUHk2RHHeSfViOrcq91EKfEzFXX8lOikXKbHVs0bYHryrjJeCheW_Z5shIimfgMqLZIIA8F7INPpAeCppkicUkStBixiCO0YGRZAdQcmI3QTBttIwd-mnBc8SQqwMfwFc9DCpwvdcdyZ6B8tRwpZuPJM1u8k2XuH17wUfeCLgkaHgczcAsWQm3T5Ldew",
			successExpected: true,
			outputExpected: []string{
				fmt.Sprintf("your access token will expire at %s", time.Unix(4102444800, 0).Format(time.RFC3339)),
				fmt.Sprintf("demo@api7.cloud"),
			},
			mockFn: func(api *cloud.MockAPI) {
				api.EXPECT().Me().Return(&types.User{
					Email: "demo@api7.cloud",
				}, nil)

			},
		},
		{
			name: "success never expire token",
			// never expire token
			token:           "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOi0xLCJqdGkiOiJ1c2VyX2lkIn0.Jf6YGTi-SD6mc8y_VZ4dJ4PiV1wfBBvUXAggWWQPr-0ZRc8hHyTLhxrKag8qsbByKBCFWO9jfPYMUrDnIgWzhKMg6s60dScYXGN1eaqmajBJFLKlHCGFSPojbJAVhah3KZjLzJDFRj_xqm0Z-AL7V4eSP0Uz4Ax7Qqu-Ubzpb4WQcLgiALURD_f47eiakMMMrIQ-ZstF2Qw4zKaWiZv-YIUhjiHRCsN2nJ2RONAU5sclqy4AlXqgOYrm_OzkN9uBH0by7QNpK2lrSTrtNBrVOg8SX-vTihEEPP4Ao_x41zwcIp_67_2YQ8uaWc_CJBjwsO6wompIu5lbn-7ghWf5-A",
			successExpected: true,
			outputExpected: []string{
				fmt.Sprintf("You are using a token that has no expiration time, please note the security risk"),
				fmt.Sprintf("demo@api7.cloud"),
			},
			mockFn: func(api *cloud.MockAPI) {
				api.EXPECT().Me().Return(&types.User{
					Email: "demo@api7.cloud",
				}, nil)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Because `os.Exit(-1)` will be triggered in the failure case, so here the test is executed using a subprocess
			//The method come from: https://talks.golang.org/2014/testing.slide#23
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				if tt.mockFn != nil {
					api := cloud.NewMockAPI(gomock.NewController(t))
					tt.mockFn(api)
					cloud.DefaultClient = api
				}
				cmd := NewCommand()
				err := cmd.Execute()
				assert.NoError(t, err)
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")
			cmd.Stdin = strings.NewReader(tt.token + "\n")

			output, err := cmd.CombinedOutput()

			if tt.successExpected {
				assert.NoError(t, err, "checking configure command execution successful")
				for _, s := range tt.outputExpected {
					assert.Contains(t, string(output), s, "checking output")
				}

				cfg, err := persistence.LoadCredential()
				assert.NoError(t, err, "checking load config error")
				assert.Equal(t, tt.token, cfg.User.AccessToken, "checking token")
			} else {
				assert.Error(t, err, "checking configure command execution failed")
				for _, s := range tt.outputExpected {
					assert.Contains(t, string(output), s, "checking output")
				}
			}
		})
	}
}
