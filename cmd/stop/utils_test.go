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

package stop

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/commands"
)

func TestStopPreRunForKubernetes(t *testing.T) {
	type testCase struct {
		name        string
		errorReason string
		mockFn      func(t *testing.T, test *testCase)
		kubectl     commands.Cmd
	}

	testCases := []testCase{
		{
			name:        "failed to delete configmap and secret on kubernetes",
			errorReason: "mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "", errors.New("mock error")).AnyTimes()
				test.kubectl = mockCmd
			},
		},
		{
			name:    "delete configmap and secret on kubernetes should succeed",
			kubectl: commands.New("kubectl", true),
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "", nil).AnyTimes()
				test.kubectl = mockCmd
			},
		},
		{
			name: "when configmap and secret not exist on kubernetes delete should succeed",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "NotFound", errors.New("mock error")).AnyTimes()
				test.kubectl = mockCmd
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn(t, &tc)

			err := stopPreRunForKubernetes(tc.kubectl)
			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "check error")
			} else {
				assert.NoError(t, err, "check error")
			}
		})
	}
}
