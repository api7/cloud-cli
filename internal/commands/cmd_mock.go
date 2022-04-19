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

// Code generated by MockGen. DO NOT EDIT.
// Source: ./types.go

// Package commands is a generated GoMock package.
package commands

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCmd is a mock of Cmd interface.
type MockCmd struct {
	ctrl     *gomock.Controller
	recorder *MockCmdMockRecorder
}

// MockCmdMockRecorder is the mock recorder for MockCmd.
type MockCmdMockRecorder struct {
	mock *MockCmd
}

// NewMockCmd creates a new mock instance.
func NewMockCmd(ctrl *gomock.Controller) *MockCmd {
	mock := &MockCmd{ctrl: ctrl}
	mock.recorder = &MockCmdMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCmd) EXPECT() *MockCmdMockRecorder {
	return m.recorder
}

// AppendArgs mocks base method.
func (m *MockCmd) AppendArgs(args ...string) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "AppendArgs", varargs...)
}

// AppendArgs indicates an expected call of AppendArgs.
func (mr *MockCmdMockRecorder) AppendArgs(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendArgs", reflect.TypeOf((*MockCmd)(nil).AppendArgs), args...)
}

// Execute mocks base method.
func (m *MockCmd) Execute(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute.
func (mr *MockCmdMockRecorder) Execute(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCmd)(nil).Execute), ctx)
}

// Run mocks base method.
func (m *MockCmd) Run(ctx context.Context) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Run indicates an expected call of Run.
func (mr *MockCmdMockRecorder) Run(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockCmd)(nil).Run), ctx)
}

// String mocks base method.
func (m *MockCmd) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockCmdMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockCmd)(nil).String))
}
