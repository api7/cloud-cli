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

package deploy

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
)

func newDockerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker [ARGS...]",
		Short: "Deploy Apache APISIX to the Docker container",
		Example: `
cloud-cli deploy docker \
		--apisix-image apisix/apisix:2.11.0 \
		--docker-run-arg --detach \
		--docker-run-arg --hostname=apisix-1`,
		Run: func(cmd *cobra.Command, args []string) {
			var docker *commands.Cmd
			opts := options.Global.Deploy.Docker
			if opts.DockerCLIPath != "" {
				docker = commands.New(opts.DockerCLIPath, options.Global.DryRun)
			} else {
				docker = commands.New("docker", options.Global.DryRun)
			}
			docker.AppendArgs("run", opts.APISIXImage)
			for _, args := range opts.DockerRunArgs {
				docker.AppendArgs(strings.Split(args, "=")...)
			}

			if options.Global.DryRun {
				output.Infof("Running:\n%s\n", docker.String())
			}

			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()

			stdout, stderr, err := docker.Run(ctx)
			if stderr != "" {
				output.Warnf(stderr)
			}
			if stdout != "" {
				output.Verbosef(stdout)
			}
			if err != nil {
				output.Errorf(err.Error())
			}
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Docker.APISIXImage, "apisix-image", "apache/apisix:2.11.0", "Specify the Apache APISIX image")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Docker.DockerCLIPath, "docker-cli-path", "", "Specify the filepath of the docker command")
	cmd.PersistentFlags().StringSliceVar(&options.Global.Deploy.Docker.DockerRunArgs, "docker-run-arg", []string{}, "Specify the arguments (in the format of name=value) for the docker run command")

	return cmd
}
