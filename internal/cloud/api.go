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
	"context"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-go-sdk"
	"github.com/pkg/errors"
	"io"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (a *api) Me() (*cloud.User, error) {
	return a.sdk.Me(context.TODO())
}

func (a *api) ListClusters(orgID cloud.ID, limit int, skip int) ([]*cloud.Cluster, error) {
	var clusters []*cloud.Cluster
	pageSize := limit
	page := int(math.Floor(float64(skip)/float64(limit))) + 1
	start := skip % limit
	end := skip%limit + limit

	iter, err := a.sdk.ListClusters(context.TODO(), &cloud.ResourceListOptions{
		Organization: &cloud.Organization{
			ID: orgID,
		},
		Pagination: &cloud.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cluster iterator")
	}

	for {
		cluster, err := iter.Next()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get next cluster")
		}
		if cluster == nil || len(clusters) == end {
			if len(clusters) <= start {
				return clusters, nil
			}
			return clusters[start:], nil
		}

		clusters = append(clusters, cluster)
	}
}

func (a *api) GetTLSBundle(clusterID cloud.ID) (*cloud.TLSBundle, error) {
	return a.sdk.GenerateGatewaySideCertificate(context.TODO(), clusterID, nil)
}

func (a *api) GetCloudLuaModule() ([]byte, error) {
	req, err := a.newRequest(http.MethodGet, a.cloudLuaModuleURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response code: %d, message: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func (a *api) GetStartupConfig(clusterID cloud.ID, configType StartupConfigType) (string, error) {
	config, err := a.sdk.GetGatewayInstanceStartupConfigTemplate(context.TODO(), clusterID, string(configType), nil)
	if err != nil {
		return "", err
	}
	return config, nil
}

func (a *api) GetDefaultOrganization() (*cloud.Organization, error) {
	user, err := a.Me()
	if err != nil {
		return nil, errors.Wrap(err, "failed to access user information")
	}

	if len(user.OrgIDs) == 0 {
		return nil, errors.New("incomplete user information, no organization")
	}

	return a.sdk.GetOrganization(context.TODO(), user.OrgIDs[0], nil)
}

func (a *api) GetDefaultCluster() (*cloud.Cluster, error) {
	user, err := a.Me()
	if err != nil {
		return nil, errors.Wrap(err, "failed to access user information")
	}

	if len(user.OrgIDs) == 0 {
		return nil, errors.New("incomplete user information, no organization")
	}
	iter, err := a.sdk.ListClusters(context.TODO(), &cloud.ResourceListOptions{
		Organization: &cloud.Organization{
			ID: user.OrgIDs[0],
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cluster iterator")
	}

	// Let's just fetch the first cluster.
	cluster, err := iter.Next()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default cluster")
	}
	if cluster == nil {
		return nil, errors.New("no cluster available")
	}

	return cluster, nil
}

func (a *api) GetClusterDetail(clusterID cloud.ID) (*cloud.Cluster, error) {
	user, err := a.Me()
	if err != nil {
		return nil, errors.Wrap(err, "failed to access user information")
	}
	cluster, err := a.sdk.GetCluster(context.TODO(), clusterID, &cloud.ResourceGetOptions{
		Organization: &cloud.Organization{
			ID: user.OrgIDs[0],
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cluster detail")
	}

	return cluster, nil
}

func (a *api) GetSSL(clusterID, sslID cloud.ID) (*cloud.CertificateDetails, error) {
	return a.sdk.GetCertificate(context.TODO(), sslID, &cloud.ResourceGetOptions{
		Cluster: &cloud.Cluster{
			ID: clusterID,
		},
	})
}

func (a *api) DeleteSSL(clusterID, sslID cloud.ID) error {
	return a.sdk.DeleteCertificate(context.TODO(), sslID, &cloud.ResourceDeleteOptions{
		Cluster: &cloud.Cluster{
			ID: clusterID,
		},
	})
}

func (a *api) ListSSL(clusterID cloud.ID, limit int, skip int) ([]*cloud.CertificateDetails, error) {
	var (
		ssl []*cloud.CertificateDetails
	)

	pageSize := 25
	firstPage := skip/pageSize + 1
	skipOnFirstPage := skip % pageSize

	iter, err := a.sdk.ListCertificates(context.TODO(), &cloud.ResourceListOptions{
		Cluster: &cloud.Cluster{
			ID: clusterID,
		},
		Pagination: &cloud.Pagination{
			Page:     firstPage,
			PageSize: pageSize,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get certificate iterator")
	}

	for {
		cert, err := iter.Next()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get next cert")
		}
		if cert == nil {
			// end
			break
		}

		if skipOnFirstPage > 0 {
			skipOnFirstPage--
			continue
		}
		ssl = append(ssl, cert)
		if len(ssl) >= limit {
			break
		}
	}
	return ssl, nil
}

func (a *api) CreateSSL(clusterID cloud.ID, ssl *cloud.Certificate) (*cloud.CertificateDetails, error) {
	return a.sdk.CreateCertificate(context.TODO(), ssl, &cloud.ResourceCreateOptions{
		Cluster: &cloud.Cluster{
			ID: clusterID,
		},
	})
}

func (a *api) UpdateSSL(clusterID cloud.ID, ssl *cloud.Certificate) (*cloud.CertificateDetails, error) {
	return a.sdk.UpdateCertificate(context.TODO(), ssl, &cloud.ResourceUpdateOptions{
		Cluster: &cloud.Cluster{
			ID: clusterID,
		},
	})
}

func (a *api) ListServices(clusterID cloud.ID, limit int, skip int) ([]*cloud.Application, error) {
	var services []*cloud.Application
	pageSize := limit
	page := int(math.Floor(float64(skip)/float64(limit))) + 1
	start := skip % limit
	end := skip%limit + limit

	iter, err := a.sdk.ListApplications(context.TODO(), &cloud.ResourceListOptions{
		Cluster: &cloud.Cluster{
			ID: clusterID,
		},
		Pagination: &cloud.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list service iterator")
	}

	for {
		service, err := iter.Next()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get next service")
		}
		if service == nil || len(services) == end {
			if len(services) <= start {
				return services, nil
			}
			return services[start:], nil
		}

		services = append(services, service)
	}
}

func (a *api) UpdateService(clusterID cloud.ID, svc *cloud.Application) (*cloud.Application, error) {
	newSvc, err := a.sdk.UpdateApplication(context.TODO(), svc,
		&cloud.ResourceUpdateOptions{
			Cluster: &cloud.Cluster{
				ID: clusterID,
			},
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update service")
	}

	return newSvc, nil
}

func (a *api) GetService(clusterID cloud.ID, appID cloud.ID) (*cloud.Application, error) {
	service, err := a.sdk.GetApplication(context.TODO(), appID, &cloud.ResourceGetOptions{
		Cluster: &cloud.Cluster{
			ID: clusterID,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get service iterator")
	}

	return service, nil
}

func (a *api) DeleteService(clusterID cloud.ID, appID cloud.ID) error {
	err := a.sdk.DeleteApplication(context.TODO(), appID, &cloud.ResourceDeleteOptions{
		Cluster: &cloud.Cluster{
			ID: clusterID,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete service")
	}
	return nil
}

func (a *api) CreateService(clusterID cloud.ID, svc *cloud.Application) (*cloud.Application, error) {
	newSvc, err := a.sdk.CreateApplication(context.TODO(), svc,
		&cloud.ResourceCreateOptions{
			Cluster: &cloud.Cluster{
				ID: clusterID,
			},
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create service")
	}

	return newSvc, nil
}

func (a *api) newRequest(method string, url *url.URL, body io.Reader) (*http.Request, error) {
	// Respect users' settings if host and scheme are not empty.
	if url.Host == "" {
		url.Host = a.host
	}
	if url.Scheme == "" {
		url.Scheme = a.scheme
	}

	request, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}

	requestDump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, err
	}
	output.Verbosef("Sending request:\n%s", requestDump)

	return request, nil
}
