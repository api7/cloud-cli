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

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/config"
	"github.com/api7/cloud-cli/internal/output"
)

// NewCommand creates the `configure` sub-command object.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure the credential for accessing API7 Cloud.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("API7 Cloud Access Key: ")
			var accessKey string

			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				accessKey = scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				output.Errorf("reading standard input: %s", err)
			}

			token, err := jwt.Parse(accessKey, nil)
			if err != nil {
				if e, ok := err.(*jwt.ValidationError); ok {
					if e.Errors != jwt.ValidationErrorUnverifiable {
						output.Errorf("invalid access key: %s", err)
					}
				} else {
					output.Errorf("access key parse error: %s", err)
				}
			}
			claim, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				output.Errorf("invalid claim type")
			}

			expireAt := int64(claim["exp"].(float64))
			if expireAt != -1 {
				if expireAt < time.Now().Unix() {
					output.Errorf("access key expired")
				}
			}

			output.Warnf("your access key will expire at %s", time.Unix(expireAt, 0).Format(time.RFC3339))

			if err := config.Save(&config.Config{
				User: config.User{
					AccessKey: accessKey,
				},
			}); err != nil {
				output.Errorf(err.Error())
			}
		},
	}

	return cmd
}
