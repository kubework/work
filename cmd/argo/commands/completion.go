package commands

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	bashCompletionFunc = `
__work_get_workflow() {
	local status="$1"
	local -a work_out
	if work_out=($(work list --status="$status" --output name 2>/dev/null)); then
		COMPREPLY+=( $( compgen -W "${work_out[*]}" -- "$cur" ) )
	fi
}

__work_get_logs() {
	# Determine if were completing a workflow or not.
	local workflow=0
	for comp_word in "${COMP_WORDS[@]}"; do
		if [[ $comp_word =~ ^(-w|--workflow)$ ]]; then
			workflow=1
			break
		fi
	done

	# If completing a workflow, call normal function.
	if [[ $workflow -eq 1 ]]; then
		__work_get_workflow && return $?
	fi

	# Otherwise, complete the list of pods
	local -a kubectl_out
	if kubectl_out=($(kubectl get pods --no-headers --label-columns=workflows.kubework.io/workflow 2>/dev/null | awk '{if ($6!="") print $1}' 2>/dev/null)); then
		COMPREPLY+=( $( compgen -W "${kubectl_out[*]}" -- "$cur" ) )
	fi
}

__work_list_files() {
	COMPREPLY+=( $( compgen -f -o plusdirs -X '!*.@(yaml|yml|json)' -- "$cur" ) )
}

__work_custom_func() {
	case ${last_command} in
		work_delete | work_get | work_resubmit)
			__work_get_workflow
			return
			;;
		work_suspend | work_terminate | work_wait | work_watch)
			__work_get_workflow "Running,Pending"
			return
			;;
		work_resume)
			__work_get_workflow "Running"
			return
			;;
		work_retry)
			__work_get_workflow "Failed"
			return
			;;
		work_logs)
			__work_get_logs
			return
			;;
		work_submit)
			__work_list_files
			return
			;;
		work_lint)
			__work_list_files
			return
			;;
		*)
			;;
	esac
}
	`
)

func NewCompletionCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "completion SHELL",
		Short: "output shell completion code for the specified shell (bash or zsh)",
		Long: `Write bash or zsh shell completion code to standard output.

For bash, ensure you have bash completions installed and enabled.
To access completions in your current shell, run
$ source <(work completion bash)
Alternatively, write it to a file and source in .bash_profile

For zsh, output to a file in a directory referenced by the $fpath shell
variable.
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			shell := args[0]
			rootCommand := NewCommand()
			rootCommand.BashCompletionFunction = bashCompletionFunc
			availableCompletions := map[string]func(io.Writer) error{
				"bash": rootCommand.GenBashCompletion,
				"zsh":  rootCommand.GenZshCompletion,
			}
			completion, ok := availableCompletions[shell]
			if !ok {
				fmt.Printf("Invalid shell '%s'. The supported shells are bash and zsh.\n", shell)
				os.Exit(1)
			}
			if err := completion(os.Stdout); err != nil {
				log.Fatal(err)
			}
		},
	}

	return command
}
