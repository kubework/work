package commands

import (
	"fmt"

	"github.com/kubework/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kubework/work/cmd/work/commands/client"
	workflowpkg "github.com/kubework/work/pkg/apiclient/workflow"
)

func NewTerminateCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "terminate WORKFLOW WORKFLOW2...",
		Short: "terminate zero or more workflows",
		Run: func(cmd *cobra.Command, args []string) {

			ctx, apiClient := client.NewAPIClient()
			serviceClient := apiClient.NewWorkflowServiceClient()
			namespace := client.Namespace()
			for _, name := range args {
				wf, err := serviceClient.TerminateWorkflow(ctx, &workflowpkg.WorkflowTerminateRequest{
					Name:      name,
					Namespace: namespace,
				})
				errors.CheckError(err)
				fmt.Printf("workflow %s terminated\n", wf.Name)
			}
		},
	}
	return command
}
