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

package apisix

import (
	"os"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// MergeConfig merge the user customized config with default settings.
func MergeConfig(config []byte, defaultConfig []byte) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if config != nil {
		if err := yaml.Unmarshal(config, &data); err != nil {
			return nil, errors.Wrap(err, "unmarshal config")
		}
	}
	defaultData := make(map[string]interface{})
	if defaultConfig != nil {
		if err := yaml.Unmarshal(defaultConfig, &defaultData); err != nil {
			return nil, errors.Wrap(err, "unmarshal default config")
		}
	}

	if err := mergo.Merge(&data, defaultData, mergo.WithOverride); err != nil {
		return nil, err
	}
	return data, nil
}

// SaveConfigToTemp saves the config to the temporary file and return its name.
func SaveConfigToTemp(config map[string]interface{}, pattern string) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", errors.Wrap(err, "marshal config")
	}
	tempFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	// Assign permission for other users so that processes inside container can read it.
	if err := tempFile.Chmod(0644); err != nil {
		return "", err
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(data); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

// SaveConfig saves the config to the specified file.
func SaveConfig(config map[string]interface{}, filepath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "marshal config")
	}
	if err = os.WriteFile(filepath, data, 0644); err != nil {
		return err
	}
	// Assign permission for other users so that processes inside container can read it.
	if err = os.Chmod(filepath, 0644); err != nil {
		return err
	}
	return nil
}
