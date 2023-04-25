package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/persistence"
	sdk "github.com/api7/cloud-go-sdk"
)

func TestServiceList(t *testing.T) {
	testCases := []struct {
		name      string
		config    *persistence.CloudConfiguration
		args      []string
		mockCloud func(api *cloud.MockAPI)
		outputs   []string
	}{
		{
			name: "list service",
			config: &persistence.CloudConfiguration{
				DefaultProfile: "prod",
				Profiles: []persistence.Profile{
					{
						Name:    "prod",
						Address: "https://prod.api7.ai",
						User: persistence.User{
							AccessToken: "prod-token",
						},
					},
				},
			},
			args: []string{"list", "--kind", "service"},
			mockCloud: func(api *cloud.MockAPI) {
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID:   123,
					Name: "API7.AI",
				}, nil)
				api.EXPECT().ListServices(sdk.ID(123), 10, 0).Return([]*sdk.Application{
					{
						ID:        123,
						ClusterID: 123,
					},
				}, nil)
			},
			outputs: []string{"[\n\t{\n\t\t\"name\": \"\",\n\t\t\"description\": \"\",\n\t\t\"path_prefix\": \"\",\n\t\t\"hosts\": null,\n\t\t\"upstreams\": null,\n\t\t\"active\": 0,\n\t\t\"id\": \"123\",\n\t\t\"cluster_id\": \"123\",\n\t\t\"status\": 0,\n\t\t\"created_at\": \"0001-01-01T00:00:00Z\",\n\t\t\"updated_at\": \"0001-01-01T00:00:00Z\",\n\t\t\"available_cert_ids\": null,\n\t\t\"canary_release_id\": null,\n\t\t\"canary_upstream_version_list\": null\n\t}\n]"},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := persistence.SaveConfiguration(tc.config)
			assert.NoError(t, err, "prepare fake cloud configuration")

			ctrl := gomock.NewController(t)
			api := cloud.NewMockAPI(ctrl)
			cloud.NewClient = func(_ string, _ string, _ bool) (cloud.API, error) {
				return api, nil
			}
			if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
				if tc.mockCloud != nil {
					tc.mockCloud(api)
				}
				cloud.DefaultClient = api
				cmd := NewCommand()
				cmd.SetArgs(tc.args)
				err := cmd.Execute()
				assert.NoError(t, err, "check if the command executed successfully")
				return
			}

			cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")

			output, _ := cmd.CombinedOutput()
			fmt.Println(string(output))
			for _, o := range tc.outputs {
				assert.Contains(t, string(output), strings.TrimSpace(o), "check output")
			}

		})
	}
}
