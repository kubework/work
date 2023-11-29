package apiclient

import (
	"context"

	"k8s.io/client-go/tools/clientcmd"

	cronworkflowpkg "github.com/kubework/work/pkg/apiclient/cronworkflow"
	workflowpkg "github.com/kubework/work/pkg/apiclient/workflow"
	workflowarchivepkg "github.com/kubework/work/pkg/apiclient/workflowarchive"
	workflowtemplatepkg "github.com/kubework/work/pkg/apiclient/workflowtemplate"
)

type Client interface {
	NewArchivedWorkflowServiceClient() (workflowarchivepkg.ArchivedWorkflowServiceClient, error)
	NewWorkflowServiceClient() workflowpkg.WorkflowServiceClient
	NewCronWorkflowServiceClient() cronworkflowpkg.CronWorkflowServiceClient
	NewWorkflowTemplateServiceClient() workflowtemplatepkg.WorkflowTemplateServiceClient
}

func NewClient(workServer string, authSupplier func() string, clientConfig clientcmd.ClientConfig) (context.Context, Client, error) {
	if workServer != "" {
		return newWorkServerClient(workServer, authSupplier())
	} else {
		return newWorkKubeClient(clientConfig)
	}
}
