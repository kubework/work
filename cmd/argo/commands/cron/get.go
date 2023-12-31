package cron

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kubework/pkg/errors"
	"github.com/kubework/pkg/humanize"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/kubework/work/cmd/work/commands/client"
	"github.com/kubework/work/pkg/apiclient/cronworkflow"
	wfv1 "github.com/kubework/work/pkg/apis/workflow/v1alpha1"
)

func NewGetCommand() *cobra.Command {
	var (
		output string
	)

	var command = &cobra.Command{
		Use:   "get CRON_WORKFLOW...",
		Short: "display details about a cron workflow",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			ctx, apiClient := client.NewAPIClient()
			serviceClient := apiClient.NewCronWorkflowServiceClient()
			namespace := client.Namespace()

			for _, arg := range args {
				cronWf, err := serviceClient.GetCronWorkflow(ctx, &cronworkflow.GetCronWorkflowRequest{
					Name:      arg,
					Namespace: namespace,
				})
				errors.CheckError(err)
				printCronWorkflow(cronWf, output)
			}
		},
	}

	command.Flags().StringVarP(&output, "output", "o", "", "Output format. One of: json|yaml|wide")
	return command
}

func printCronWorkflow(wf *wfv1.CronWorkflow, outFmt string) {
	switch outFmt {
	case "name":
		fmt.Println(wf.ObjectMeta.Name)
	case "json":
		outBytes, _ := json.MarshalIndent(wf, "", "    ")
		fmt.Println(string(outBytes))
	case "yaml":
		outBytes, _ := yaml.Marshal(wf)
		fmt.Print(string(outBytes))
	case "wide", "":
		printCronWorkflowTemplate(wf)
	default:
		log.Fatalf("Unknown output format: %s", outFmt)
	}
}

func printCronWorkflowTemplate(wf *wfv1.CronWorkflow) {
	const fmtStr = "%-30s %v\n"
	fmt.Printf(fmtStr, "Name:", wf.ObjectMeta.Name)
	fmt.Printf(fmtStr, "Namespace:", wf.ObjectMeta.Namespace)
	fmt.Printf(fmtStr, "Created:", humanize.Timestamp(wf.ObjectMeta.CreationTimestamp.Time))
	fmt.Printf(fmtStr, "Schedule:", wf.Spec.Schedule)
	fmt.Printf(fmtStr, "Suspended:", wf.Spec.Suspend)
	if wf.Spec.Timezone != "" {
		fmt.Printf(fmtStr, "Timezone:", wf.Spec.Timezone)
	}
	if wf.Spec.StartingDeadlineSeconds != nil {
		fmt.Printf(fmtStr, "StartingDeadlineSeconds:", *wf.Spec.StartingDeadlineSeconds)
	}
	if wf.Spec.ConcurrencyPolicy != "" {
		fmt.Printf(fmtStr, "ConcurrencyPolicy:", wf.Spec.ConcurrencyPolicy)
	}
	if wf.Status.LastScheduledTime != nil {
		fmt.Printf(fmtStr, "LastScheduledTime:", humanize.Timestamp(wf.Status.LastScheduledTime.Time))
	}
	if len(wf.Status.Active) > 0 {
		var activeWfNames []string
		for _, activeWf := range wf.Status.Active {
			activeWfNames = append(activeWfNames, activeWf.Name)
		}
		fmt.Printf(fmtStr, "Active Workflows:", strings.Join(activeWfNames, ", "))
	}
}
