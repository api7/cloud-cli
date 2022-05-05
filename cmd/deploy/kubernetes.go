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
	"io/ioutil"
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

const (
	// defaultHelmChartsUrl is the default url for helm charts
	defaultHelmChartsUrl = "https://charts.apiseven.com"
)

func newKubernetesCommand() *cobra.Command {
	var (
		ctx     deployContext
		kubectl commands.Cmd
	)

	cmd := &cobra.Command{
		Use:   "kubernetes [ARGS...]",
		Short: "Deploy Apache APISIX on Kubernetes",
		Example: `
cloud-cli deploy kubernetes \
		--name apisix \
		--namespace apisix \
		--apisix-image apache/apisix:2.11.0-centos \
		--helm-install-arg --output=table \
		--helm-install-arg --wait`,
		PreRun: func(cmd *cobra.Command, args []string) {
			opts := &options.Global.Deploy.Kubernetes
			if opts.KubectlCLIPath == "" {
				opts.KubectlCLIPath = "kubectl"
			}
			kubectl = commands.New(opts.KubectlCLIPath, options.Global.DryRun)

			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
			if err := deployPreRunForKubernetes(&ctx, kubectl); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				data         []byte
				mergedConfig map[string]interface{}
				configFile   string
				err          error
			)

			if ctx.KubernetesOpts.HelmCLIPath == "" {
				ctx.KubernetesOpts.HelmCLIPath = "helm"
			}
			helm := commands.New(ctx.KubernetesOpts.HelmCLIPath, options.Global.DryRun)

			newCtx, cancel := context.WithTimeout(context.TODO(), consts.DefaultHelmTimeout)
			defer cancel()
			go utils.WaitForSignal(func() {
				cancel()
			})

			{
				helm.AppendArgs("repo", "add", "apisix", defaultHelmChartsUrl)
				helmRun(newCtx, helm)
			}

			{
				helm = commands.New(ctx.KubernetesOpts.HelmCLIPath, options.Global.DryRun)
				helm.AppendArgs("repo", "update")
				helmRun(newCtx, helm)
			}

			{
				helm = commands.New(ctx.KubernetesOpts.HelmCLIPath, options.Global.DryRun)
				helm.AppendArgs("install", options.Global.Deploy.Name, "apisix/apisix")
				helm.AppendArgs("--namespace", ctx.KubernetesOpts.NameSpace)

				var customizeValues string
				for _, args := range ctx.KubernetesOpts.HelmInstallArgs {
					if strings.Contains(args, "--values=") {
						kv := strings.Split(args, "=")
						if len(kv) != 2 {
							output.Errorf("invalid --values option")
						}
						customizeValues = kv[1]
						continue
					}
					helm.AppendArgs(strings.SplitN(args, "=", 2)...)
				}

				if customizeValues != "" {
					if data, err = ioutil.ReadFile(customizeValues); err != nil {
						output.Errorf("invalid --apisix-config-file option: %s", err)
					}
				}

				if mergedConfig, err = apisix.MergeConfig(data, ctx.essentialConfig); err != nil {
					output.Errorf(err.Error())
				}
				if configFile, err = apisix.SaveConfig(mergedConfig, "helm-values-*.yaml"); err != nil {
					output.Errorf(err.Error())
				}
				helm.AppendArgs("--values", configFile)

				helmRun(newCtx, helm)
				printInstallDetailForKubernetes(kubectl)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.NameSpace, "namespace", "apisix", "Specify the Kubernetes name space")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.APISIXImage, "apisix-image", "apache/apisix:2.11.0-centos", "Specify the Apache APISIX image")
	cmd.PersistentFlags().UintVar(&options.Global.Deploy.Kubernetes.ReplicaCount, "replica-count", 1, "Specify the pod replica count")
	cmd.PersistentFlags().StringSliceVar(&options.Global.Deploy.Kubernetes.HelmInstallArgs, "helm-install-arg", []string{}, "Specify the arguments (in the format of name=value) for the helm install command")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.KubectlCLIPath, "kubectl-cli-path", "", "Specify the filepath of the kubectl command")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.HelmCLIPath, "helm-cli-path", "", "Specify the filepath of the helm command")

	return cmd
}

func helmRun(ctx context.Context, helm commands.Cmd) {
	if options.Global.DryRun {
		output.Infof("Running:\n%s\n", helm.String())
	} else {
		output.Verbosef("Running:\n%s\n", helm.String())
	}

	stdout, stderr, err := helm.Run(ctx)
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
