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
	"time"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/apisix"
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

func newBareCommand() *cobra.Command {
	var (
		cloudLuaModuleDir string
	)
	cmd := &cobra.Command{
		Use:   "bare [ARGS...]",
		Short: "Deploy Apache APISIX on bare metal (only CentOS 7) ",
		Example: `
cloud-cli deploy bare \
		--apisix-version 2.11.0`,
		PreRun: func(cmd *cobra.Command, args []string) {
			var (
				err error
			)
			if err := persistence.PrepareCertificate(); err != nil {
				output.Errorf("Failed to prepare certificate: %s", err)
				return
			}
			cloudLuaModuleDir, err = persistence.SaveCloudLuaModule()
			if err != nil {
				output.Errorf("Failed to save cloud lua module: %s", err)
				return
			}
			output.Verbosef("Saved cloud lua module to: %s", cloudLuaModuleDir)
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
			defer cancel()

			opts := options.Global.Deploy.Bare
			path, err := persistence.DownloadRPM(ctx, opts.APISIXVersion)
			if err != nil {
				output.Errorf(err.Error())
				return
			}
			if path != "" {
				installer := commands.New("yum", options.Global.DryRun)
				installer.AppendArgs("install", "-y", path+"/*.rpm")
				if err := installer.Execute(ctx); err != nil {
					return
				}
			}

			var configFile string
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
					configFile, err = apisix.SaveConfig(mergedConfig)
					if err != nil {
						output.Errorf(err.Error())
						return
					}
				}
			}

			bare := commands.New("apisix", options.Global.DryRun)
			bare.AppendArgs("start")

			if configFile != "" {
				bare.AppendArgs("-c", configFile)
			}
			if err = bare.Execute(ctx); err != nil {
				return
			}
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Bare.APISIXVersion, "apisix-version", "2.11.0", "Specifies the APISIX version, default value is 2.11.0")

	return cmd
}
