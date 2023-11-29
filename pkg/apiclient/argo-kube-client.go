package apiclient

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kubework/work/persist/sqldb"
	"github.com/kubework/work/pkg/apiclient/cronworkflow"
	workflowpkg "github.com/kubework/work/pkg/apiclient/workflow"
	workflowarchivepkg "github.com/kubework/work/pkg/apiclient/workflowarchive"
	"github.com/kubework/work/pkg/apiclient/workflowtemplate"
	"github.com/kubework/work/pkg/client/clientset/versioned"
	"github.com/kubework/work/server/auth"
	cronworkflowserver "github.com/kubework/work/server/cronworkflow"
	workflowserver "github.com/kubework/work/server/workflow"
	workflowtemplateserver "github.com/kubework/work/server/workflowtemplate"
	"github.com/kubework/work/util/help"
)

type workKubeClient struct {
}

func newWorkKubeClient(clientConfig clientcmd.ClientConfig) (context.Context, Client, error) {
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	wfClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}
	gatekeeper := auth.NewGatekeeper(auth.Server, wfClient, kubeClient, restConfig)
	ctx, err := gatekeeper.Context(context.Background())
	if err != nil {
		return nil, nil, err
	}
	return ctx, &workKubeClient{}, nil
}

func (a *workKubeClient) NewWorkflowServiceClient() workflowpkg.WorkflowServiceClient {
	return &workKubeWorkflowServiceClient{workflowserver.NewWorkflowServer(sqldb.ExplosiveOffloadNodeStatusRepo)}
}

func (a *workKubeClient) NewCronWorkflowServiceClient() cronworkflow.CronWorkflowServiceClient {
	return &workKubeCronWorkflowServiceClient{cronworkflowserver.NewCronWorkflowServer()}
}
func (a *workKubeClient) NewWorkflowTemplateServiceClient() workflowtemplate.WorkflowTemplateServiceClient {
	return &workKubeWorkflowTemplateServiceClient{workflowtemplateserver.NewWorkflowTemplateServer()}
}

func (a *workKubeClient) NewArchivedWorkflowServiceClient() (workflowarchivepkg.ArchivedWorkflowServiceClient, error) {
	return nil, fmt.Errorf("it is impossible to interact with the workflow archive if you are not using the Work Server, see " + help.CLI)
}
