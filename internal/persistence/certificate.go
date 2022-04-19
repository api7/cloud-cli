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

package persistence

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/utils"
)

// PrepareCertificate downloads the client certificate and key from API7 Cloud.
// This certificate is used for the communication between APISIX and API7 Cloud.
func PrepareCertificate(cpID string) error {
	certFilename := filepath.Join(TLSDir, "tls.crt")
	if available, err := checkIfCertificateAvailable(certFilename); err != nil {
		return errors.Wrap(err, "check certificate availability")
	} else if available {
		return nil
	}

	output.Verbosef("Downloading tls bundle from API7 Cloud")

	// Currently, only one control plane is supported for an organization.
	bundle, err := cloud.DefaultClient.GetTLSBundle(cpID)
	if err != nil {
		return errors.Wrap(err, "download tls bundle")
	}

	err = ioutil.WriteFile(certFilename, []byte(bundle.Certificate), 0644)
	if err != nil {
		return errors.Wrap(err, "save certificate")
	}

	certKeyFilename := filepath.Join(TLSDir, "tls.key")
	err = ioutil.WriteFile(certKeyFilename, []byte(bundle.PrivateKey), 0644)
	if err != nil {
		return errors.Wrap(err, "save private key")
	}

	certCAFilename := filepath.Join(TLSDir, "ca.crt")
	err = ioutil.WriteFile(certCAFilename, []byte(bundle.CACertificate), 0644)
	if err != nil {
		return errors.Wrap(err, "save ca certificate")
	}
	return nil
}

func checkIfCertificateAvailable(certFilename string) (bool, error) {
	data, err := ioutil.ReadFile(certFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	expired, err := utils.CheckIfCertificateIsExpired(data)
	if err != nil {
		return false, err
	}
	return !expired, nil
}
