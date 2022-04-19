//  Copyright 2022 API7.ai, Inc under one or more contributor license
//  agreements.  See the NOTICE file distributed with this work for
//  additional information regarding copyright ownership.
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

package apisix

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		config      []byte
		essential   []byte
		result      map[string]interface{}
		errorReason string
	}{
		{
			name: "bad config",
			config: []byte(`
apisix123
`),
			essential: []byte(`
apisix:
  enable_admin: false
  ssl:
    ssl_trusted_certificate: /path/to/ca.crt
  lua_module_hook: cloud
  extra_lua_path: /opt/cloud_module_beta/?.ljbc;
nginx_config:
  http:
    custom_lua_shared_dict:
      cloud: 1m
etcd:
  host:
    - "https://default.foo.api7.cloud:443"
  tls:
    cert: /cloud/tls/tls.crt
    key: /cloud/tls/tls.key
    sni: default.foo.api7.cloud
    verify: true
`),

			errorReason: "unmarshal config",
		},
		{
			name:        "bad default config",
			config:      []byte(`apisix: 123`),
			essential:   []byte(`apisix123`),
			errorReason: "unmarshal default config",
		},
		{
			name: "success",
			config: []byte(`
apisix:
  enable_admin: true
  ssl:
    ssl_trusted_certificate: /path/to/ca.crt
`),
			essential: []byte(`
apisix:
  enable_admin: false
  ssl:
    ssl_trusted_certificate: /path/to/ca2.crt
`),
			result: map[string]interface{}{
				"apisix": map[string]interface{}{
					"enable_admin": false,
					"ssl": map[string]interface{}{
						"ssl_trusted_certificate": "/path/to/ca2.crt",
					},
				},
			},
		},
		{
			name:   "empty config",
			config: nil,
			essential: []byte(`
apisix:
  enable_admin: false
  ssl:
    ssl_trusted_certificate: /path/to/ca2.crt
`),
			result: map[string]interface{}{
				"apisix": map[string]interface{}{
					"enable_admin": false,
					"ssl": map[string]interface{}{
						"ssl_trusted_certificate": "/path/to/ca2.crt",
					},
				},
			},
		},
		{
			name: "empty default config",
			config: []byte(`
apisix:
  enable_admin: false
  ssl:
    ssl_trusted_certificate: /path/to/ca2.crt
`),
			result: map[string]interface{}{
				"apisix": map[string]interface{}{
					"enable_admin": false,
					"ssl": map[string]interface{}{
						"ssl_trusted_certificate": "/path/to/ca2.crt",
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := MergeConfig(tc.config, tc.essential)
			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "check error reason")
			} else {
				assert.Nil(t, err, "check if reason is nil")
				json1, _ := json.Marshal(tc.result)
				json2, _ := json.Marshal(result)
				assert.JSONEq(t, string(json1), string(json2), "check result")
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		config          map[string]interface{}
		filenamePattern string
		errorReason     string
	}{
		{
			name: "success",
			config: map[string]interface{}{
				"apisix": map[string]interface{}{
					"enable_admin": false,
					"ssl": map[string]interface{}{
						"ssl_trusted_certificate": "/path/to/ca2.crt",
					},
				},
			},
			filenamePattern: `apisix-config-\d+\.yaml`,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			filename, err := SaveConfig(tc.config, "apisix-config-*.yaml")
			if tc.errorReason != "" {
				assert.NotNil(t, err, "check if err is not nil")
				assert.Equal(t, tc.errorReason, err.Error(), "check error")
			} else {
				assert.Nil(t, err, "check if err is nil")
				assert.Regexp(t, regexp.MustCompile(tc.filenamePattern), filename, "check filename")
			}
		})
	}
}
