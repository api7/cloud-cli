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

package deploy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/apisix"
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/api7/cloud-cli/internal/utils"
)

func newDockerCommand() *cobra.Command {
	var (
		ctx deployContext
	)
	cmd := &cobra.Command{
		Use:   "docker [ARGS...]",
		Short: "Deploy Apache APISIX to the Docker container",
		Example: `
cloud-cli deploy docker \
		--name apisix-0 \
		--apisix-image apache/apisix:2.15.0-centos \
		--docker-run-arg --detach \
		--docker-run-arg --mount=type=bind,source=/etc/hosts,target=/etc/hosts,readonly \
		--docker-run-arg --hostname=apisix-1`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := options.Global.Deploy.Docker.Validate(); err != nil {
				output.Errorf(err.Error())
				return
			}

			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
			if err := deployPreRunForDocker(&ctx); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				data []byte
				err  error
			)
			opts := options.Global.Deploy.Docker
			docker := getDockerCommand()
			docker.AppendArgs("run")
			for _, args := range opts.DockerRunArgs {
				docker.AppendArgs(strings.SplitN(args, "=", 2)...)
			}
			docker.AppendArgs("--detach")

			if options.Global.Deploy.APISIXConfigFile != "" {
				data, err = os.ReadFile(options.Global.Deploy.APISIXConfigFile)
				if err != nil {
					output.Errorf("invalid --apisix-config-file option: %s", err)
					return
				}
			}
			mergedConfig, err := apisix.MergeConfig(data, ctx.essentialConfig)
			if err != nil {
				output.Errorf(err.Error())
				return
			}
			if len(mergedConfig) > 0 {
				configFile, err := apisix.SaveConfigToTemp(mergedConfig, "apisix-config-*.yaml")
				if err != nil {
					output.Errorf(err.Error())
					return
				}
				docker.AppendArgs("--mount", "type=bind,source="+configFile+",target=/usr/local/apisix/conf/config.yaml,readonly")
			}
			docker.AppendArgs("--mount", "type=bind,source="+ctx.cloudLuaModuleDir+",target=/cloud_lua_module,readonly")
			docker.AppendArgs("--mount", "type=bind,source="+ctx.tlsDir+",target=/cloud/tls,readonly")
			docker.AppendArgs("--mount", "type=bind,source="+ctx.apisixIDFile+",target=/usr/local/apisix/conf/apisix.uid,readonly")
			// For cloud_lua_module/apisix/cli/*.lua, we have to mount them to the /usr/local/apisix/apisix/cli/ directory. Or
			// they cannot be loaded as the /usr/local/apisix/apisix/cli/apisix.lua puts /usr/local/apisix as the first item
			// in package.path.
			{
				sourcePath := filepath.Join(ctx.cloudLuaModuleDir, "apisix", "cli", "etcd.ljbc")
				docker.AppendArgs("--mount", "type=bind,source="+sourcePath+",target=/usr/local/apisix/apisix/cli/etcd.lua,readonly")
			}
			{
				sourcePath := filepath.Join(ctx.cloudLuaModuleDir, "apisix", "cli", "local_storage.ljbc")
				docker.AppendArgs("--mount", "type=bind,source="+sourcePath+",target=/usr/local/apisix/apisix/cli/local_storage.lua,readonly")
			}

			if options.Global.Deploy.Docker.LocalCacheBindPath != "" {
				docker.AppendArgs("--mount", "type=bind,source="+options.Global.Deploy.Docker.LocalCacheBindPath+",target=/usr/local/apisix/conf/apisix.data")
			}

			// TODO support customization of the HTTP and HTTPS ports.
			docker.AppendArgs("-p", fmt.Sprintf("%d:9080", options.Global.Deploy.Docker.HTTPHostPort))
			docker.AppendArgs("-p", fmt.Sprintf("%d:9443", options.Global.Deploy.Docker.HTTPSHostPort))

			if options.Global.Deploy.Name != "" {
				docker.AppendArgs("--name", options.Global.Deploy.Name)
				docker.AppendArgs("--hostname", options.Global.Deploy.Name)
			} else {
				docker.AppendArgs("--name", consts.DefaultDeploymentName)
				docker.AppendArgs("--hostname", consts.DefaultDeploymentName)
			}

			docker.AppendArgs(opts.APISIXImage)

			if options.Global.DryRun {
				output.Infof("Running:\n%s\n", docker.String())
			} else {
				output.Verbosef("Running:\n%s\n", docker.String())
			}

			newctx, cancel := context.WithCancel(context.TODO())
			go utils.WaitForSignal(func() {
				cancel()
			})

			stdout, stderr, err := docker.Run(newctx)
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
			fmt.Println("Congratulations! Your APISIX instance was deployed successfully")

			docker = getDockerCommand()
			containerID, err := getDockerContainerIDByName(newctx, docker, options.Global.Deploy.Name)
			if err != nil {
				message := fmt.Sprintf("failed to get APISIX container ID: %s\nPlease check it via docker ps command\n", err)
				output.Errorf(message)
				return
			}
			fmt.Printf("Container ID: %s\n", containerID)
			fmt.Printf("APISIX ID: %s\n", ctx.apisixID)
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Docker.APISIXImage, "apisix-image", "apache/apisix:2.15.0-centos", "Specify the Apache APISIX image")
	cmd.PersistentFlags().IntVar(&options.Global.Deploy.Docker.HTTPHostPort, "http-host-port", 9080, "Specify the host port for HTTP")
	cmd.PersistentFlags().IntVar(&options.Global.Deploy.Docker.HTTPSHostPort, "https-host-port", 9443, "Specify the host port for HTTPS")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Docker.DockerCLIPath, "docker-cli-path", "", "Specify the filepath of the docker command")
	cmd.PersistentFlags().StringSliceVar(&options.Global.Deploy.Docker.DockerRunArgs, "docker-run-arg", []string{}, "Specify the arguments (in the format of name=value, e.g. --mount=type=bind,source=/etc/hosts,target=/etc/hosts,readonly) for the docker run command")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Docker.LocalCacheBindPath, "local-cache-bind-path", "", "Specify the path to bind to the local cache directory in the container")

	return cmd
}

func getDockerCommand() commands.Cmd {
	opts := options.Global.Deploy.Docker
	if opts.DockerCLIPath != "" {
		return commands.New(opts.DockerCLIPath, options.Global.DryRun)
	}
	return commands.New("docker", options.Global.DryRun)
}
