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

package persistence

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/output"
)

func init() {
	if err := Init(); err != nil {
		panic(err)
	}
}

// SaveCredential to file for persistence
func SaveCredential(credential *Credential) error {
	data, err := yaml.Marshal(credential)
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(credentialDir)
	if _, err = os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return fmt.Errorf("failed to create config directory in %s: %s", dir, err)
		}
	}

	file, err := os.Create(credentialDir)
	if err != nil {
		return fmt.Errorf("failed to create file in %s for credential: %s", credentialDir, err)
	}

	write, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write credential to %s, %s", credentialDir, err)
	}
	if write != len(data) {
		return fmt.Errorf("failed to write credential to %s", credentialDir)
	}

	return nil
}

// LoadCredential from file
func LoadCredential() (*Credential, error) {
	file, err := os.Open(credentialDir)
	if err != nil {
		return nil, err
	}

	var credential Credential
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&credential)
	if err != nil {
		return nil, fmt.Errorf("failed to decode credential, %s", err)
	}

	return &credential, nil
}

// CheckAccessTokenAndInitCloudClient checks if cloud-cli configured the token correctly.
// Then use this token to initialize the cloud client.
func CheckAccessTokenAndInitCloudClient() error {
	accessToken := os.Getenv(consts.Api7CloudAccessTokenEnv)
	if accessToken == "" {
		credential, err := LoadCredential()
		if err != nil {
			return fmt.Errorf("Failed to load credential: %s,\nPlease run 'cloud-cli configure' first, access token can be created from https://console.api7.cloud", err)
		}
		accessToken = credential.User.AccessToken
	}

	output.Verbosef("Loaded access token: %s", accessToken)

	if err := cloud.InitDefaultClient(accessToken); err != nil {
		return fmt.Errorf("Failed to init api7 cloud client: %s", err)
	}
	return nil
}
