package apiclient

import (
	"context"

	"google.golang.org/grpc"

	workflowtemplatepkg "github.com/kubework/work/pkg/apiclient/workflowtemplate"
	"github.com/kubework/work/pkg/apis/workflow/v1alpha1"
)

type workKubeWorkflowTemplateServiceClient struct {
	delegate workflowtemplatepkg.WorkflowTemplateServiceServer
}

func (a workKubeWorkflowTemplateServiceClient) CreateWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateCreateRequest, _ ...grpc.CallOption) (*v1alpha1.WorkflowTemplate, error) {
	return a.delegate.CreateWorkflowTemplate(ctx, req)
}

func (a workKubeWorkflowTemplateServiceClient) GetWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateGetRequest, _ ...grpc.CallOption) (*v1alpha1.WorkflowTemplate, error) {
	return a.delegate.GetWorkflowTemplate(ctx, req)
}

func (a workKubeWorkflowTemplateServiceClient) ListWorkflowTemplates(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateListRequest, _ ...grpc.CallOption) (*v1alpha1.WorkflowTemplateList, error) {
	return a.delegate.ListWorkflowTemplates(ctx, req)
}

func (a workKubeWorkflowTemplateServiceClient) UpdateWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateUpdateRequest, _ ...grpc.CallOption) (*v1alpha1.WorkflowTemplate, error) {
	return a.delegate.UpdateWorkflowTemplate(ctx, req)
}

func (a workKubeWorkflowTemplateServiceClient) DeleteWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateDeleteRequest, _ ...grpc.CallOption) (*workflowtemplatepkg.WorkflowTemplateDeleteResponse, error) {
	return a.delegate.DeleteWorkflowTemplate(ctx, req)
}

func (a workKubeWorkflowTemplateServiceClient) LintWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateLintRequest, _ ...grpc.CallOption) (*v1alpha1.WorkflowTemplate, error) {
	return a.delegate.LintWorkflowTemplate(ctx, req)
}
