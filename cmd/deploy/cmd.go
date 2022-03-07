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
)

// NewCommand creates the deploy sub-command object.
func NewCommand(opts *options.DeployOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy [COMMAND] [ARG...]",
		Short: "Deploy Apache APISIX with being connected to API7 Cloud.",
	}

	cmd.PersistentFlags().StringVar(&opts.APISIXConfigFile, "apisix-config", "", "Specify the custom APISIX configuration file")
	cmd.PersistentFlags().StringVar(&opts.APISIXInstanceID, "apisix-id", "", "Specify the custom APISIX instance ID")

	cmd.AddCommand(newDockerCommand(&opts.Docker))

	return cmd
}