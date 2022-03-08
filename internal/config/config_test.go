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
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	id := uuid.NewString()
	configFileLocation = fmt.Sprintf("/tmp/%s/config", id)
	err := Save(&Config{
		User: User{AccessToken: "test-0"},
	})
	assert.NoError(t, err, "save to file in not exist dir should be success")

	config, err := Load()
	assert.NoError(t, err, "load from file should be success")
	assert.Equal(t, "test-0", config.User.AccessToken, "access token should be test-0")

	id = uuid.NewString()
	err = os.MkdirAll(fmt.Sprintf("/tmp/%s", id), fs.ModePerm)
	assert.NoError(t, err, "create dir should be success")

	configFileLocation = fmt.Sprintf("/tmp/%s/config", id)
	err = Save(&Config{
		User: User{AccessToken: "test-1"},
	})
	assert.NoError(t, err, "save to file in exist dir should be success")

	config, err = Load()
	assert.NoError(t, err, "load from file should be success")
	assert.Equal(t, "test-1", config.User.AccessToken, "access token should be test-1")

	err = Save(&Config{
		User: User{AccessToken: "test-2"},
	})
	assert.NoError(t, err, "overwrite config file should be success")

	config, err = Load()
	assert.NoError(t, err, "load from file should be success")
	assert.Equal(t, "test-2", config.User.AccessToken, "access token should be test-2")
}

func TestLoad(t *testing.T) {
	id := uuid.NewString()
	configFileLocation = fmt.Sprintf("/tmp/%s/config", id)

	_, err := Load()
	assert.Contains(t, err.Error(), "no such file or directory", "load from file should be failed")

	dir := filepath.Dir(configFileLocation)
	err = os.MkdirAll(dir, 0750)
	assert.NoError(t, err, "create dir should be success")

	err = os.WriteFile(configFileLocation, []byte("invalid config"), fs.ModePerm)
	assert.NoError(t, err, "write config file should be success")

	_, err = Load()
	assert.Contains(t, err.Error(), "failed to decode config file", "load from file should be failed")
}
