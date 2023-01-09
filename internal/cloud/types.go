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
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/api7/cloud-go-sdk"
	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/consts"
)

var (
	DefaultClient API
)

// InitDefaultClient initializes the default client with the given configuration
func InitDefaultClient(cloudAddr, accessToken string) (err error) {
	if DefaultClient != nil {
		return nil
	}
	DefaultClient, err = newClient(cloudAddr, accessToken)
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
	defaultApi7CloudLuaModuleURL = "https://github.com/api7/cloud-scripts/raw/main/assets/cloud_module_beta.tar.gz"
)

// StartupConfigType is type of gateway startup config
type StartupConfigType string

const (
	APISIX StartupConfigType = "apisix"
	HELM   StartupConfigType = "helm"
)

// API warp API7 Cloud REST API
type API interface {
	// Me returns the current user information
	Me() (*cloud.User, error)
	// ListControlPlanes returns the list of control planes in organization
	ListControlPlanes(orgID cloud.ID) ([]*cloud.ControlPlane, error)
	// GetTLSBundle gets the tls bundle used to communicate with API7 Cloud. returns the control plane with the given ID
	GetTLSBundle(cpID cloud.ID) (*cloud.TLSBundle, error)
	// GetCloudLuaModule returns the Cloud Lua code (in the tar.gz format)
	GetCloudLuaModule() ([]byte, error)
	// GetStartupConfig gets the startup configuration from API7 Cloud for deploy APISIX by specify config type.
	GetStartupConfig(cpID cloud.ID, configType StartupConfigType) (string, error)
	// GetDefaultOrganization returns the default organization for the current user.
	GetDefaultOrganization() (*cloud.Organization, error)
	// GetDefaultControlPlane returns the default control plane for the current organization.
	GetDefaultControlPlane() (*cloud.ControlPlane, error)
	// DebugShowConfig returns the translated Apache APISIX object with the given API7 Cloud resource type and id.
	DebugShowConfig(cpID cloud.ID, resource string, id string) (string, error)
}

type api struct {
	sdk               cloud.Interface
	host              string
	scheme            string
	accessToken       string
	httpClient        *http.Client
	cloudLuaModuleURL *url.URL
}

var (
	// NewClient is a function to create a new API7 Cloud API Client
	NewClient = newClient
)

// newClient returns a new API7 Cloud API Client
func newClient(cloudAddr, accessToken string) (API, error) {
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

	sdk, err := cloud.NewInterface(&cloud.Options{
		ServerAddr:  cloudAddr,
		Token:       accessToken,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "initialize cloud go sdk")
	}

	return &api{
		sdk:               sdk,
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
