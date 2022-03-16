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
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

func newKubernetesCommand() *cobra.Command {
	var (
		ctx deployContext
	)
	cmd := &cobra.Command{
		Use:   "kubernetes [ARGS...]",
		Short: "Deploy Apache APISIX on Kubernetes",
		Example: `
cloud-cli deploy kubernetes \
		--name apisix \
		--namespace apisix \
		--apisix-image apisix/apisix:2.11.0 \
		--helm-install-arg --output=table \
		--helm-install-arg --wait`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
			if err := deployPreRunForKubernetes(&ctx); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
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
