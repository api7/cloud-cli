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

package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/cmd/config"
	"github.com/api7/cloud-cli/cmd/configure"
	"github.com/api7/cloud-cli/cmd/debug"
	"github.com/api7/cloud-cli/cmd/deploy"
	"github.com/api7/cloud-cli/cmd/stop"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/utils"
	"github.com/api7/cloud-cli/internal/version"
)

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cloud-cli [OPTIONS] COMMANDS",
		Short:   "Universal command line interface for API7 Cloud",
		Version: version.V.String(),
	}
	cmd.PersistentFlags().BoolVar(&options.Global.Verbose, "verbose", false, "Enable verbose output")
	cmd.PersistentFlags().BoolVar(&options.Global.DryRun, "dry-run", false, "Enable dry run mode")

	cmd.AddCommand(deploy.NewCommand())
	cmd.AddCommand(configure.NewCommand())
	cmd.AddCommand(stop.NewStopCommand())
	cmd.AddCommand(debug.NewCommand())
	cmd.AddCommand(config.NewCommand())

	return cmd
}

func main() {
	defer func() {
		if !options.Global.Verbose {
			return
		}
		// waiting for trace logger done
		utils.WantExit = true
		utils.VerboseWg.Wait()
	}()

	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
