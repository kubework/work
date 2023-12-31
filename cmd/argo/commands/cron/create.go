package cron

import (
	"log"
	"os"

	"github.com/kubework/pkg/json"
	"github.com/spf13/cobra"

	"github.com/kubework/work/cmd/work/commands/client"
	cronworkflowpkg "github.com/kubework/work/pkg/apiclient/cronworkflow"

	wfv1 "github.com/kubework/work/pkg/apis/workflow/v1alpha1"
	"github.com/kubework/work/workflow/common"
	"github.com/kubework/work/workflow/util"
)

type cliCreateOpts struct {
	output string // --output
	strict bool   // --strict
}

type cronWorkflowSubmitOpts struct {
	instanceId string
}

func NewCreateCommand() *cobra.Command {
	var (
		cliCreateOpts cliCreateOpts
		submitOpts    cronWorkflowSubmitOpts
	)
	var command = &cobra.Command{
		Use:   "create FILE1 FILE2...",
		Short: "create a cron workflow",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			CreateCronWorkflows(args, &cliCreateOpts, &submitOpts)
		},
	}
	command.Flags().StringVarP(&cliCreateOpts.output, "output", "o", "", "Output format. One of: name|json|yaml|wide")
	command.Flags().StringVar(&submitOpts.instanceId, "instanceid", "", "submit with a specific controller's instance id label")
	command.Flags().BoolVar(&cliCreateOpts.strict, "strict", true, "perform strict workflow validation")
	return command
}

func CreateCronWorkflows(filePaths []string, cliOpts *cliCreateOpts, submitOpts *cronWorkflowSubmitOpts) {

	ctx, apiClient := client.NewAPIClient()
	serviceClient := apiClient.NewCronWorkflowServiceClient()
	namespace := client.Namespace()

	fileContents, err := util.ReadManifest(filePaths...)
	if err != nil {
		log.Fatal(err)
	}

	var cronWorkflows []wfv1.CronWorkflow
	for _, body := range fileContents {
		cronWfs := unmarshalCronWorkflows(body, cliOpts.strict)
		cronWorkflows = append(cronWorkflows, cronWfs...)
	}

	if len(cronWorkflows) == 0 {
		log.Println("No CronWorkflows found in given files")
		os.Exit(1)
	}

	for _, cronWf := range cronWorkflows {
		applySubmitOpts(&cronWf, submitOpts)
		created, err := serviceClient.CreateCronWorkflow(ctx, &cronworkflowpkg.CreateCronWorkflowRequest{
			Namespace:    namespace,
			CronWorkflow: &cronWf,
		})
		if err != nil {
			log.Fatalf("Failed to create workflow template: %v", err)
		}
		printCronWorkflowTemplate(created)
	}
}

// unmarshalCronWorkflows unmarshals the input bytes as either json or yaml
func unmarshalCronWorkflows(wfBytes []byte, strict bool) []wfv1.CronWorkflow {
	var cronWf wfv1.CronWorkflow
	var jsonOpts []json.JSONOpt
	if strict {
		jsonOpts = append(jsonOpts, json.DisallowUnknownFields)
	}
	err := json.Unmarshal(wfBytes, &cronWf, jsonOpts...)
	if err == nil {
		return []wfv1.CronWorkflow{cronWf}
	}
	yamlWfs, err := common.SplitCronWorkflowYAMLFile(wfBytes, strict)
	if err == nil {
		return yamlWfs
	}
	log.Fatalf("Failed to parse workflow template: %v", err)
	return nil
}

func applySubmitOpts(cwf *wfv1.CronWorkflow, submitOpts *cronWorkflowSubmitOpts) {
	labels := cwf.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	if submitOpts.instanceId != "" {
		labels[common.LabelKeyControllerInstanceID] = submitOpts.instanceId
	}
	cwf.SetLabels(labels)
}
