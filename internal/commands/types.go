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

package commands

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Cmd wraps the os/exec.Cmd object.
type Cmd struct {
	name   string
	args   []string
	stdout *bytes.Buffer
	stderr *bytes.Buffer

	dryrun bool
}

// New creates a Cmd object.
func New(name string, dryrun bool) *Cmd {
	if dryrun {
		return &Cmd{
			name:   name,
			dryrun: true,
		}
	}
	return &Cmd{
		name:   name,
		stdout: bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
	}
}

// AppendArgs appends a couple of args to the Cmd object.
func (c *Cmd) AppendArgs(args ...string) *Cmd {
	c.args = append(c.args, args...)
	return c
}

// String prints the command.
func (c *Cmd) String() string {
	return c.name + " " + strings.Join(c.args, " ")
}

// Run launches the command and return the stdout, stderr.
func (c *Cmd) Run(ctx context.Context) (string, string, error) {
	if c.dryrun {
		return "", "", nil
	}
	cmd := exec.CommandContext(ctx, c.name, c.args...)
	cmd.Env = os.Environ()
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	err := cmd.Run()
	if err != nil {
		return c.stdout.String(), c.stderr.String(), errors.Wrap(err, c.name)
	}
	return c.stdout.String(), c.stderr.String(), nil
}
