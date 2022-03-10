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
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/apisix"
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
)

var _rpmPackageFilePath = filepath.Join(os.Getenv("HOME"), ".api7cloud/rpm")

func newBareCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bare [ARGS...]",
		Short: "Deploy Apache APISIX on the Linux(CentOS 7)",
		Example: `
cloud-cli deploy bare \
		--apisix-version 2.11.0`,
		Run: func(cmd *cobra.Command, args []string) {
			opts := options.Global.Deploy.Bare

			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()

			path := filepath.Join(_rpmPackageFilePath, opts.APISIXVersion)
			_ = os.Mkdir(path, 750)
			download(ctx, path, opts.APISIXVersion)

			installer := commands.New("yum", options.Global.DryRun)
			installer.AppendArgs("install", "-y", path+"/*.rpm")
			execute(ctx, installer)

			if options.Global.Deploy.APISIXConfigFile != "" {
				data, err := ioutil.ReadFile(options.Global.Deploy.APISIXConfigFile)
				if err != nil {
					output.Errorf("invalid --apisix-config-file option: %s", err)
					return
				}
				mergedConfig, err := apisix.MergeConfig(data, nil)
				if err != nil {
					output.Errorf(err.Error())
					return
				}
				if len(mergedConfig) > 0 {
					configFile, err := apisix.SaveConfig(mergedConfig)
					if err != nil {
						output.Errorf(err.Error())
						return
					}
					cmd := commands.New("cp", options.Global.DryRun)
					cmd.AppendArgs("-f", configFile, "/usr/local/apisix/conf/config.yaml")
					execute(ctx, cmd)
				}
			}

			bare := commands.New("apisix", options.Global.DryRun)
			bare.AppendArgs("start")
			execute(ctx, bare)
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Bare.APISIXVersion, "apisix-version", "2.11.0", "Specifies the APISIX version, default value is 2.11.0")

	return cmd
}

func download(context context.Context, rpmFilePath, version string) {
	// install the repositories of OpenResty
	cmd := commands.New("yum", options.Global.DryRun)
	cmd.AppendArgs("install", "-y", "https://repos.apiseven.com/packages/centos/apache-apisix-repo-1.0-1.noarch.rpm")
	execute(context, cmd)

	// install the repositories of Apache APISIX.
	cmd = commands.New("yum-config-manager", options.Global.DryRun)
	cmd.AppendArgs("--add-repo", "https://repos.apiseven.com/packages/centos/apache-apisix.repo")
	execute(context, cmd)

	// download apisix rpm
	cmd = commands.New("yum", options.Global.DryRun)
	cmd.AppendArgs("install", "-y", "--downloadonly")
	cmd.AppendArgs("--downloaddir=" + rpmFilePath)
	cmd.AppendArgs("apisix-" + version)
	execute(context, cmd)
}

func execute(context context.Context, cmd *commands.Cmd) {
	if options.Global.DryRun {
		output.Infof(cmd.String())
	}

	stdout, stderr, err := cmd.Run(context)
	if stderr != "" {
		output.Warnf(stderr)
	}
	if stdout != "" {
		output.Verbosef(stdout)
	}
	if err != nil {
		output.Errorf(err.Error())
	}
}
