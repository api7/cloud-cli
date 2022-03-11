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

package cloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/types"
)

func TestMe(t *testing.T) {
	tests := []struct {
		name      string
		code      int
		body      string
		want      *types.User
		wantErr   bool
		errReason string
	}{
		{
			name:      "server internal error",
			code:      http.StatusInternalServerError,
			want:      nil,
			wantErr:   true,
			errReason: "Server internal error, please try again later",
		},
		{
			name:      "http not found",
			code:      http.StatusNotFound,
			want:      nil,
			wantErr:   true,
			errReason: "Resource not found",
		},
		{
			name:      "malformed json",
			code:      http.StatusOK,
			body:      `invalid json`,
			want:      nil,
			wantErr:   true,
			errReason: "Got a malformed response from server",
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
			want: &types.User{
				ID:         "321",
				FirstName:  "first",
				LastName:   "last",
				Email:      "demo@api7.ai",
				OrgIDs:     []string{"123"},
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

			err := os.Setenv(consts.Api7CloudAddrEnv, server.URL)
			assert.NoError(t, err, "checking env setup")

			api, err := newClient("test-token")
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

func TestListControlPlanes(t *testing.T) {
	tests := []struct {
		name      string
		orgID     string
		code      int
		body      string
		want      *types.ControlPlaneSummary
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
			want: &types.ControlPlaneSummary{
				ControlPlane: types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID:   "392306215415186327",
						Name: "default",
					},
					OrganizationID: "392306215398409111",
					RegionID:       "56523356",
					Status:         3,
					Domain:         "default.xvlbzgs4bqbjdmybyqk.api7.cloud",
					ConfigPayload:  "",
				},
				OrgName: "XVlBzg",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), fmt.Sprintf("/api/v1/orgs/%s/controlplanes", tt.orgID))
				assert.Equal(t, req.Header.Get("Authorization"), "Bearer test-token")

				rw.WriteHeader(tt.code)
				if tt.body != "" {
					write, err := rw.Write([]byte(tt.body))
					assert.NoError(t, err, "send mock response")
					assert.Equal(t, len(tt.body), write, "write mock response")
				}
			}))

			defer server.Close()

			err := os.Setenv(consts.Api7CloudAddrEnv, server.URL)
			assert.NoError(t, err, "checking env setup")

			api, err := newClient("test-token")
			assert.NoError(t, err, "checking new cloud api client")

			result, err := api.ListControlPlanes(tt.orgID)

			if tt.wantErr {
				assert.Contains(t, err.Error(), tt.errReason, "checking error reason")
			} else {
				assert.NoError(t, err, "checking error")
				assert.Len(t, result, 1, "checking control planes count")
				assert.Equal(t, tt.want, result[0], "checking result")
			}
		})
	}
}

func TestGetTLSBundle(t *testing.T) {
	tests := []struct {
		name      string
		cpID      string
		code      int
		body      string
		want      *types.TLSBundle
		wantErr   bool
		errReason string
	}{
		{
			name: "success",
			code: http.StatusOK,
			cpID: "1",
			want: &types.TLSBundle{
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
			cpID:      "1",
			errReason: "Server internal error, please try again later",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, req.URL.String(), fmt.Sprintf("/api/v1/controlplanes/%s/dp_certificate", tt.cpID))
				assert.Equal(t, req.Header.Get("Authorization"), "Bearer test-token")

				rw.WriteHeader(tt.code)
				if tt.body != "" {
					write, err := rw.Write([]byte(tt.body))
					assert.NoError(t, err, "send mock response")
					assert.Equal(t, len(tt.body), write, "write mock response")
				}
			}))

			defer server.Close()

			err := os.Setenv(consts.Api7CloudAddrEnv, server.URL)
			assert.NoError(t, err, "checking env setup")

			api, err := newClient("test-token")
			assert.NoError(t, err, "checking new cloud api client")

			bundle, err := api.GetTLSBundle(tt.cpID)

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

			api, err := newClient("test-token")
			assert.NoError(t, err, "checking new cloud api client")

			data, err := api.GetCloudLuaModule()

			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "checking error reason")
			} else {
				assert.NoError(t, err, "checking error")
				assert.Equal(t, tc.body, string(data), "check the tls bundle")
			}
		})
	}
}
