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
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// User is credential for authentication.
type User struct {
	AccessToken string `json:"access_token" yaml:"access_token"`
}

// Credential is the top-level credential for the cloud cli.
type Credential struct {
	User User `json:"user" yaml:"user"`
}

var (
	// HomeDir is the home directory of the api7 cloud.
	HomeDir = filepath.Join(os.Getenv("HOME"), ".api7cloud")
	// TLSDir is the directory to store TLS certificates.
	TLSDir string
)

// Init initializes the persistence context.
func Init() error {
	TLSDir = filepath.Join(HomeDir, "tls")
	if err := os.MkdirAll(TLSDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create tls directory")
	}
	return nil
}
