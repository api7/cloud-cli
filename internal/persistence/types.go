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

package persistence

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// User is credential for authentication.
type User struct {
	AccessToken string `json:"-" yaml:"access_token"`
}

// Profile represents a configuration profile.
type Profile struct {
	// Name is the name of the profile.
	Name string `json:"name" yaml:"name"`
	// Address is the address of API7 Cloud server.
	Address string `json:"address" yaml:"address"`
	// User is the user credential.
	User User `json:"user" yaml:"user"`
}

// CloudConfiguration is the configuration for the cloud cli.
type CloudConfiguration struct {
	// DefaultProfile is the active profile.
	DefaultProfile string `json:"default_profile" yaml:"default_profile"`
	// Profiles is the list of profiles.
	Profiles []Profile `json:"profiles" yaml:"profiles"`
}

// ConfigureProfile adds a profile to the configuration if not exists, otherwise update profile by name.
func (c *CloudConfiguration) ConfigureProfile(profile Profile) {
	for i, p := range c.Profiles {
		if p.Name == profile.Name {
			c.Profiles[i] = profile
			return
		}
	}
	c.Profiles = append(c.Profiles, profile)
}

// GetProfile returns the profile by name.
func (c *CloudConfiguration) GetProfile(name string) (*Profile, error) {
	for _, p := range c.Profiles {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, errors.Errorf("profile %s not found", name)
}

// GetDefaultProfile returns the default profile.
func (c *CloudConfiguration) GetDefaultProfile() (*Profile, error) {
	profile, err := c.GetProfile(c.DefaultProfile)
	if err != nil {
		return nil, errors.New("default profile not found")
	}
	return profile, nil
}

// Validate validates the configuration.
func (c *CloudConfiguration) Validate() error {
	if _, err := c.GetDefaultProfile(); err != nil {
		return err
	}

	return nil
}

var (
	// HomeDir is the home directory of the api7 cloud.
	HomeDir = filepath.Join(os.Getenv("HOME"), ".api7cloud")
	// TLSDir is the directory to store TLS certificates.
	TLSDir string
	// APISIXConfigDir is the directory to store APISIX configuration file.
	APISIXConfigDir string
	configDir       string
)

// Init initializes the persistence context.
func Init() error {
	configDir = filepath.Join(HomeDir, "config")

	TLSDir = filepath.Join(HomeDir, "tls")
	if err := os.MkdirAll(TLSDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create tls directory")
	}
	if err := os.Chmod(TLSDir, 0755); err != nil {
		return errors.Wrap(err, "change tls directory permission")
	}

	APISIXConfigDir = filepath.Join(HomeDir, "apisix")
	if err := os.MkdirAll(APISIXConfigDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create apisix config directory")
	}
	if err := os.Chmod(APISIXConfigDir, 0755); err != nil {
		return errors.Wrap(err, "change apisix config directory permission")
	}

	return nil
}
