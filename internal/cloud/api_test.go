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

package cloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	sdk "github.com/api7/cloud-go-sdk"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
)

func TestMe(t *testing.T) {
	tests := []struct {
		name      string
		code      int
		body      string
		want      *sdk.User
		wantErr   bool
		errReason string
	}{
		{
			name:      "server internal error",
			code:      http.StatusInternalServerError,
			body:      "internal server error",
			want:      nil,
			wantErr:   true,
			errReason: "status code: 500, message: internal server error",
		},
		{
			name:      "http not found",
			code:      http.StatusNotFound,
			want:      nil,
			wantErr:   true,
			body:      `{"status": {"code": 4, "message": "not found"}, "error": "deliberated not found"}`,
			errReason: "status code: 404, error code: 4, error reason: not found, details: deliberated not found",
		},
		{
			name:      "malformed json",
			code:      http.StatusOK,
			body:      `invalid json`,
			want:      nil,
			wantErr:   true,
			errReason: "decode response body: invalid character 'i' looking for beginning of value",
		},
		{
			name:      "token expired",
			code:      http.StatusUnauthorized,
			body:      `{"status":{"code":6,"message":"unauthorized"},"error":"Token is expired"}`,
			want:      nil,
			wantErr:   true,
			errReason: "Token is expired",
		},
		{
			name: "success",
			code: http.StatusOK,
			body: `
			{
				"payload": {
					"id": "321",
					"first_name": "first",
					"last_name": "last",
					"email": "demo@api7.ai",
					"org_ids": [
						"123"
					],
					"connection": "",
					"avatar_url": "https://api7.ai/avatar/default.png"
				},
				"status": {
					"code": 0,
					"message": "OK"
				}
			}`,
			want: &sdk.User{
				ID:         "321",
				FirstName:  "first",
				LastName:   "last",
				Email:      "demo@api7.ai",
				OrgIDs:     []sdk.ID{123},
				Connection: "",
				AvatarURL:  "https://api7.ai/avatar/default.png",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options.Global.Verbose = true
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), "/api/v1/user/me")
				assert.Equal(t, req.Header.Get("Authorization"), "Bearer test-token")

				rw.WriteHeader(tt.code)
				if tt.body != "" {
					write, err := rw.Write([]byte(tt.body))
					assert.NoError(t, err, "send mock response")
					assert.Equal(t, len(tt.body), write, "write mock response")
				}
			}))

			defer server.Close()

			api, err := newClient(server.URL, "test-token")
			assert.NoError(t, err, "checking new cloud api client")

			result, err := api.Me()

			if tt.wantErr {
				assert.Contains(t, err.Error(), tt.errReason, "checking error reason")
			} else {
				assert.NoError(t, err, "checking error")
				assert.Equal(t, tt.want, result, "checking result")
			}
		})
	}
}

func TestClusters(t *testing.T) {
	tests := []struct {
		name      string
		orgID     sdk.ID
		code      int
		body      string
		want      *sdk.Cluster
		wantErr   bool
		errReason string
	}{
		{
			name: "success",
			code: http.StatusOK,
			body: `
			{
				"payload": {
					"count": 1,
					"list": [
						{
							"id": "392306215415186327",
							"name": "default",
							"org_id": "392306215398409111",
							"region_id": "56523356",
							"status": 3,
							"domain": "default.xvlbzgs4bqbjdmybyqk.api7.cloud",
							"config_payload": "",
							"org_name": "XVlBzg"
						}
					]
				},
				"status": {
					"code": 0,
					"message": "OK"
				}
			}`,
			want: &sdk.Cluster{
				ID:   392306215415186327,
				Name: "default",
				ClusterSpec: sdk.ClusterSpec{
					OrganizationID: 392306215398409111,
					RegionID:       56523356,
					Status:         3,
					Domain:         "default.xvlbzgs4bqbjdmybyqk.api7.cloud",
					ConfigPayload:  "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Contains(t, req.URL.String(), fmt.Sprintf("/api/v1/orgs/%s/clusters", tt.orgID))
				assert.Equal(t, req.Header.Get("Authorization"), "Bearer test-token")

				rw.WriteHeader(tt.code)
				if req.URL.Query().Get("page") == "2" {
					// to avoid dead loop, return empty list when reaching page 2.
					_, err := rw.Write([]byte(
						`
			{
				"payload": {
					"count": 0,
					"list": []
				},
				"status": {
					"code": 0,
					"message": "OK"
				}
			}`))
					assert.NoError(t, err, "send mock response")

				} else if tt.body != "" {
					write, err := rw.Write([]byte(tt.body))
					assert.NoError(t, err, "send mock response")
					assert.Equal(t, len(tt.body), write, "write mock response")
				}
			}))

			defer server.Close()

			api, err := newClient(server.URL, "test-token")
			assert.NoError(t, err, "checking new cloud api client")

			result, err := api.ListClusters(tt.orgID)

			if tt.wantErr {
				assert.Contains(t, err.Error(), tt.errReason, "checking error reason")
			} else {
				assert.NoError(t, err, "checking error")
				assert.Len(t, result, 1, "checking clusters count")
				assert.Equal(t, tt.want, result[0], "checking result")
			}
		})
	}
}

func TestGetTLSBundle(t *testing.T) {
	tests := []struct {
		name      string
		clusterID sdk.ID
		code      int
		body      string
		want      *sdk.TLSBundle
		wantErr   bool
		errReason string
	}{
		{
			name:      "success",
			code:      http.StatusOK,
			clusterID: 1,
			want: &sdk.TLSBundle{
				Certificate:   "1",
				PrivateKey:    "1",
				CACertificate: "1",
			},
			body: `
				{
					"code": 0,
					"payload": {
						"certificate": "1",
						"private_key": "1",
						"ca_certificate": "1"
					}
				}
			`,
		},
		{
			name:      "internal server error",
			code:      http.StatusInternalServerError,
			clusterID: 1,
			body:      "internal server error",
			errReason: "status code: 500, message: internal server error",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), fmt.Sprintf("/api/v1/clusters/%s/dp_certificate", tt.clusterID))
				assert.Equal(t, req.Header.Get("Authorization"), "Bearer test-token")

				rw.WriteHeader(tt.code)
				if tt.body != "" {
					write, err := rw.Write([]byte(tt.body))
					assert.NoError(t, err, "send mock response")
					assert.Equal(t, len(tt.body), write, "write mock response")
				}
			}))

			defer server.Close()

			api, err := newClient(server.URL, "test-token")
			assert.NoError(t, err, "checking new cloud api client")

			bundle, err := api.GetTLSBundle(tt.clusterID)

			if tt.wantErr {
				assert.Contains(t, err.Error(), tt.errReason, "checking error reason")
			} else {
				assert.NoError(t, err, "checking error")
				assert.Equal(t, tt.want, bundle, "check the tls bundle")
			}
		})
	}
}

func TestGetCloudLuaModule(t *testing.T) {
	testCases := []struct {
		name        string
		errorReason string
		code        int
		body        string
	}{
		{
			name:        "bad code 400",
			errorReason: "unexpected response code: 400, message",
			code:        http.StatusBadRequest,
		},
		{
			name: "success",
			code: http.StatusOK,
			body: "the lua module",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(tc.code)
				if tc.body != "" {
					write, err := rw.Write([]byte(tc.body))
					assert.NoError(t, err, "send mock response")
					assert.Equal(t, len(tc.body), write, "write mock response")
				}
			}))

			defer server.Close()

			err := os.Setenv(consts.Api7CloudLuaModuleURL, server.URL+"/")
			assert.NoError(t, err, "checking env setup")

			api, err := newClient(server.URL, "test-token")
			assert.NoError(t, err, "checking new cloud api client")

			data, err := api.GetCloudLuaModule()

			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "checking error reason")
			} else {
				assert.NoError(t, err, "checking error")
				assert.Equal(t, tc.body, string(data), "check the lua module")
			}
		})
	}
}

func TestGetStartupConfig(t *testing.T) {
	testCases := []struct {
		name         string
		configType   StartupConfigType
		errorReason  string
		code         int
		expectedBody string
		body         string
	}{
		{
			name:        "bad code 400",
			configType:  APISIX,
			errorReason: "Error Code: 4, Error Reason",
			code:        http.StatusBadRequest,
			body: `
				{
					"status": {
						"code": 4
					},
					"error_reason": "400"
				}`,
		},
		{
			name:       "success (apisix)",
			configType: APISIX,
			code:       http.StatusOK,
			body: `
				{
					"status": {
						"code": 0
					},
					"payload": {
						"configuration": "abc"
					}
				}
			`,
			expectedBody: "abc",
		},
		{
			name:       "success (helm)",
			configType: HELM,
			code:       http.StatusOK,
			body: `
				{
					"status": {
						"code": 0
					},
					"payload": {
						"configuration": "abc"
					}
				}
			`,
			expectedBody: "abc",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), fmt.Sprintf("/api/v1/clusters/%s/startup_config_tpl/%s", "1", tc.configType))
				rw.WriteHeader(tc.code)
				if tc.body != "" {
					_, err := rw.Write([]byte(tc.body))
					assert.NoError(t, err, "send mock response")
				}
			}))

			defer server.Close()

			api, err := newClient(server.URL, "test-token")
			assert.NoError(t, err, "checking new cloud api client")

			data, err := api.GetStartupConfig(1, tc.configType)

			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "checking error reason")
			} else {
				assert.NoError(t, err, "checking error")
				assert.Equal(t, tc.expectedBody, data, "check the startup config")
			}
		})
	}
}
