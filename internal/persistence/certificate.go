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
	"os"
	"path/filepath"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/utils"
)

// PrepareCertificate downloads the client certificate and key from API7 Cloud.
// This certificate is used for the communication between APISIX and API7 Cloud.
func PrepareCertificate(clusterID sdk.ID) error {
	clusterTLSDir := filepath.Join(TLSDir, clusterID.String())

	certFilename := filepath.Join(clusterTLSDir, "tls.crt")
	if available, err := checkIfCertificateAvailable(certFilename); err != nil {
		return errors.Wrap(err, "check certificate availability")
	} else if available {
		return nil
	}

	output.Verbosef("Downloading tls bundle from API7 Cloud")

	// Currently, only one cluster is supported for an organization.
	bundle, err := cloud.DefaultClient.GetTLSBundle(clusterID)
	if err != nil {
		return errors.Wrap(err, "download tls bundle")
	}

	// Make cluster tls dir
	if err = os.MkdirAll(clusterTLSDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create tls directory")
	}
	if err = os.Chmod(clusterTLSDir, 0755); err != nil {
		return errors.Wrap(err, "change tls directory permission")
	}

	err = os.WriteFile(certFilename, []byte(bundle.Certificate), 0644)
	if err != nil {
		return errors.Wrap(err, "save certificate")
	}
	// permission in WriteFile is before umask, so we need to chmod it
	if err = os.Chmod(certFilename, 0644); err != nil {
		return errors.Wrap(err, "change certificate permission")
	}

	certKeyFilename := filepath.Join(clusterTLSDir, "tls.key")
	err = os.WriteFile(certKeyFilename, []byte(bundle.PrivateKey), 0644)
	if err != nil {
		return errors.Wrap(err, "save private key")
	}
	if err = os.Chmod(certKeyFilename, 0644); err != nil {
		return errors.Wrap(err, "change private key permission")
	}

	certCAFilename := filepath.Join(clusterTLSDir, "ca.crt")
	err = os.WriteFile(certCAFilename, []byte(bundle.CACertificate), 0644)
	if err != nil {
		return errors.Wrap(err, "save ca certificate")
	}
	if err = os.Chmod(certCAFilename, 0644); err != nil {
		return errors.Wrap(err, "change ca certificate permission")
	}

	return nil
}

func checkIfCertificateAvailable(certFilename string) (bool, error) {
	data, err := os.ReadFile(certFilename)
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
