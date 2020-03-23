// Copyright 2020 Teserakt AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/teserakt-io/automation-engine/internal/pb"
)

const (
	// CompletionFuncNameAction holds the name of the bash function used to autocomplete action flag
	CompletionFuncNameAction = "__c2ae_autocomplete_actions"
	// CompletionFuncNameTriggerType holds the name of the bash function used to autocomplete trigger type flag
	CompletionFuncNameTriggerType = "__c2ae_autocomplete_trigger_types"
	// CompletionFuncNameTargetType holds the name of the bash function used to autocomplete target type flag
	CompletionFuncNameTargetType = "__c2ae_autocomplete_target_types"
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

. <(c2ae-cli completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(c2ae-cli completion)`,
		RunE: completionCmd.run,
	}

	cobraCmd.Flags().BoolVar(
		&completionCmd.flags.IsZsh,
		"zsh",
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

	var actionNames []string
	for _, name := range pb.ActionType_name {
		actionNames = append(actionNames, name)
	}

	var triggerTypes []string
	for _, t := range pb.TriggerType_name {
		triggerTypes = append(triggerTypes, t)
	}

	var targetTypes []string
	for _, t := range pb.TargetType_name {
		targetTypes = append(targetTypes, t)
	}

	out += c.generateCompletionFunc(CompletionFuncNameAction, actionNames)
	out += c.generateCompletionFunc(CompletionFuncNameTriggerType, triggerTypes)
	out += c.generateCompletionFunc(CompletionFuncNameTargetType, targetTypes)

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
