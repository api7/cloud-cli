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

package commands

import (
	"bytes"
	"context"
)

// Cmd is the command constructor and runner.
type Cmd interface {
	// AppendArgs appends a couple of args to the Cmd.
	AppendArgs(args ...string)
	// String prints the command.
	String() string
	// Run launches the command and return the stdout, stderr.
	Run(ctx context.Context) (string, string, error)
	// Execute launches the command and return error, and print stdout and stderr to console.
	Execute(ctx context.Context) error
}

// Cmd wraps the os/exec.Cmd object.
type cmd struct {
	name   string
	args   []string
	stdout *bytes.Buffer
	stderr *bytes.Buffer

	dryrun bool
}

// New creates a Cmd object.
func New(name string, dryrun bool) Cmd {
	if dryrun {
		return &cmd{
			name:   name,
			dryrun: true,
		}
	}
	return &cmd{
		name:   name,
		stdout: bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
	}
}
