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

package configure

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

// NewCommand creates the `configure` sub-command object.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure the credential for accessing API7 Cloud.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("API7 Cloud Access Token: ")
			var accessToken string

			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				accessToken = scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				output.Errorf("reading standard input: %s", err)
			}

			token, err := jwt.Parse(accessToken, nil)
			if err != nil {
				if e, ok := err.(*jwt.ValidationError); ok {
					if e.Errors != jwt.ValidationErrorUnverifiable {
						output.Errorf("invalid access token: %s", err)
					}
				} else {
					output.Errorf("access token parse error: %s", err)
				}
			}
			claim, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				output.Errorf("invalid claim type")
			}

			if claim["exp"] != nil {
				expireAt := int64(claim["exp"].(float64))
				if expireAt < time.Now().Unix() {
					output.Errorf("access token expired")
				}
				output.Warnf("your access token will expire at %s", time.Unix(expireAt, 0).Format(time.RFC3339))
			} else {
				output.Warnf("You are using a token that has no expiration time, please note the security risk")
			}

			if err := persistence.SaveCredential(&persistence.Credential{
				User: persistence.User{
					AccessToken: accessToken,
				},
			}); err != nil {
				output.Errorf(err.Error())
			}

			err = cloud.InitDefaultClient(accessToken)
			if err != nil {
				output.Errorf("failed to initialize api7 cloud client: %s", err)
			}

			me, err := cloud.Client().Me()
			if err != nil {
				output.Errorf("failed to request api7 cloud: %s", err)
			}

			output.Infof("successfully configured api7 cloud access token, your account is %s", me.Email)
		},
	}

	return cmd
}
