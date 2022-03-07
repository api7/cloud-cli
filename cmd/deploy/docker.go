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
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/options"
)

func newDockerCommand(opts *options.DockerDeployOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker [ARGS...]",
		Short: "Deploy Apache APISIX to the Docker container",
		Example: `
cloud-cli deploy docker \
		--apisix-image apisix/apisix:2.11.0 \
		--docker-run-arg --detach \
		--docker-run-arg --hostname=apisix-1`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	cmd.PersistentFlags().StringVar(&opts.APISIXImage, "apisix-image", "apache/apisix:2.11.0", "Specify the Apache APISIX image")
	cmd.PersistentFlags().StringVar(&opts.DockerCLIPath, "docker-cli-path", "", "Specify the filepath of the docker command")
	cmd.PersistentFlags().StringSliceVar(&opts.DockerRunArgs, "docker-run-arg", []string{}, "Specify the arguments (in the format of name=value) for the docker run command")

	return cmd
}
