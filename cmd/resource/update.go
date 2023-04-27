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
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/spf13/cobra"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
)

var (
	_resourceUpdateHandler = map[string]func(sdk.ID) interface{}{
		"service": func(id sdk.ID) interface{} {
			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to list services: %s", err.Error())
			}
			svc, err := readServiceFromFile(options.Global.Resource.Update.FromFile)
			if err != nil {
				output.Errorf("Failed to read service from file: %s", err.Error())
			}
			svc.ID = id
			newSvc, err := cloud.DefaultClient.UpdateService(cluster.ID, svc)
			if err != nil {
				output.Errorf("Failed to list services: %s", err.Error())
			}
			return newSvc
		},
		"consumer": func(id sdk.ID) interface{} {
			cluster, err := cloud.Client().GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to list consumer: %s", err.Error())
			}
			consumer, err := readConsumerFromFile(options.Global.Resource.Update.FromFile)
			if err != nil {
				output.Errorf("Failed to read consumer from file: %s", err.Error())
			}
			consumer.ID = id
			newConsumer, err := cloud.DefaultClient.UpdateConsumer(cluster.ID, consumer)
			if err != nil {
				output.Errorf("Failed to update consumers: %s", err.Error())
			}
			return newConsumer
		},
		"ssl": func(id sdk.ID) interface{} {
			var (
				caCert []byte
			)

			cluster, err := cloud.DefaultClient.GetDefaultCluster()
			if err != nil {
				output.Errorf("Failed to get default cluster: %s", err)
			}

			sslOptions := options.Global.Resource.Update.SSL
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
				ID: id,
				CertificateSpec: sdk.CertificateSpec{
					Certificate: string(cert),
					PrivateKey:  string(pkey),
					Labels:      options.Global.Resource.Update.Labels,
					Type:        sslOptions.Type,
				},
			}
			if caCert != nil {
				certificate.CACertificate = string(caCert)
			}

			details, err := cloud.DefaultClient.UpdateSSL(cluster.ID, certificate)
			if err != nil {
				output.Errorf("Failed to update certificate: %s", err)
			}
			return details
		},
	}
)

func readServiceFromFile(filename string) (*sdk.Application, error) {
	var (
		app *sdk.Application
		err error
	)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	switch filepath.Ext(filename) {
	case ".json":
		err = json.Unmarshal(data, &app)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &app)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal service: %s", err)
	}

	return app, nil
}

func readConsumerFromFile(filename string) (*sdk.Consumer, error) {
	var (
		consumer *sdk.Consumer
		err      error
	)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	switch filepath.Ext(filename) {
	case ".json":
		err = json.Unmarshal(data, &consumer)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &consumer)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal consumer: %s", err)
	}

	return consumer, nil
}

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update a resource",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := options.Global.Resource.Update.Validate(); err != nil {
				output.Errorf(err.Error())
				return
			}
			if err := persistence.Init(); err != nil {
				output.Errorf(err.Error())
				return
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			id := options.Global.Resource.Update.ID
			kind := options.Global.Resource.Update.Kind

			handler, ok := _resourceUpdateHandler[kind]
			if !ok {
				output.Errorf("This kind of resource is not supported")
			} else {
				uint64ID, _ := strconv.ParseUint(id, 10, 64)
				resource := handler(sdk.ID(uint64ID))
				text, _ := json.MarshalIndent(resource, "", "\t")
				fmt.Println(string(text))
			}
		},
	}
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Update.FromFile, "from-file", "", "Specify the resource definition file")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Update.Kind, "kind", "", "Specify the resource kind")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Update.ID, "id", "", "Specify the resource ID")
	cmd.PersistentFlags().StringSliceVar(&options.Global.Resource.Update.Labels, "label", nil, "Add label for this resource")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Update.SSL.CertFile, "cert", "", "Specify the certificate file (this option is only useful when kind is ssl)")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Update.SSL.PKeyFile, "pkey", "", "Specify the private key file (this option is only useful when kind is ssl)")
	cmd.PersistentFlags().StringVar(&options.Global.Resource.Update.SSL.CACertFile, "cacert", "", "Specify the CA certificate key file (this option is only useful when kind is ssl)")
	cmd.PersistentFlags().StringVar((*string)(&options.Global.Resource.Update.SSL.Type), "ssl-type", "", "Specify the SSL type (optional value can be \"server\", \"client\"")
	return cmd
}
