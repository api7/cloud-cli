package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/persistence"
)

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
