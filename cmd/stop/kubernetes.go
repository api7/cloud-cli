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

package stop

import (
	"context"
	"github.com/api7/cloud-cli/internal/consts"
	"github.com/spf13/cobra"
	"strings"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/utils"
)

func newStopKubernetesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes [ARG...]",
		Short: "Stop Apache APISIX on Kubernetes",
		Example: `
cloud-cli stop kubernetes \
		--name apisix \
		--namespace apisix \
		--helm-uninstall-arg --keep-history \
		--helm-uninstall-arg --wait`,
		PreRun: func(cmd *cobra.Command, args []string) {
			opts := options.Global.Stop.Kubernetes
			if opts.KubectlCLIPath == "" {
				opts.KubectlCLIPath = "kubectl"
			}
			kubectl := commands.New(opts.KubectlCLIPath, options.Global.DryRun)

			if err := stopPreRunForKubernetes(kubectl); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.TODO(), consts.DefaultHelmTimeout)
			defer cancel()
			go utils.WaitForSignal(func() {
				cancel()
			})

			opts := options.Global.Stop.Kubernetes
			if opts.HelmCLIPath == "" {
				opts.HelmCLIPath = "helm"
			}

			helm := commands.New(opts.HelmCLIPath, options.Global.DryRun)
			helm.AppendArgs("uninstall", options.Global.Stop.Name, "--namespace", opts.NameSpace)

			for _, args := range opts.HelmUnInstallArgs {
				helm.AppendArgs(strings.Split(args, "=")...)
			}

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
		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Stop.Kubernetes.NameSpace, "namespace", "apisix", "Specify the Kubernetes name space")
	cmd.PersistentFlags().StringVar(&options.Global.Stop.Kubernetes.HelmCLIPath, "helm-cli-path", "", "Specify the filepath of the helm command")
	cmd.PersistentFlags().StringSliceVar(&options.Global.Stop.Kubernetes.HelmUnInstallArgs, "helm-uninstall-arg", []string{}, "Specify the arguments (in the format of name=value) for the helm uninstall command")
	cmd.PersistentFlags().StringVar(&options.Global.Stop.Kubernetes.KubectlCLIPath, "kubectl-cli-path", "", "Specify the filepath of the kubectl command")
	return cmd
}
