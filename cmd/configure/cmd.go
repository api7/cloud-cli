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

package configure

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

// NewCommand creates the `configure` sub-command object.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure the credential for accessing API7 Cloud.",
		Run: func(cmd *cobra.Command, args []string) {
			if options.Global.Configure.AccessToken == "" {
				fmt.Printf("API7 Cloud Access Token: ")

				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					options.Global.Configure.AccessToken = scanner.Text()
				}
				if err := scanner.Err(); err != nil {
					output.Errorf("reading standard input: %s", err)
				}
			}

			token, err := jwt.Parse(options.Global.Configure.AccessToken, nil)
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

			err = cloud.InitDefaultClient(options.Global.Configure.Addr, options.Global.Configure.AccessToken)
			if err != nil {
				output.Errorf("failed to initialize api7 cloud client: %s", err)
			}

			me, err := cloud.Client().Me()
			if err != nil {
				output.Errorf("failed to request api7 cloud: %s", err)
			}

			profileName := options.Global.Configure.Profile

			configuration, err := persistence.LoadConfiguration()
			if err != nil {
				output.Verbosef("there is no configuration file, create a new one")
				configuration = &persistence.CloudConfiguration{}
				// generate a random profile name if not specified at first time
				if profileName == "" {
					profileName = namesgenerator.GetRandomName(0)
				}
			} else {
				if profileName == "" {
					profileName = configuration.DefaultProfile
				}
			}

			newProfile := persistence.Profile{
				Name:    profileName,
				Address: options.Global.Configure.Addr,
				User: persistence.User{
					AccessToken: options.Global.Configure.AccessToken,
				},
			}
			configuration.ConfigureProfile(newProfile)

			if options.Global.Configure.Default || len(configuration.Profiles) == 1 {
				configuration.DefaultProfile = newProfile.Name
			}

			if err := persistence.SaveConfiguration(configuration); err != nil {
				output.Errorf(err.Error())
			}

			output.Infof("successfully configured api7 cloud access token, your account is %s", me.Email)
		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Configure.Addr, "addr", "https://api.api7.cloud", "Specify the API7 Cloud server address")
	cmd.PersistentFlags().StringVar(&options.Global.Configure.Profile, "profile", "", "Specify the profile name")
	cmd.PersistentFlags().BoolVar(&options.Global.Configure.Default, "set-default", true, "Set the profile as default")
	cmd.PersistentFlags().StringVar(&options.Global.Configure.AccessToken, "token", "", "Specify the access token")

	return cmd
}
