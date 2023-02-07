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

// only for testing
package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/persistence"
)

// PrepareFakeConfiguration prepares a fake configuration for testing.
func PrepareFakeConfiguration(t *testing.T) {
	err := persistence.SaveConfiguration(&persistence.CloudConfiguration{
		DefaultProfile: "default",
		Profiles: []persistence.Profile{
			{
				Name:    "default",
				Address: "https://api.api7.cloud",
				User:    persistence.User{AccessToken: "test-token"},
			},
		},
	})
	assert.NoError(t, err, "prepare fake cloud configuration")
}
