package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/teserakt/c2se/internal/cli/grpc"
	"gitlab.com/teserakt/c2se/internal/pb"
)

type createCommand struct {
	cobraCmd          *cobra.Command
	c2seClientFactory grpc.ClientFactory
	flags             createCommandFlags
}

type createCommandFlags struct {
	Description string
	Action      string
}

var _ Command = &createCommand{}

// NewCreateCommand creates a new command to create a new rule
func NewCreateCommand(c2seClientFactory grpc.ClientFactory) Command {
	createCmd := &createCommand{
		c2seClientFactory: c2seClientFactory,
	}

	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new rule",
		RunE:  createCmd.run,
	}

	cobraCmd.Flags().StringVarP(&createCmd.flags.Description, "description", "", "", "short description of the rule")
	cobraCmd.Flags().StringVarP(&createCmd.flags.Action, "action", "", "", "action to be performed when the rule will trigger")

	cobraCmd.MarkFlagCustom("action", CompletionFuncNameAction)

	cobraCmd.MarkFlagRequired("description")
	cobraCmd.MarkFlagRequired("action")

	createCmd.cobraCmd = cobraCmd

	return createCmd
}

func (c *createCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *createCommand) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	action, ok := pb.ActionType_value[c.flags.Action]
	if !ok {
		return fmt.Errorf("unknown action %s", c.flags.Action)
	}

	req := &pb.AddRuleRequest{
		Description: c.flags.Description,
		Action:      pb.ActionType(action),
	}

	client, err := c.c2seClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}

	resp, err := client.AddRule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add rule: %s", err)
	}

	fmt.Printf("Rule #%d created!\n", resp.Rule.Id)

	return nil
}
