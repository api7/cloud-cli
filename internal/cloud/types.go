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
	"os"

	"github.com/api7/cloud-cli/internal/types"
)

const (
	DefaultCloudApiServer       = "console.api7.cloud"
	DefaultCloudApiServerScheme = "https"
)

// API warp API7 Cloud REST API
type API interface {
	// Me returns the current user information
	Me() (*types.User, error)
	// ListControlPlanes returns the list of control planes in organization
	ListControlPlanes(orgID string) ([]*types.ControlPlaneSummary, error)
}

type api struct {
	apiServer   string
	scheme      string
	accessToken string
	httpClient  *http.Client
}

func New(accessToken string) API {
	apiServer := os.Getenv("CLOUD_API_SERVER")
	if apiServer == "" {
		apiServer = DefaultCloudApiServer
	}

	scheme := os.Getenv("CLOUD_API_SERVER_SCHEME")
	if scheme == "" {
		scheme = DefaultCloudApiServerScheme
	}

	return &api{
		apiServer:   apiServer,
		scheme:      scheme,
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}
}
