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
	"time"

	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/utils"
)

func newStopBareCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bare",
		Short: "Stop Apache APISIX on bare metal (only CentOS 7)",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Minute)
			go utils.WaitForSignal(func() {
				cancel()
			})

			bare := commands.New("apisix", options.Global.DryRun)
			bare.AppendArgs("stop")
			if options.Global.DryRun {
				output.Infof(bare.String())
				return
			}
			stdout, stderr, err := bare.Run(ctx)
			if stderr != "" {
				output.Warnf(stderr)
			}
			if stdout != "" {
				output.Verbosef(stdout)
			}
			if err != nil {
				output.Errorf(err.Error())
				return
			}
		},
	}

	return cmd
}
