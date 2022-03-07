// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
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
package main

import (
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/cmd/deploy"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/version"
)

var (
	_globalOptions options.GlobalOptions
)

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cloud-cli [OPTIONS] COMMANDS",
		Short:   "Universal command line interface for API7 Cloud",
		Version: version.V.String(),
	}
	cmd.PersistentFlags().BoolVar(&_globalOptions.Verbose, "verbose", false, "Enable verbose output")
	cmd.PersistentFlags().BoolVar(&_globalOptions.DryRun, "dry-run", false, "Enable dry run mode")

	cmd.AddCommand(deploy.NewCommand(&_globalOptions.Deploy))

	return cmd
}

func main() {
	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
