package commands

import (
	"github.com/spf13/cobra"

	"gitlab.com/teserakt/c2se/internal/cli/grpc"
)

// Command defines a cli Command
type Command interface {
	CobraCmd() *cobra.Command
}

type rootCommandFlags struct {
	Endpoint string
}

type rootCommand struct {
	cobraCmd *cobra.Command
	flags    rootCommandFlags
}

var _ Command = &rootCommand{}

// NewRootCommand creates and configure a new cli root command
func NewRootCommand(c2seClientFactory grpc.ClientFactory, version string) Command {

	rootCmd := &rootCommand{}

	listCmd := NewListCommand(c2seClientFactory)
	createCmd := NewCreateCommand(c2seClientFactory)
	addTriggerCmd := NewAddTriggerCommand(c2seClientFactory)
	addTargetCmd := NewAddTargetCommand(c2seClientFactory)
	showCmd := NewShowCommand(c2seClientFactory)
	deleteCmd := NewDeleteCommand(c2seClientFactory)

	completionCmd := NewCompletionCommand(rootCmd)

	cobraCmd := &cobra.Command{
		Use:                    "c2se-cli",
		BashCompletionFunction: completionCmd.GenerateCustomCompletionFuncs(),
		Version:                version,
		SilenceUsage:           true,
		SilenceErrors:          true,
	}

	cobraCmd.PersistentFlags().StringVarP(
		&rootCmd.flags.Endpoint,
		grpc.EndpointFlag,
		"e",
		"127.0.0.1:5556", "url to the c2se grpc api",
	)

	cobraCmd.AddCommand(
		listCmd.CobraCmd(),
		createCmd.CobraCmd(),
		addTriggerCmd.CobraCmd(),
		addTargetCmd.CobraCmd(),
		showCmd.CobraCmd(),
		deleteCmd.CobraCmd(),

		// Autocompletion script generation command
		completionCmd.CobraCmd(),
	)

	cobraCmd.SetVersionTemplate(`{{printf "%s" .Version}}`)

	rootCmd.cobraCmd = cobraCmd

	return rootCmd
}

func (c *rootCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}
