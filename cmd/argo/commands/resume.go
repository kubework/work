package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/kubework/work/cmd/work/commands/client"
	workflowpkg "github.com/kubework/work/pkg/apiclient/workflow"
)

func NewResumeCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "resume WORKFLOW1 WORKFLOW2...",
		Short: "resume zero or more workflows",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, apiClient := client.NewAPIClient()
			serviceClient := apiClient.NewWorkflowServiceClient()
			namespace := client.Namespace()
			for _, wfName := range args {
				_, err := serviceClient.ResumeWorkflow(ctx, &workflowpkg.WorkflowResumeRequest{
					Name:      wfName,
					Namespace: namespace,
				})
				if err != nil {
					log.Fatalf("Failed to resume %s: %+v", wfName, err)
				}
				fmt.Printf("workflow %s resumed\n", wfName)
			}

		},
	}
	return command
}
