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
	"github.com/api7/cloud-cli/internal/options"
)

func init() {
	if err := Init(); err != nil {
		panic(err)
	}
}

// SaveConfiguration to file for persistence
func SaveConfiguration(config *CloudConfiguration) error {
	if err := config.Validate(); err != nil {
		return err
	}
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(configDir)
	if _, err = os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return fmt.Errorf("failed to create config directory in %s: %s", dir, err)
		}
	}

	file, err := os.Create(configDir)
	if err != nil {
		return fmt.Errorf("failed to create file in %s for api7 cloud configuration: %s", configDir, err)
	}

	write, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write configuration to %s, %s", configDir, err)
	}
	if write != len(data) {
		return fmt.Errorf("failed to write configuration to %s, truncated write", configDir)
	}

	return nil
}

// LoadConfiguration from file
func LoadConfiguration() (*CloudConfiguration, error) {
	file, err := os.Open(configDir)
	if err != nil {
		return nil, err
	}

	var config CloudConfiguration
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode configuration, %s", err)
	}

	return &config, nil
}

// CheckConfigurationAndInitCloudClient checks if cloud-cli configured the server address and token correctly.
// Then use this token to initialize the cloud client.
func CheckConfigurationAndInitCloudClient() error {
	configuration, err := LoadConfiguration()
	if err != nil {
		return fmt.Errorf("Failed to load configuration: %s,\nPlease run 'cloud-cli configure' first, access token can be created from API7 WEB Console.", err)
	}

	profileName := os.Getenv(consts.Api7CloudProfile)
	if profileName == "" {
		profileName = options.Global.Profile
	}
	if profileName == "" {
		profileName = configuration.DefaultProfile
	}

	profile, err := configuration.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("Failed to get %s profile, Please check your configuration file: %s", profileName, configDir)
	}

	if err := cloud.InitDefaultClient(profile.Address, profile.User.AccessToken); err != nil {
		return fmt.Errorf("Failed to init api7 cloud client: %s", err)
	}
	return nil
}
