package apiclient

import (
	"context"

	"google.golang.org/grpc"

	workflowpkg "github.com/kubework/work/pkg/apiclient/workflow"
	"github.com/kubework/work/pkg/apis/workflow/v1alpha1"
)

type workKubeWorkflowServiceClient struct {
	delegate workflowpkg.WorkflowServiceServer
}

func (c workKubeWorkflowServiceClient) CreateWorkflow(ctx context.Context, req *workflowpkg.WorkflowCreateRequest, _ ...grpc.CallOption) (*v1alpha1.Workflow, error) {
	return c.delegate.CreateWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) GetWorkflow(ctx context.Context, req *workflowpkg.WorkflowGetRequest, _ ...grpc.CallOption) (*v1alpha1.Workflow, error) {
	return c.delegate.GetWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) ListWorkflows(ctx context.Context, req *workflowpkg.WorkflowListRequest, _ ...grpc.CallOption) (*v1alpha1.WorkflowList, error) {
	return c.delegate.ListWorkflows(ctx, req)
}

func (c workKubeWorkflowServiceClient) WatchWorkflows(ctx context.Context, req *workflowpkg.WatchWorkflowsRequest, _ ...grpc.CallOption) (workflowpkg.WorkflowService_WatchWorkflowsClient, error) {
	panic("not implemented")
}

func (c workKubeWorkflowServiceClient) DeleteWorkflow(ctx context.Context, req *workflowpkg.WorkflowDeleteRequest, _ ...grpc.CallOption) (*workflowpkg.WorkflowDeleteResponse, error) {
	return c.delegate.DeleteWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) RetryWorkflow(ctx context.Context, req *workflowpkg.WorkflowRetryRequest, _ ...grpc.CallOption) (*v1alpha1.Workflow, error) {
	return c.delegate.RetryWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) ResubmitWorkflow(ctx context.Context, req *workflowpkg.WorkflowResubmitRequest, _ ...grpc.CallOption) (*v1alpha1.Workflow, error) {
	return c.delegate.ResubmitWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) ResumeWorkflow(ctx context.Context, req *workflowpkg.WorkflowResumeRequest, _ ...grpc.CallOption) (*v1alpha1.Workflow, error) {
	return c.delegate.ResumeWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) SuspendWorkflow(ctx context.Context, req *workflowpkg.WorkflowSuspendRequest, _ ...grpc.CallOption) (*v1alpha1.Workflow, error) {
	return c.delegate.SuspendWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) TerminateWorkflow(ctx context.Context, req *workflowpkg.WorkflowTerminateRequest, _ ...grpc.CallOption) (*v1alpha1.Workflow, error) {
	return c.delegate.TerminateWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) LintWorkflow(ctx context.Context, req *workflowpkg.WorkflowLintRequest, _ ...grpc.CallOption) (*v1alpha1.Workflow, error) {
	return c.delegate.LintWorkflow(ctx, req)
}

func (c workKubeWorkflowServiceClient) PodLogs(ctx context.Context, req *workflowpkg.WorkflowLogRequest, _ ...grpc.CallOption) (workflowpkg.WorkflowService_PodLogsClient, error) {
	panic("not implemented")
}
