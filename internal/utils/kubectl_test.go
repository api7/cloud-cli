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
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "", errors.New("mock error"))
				test.kubectl = mockCmd
			},
		},
		{
			name:        "container not found",
			errorReason: "mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "container not found", errors.New("mock error")).Times(getAPISIXIDRetry)
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
				mockCmd.EXPECT().Run(gomock.Any()).Return("aaa-id", "", nil)
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
				assert.Equal(t, tc.wantRsp, rsp, "check response")
			}
		})
	}
}
