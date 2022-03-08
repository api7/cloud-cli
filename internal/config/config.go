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

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	configFileLocation string
)

func init() {
	configFileLocation = filepath.Join(os.Getenv("HOME"), ".api7cloud/config")
}

// Save config to file for persistence
func Save(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(configFileLocation)
	if _, err = os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return fmt.Errorf("failed to create config directory in $HOME/.api7cloud: %s", err)
		}
	}

	file, err := os.Create(configFileLocation)
	if err != nil {
		return fmt.Errorf("failed create file in $HOME/.api7cloud/config for credential: %s", err)
	}

	write, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("failed write credential to %s, %s", configFileLocation, err)
	}
	if write != len(data) {
		return fmt.Errorf("failed write credential to %s", configFileLocation)
	}

	return nil
}

// Load config from file
func Load() (*Config, error) {
	file, err := os.Open(configFileLocation)
	if err != nil {
		return nil, err
	}

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file, %s", err)
	}

	return &config, nil
}
