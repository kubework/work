package archive

import (
	"fmt"

	"github.com/kubework/pkg/errors"
	"github.com/spf13/cobra"

	client "github.com/kubework/work/cmd/work/commands/client"
	workflowarchivepkg "github.com/kubework/work/pkg/apiclient/workflowarchive"
)

func NewDeleteCommand() *cobra.Command {
	var command = &cobra.Command{
		Use: "delete UID...",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, apiClient := client.NewAPIClient()
			serviceClient, err := apiClient.NewArchivedWorkflowServiceClient()
			errors.CheckError(err)
			for _, uid := range args {
				_, err = serviceClient.DeleteArchivedWorkflow(ctx, &workflowarchivepkg.DeleteArchivedWorkflowRequest{Uid: uid})
				errors.CheckError(err)
				fmt.Printf("Archived workflow '%s' deleted\n", uid)
			}
		},
	}
	return command
}
