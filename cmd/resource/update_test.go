package resource

import (
	"fmt"
	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/persistence"
	sdk "github.com/api7/cloud-go-sdk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestServiceUpdate(t *testing.T) {

	testCases := []struct {
		name       string
		config     *persistence.CloudConfiguration
		args       []string
		mockCloud  func(api *cloud.MockAPI)
		outputs    string
		testConfig string
	}{
		{
			name: "update service",
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
			args:       []string{"update", "--kind", "service", "--config", os.TempDir() + "config.json"},
			testConfig: os.TempDir() + "config.json",
			mockCloud: func(api *cloud.MockAPI) {
				api.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 123,
				}, nil)
				api.EXPECT().UpdateService(sdk.ID(123), os.TempDir()+"config.json").Return(&sdk.Application{
					ID:        sdk.ID(123),
					ClusterID: sdk.ID(123),
					ApplicationSpec: sdk.ApplicationSpec{
						Description: "456",
					},
				}, nil)
			},
			outputs: "{\n\t\"name\": \"\",\n\t\"description\": \"456\",\n\t\"path_prefix\": \"\",\n\t\"hosts\": null,\n\t\"upstreams\": null,\n\t\"active\": 0,\n\t\"id\": \"123\",\n\t\"cluster_id\": \"123\",\n\t\"status\": 0,\n\t\"created_at\": \"0001-01-01T00:00:00Z\",\n\t\"updated_at\": \"0001-01-01T00:00:00Z\",\n\t\"available_cert_ids\": null,\n\t\"canary_release_id\": null,\n\t\"canary_upstream_version_list\": null\n}",
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
			testFile, _ := os.Create(tc.testConfig)
			defer testFile.Close()
			_, _ = testFile.Write([]byte(tc.outputs))

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

			assert.Contains(t, string(output), tc.outputs, "check output")

			os.Remove(tc.testConfig)
		})
	}
}
