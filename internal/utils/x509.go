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

package utils

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/pkg/errors"
)

// CheckIfCertificateIsExpired checks whether the certificate is expired or not.
// Note this function accepts the certificate as a byte array (in PEM format).
func CheckIfCertificateIsExpired(data []byte) (bool, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return false, errors.New("failed to decode certificate from PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse certificate")
	}
	if cert.NotAfter.Before(time.Now()) {
		return true, nil
	}
	return false, nil
}
