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

package stop

import (
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/options"
)

// NewStopCommand creates the stop sub-command object.
func NewStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [COMMAND] [ARG...]",
		Short: "Stop Apache APISIX instance.",
	}

	cmd.PersistentFlags().StringVar(&options.Global.Stop.Name, "name", "apisix", "The identifier of this deployment, it would be the container name (on Docker), the helm release (on Kubernetes) and it's useless if APISIX is deployed on bare metal")
	cmd.PersistentFlags().BoolVar(&options.Global.Stop.Remove, "rm", false, "The identifier of this deployment, which will force the removal of containers in docker")

	cmd.AddCommand(newStopBareCommand())
	cmd.AddCommand(newStopDockerCommand())
	cmd.AddCommand(newStopKubernetesCommand())
	return cmd
}
