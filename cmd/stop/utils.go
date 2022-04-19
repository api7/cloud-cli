//  Copyright 2022 API7.ai, Inc under one or more contributor license
//  agreements.  See the NOTICE file distributed with this work for
//  additional information regarding copyright ownership.
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

package stop

import (
	"context"
	"fmt"
	"strings"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/types"
	"github.com/api7/cloud-cli/internal/utils"
)

func stopPreRunForKubernetes(kubectl commands.Cmd) error {
	var err error
	if err = deleteOnKubernetes(kubectl, types.ConfigMap); err != nil {
		return fmt.Errorf("Failed to delete configmap on kubernetes: %s", err.Error())
	}
	if err = deleteOnKubernetes(kubectl, types.Secret); err != nil {
		return fmt.Errorf("Failed to delete secret on kubernetes: %s", err.Error())
	}

	return nil
}

func deleteOnKubernetes(kubectl commands.Cmd, k types.K8sResourceKind) error {
	opts := options.Global.Stop.Kubernetes

	kubectl.AppendArgs("delete")
	switch k {
	case types.ConfigMap:
		kubectl.AppendArgs("configmap", consts.DefaultConfigMapName)
	case types.Secret:
		kubectl.AppendArgs("secret", consts.DefaultSecretName)
	default:
		panic(fmt.Sprintf("invaild kind: %d", k))
	}

	kubectl.AppendArgs("--namespace", opts.NameSpace)

	if options.Global.DryRun {
		output.Infof("Running:\n%s\n", kubectl.String())
	} else {
		output.Verbosef("Running:\n%s\n", kubectl.String())
	}

	newCtx, cancel := context.WithTimeout(context.TODO(), consts.DefaultKubectlTimeout)
	defer cancel()
	go utils.WaitForSignal(func() {
		cancel()
	})

	stdout, stderr, err := kubectl.Run(newCtx)
	if stderr != "" {
		output.Warnf(stderr)
		// if stderr contains NotFound flag, we don't think it's an error
		if strings.Contains(stderr, "NotFound") {
			err = nil
		}
	}
	if stdout != "" {
		output.Verbosef(stdout)
	}
	if err != nil {
		return err
	}

	return nil
}
