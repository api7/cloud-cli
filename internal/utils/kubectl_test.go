//  Copyright 2022 API7.ai, Inc under one or more contributor license
//  agreements.  See the NOTICE file distributed with this work for
//  additional information regarding copyright ownership.
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

package utils

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/commands"
)

func TestGetDeploymentName(t *testing.T) {
	type testCase struct {
		name        string
		errorReason string
		mockFn      func(t *testing.T, test *testCase)
		kubectl     commands.Cmd
		wantRsp     string
	}

	testCases := []testCase{
		{
			name:        "execute kubectl command error",
			errorReason: "mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "", errors.New("mock error"))
				test.kubectl = mockCmd
			},
		},
		{
			name: "get deployment name succeed",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("xxx-apisix", "", nil)
				test.kubectl = mockCmd
			},
			wantRsp: "xxx-apisix",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn(t, &tc)
			rsp, err := GetDeploymentName(tc.kubectl)
			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "check error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantRsp, rsp, "check response")
			}
		})
	}
}

func TestGetPodsNames(t *testing.T) {
	type testCase struct {
		name        string
		errorReason string
		mockFn      func(t *testing.T, test *testCase)
		kubectl     commands.Cmd
		wantRsp     []string
	}

	testCases := []testCase{
		{
			name:        "execute kubectl command error",
			errorReason: "mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "", errors.New("mock error"))
				test.kubectl = mockCmd
			},
		},
		{
			name: "get pod name succeed",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("aaa bbb", "", nil)
				test.kubectl = mockCmd
			},
			wantRsp: []string{"aaa", "bbb"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn(t, &tc)
			rsp, err := GetPodsNames(tc.kubectl)
			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "check error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantRsp, rsp, "check response")
			}
		})
	}
}

func TestGetAPISIXID(t *testing.T) {
	type testCase struct {
		name        string
		errorReason string
		mockFn      func(t *testing.T, test *testCase)
		kubectl     commands.Cmd
		wantRsp     string
	}

	testCases := []testCase{
		{
			name:        "execute kubectl command error",
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
			name: "get APISIX ID succeed",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("aaa-id", "", nil).AnyTimes()
				test.kubectl = mockCmd
			},
			wantRsp: "aaa-id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn(t, &tc)
			rsp, err := GetAPISIXID(tc.kubectl, "aaa")
			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "check error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantRsp, rsp, "check response")
			}
		})
	}
}

func TestGetServiceName(t *testing.T) {
	type testCase struct {
		name        string
		errorReason string
		mockFn      func(t *testing.T, test *testCase)
		kubectl     commands.Cmd
		wantRsp     string
	}

	testCases := []testCase{
		{
			name:        "execute kubectl command error",
			errorReason: "mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "", errors.New("mock error"))
				test.kubectl = mockCmd
			},
		},
		{
			name: "get service name succeed",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("aaa-apisix-gateway", "", nil)
				test.kubectl = mockCmd
			},
			wantRsp: "aaa-apisix-gateway",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFn(t, &tc)
			rsp, err := GetServiceName(tc.kubectl)
			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "check error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantRsp, rsp, "check response")
			}
		})
	}
}
