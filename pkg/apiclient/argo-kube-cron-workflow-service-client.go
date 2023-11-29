package apiclient

import (
	"context"

	"google.golang.org/grpc"

	cronworkflowpkg "github.com/kubework/work/pkg/apiclient/cronworkflow"
	"github.com/kubework/work/pkg/apis/workflow/v1alpha1"
)

type workKubeCronWorkflowServiceClient struct {
	delegate cronworkflowpkg.CronWorkflowServiceServer
}

func (c workKubeCronWorkflowServiceClient) LintCronWorkflow(ctx context.Context, req *cronworkflowpkg.LintCronWorkflowRequest, _ ...grpc.CallOption) (*v1alpha1.CronWorkflow, error) {
	return c.delegate.LintCronWorkflow(ctx, req)
}

func (c workKubeCronWorkflowServiceClient) CreateCronWorkflow(ctx context.Context, req *cronworkflowpkg.CreateCronWorkflowRequest, _ ...grpc.CallOption) (*v1alpha1.CronWorkflow, error) {
	return c.delegate.CreateCronWorkflow(ctx, req)
}

func (c workKubeCronWorkflowServiceClient) ListCronWorkflows(ctx context.Context, req *cronworkflowpkg.ListCronWorkflowsRequest, _ ...grpc.CallOption) (*v1alpha1.CronWorkflowList, error) {
	return c.delegate.ListCronWorkflows(ctx, req)
}

func (c workKubeCronWorkflowServiceClient) GetCronWorkflow(ctx context.Context, req *cronworkflowpkg.GetCronWorkflowRequest, _ ...grpc.CallOption) (*v1alpha1.CronWorkflow, error) {
	return c.delegate.GetCronWorkflow(ctx, req)
}

func (c workKubeCronWorkflowServiceClient) UpdateCronWorkflow(ctx context.Context, req *cronworkflowpkg.UpdateCronWorkflowRequest, _ ...grpc.CallOption) (*v1alpha1.CronWorkflow, error) {
	return c.delegate.UpdateCronWorkflow(ctx, req)
}

func (c workKubeCronWorkflowServiceClient) DeleteCronWorkflow(ctx context.Context, req *cronworkflowpkg.DeleteCronWorkflowRequest, _ ...grpc.CallOption) (*cronworkflowpkg.CronWorkflowDeletedResponse, error) {
	return c.delegate.DeleteCronWorkflow(ctx, req)
}
