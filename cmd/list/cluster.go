package list

import (
	"fmt"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/spf13/cobra"
)

func newClustersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clusters",
		Short: "List the clusters obtained from the API7 Cloud.",
		Example: `
cloud-cli list clusters --count 10 --skip 1 `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(options.Global.List)
			user, err := cloud.Client().Me()
			if err != nil {
				output.Errorf(err.Error())
			}
			clustersList, err := cloud.DefaultClient.ListClusters(user.OrgIDs[0], options.Global.List.Clusters.Count, options.Global.List.Clusters.Skip)
			if err != nil {
				output.Errorf("Failed to list clusters: %s", err.Error())
			}
			fmt.Println(clustersList)
		},
	}

	cmd.PersistentFlags().IntVar(&options.Global.List.Clusters.Count, "count", 10, "Specify the amount of data to be listed")
	cmd.PersistentFlags().IntVar(&options.Global.List.Clusters.Skip, "skip", 1, "Specifies how much data to skip ahead")

	return cmd
}
