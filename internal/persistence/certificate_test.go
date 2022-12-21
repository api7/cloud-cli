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
	"testing"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
)

func TestPrepareCertificate(t *testing.T) {
	testCases := []struct {
		name           string
		cpID           sdk.ID
		preparedCert   string
		preparedKey    string
		preparedCACert string
		expectedCert   string
		expectedKey    string
		expectedCACert string
		errorReason    string
		mockFn         func(t *testing.T)
	}{
		{
			name:           "bad tls bundle",
			preparedCert:   "abcdef",
			preparedKey:    "abcdef",
			preparedCACert: "abcdef",
			errorReason:    "check certificate availability: failed to decode certificate from PEM",
			mockFn: func(t *testing.T) {

			},
		},
		{
			name:        "empty cert but get tls bundle failed (mock error)",
			errorReason: "download tls bundle: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name: "success",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)
				cloud.DefaultClient = mockClient
			},
			expectedCert:   "1",
			expectedKey:    "1",
			expectedCACert: "1",
		},
		{
			name: "success with cp id",
			cpID: 123456,
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)
				cloud.DefaultClient = mockClient
			},
			expectedCert:   "1",
			expectedKey:    "1",
			expectedCACert: "1",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.cpID == 0 {
				tc.cpID = 1
			}
			cpTLSDir := filepath.Join(TLSDir, tc.cpID.String())

			if tc.preparedCert != "" {
				// create cp tls dir
				if err := os.MkdirAll(cpTLSDir, 0755); err != nil {
					assert.Nil(t, err, "failed to create tls directory")
				}
				if err := os.Chmod(cpTLSDir, 0755); err != nil {
					assert.Nil(t, err, "change tls directory permission")
				}
				defer os.Remove(cpTLSDir)

				certFilename := filepath.Join(cpTLSDir, "tls.crt")
				certKeyFilename := filepath.Join(cpTLSDir, "tls.key")
				certCAFilename := filepath.Join(cpTLSDir, "ca.crt")

				err := os.WriteFile(certFilename, []byte(tc.preparedCert), 0600)
				assert.Nil(t, err, "check if cert is saved")
				defer os.Remove(certFilename)

				err = os.WriteFile(certKeyFilename, []byte(tc.preparedKey), 0600)
				assert.Nil(t, err, "check if pkey is saved")
				defer os.Remove(certKeyFilename)

				err = os.WriteFile(certCAFilename, []byte(tc.preparedCACert), 0600)
				assert.Nil(t, err, "check if ca cert is saved")
				defer os.Remove(certCAFilename)
			}

			tc.mockFn(t)

			err := PrepareCertificate(tc.cpID)
			defer os.Remove(cpTLSDir)
			if tc.errorReason == "" {
				assert.Nil(t, err, "check if err is nil")
				certFilename := filepath.Join(cpTLSDir, "tls.crt")
				certKeyFilename := filepath.Join(cpTLSDir, "tls.key")
				certCAFilename := filepath.Join(cpTLSDir, "ca.crt")
				defer os.Remove(certFilename)
				defer os.Remove(certKeyFilename)
				defer os.Remove(certCAFilename)
				cert, err := os.ReadFile(certFilename)
				assert.Nil(t, err, "read cert")
				pkey, err := os.ReadFile(certKeyFilename)
				assert.Nil(t, err, "read pkey")
				ca, err := os.ReadFile(certCAFilename)
				assert.Nil(t, err, "read ca cert")

				assert.Equal(t, tc.expectedCert, string(cert), "check cert")
				assert.Equal(t, tc.expectedKey, string(pkey), "check pkey")
				assert.Equal(t, tc.expectedCACert, string(ca), "check ca cert")
			} else {
				assert.Equal(t, tc.errorReason, err.Error(), "check if err is correct")
			}
		})
	}
}
