package workflowtemplate

import (
	"context"
	"fmt"
	"sort"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	workflowtemplatepkg "github.com/kubework/work/pkg/apiclient/workflowtemplate"
	"github.com/kubework/work/pkg/apis/workflow/v1alpha1"
	"github.com/kubework/work/server/auth"
	"github.com/kubework/work/workflow/templateresolution"
	"github.com/kubework/work/workflow/validate"
)

type WorkflowTemplateServer struct {
}

func NewWorkflowTemplateServer() workflowtemplatepkg.WorkflowTemplateServiceServer {
	return &WorkflowTemplateServer{}
}

func (wts *WorkflowTemplateServer) CreateWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateCreateRequest) (*v1alpha1.WorkflowTemplate, error) {
	wfClient := auth.GetWfClient(ctx)
	if req.Template == nil {
		return nil, fmt.Errorf("workflow template was not found in the request body")
	}
	wftmplGetter := templateresolution.WrapWorkflowTemplateInterface(wfClient.KubeworkV1alpha1().WorkflowTemplates(req.Namespace))

	err := validate.ValidateWorkflowTemplate(wftmplGetter, req.Template)
	if err != nil {
		return nil, err
	}

	return wfClient.KubeworkV1alpha1().WorkflowTemplates(req.Namespace).Create(req.Template)

}

func (wts *WorkflowTemplateServer) GetWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateGetRequest) (*v1alpha1.WorkflowTemplate, error) {
	wfClient := auth.GetWfClient(ctx)

	wfTmpl, err := wfClient.KubeworkV1alpha1().WorkflowTemplates(req.Namespace).Get(req.Name, v1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return wfTmpl, err
}

func (wts *WorkflowTemplateServer) ListWorkflowTemplates(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateListRequest) (*v1alpha1.WorkflowTemplateList, error) {
	wfClient := auth.GetWfClient(ctx)
	options := v1.ListOptions{}
	if req.ListOptions != nil {
		options = *req.ListOptions
	}
	wfList, err := wfClient.KubeworkV1alpha1().WorkflowTemplates(req.Namespace).List(options)
	if err != nil {
		return nil, err
	}

	sort.Sort(wfList.Items)

	return wfList, nil
}

func (wts *WorkflowTemplateServer) DeleteWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateDeleteRequest) (*workflowtemplatepkg.WorkflowTemplateDeleteResponse, error) {
	wfClient := auth.GetWfClient(ctx)

	err := wfClient.KubeworkV1alpha1().WorkflowTemplates(req.Namespace).Delete(req.Name, &v1.DeleteOptions{})
	if err != nil {
		return nil, err
	}

	return &workflowtemplatepkg.WorkflowTemplateDeleteResponse{}, nil
}

func (wts *WorkflowTemplateServer) LintWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateLintRequest) (*v1alpha1.WorkflowTemplate, error) {
	wfClient := auth.GetWfClient(ctx)

	wftmplGetter := templateresolution.WrapWorkflowTemplateInterface(wfClient.KubeworkV1alpha1().WorkflowTemplates(req.Namespace))

	err := validate.ValidateWorkflowTemplate(wftmplGetter, req.Template)
	if err != nil {
		return nil, err
	}

	return req.Template, nil
}

func (wts *WorkflowTemplateServer) UpdateWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateUpdateRequest) (*v1alpha1.WorkflowTemplate, error) {
	if req.Template == nil {
		return nil, fmt.Errorf("WorkflowTemplate is not found in Request body")
	}
	wfClient := auth.GetWfClient(ctx)
	wftmplGetter := templateresolution.WrapWorkflowTemplateInterface(wfClient.KubeworkV1alpha1().WorkflowTemplates(req.Namespace))

	err := validate.ValidateWorkflowTemplate(wftmplGetter, req.Template)
	if err != nil {
		return nil, err
	}

	res, err := wfClient.KubeworkV1alpha1().WorkflowTemplates(req.Namespace).Update(req.Template)
	return res, err
}
