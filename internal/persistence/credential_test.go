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

func TestSaveCredential(t *testing.T) {
	id := uuid.NewString()
	credentialFileLocation = fmt.Sprintf("/tmp/%s/credential", id)
	err := SaveCredential(&Credential{
		User: User{AccessToken: "test-0"},
	})
	assert.NoError(t, err, "save to file in not exist dir should be success")

	credential, err := LoadCredential()
	assert.NoError(t, err, "load from file should be success")
	assert.Equal(t, "test-0", credential.User.AccessToken, "access token should be test-0")

	id = uuid.NewString()
	err = os.MkdirAll(fmt.Sprintf("/tmp/%s", id), fs.ModePerm)
	assert.NoError(t, err, "create dir should be success")

	credentialFileLocation = fmt.Sprintf("/tmp/%s/credential", id)
	err = SaveCredential(&Credential{
		User: User{AccessToken: "test-1"},
	})
	assert.NoError(t, err, "save to file in exist dir should be success")

	credential, err = LoadCredential()
	assert.NoError(t, err, "load from file should be success")
	assert.Equal(t, "test-1", credential.User.AccessToken, "access token should be test-1")

	err = SaveCredential(&Credential{
		User: User{AccessToken: "test-2"},
	})
	assert.NoError(t, err, "overwrite credential file should be success")

	credential, err = LoadCredential()
	assert.NoError(t, err, "load from file should be success")
	assert.Equal(t, "test-2", credential.User.AccessToken, "access token should be test-2")
}

func TestLoad(t *testing.T) {
	id := uuid.NewString()
	credentialFileLocation = fmt.Sprintf("/tmp/%s/credential", id)

	_, err := LoadCredential()
	assert.Contains(t, err.Error(), "no such file or directory", "load from file should be failed")

	dir := filepath.Dir(credentialFileLocation)
	err = os.MkdirAll(dir, 0750)
	assert.NoError(t, err, "create dir should be success")

	err = os.WriteFile(credentialFileLocation, []byte("invalid credential"), fs.ModePerm)
	assert.NoError(t, err, "write credential file should be success")

	_, err = LoadCredential()
	assert.Contains(t, err.Error(), "failed to decode credential", "load from file should be failed")
}
