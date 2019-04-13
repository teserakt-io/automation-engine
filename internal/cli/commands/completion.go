package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/teserakt/c2se/internal/pb"
)

const (
	// CompletionFuncNameAction holds the name of the bash function used to autocomplete action flag
	CompletionFuncNameAction = "__c2se_autocomplete_actions"
	// CompletionFuncNameTriggerType holds the name of the bash function used to autocomplet trigger type flag
	CompletionFuncNameTriggerType = "__c2se_autocomplete_trigger_types"
)

// CompletionCommand defines a custom Command to deal with auto completion
type CompletionCommand struct {
	cobraCmd *cobra.Command
	rootCmd  Command
	flags    completionCommandFlags
}

type completionCommandFlags struct {
	IsZsh bool
}

var _ Command = &CompletionCommand{}

// NewCompletionCommand returns the cobra command used to generate the autocompletion
func NewCompletionCommand(rootCommand Command) *CompletionCommand {

	completionCmd := &CompletionCommand{
		rootCmd: rootCommand,
	}

	cobraCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

. <(c2se-cli completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(c2se-cli completion)`,
		RunE: completionCmd.run,
	}

	cobraCmd.Flags().BoolVarP(
		&completionCmd.flags.IsZsh,
		"zsh",
		"",
		false,
		"Generate zsh completion script (default: bash)",
	)

	completionCmd.cobraCmd = cobraCmd

	return completionCmd
}

// CobraCmd returns the cobra command
func (c *CompletionCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *CompletionCommand) run(cmd *cobra.Command, args []string) error {
	if c.flags.IsZsh {
		c.rootCmd.CobraCmd().GenZshCompletion(os.Stdout)

		return nil
	}

	c.rootCmd.CobraCmd().GenBashCompletion(os.Stdout)

	return nil
}

// GenerateCustomCompletionFuncs returns the bash script snippets to use for custom autocompletion
func (c *CompletionCommand) GenerateCustomCompletionFuncs() string {
	var out string

	// Autocomplete for pb.ActionType
	var actionNames []string
	for _, name := range pb.ActionType_name {
		actionNames = append(actionNames, name)
	}

	var triggerTypes []string
	for _, t := range pb.TriggerType_name {
		triggerTypes = append(triggerTypes, t)
	}

	out += c.generateCompletionFunc(CompletionFuncNameAction, actionNames)
	out += c.generateCompletionFunc(CompletionFuncNameTriggerType, triggerTypes)

	return out
}

func (c *CompletionCommand) generateCompletionFunc(funcName string, suggestions []string) string {
	return fmt.Sprintf(`
	%s()
	{
		COMPREPLY=( $(compgen -W "%s" -- "$cur") )
	}
	`,
		funcName,
		strings.Join(suggestions, " "),
	)
}
