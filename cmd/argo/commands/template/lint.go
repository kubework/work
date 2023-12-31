package template

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kubework/work/cmd/work/commands/client"
	workflowtemplatepkg "github.com/kubework/work/pkg/apiclient/workflowtemplate"
	wfv1 "github.com/kubework/work/pkg/apis/workflow/v1alpha1"
	cmdutil "github.com/kubework/work/util/cmd"
	"github.com/kubework/work/workflow/validate"
)

func NewLintCommand() *cobra.Command {
	var (
		strict bool
	)
	var command = &cobra.Command{
		Use:   "lint (DIRECTORY | FILE1 FILE2 FILE3...)",
		Short: "validate a file or directory of workflow template manifests",
		Run: func(cmd *cobra.Command, args []string) {
			err := ServerSideLint(args, strict)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("WorkflowTemplate manifests validated\n")
		},
	}
	command.Flags().BoolVar(&strict, "strict", true, "perform strict workflow validation")
	return command
}

func ServerSideLint(args []string, strict bool) error {
	validateDir := cmdutil.MustIsDir(args[0])

	ctx, apiClient := client.NewAPIClient()
	serviceClient := apiClient.NewWorkflowTemplateServiceClient()
	namespace := client.Namespace()

	if validateDir {
		if len(args) > 1 {
			fmt.Printf("Validation of a single directory supported")
			os.Exit(1)
		}
		walkFunc := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info == nil || info.IsDir() {
				return nil
			}
			fileExt := filepath.Ext(info.Name())
			switch fileExt {
			case ".yaml", ".yml", ".json":
			default:
				return nil
			}
			wfTmpls, err := validate.ParseWfTmplFromFile(path, strict)
			if err != nil {
				log.Error(err)
			}
			for _, wfTmpl := range wfTmpls {
				err := ServerLintValidation(ctx, serviceClient, wfTmpl, namespace)
				if err != nil {
					log.Error(err)
				}
			}
			return nil
		}
		return filepath.Walk(args[0], walkFunc)
	} else {
		for _, arg := range args {
			wfTmpls, err := validate.ParseWfTmplFromFile(arg, strict)
			if err != nil {
				log.Error(err)
			}
			for _, wfTmpl := range wfTmpls {
				err := ServerLintValidation(ctx, serviceClient, wfTmpl, namespace)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	return nil
}

func ServerLintValidation(ctx context.Context, client workflowtemplatepkg.WorkflowTemplateServiceClient, wfTmpl wfv1.WorkflowTemplate, ns string) error {
	wfTmplReq := workflowtemplatepkg.WorkflowTemplateLintRequest{
		Namespace: ns,
		Template:  &wfTmpl,
	}
	_, err := client.LintWorkflowTemplate(ctx, &wfTmplReq)
	return err
}
