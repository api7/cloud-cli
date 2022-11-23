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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/types"
)

func (a *api) Me() (*types.User, error) {
	var user types.User

	if err := a.makeGetRequest(&url.URL{
		Path: "/api/v1/user/me",
	}, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *api) ListControlPlanes(orgID string) ([]*types.ControlPlaneSummary, error) {
	var response types.GetOrganizationControlPlanesResponsePayload

	if err := a.makeGetRequest(&url.URL{
		Path: fmt.Sprintf("/api/v1/orgs/%s/controlplanes", orgID),
	}, &response); err != nil {
		return nil, err
	}

	return response.List, nil
}

func (a *api) GetTLSBundle(cpID string) (*types.TLSBundle, error) {
	var bundle types.TLSBundle

	if err := a.makeGetRequest(&url.URL{
		Path: fmt.Sprintf("/api/v1/controlplanes/%s/dp_certificate", cpID),
	}, &bundle); err != nil {
		return nil, err
	}
	return &bundle, nil
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

func (a *api) GetStartupConfig(cpID string, configType StartupConfigType) (string, error) {
	var response types.ControlPlaneStartupConfigResponsePayload

	if err := a.makeGetRequest(&url.URL{
		Path: fmt.Sprintf("/api/v1/controlplanes/%s/startup_config_tpl/%s", cpID, configType),
	}, &response); err != nil {
		return "", err
	}
	return response.Configuration, nil
}

func (a *api) GetDefaultOrganization() (*types.Organization, error) {
	user, err := a.Me()
	if err != nil {
		return nil, errors.Wrap(err, "failed to access user information")
	}

	if len(user.OrgIDs) == 0 {
		return nil, errors.New("incomplete user information, no organization")
	}

	var org types.Organization
	if err := a.makeGetRequest(&url.URL{
		Path: fmt.Sprintf("/api/v1/orgs/%s", user.OrgIDs[0]),
	}, &org); err != nil {
		return nil, err
	}
	return &org, nil
}

func (a *api) GetDefaultControlPlane() (*types.ControlPlane, error) {
	user, err := a.Me()
	if err != nil {
		return nil, errors.Wrap(err, "failed to access user information")
	}

	if len(user.OrgIDs) == 0 {
		return nil, errors.New("incomplete user information, no organization")
	}
	controlPlanes, err := a.ListControlPlanes(user.OrgIDs[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to list control planes")
	}
	if len(controlPlanes) == 0 {
		return nil, errors.New("no control plane available")
	}
	return &controlPlanes[0].ControlPlane, nil
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
