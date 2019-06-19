package commands

import (
	"github.com/spf13/cobra"

	"gitlab.com/teserakt/c2ae/internal/cli"
)

// Command defines a cli Command
type Command interface {
	CobraCmd() *cobra.Command
}

type rootCommandFlags struct {
	Endpoint string
	Cert     string
}

type rootCommand struct {
	cobraCmd *cobra.Command
	flags    rootCommandFlags
}

var _ Command = &rootCommand{}

// NewRootCommand creates and configure a new cli root command
func NewRootCommand(c2aeClientFactory cli.APIClientFactory, version string) Command {

	rootCmd := &rootCommand{}

	listCmd := NewListCommand(c2aeClientFactory)
	createCmd := NewCreateCommand(c2aeClientFactory)
	addTriggerCmd := NewAddTriggerCommand(c2aeClientFactory)
	addTargetCmd := NewAddTargetCommand(c2aeClientFactory)
	showCmd := NewShowCommand(c2aeClientFactory)
	deleteCmd := NewDeleteCommand(c2aeClientFactory)

	completionCmd := NewCompletionCommand(rootCmd)

	cobraCmd := &cobra.Command{
		Use:                    "c2ae-cli",
		BashCompletionFunction: completionCmd.GenerateCustomCompletionFuncs(),
		Version:                version,
		SilenceUsage:           true,
		SilenceErrors:          true,
	}

	cobraCmd.PersistentFlags().StringVarP(
		&rootCmd.flags.Endpoint,
		cli.EndpointFlag,
		"e",
		"127.0.0.1:5556", "url to the c2ae grpc api",
	)

	cobraCmd.PersistentFlags().StringVarP(
		&rootCmd.flags.Cert,
		cli.CertFlag,
		"c",
		"configs/c2ae-cert.pem", "path to the c2ae grpc api certificate",
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
