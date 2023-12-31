package cron

import (
	"github.com/kubework/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kubework/work/cmd/work/commands/client"
	cronworkflowpkg "github.com/kubework/work/pkg/apiclient/cronworkflow"
)

// NewDeleteCommand returns a new instance of an `work delete` command
func NewDeleteCommand() *cobra.Command {
	var (
		all bool
	)

	var command = &cobra.Command{
		Use:   "delete [CRON_WORKFLOW... | --all]",
		Short: "delete a cron workflow",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, apiClient := client.NewAPIClient()
			serviceClient := apiClient.NewCronWorkflowServiceClient()
			if all {
				cronWfList, err := serviceClient.ListCronWorkflows(ctx, &cronworkflowpkg.ListCronWorkflowsRequest{
					Namespace: client.Namespace(),
				})
				errors.CheckError(err)
				for _, cronWf := range cronWfList.Items {
					args = append(args, cronWf.Name)
				}
			}
			for _, name := range args {
				_, err := serviceClient.DeleteCronWorkflow(ctx, &cronworkflowpkg.DeleteCronWorkflowRequest{
					Name:      name,
					Namespace: client.Namespace(),
				})
				errors.CheckError(err)
			}
		},
	}

	command.Flags().BoolVar(&all, "all", false, "Delete all workflow templates")
	return command
}
