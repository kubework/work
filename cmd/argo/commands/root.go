package commands

import (
	"fmt"

	"github.com/kubework/pkg/cli"
	"github.com/spf13/cobra"

	"github.com/kubework/work/cmd/work/commands/auth"
	"github.com/kubework/work/cmd/work/commands/cron"
	"github.com/kubework/work/util/help"

	"github.com/kubework/work/cmd/work/commands/archive"
	"github.com/kubework/work/cmd/work/commands/client"
	"github.com/kubework/work/cmd/work/commands/template"
	"github.com/kubework/work/util/cmd"
)

const (
	// CLIName is the name of the CLI
	CLIName = "work"
)

// NewCommand returns a new instance of an work command
func NewCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CLIName,
		Short: "work is the command line interface to Work",
		Example: fmt.Sprintf(`
If you're using the Work Server (e.g. because you need large workflow support or workflow archive), please read %s.`, help.CLI),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewCompletionCommand())
	command.AddCommand(NewDeleteCommand())
	command.AddCommand(NewGetCommand())
	command.AddCommand(NewLintCommand())
	command.AddCommand(NewListCommand())
	command.AddCommand(NewLogsCommand())
	command.AddCommand(NewResubmitCommand())
	command.AddCommand(NewResumeCommand())
	command.AddCommand(NewRetryCommand())
	command.AddCommand(NewServerCommand())
	command.AddCommand(NewSubmitCommand())
	command.AddCommand(NewSuspendCommand())
	command.AddCommand(auth.NewAuthCommand())
	command.AddCommand(NewWaitCommand())
	command.AddCommand(NewWatchCommand())
	command.AddCommand(NewTerminateCommand())
	command.AddCommand(archive.NewArchiveCommand())
	command.AddCommand(cmd.NewVersionCmd(CLIName))
	command.AddCommand(template.NewTemplateCommand())
	command.AddCommand(cron.NewCronWorkflowCommand())
	client.AddKubectlFlagsToCmd(command)
	client.AddWorkServerFlagsToCmd(command)

	// global log level
	var logLevel string
	command.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cli.SetLogLevel(logLevel)
	}
	command.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "Set the logging level. One of: debug|info|warn|error")

	return command
}
