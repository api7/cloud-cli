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

package config

import (
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

func newRenewCertificateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew-cert",
		Short: "Renew the Certificate for communicating with API7 Cloud",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := persistence.CheckConfigurationAndInitCloudClient(); err != nil {
				output.Errorf(err.Error())
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}

			defaultCluster, err := cloud.DefaultClient.GetDefaultCluster()
			if err != nil {
				output.Errorf(err.Error())
				return
			}

			if err = persistence.DownloadNewCertificate(defaultCluster.ID); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
	}
	return cmd
}
