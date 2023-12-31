package cron

import (
	"fmt"

	"github.com/kubework/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kubework/work/cmd/work/commands/client"
	cronworkflowpkg "github.com/kubework/work/pkg/apiclient/cronworkflow"
)

// NewSuspendCommand returns a new instance of an `work suspend` command
func NewSuspendCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "suspend CRON_WORKFLOW...",
		Short: "suspend zero or more cron workflows",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, apiClient := client.NewAPIClient()
			serviceClient := apiClient.NewCronWorkflowServiceClient()
			namespace := client.Namespace()
			for _, name := range args {
				cronWf, err := serviceClient.GetCronWorkflow(ctx, &cronworkflowpkg.GetCronWorkflowRequest{
					Name:      name,
					Namespace: namespace,
				})
				errors.CheckError(err)
				cronWf.Spec.Suspend = true
				_, err = serviceClient.UpdateCronWorkflow(ctx, &cronworkflowpkg.UpdateCronWorkflowRequest{
					Name:         name,
					Namespace:    namespace,
					CronWorkflow: cronWf,
				})
				errors.CheckError(err)
				fmt.Printf("CronWorkflow '%s' suspended\n", name)
			}
		},
	}
	return command
}
