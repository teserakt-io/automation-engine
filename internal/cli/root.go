package cli

import (
	"github.com/spf13/cobra"

	"gitlab.com/teserakt/c2se/internal/pb"
)

// Command defines a cli Command
type Command interface {
	Execute() error
	CobraCmd() *cobra.Command
}

type rootCommand struct {
	cobraCmd *cobra.Command
}

var _ Command = &rootCommand{}

// NewRootCommand creates and configure a new cli root command
func NewRootCommand(c2seClient pb.C2ScriptEngineClient) Command {

	cobraCmd := &cobra.Command{
		Use: "C2 script-engine cli",
	}

	listCmd := NewListCommand(c2seClient)
	createCmd := NewCreateCommand(c2seClient)

	cobraCmd.AddCommand(
		listCmd.CobraCmd(),
		//	showCmd,
		createCmd.CobraCmd(),
	//	addTriggerCmd,
	//	addTargetCmd,
	)

	return &rootCommand{
		cobraCmd: cobraCmd,
	}
}

func (c *rootCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *rootCommand) Execute() error {
	return c.CobraCmd().Execute()
}
