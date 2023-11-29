package apiclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	cronworkflowpkg "github.com/kubework/work/pkg/apiclient/cronworkflow"
	workflowpkg "github.com/kubework/work/pkg/apiclient/workflow"
	workflowarchivepkg "github.com/kubework/work/pkg/apiclient/workflowarchive"
	workflowtemplatepkg "github.com/kubework/work/pkg/apiclient/workflowtemplate"
)

type workServerClient struct {
	*grpc.ClientConn
}

func newWorkServerClient(workServer, auth string) (context.Context, Client, error) {
	conn, err := NewClientConn(workServer)
	if err != nil {
		return nil, nil, err
	}
	return newContext(auth), &workServerClient{conn}, nil
}

func (a *workServerClient) NewWorkflowServiceClient() workflowpkg.WorkflowServiceClient {
	return workflowpkg.NewWorkflowServiceClient(a.ClientConn)
}

func (a *workServerClient) NewCronWorkflowServiceClient() cronworkflowpkg.CronWorkflowServiceClient {
	return cronworkflowpkg.NewCronWorkflowServiceClient(a.ClientConn)
}

func (a *workServerClient) NewWorkflowTemplateServiceClient() workflowtemplatepkg.WorkflowTemplateServiceClient {
	return workflowtemplatepkg.NewWorkflowTemplateServiceClient(a.ClientConn)
}

func (a *workServerClient) NewArchivedWorkflowServiceClient() (workflowarchivepkg.ArchivedWorkflowServiceClient, error) {
	return workflowarchivepkg.NewArchivedWorkflowServiceClient(a.ClientConn), nil
}

func NewClientConn(workServer string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(workServer, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DEPRECATED
func NewContext(auth string) context.Context {
	return newContext(auth)
}

func newContext(auth string) context.Context {
	if auth == "" {
		return context.Background()
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", auth))
}
