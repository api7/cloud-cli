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
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/spf13/cobra"
)

func newKubernetesCommand()  *cobra.Command {
	var (
		cloudLuaModuleDir string
	)
	cmd := &cobra.Command{
		Use:   "kubernetes [ARGS...]",
		Short: "Deploy Apache APISIX on kubernetes",
		Example: `
cloud-cli deploy kubernetes \
		--name apisix \
		--namespace apisix \
		--apisix-image apisix/apisix:2.11.0 \
		--secret-name cloud-ssl \
		--helm-install-arg --output=table \
		--helm-install-arg --description=this is a description`,
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error

			if err = persistence.PrepareCertificate(); err != nil {
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

		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.NameSpace, "namespace", "apisix", "Specify the Kubernetes nameSpace")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.APISIXImage, "apisix-image", "apache/apisix:2.11.0", "Specify the Apache APISIX image")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.SecretName, "secret-name", "cloud-ssl", "Specify the kubernetes secret name")
	cmd.PersistentFlags().StringSliceVar(&options.Global.Deploy.Kubernetes.HelmInstallArgs, "helm-install-arg", []string{}, "Specify the arguments (in the format of name=value) for the helm install command")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.KubectlCLIPath, "kubectl-cli-path", "", "Specify the filepath of the kubectl command")
	cmd.PersistentFlags().StringVar(&options.Global.Deploy.Kubernetes.HelmCLIPath, "helm-cli-path", "", "Specify the filepath of the helm command")

	return cmd
}

func createSecretOnK8s() error {
	var (
		kubectl *commands.Cmd
		data   []byte
		err    error
	)

	opts := options.Global.Deploy.Kubernetes
	if opts.KubectlCLIPath != "" {
		kubectl = commands.New(opts.KubectlCLIPath, options.Global.DryRun)
	} else {
		kubectl = commands.New("kubectl", options.Global.DryRun)
	}

	kubectl.AppendArgs("create","secret","generic",options.Global.Deploy.Kubernetes.SecretName)
	kubectl.AppendArgs("--form-file",)
}

func deployByHelm() error {

}