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

package cloud

import (
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/types"
)

var (
	DefaultClient API
)

// InitDefaultClient initializes the default client with the given credentials
func InitDefaultClient(accessToken string) (err error) {
	if DefaultClient != nil {
		return nil
	}
	DefaultClient, err = newClient(accessToken)
	return
}

// Client return default client to access API7 Cloud API
func Client() API {
	if DefaultClient == nil {
		panic("default client is not initialized")
	}
	return DefaultClient
}

const (
	defaultApi7CloudAddr         = "https://console.api7.cloud"
	defaultApi7CloudLuaModuleURL = "https://github.com/api7/cloud-scripts/raw/main/assets/cloud_module_beta.tar.gz"
)

// StartupConfigType is type of data plane startup config
type StartupConfigType string

const (
	APISIX StartupConfigType = "apisix"
	HELM   StartupConfigType = "helm"
)

// API warp API7 Cloud REST API
type API interface {
	// Me returns the current user information
	Me() (*types.User, error)
	// ListControlPlanes returns the list of control planes in organization
	ListControlPlanes(orgID string) ([]*types.ControlPlaneSummary, error)
	// GetTLSBundle gets the tls bundle used to communicate with API7 Cloud. returns the control plane with the given ID
	GetTLSBundle(cpID string) (*types.TLSBundle, error)
	// GetCloudLuaModule returns the Cloud Lua code (in the tar.gz format)
	GetCloudLuaModule() ([]byte, error)
	// GetStartupConfig gets the startup configuration from API7 Cloud for deploy APISIX by specify config type.
	GetStartupConfig(cpID string, configType StartupConfigType) (string, error)
	// GetDefaultControlPlane returns the default control plane for the current organization.
	GetDefaultControlPlane() (*types.ControlPlane, error)
}

type api struct {
	host              string
	scheme            string
	accessToken       string
	httpClient        *http.Client
	cloudLuaModuleURL *url.URL
}

// newClient returns a new API7 Cloud API Client
func newClient(accessToken string) (API, error) {
	cloudAddr := os.Getenv(consts.Api7CloudAddrEnv)
	if cloudAddr == "" {
		cloudAddr = defaultApi7CloudAddr
	}

	u, err := url.Parse(cloudAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse API7 Cloud server URL")
	}

	if u.Host == "" || u.Scheme == "" {
		return nil, errors.New("invalid API7 Cloud server URL")
	}

	rawCloudModuleURL := os.Getenv(consts.Api7CloudLuaModuleURL)
	if rawCloudModuleURL == "" {
		rawCloudModuleURL = defaultApi7CloudLuaModuleURL
	}
	cloudModuleURL, err := url.Parse(rawCloudModuleURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse API7 Cloud Lua module URL")
	}
	if cloudModuleURL.Host == "" || cloudModuleURL.Scheme == "" {
		return nil, errors.New("invalid API7 Cloud Lua Module URL")
	}

	return &api{
		host:              u.Host,
		scheme:            u.Scheme,
		cloudLuaModuleURL: cloudModuleURL,
		accessToken:       accessToken,
		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}, nil
}
