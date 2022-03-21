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
	"bytes"
	"context"
	_ "embed"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/apisix"
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

var (
	//go:embed manifest/install.sh
	_installScript string
	_installer     *template.Template
	_apisixRepoURL = "https://repos.apiseven.com/packages/centos/apache-apisix-repo-1.0-1.noarch.rpm"
)

type installContext struct {
	APISIXRepoURL string
	TLSDir        string
	ConfigFile    string
	Version       string
	InstanceID    string
}

func init() {
	_installer = template.Must(template.New("install script").Parse(_installScript))
}

func newBareCommand() *cobra.Command {
	var (
		ctx deployContext
	)
	cmd := &cobra.Command{
		Use:   "bare [ARGS...]",
		Short: "Deploy Apache APISIX on bare metal (only CentOS 7)",
		Example: `
cloud-cli deploy bare \
		--apisix-version 2.11.0`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
			if err := deployPreRunForBare(&ctx); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			context, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
			defer cancel()

			var (
				err  error
				data []byte
			)
			opts := options.Global.Deploy.Bare

			if options.Global.Deploy.APISIXConfigFile != "" {
				data, err = ioutil.ReadFile(options.Global.Deploy.APISIXConfigFile)
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

			var configFile string
			if len(mergedConfig) > 0 {
				configFile, err = apisix.SaveConfig(mergedConfig, "apisix-config-*.yaml")
				if err != nil {
					output.Errorf(err.Error())
					return
				}
			}

			buf := bytes.NewBuffer(nil)
			err = _installer.Execute(buf, &installContext{
				APISIXRepoURL: _apisixRepoURL,
				TLSDir:        ctx.tlsDir,
				ConfigFile:    configFile,
				Version:       opts.APISIXVersion,
				InstanceID:    options.Global.Deploy.APISIXInstanceID,
			})
			if err != nil {
				output.Errorf(err.Error())
				return
			}
			installerPath := filepath.Join(persistence.HomeDir, "scripts")
			err = os.Mkdir(installerPath, 0755)
			if err != nil {
				if !os.IsExist(err) {
					output.Errorf(err.Error())
					return
				}
			}

			installerFile := filepath.Join(installerPath, "install.sh")
			err = os.WriteFile(installerFile, buf.Bytes(), 0755)
			if err != nil {
				output.Errorf(err.Error())
				return
			}

			bare := commands.New("/usr/bin/bash", options.Global.DryRun)
			bare.AppendArgs("-C")
			bare.AppendArgs(installerFile)

			if err = bare.Execute(context); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Bare.APISIXVersion, "apisix-version", "2.11.0", "Specifies the APISIX version, default value is 2.11.0")

	return cmd
}
