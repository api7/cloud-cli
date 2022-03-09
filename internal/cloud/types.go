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
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/types"
)

const (
	api7CloudAddr        = "API&_CLOUD_ADDR"
	DefaultApi7CloudAddr = "https://console.api7.cloud"
)

// API warp API7 Cloud REST API
type API interface {
	// Me returns the current user information
	Me() (*types.User, error)
	// ListControlPlanes returns the list of control planes in organization
	ListControlPlanes(orgID string) ([]*types.ControlPlaneSummary, error)
}

type api struct {
	host        string
	scheme      string
	accessToken string
	httpClient  *http.Client
}

// New returns a new API7 Cloud API Client
func New(accessToken string) (API, error) {
	cloudAddr := os.Getenv(api7CloudAddr)
	if cloudAddr == "" {
		cloudAddr = DefaultApi7CloudAddr
	}

	u, err := url.Parse(cloudAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse API7 Cloud server URL")
	}

	if u.Host == "" || u.Scheme == "" {
		return nil, errors.New("invalid API7 Cloud server URL")
	}

	return &api{
		host:        u.Host,
		scheme:      u.Scheme,
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}, nil
}
