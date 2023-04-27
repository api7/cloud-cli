// Copyright 2023 API7.ai, Inc
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

package resource

import (
	"encoding/json"
	"fmt"
	"os"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

var (
	_resourceCreateHandler = map[string]func() interface{}{
		"ssl": func() interface{} {
			var (
				caCert []byte
			)

			cluster, err := cloud.DefaultClient.GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get default cluster: %s", err)
			}

			sslOptions := options.Global.Resource.Create.SSL
			cert, err := os.ReadFile(sslOptions.CertFile)
			if err != nil {
				output.Errorf("Failed to read certificate file: %s", err)
			}
			pkey, err := os.ReadFile(sslOptions.PKeyFile)
			if err != nil {
				output.Errorf("Failed to read private key file: %s", err)
			}
			if sslOptions.CACertFile != "" {
				caCert, err = os.ReadFile(sslOptions.CACertFile)
				if err != nil {
					output.Errorf("Failed to read CA certificate file: %s", err)
				}
			}

			certificate := &sdk.Certificate{
				CertificateSpec: sdk.CertificateSpec{
					Certificate: string(cert),
					PrivateKey:  string(pkey),
					Labels:      options.Global.Resource.Create.Labels,
					Type:        sslOptions.Type,
				},
			}
			if caCert != nil {
				certificate.CACertificate = string(caCert)
			}

			details, err := cloud.DefaultClient.CreateSSL(cluster.ID, certificate)
			if err != nil {
				output.Errorf("Failed to create certificate: %s", err)
			}
			return details
		},
		"service": func() interface{} {
			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get default cluster: %s", err.Error())
			}
			svc, err := readServiceFromFile(options.Global.Resource.Create.FromFile)
			if err != nil {
				output.Errorf("Failed to read service from file: %s", err.Error())
			}
			newSvc, err := cloud.DefaultClient.CreateService(cluster.ID, svc)
			if err != nil {
				output.Errorf("Failed to create service: %s", err.Error())
			}
			return newSvc
		},
	}
)

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a resource",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := options.Global.Resource.Create.Validate(); err != nil {
				output.Errorf(err.Error())
				return
			}
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			kind := options.Global.Resource.Create.Kind
			handler, ok := _resourceCreateHandler[kind]
			if !ok {
				output.Errorf("This kind of resource is not supported")
			} else {
				resource := handler()
				text, _ := json.MarshalIndent(resource, "", "\t")
				fmt.Println(string(text))
			}
		},
	}

	cmd.PersistentFlags().StringVar(&options.Global.Resource.Create.Kind, "kind", "", "Specify the resource kind")
	cmd.PersistentFlags().StringSliceVar(&options.Global.Resource.Create.Labels, "label", nil, "Add label for this resource")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Create.SSL.CertFile, "cert", "", "Specify the certificate file (this option is only useful when kind is ssl)")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Create.SSL.PKeyFile, "pkey", "", "Specify the private key file (this option is only useful when kind is ssl)")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Create.SSL.CACertFile, "cacert", "", "Specify the CA certificate key file (this option is only useful when kind is ssl)")
	cmd.PersistentFlags().StringVar((*string)(&options.Global.Resource.Create.SSL.Type), "ssl-type", "", "Specify the SSL type (optional value can be \"server\", \"client\"")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Create.FromFile, "from-file", "", "Specify the resource definition file")
	return cmd
}
