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

package deploy

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/apisix"
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
	Upgrade       bool
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
		--apisix-version 2.15.0`,
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

			var configFile string
			if len(mergedConfig) > 0 {
				configFile = filepath.Join(ctx.apisixConfigDir, "apisix-config-cloud.yaml")
				if err = apisix.SaveConfig(mergedConfig, configFile); err != nil {
					output.Errorf(err.Error())
					return
				}
			}

			if options.Global.Deploy.Bare.Reload {
				if err = apisix.Reload(context, ctx.tlsDir); err != nil {
					output.Errorf(err.Error())
				}
				return
			}

			deployOnBareMetal(context, &ctx, &opts, configFile)
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Bare.APISIXVersion, "apisix-version", "2.15.0", "Specifies the APISIX version, default value is 2.15.0")
	cmd.PersistentFlags().BoolVar(&options.Global.Deploy.Bare.Reload, "reload", false, "Skip deployment, only update configurations and reload APISIX")
	cmd.PersistentFlags().BoolVar(&options.Global.Deploy.Bare.Upgrade, "upgrade", false, "Skip deployment, try to upgrade APISIX version")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Bare.APISIXBinPath, "apisix-bin-path", "/usr/bin/apisix", "APISIX binary file path")

	return cmd
}
