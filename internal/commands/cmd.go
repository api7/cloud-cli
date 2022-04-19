// Copyright 2022 API7.ai, Inc under one or more contributor license
// agreements.  See the NOTICE file distributed with this work for
// additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
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
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/output"
)

// AppendArgs appends a couple of args to the Cmd object.
func (c *cmd) AppendArgs(args ...string) {
	c.args = append(c.args, args...)
}

// String prints the command.
func (c *cmd) String() string {
	return c.name + " " + strings.Join(c.args, " ")
}

// Run launches the command and return the stdout, stderr.
func (c *cmd) Run(ctx context.Context) (string, string, error) {
	defer func() {
		c.args = nil
		c.stdout = bytes.NewBuffer(nil)
		c.stderr = bytes.NewBuffer(nil)
	}()
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

func (c *cmd) Execute(ctx context.Context) error {
	if c.dryrun {
		output.Infof(c.String())
		return nil
	}

	stdout, stderr, err := c.Run(ctx)
	if stderr != "" {
		output.Warnf(stderr)
	}
	if stdout != "" {
		output.Verbosef(stdout)
	}
	if err != nil {
		output.Errorf(err.Error())
	}
	return err
}
