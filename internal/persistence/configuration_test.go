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
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSaveConfiguration(t *testing.T) {
	id := uuid.NewString()
	configDir = fmt.Sprintf("/tmp/%s/config", id)
	err := SaveConfiguration(&CloudConfiguration{
		DefaultProfile: "prod",
		Profiles: []Profile{
			{
				Name:    "dev",
				Address: "https://api.api7.cloud",
				User:    User{AccessToken: "dev-token"},
			},
		},
	})
	assert.Contains(t, err.Error(), "default profile not found", "save configuration should be failed")

	err = SaveConfiguration(&CloudConfiguration{
		DefaultProfile: "prod",
		Profiles: []Profile{
			{
				Name:    "prod",
				Address: "https://api.api7.cloud",
				User:    User{AccessToken: "prod-token"},
			},
			{
				Name:    "dev",
				Address: "https://api.api7.cloud",
				User:    User{AccessToken: "dev-token"},
			},
		},
	})
	assert.NoError(t, err, "save to file in not exist dir should be success")

	config, err := LoadConfiguration()
	assert.NoError(t, err, "load from file should be success")

	profile, err := config.GetDefaultProfile()
	assert.NoError(t, err, "get default profile")
	assert.Equal(t, "prod-token", profile.User.AccessToken, "access token should be prod-token")

	profile, err = config.GetProfile("dev")
	assert.NoError(t, err, "get dev profile")
	assert.Equal(t, "dev-token", profile.User.AccessToken, "access token should be dev-token")

	id = uuid.NewString()
	err = os.MkdirAll(fmt.Sprintf("/tmp/%s", id), fs.ModePerm)
	assert.NoError(t, err, "create dir should be success")

	configDir = fmt.Sprintf("/tmp/%s/config", id)
	err = SaveConfiguration(&CloudConfiguration{
		DefaultProfile: "prod",
		Profiles: []Profile{
			{
				Name:    "prod",
				Address: "https://api.api7.cloud",
				User:    User{AccessToken: "prod-token-2"},
			},
			{
				Name:    "dev",
				Address: "https://api.api7.cloud",
				User:    User{AccessToken: "dev-token-2"},
			},
		},
	})
	assert.NoError(t, err, "save to file in exist dir should be success")

	config, err = LoadConfiguration()
	assert.NoError(t, err, "load from file should be success")

	profile, err = config.GetDefaultProfile()
	assert.NoError(t, err, "get default profile")
	assert.Equal(t, "prod-token-2", profile.User.AccessToken, "access token should be prod-token-2")
}

func TestLoad(t *testing.T) {
	id := uuid.NewString()
	configDir = fmt.Sprintf("/tmp/%s/config", id)

	_, err := LoadConfiguration()
	assert.Contains(t, err.Error(), "no such file or directory", "load from file should be failed")

	dir := filepath.Dir(configDir)
	err = os.MkdirAll(dir, 0750)
	assert.NoError(t, err, "create dir should be success")

	err = os.WriteFile(configDir, []byte("invalid configuration"), fs.ModePerm)
	assert.NoError(t, err, "write fake configuration file should be success")

	_, err = LoadConfiguration()
	assert.Contains(t, err.Error(), "failed to decode configuration", "load from file should be failed")
}
