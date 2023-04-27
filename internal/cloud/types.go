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

package cloud

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/api7/cloud-go-sdk"
	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/utils"
)

var (
	DefaultClient API
)

// InitDefaultClient initializes the default client with the given configuration
func InitDefaultClient(cloudAddr, accessToken string, trace bool) (err error) {
	if DefaultClient != nil {
		return nil
	}
	DefaultClient, err = newClient(cloudAddr, accessToken, trace)
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
	// ListClusters returns the list of clusters in organization
	ListClusters(orgID cloud.ID, limit int, skip int) ([]*cloud.Cluster, error)
	// GetTLSBundle gets the tls bundle used to communicate with API7 Cloud. returns the cluster with the given ID
	GetTLSBundle(clusterID cloud.ID) (*cloud.TLSBundle, error)
	// GetCloudLuaModule returns the Cloud Lua code (in the tar.gz format)
	GetCloudLuaModule() ([]byte, error)
	// GetStartupConfig gets the startup configuration from API7 Cloud for deploy APISIX by specify config type.
	GetStartupConfig(clusterID cloud.ID, configType StartupConfigType) (string, error)
	// GetDefaultOrganization returns the default organization for the current user.
	GetDefaultOrganization() (*cloud.Organization, error)
	// GetDefaultCluster returns the default cluster for the current organization.
	GetDefaultCluster() (*cloud.Cluster, error)
	// GetClusterDetail returns the detail cluster for the specify cluster.
	GetClusterDetail(clusterID cloud.ID) (*cloud.Cluster, error)
	// GetSSL returns the detail of the Certificate (SSL) object.
	GetSSL(clusterID, sslID cloud.ID) (*cloud.CertificateDetails, error)
	// DeleteSSL deletes the specified SSL object.
	DeleteSSL(clusterID, sslID cloud.ID) error
	// ListSSL lists up to *limits* SSL objects in the specified cluster, and it'll
	// skip the first *skip* objects.
	ListSSL(clusterID cloud.ID, limit int, skip int) ([]*cloud.CertificateDetails, error)
	// CreateSSL creates an SSL object according to the given spec.
	CreateSSL(clusterID cloud.ID, ssl *cloud.Certificate) (*cloud.CertificateDetails, error)
	// DebugShowConfig returns the translated Apache APISIX object with the given API7 Cloud resource type and id.
	DebugShowConfig(clusterID cloud.ID, resource string, id cloud.ID) (string, error)
	// ListServices return the list of services in application
	ListServices(clusterID cloud.ID, limit int, skip int) ([]*cloud.Application, error)
	// UpdateService return the configuration after the service update
	UpdateService(clusterID cloud.ID, config string) (*cloud.Application, error)
	// GetService return the service in line with id in application
	GetService(clusterID cloud.ID, appID cloud.ID) (*cloud.Application, error)
	// DeleteService return the service delete success or fail
	DeleteService(clusterID cloud.ID, appID cloud.ID) error
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
func newClient(cloudAddr, accessToken string, trace bool) (API, error) {
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
		ServerAddr:      cloudAddr,
		Token:           accessToken,
		DialTimeout:     5 * time.Second,
		EnableHTTPTrace: trace,
	})
	if err != nil {
		return nil, errors.Wrap(err, "initialize cloud go sdk")
	}

	if trace {
		go utils.VerboseGoroutine(sdk.TraceChan())
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
