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

package commands

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	cmd := New("echo", false)
	cmd.AppendArgs("hello world")

	assert.Equal(t, "echo hello world", cmd.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stdout, stderr, err := cmd.Run(ctx)
	assert.Nil(t, err, "check cmd run error")
	assert.Equal(t, "hello world\n", stdout, "check cmd stdout")
	assert.Empty(t, stderr, "check cmd stderr")
}

func TestCmdRunTimeout(t *testing.T) {
	cmd := New("sleep", false)
	cmd.AppendArgs("5")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stdout, stderr, err := cmd.Run(ctx)
	assert.NotNil(t, err, "check err is not nil")
	assert.Equal(t, "sleep: signal: killed", err.Error(), "check cmd run error")
	assert.Empty(t, stdout, "check cmd stdout")
	assert.Empty(t, stderr, "check cmd stderr")
}

func TestCmdDryRun(t *testing.T) {
	cmd := New("sleep", true)
	cmd.AppendArgs("5")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stdout, stderr, err := cmd.Run(ctx)
	assert.Nil(t, err, "check err is nil")
	assert.Empty(t, stdout, "check cmd stdout")
	assert.Empty(t, stderr, "check cmd stderr")
}
