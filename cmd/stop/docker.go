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

package stop

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/utils"
)

func newStopDockerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker [ARG...]",
		Short: "Stop Apache APISIX on Docker",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.TODO())
			go utils.WaitForSignal(func() {
				cancel()
			})

			var (
				docker commands.Cmd
				err    error
			)
			opts := options.Global.Stop.Docker
			if opts.DockerCLIPath != "" {
				docker = commands.New(opts.DockerCLIPath, options.Global.DryRun)
			} else {
				docker = commands.New("docker", options.Global.DryRun)
			}
			if options.Global.Stop.Remove {
				docker.AppendArgs("rm")
				docker.AppendArgs("-f")
			} else {
				docker.AppendArgs("stop")
			}

			if options.Global.Stop.Name != "" {
				docker.AppendArgs(options.Global.Stop.Name)
			}
			if options.Global.DryRun {
				output.Infof(docker.String())
				return
			}
			stdout, stderr, err := docker.Run(ctx)
			if stderr != "" {
				output.Warnf(stderr)
			}
			if stdout != "" {
				output.Verbosef(stdout)
			}
			if err != nil {
				output.Errorf(err.Error())
				return
			}
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Stop.Docker.DockerCLIPath, "docker-cli-path", "", "Specify the filepath of the docker command")

	return cmd
}
