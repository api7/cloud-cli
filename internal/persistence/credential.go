//  Copyright 2022 API7.ai, Inc under one or more contributor license
//  agreements.  See the NOTICE file distributed with this work for
//  additional information regarding copyright ownership.
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

package persistence

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	credentialFileLocation string
)

func init() {
	credentialFileLocation = filepath.Join(os.Getenv("HOME"), ".api7cloud/credentials")
}

// SaveCredential to file for persistence
func SaveCredential(credential *Credential) error {
	data, err := yaml.Marshal(credential)
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(credentialFileLocation)
	if _, err = os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			return fmt.Errorf("failed to create config directory in %s: %s", dir, err)
		}
	}

	file, err := os.Create(credentialFileLocation)
	if err != nil {
		return fmt.Errorf("failed create file in %s for credential: %s", credentialFileLocation, err)
	}

	write, err := file.Write(data)
	if err != nil {
		return fmt.Errorf("failed write credential to %s, %s", credentialFileLocation, err)
	}
	if write != len(data) {
		return fmt.Errorf("failed write credential to %s", credentialFileLocation)
	}

	return nil
}

// LoadCredential from file
func LoadCredential() (*Credential, error) {
	file, err := os.Open(credentialFileLocation)
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
