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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/api7/cloud-go-sdk"
	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/types"
)

func (a *api) Me() (*cloud.User, error) {
	return a.sdk.Me(context.TODO())
}

func (a *api) ListControlPlanes(orgID cloud.ID) ([]*cloud.ControlPlane, error) {
	var controlPlanes []*cloud.ControlPlane

	iter, err := a.sdk.ListControlPlanes(context.TODO(), &cloud.ResourceListOptions{
		Organization: &cloud.Organization{
			ID: orgID,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create control plane iterator")
	}

	for {
		cp, err := iter.Next()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get next control plane")
		}
		if cp == nil {
			return controlPlanes, nil
		}
		controlPlanes = append(controlPlanes, cp)
	}
}

func (a *api) GetTLSBundle(cpID cloud.ID) (*cloud.TLSBundle, error) {
	return a.sdk.GenerateGatewaySideCertificate(context.TODO(), cpID, nil)
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

func (a *api) GetStartupConfig(cpID cloud.ID, configType StartupConfigType) (string, error) {
	var response types.ControlPlaneStartupConfigResponsePayload

	if err := a.makeGetRequest(&url.URL{
		Path: fmt.Sprintf("/api/v1/controlplanes/%s/startup_config_tpl/%s", cpID.String(), configType),
	}, &response); err != nil {
		return "", err
	}
	return response.Configuration, nil
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

func (a *api) GetDefaultControlPlane() (*cloud.ControlPlane, error) {
	user, err := a.Me()
	if err != nil {
		return nil, errors.Wrap(err, "failed to access user information")
	}

	if len(user.OrgIDs) == 0 {
		return nil, errors.New("incomplete user information, no organization")
	}
	iter, err := a.sdk.ListControlPlanes(context.TODO(), &cloud.ResourceListOptions{
		Organization: &cloud.Organization{
			ID: user.OrgIDs[0],
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create control plane iterator")
	}

	// Let's just fetch the first control plane.
	cp, err := iter.Next()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default control plane")
	}
	if cp == nil {
		return nil, errors.New("no control plane available")
	}

	return cp, nil
}

func (a *api) makeGetRequest(u *url.URL, response interface{}) error {
	req, err := a.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}

	return decodeResponse(resp, response)
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

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

	return request, nil
}

func decodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	output.Verbosef("Receiving response:\n%s", responseDump)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= http.StatusInternalServerError {
			return errors.New("Server internal error, please try again later")
		}

		if resp.StatusCode == http.StatusNotFound {
			return errors.New("Resource not found")
		}

		var rw types.ResponseWrapper
		err := json.NewDecoder(resp.Body).Decode(&rw)
		if err != nil {
			return errors.Wrap(err, "Got a malformed response from server")
		}
		return errors.New(fmt.Sprintf("Error Code: %d, Error Reason: %s", rw.Status.Code, rw.ErrorReason))
	}
	var rw types.ResponseWrapper
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	err = dec.Decode(&rw)
	if err != nil {
		return errors.Wrap(err, "Got a malformed response from server")
	}

	if v != nil {
		data, err := json.Marshal(rw.Payload)
		if err != nil {
			return errors.Wrap(err, "Got a malformed response from server")
		}

		err = json.Unmarshal(data, v)
		if err != nil {
			return errors.Wrap(err, "Got a malformed response from server")
		}
	}
	return nil
}
