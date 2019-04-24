package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/teserakt/c2se/internal/cli"
	"gitlab.com/teserakt/c2se/internal/pb"
)

type deleteCommand struct {
	cobraCmd          *cobra.Command
	c2seClientFactory cli.APIClientFactory
	flags             deleteCommandFlags
}

type deleteCommandFlags struct {
	RuleID int32
}

var _ Command = &deleteCommand{}

// NewDeleteCommand creates a new command to delete rules
func NewDeleteCommand(c2seClientFactory cli.APIClientFactory) Command {
	deleteCmd := &deleteCommand{
		c2seClientFactory: c2seClientFactory,
	}

	cobraCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a rule",
		RunE:  deleteCmd.run,
	}

	cobraCmd.Flags().Int32Var(&deleteCmd.flags.RuleID, "rule", 0, "The ruleID to show")

	cobraCmd.MarkFlagRequired("rule")

	deleteCmd.cobraCmd = cobraCmd

	return deleteCmd
}

func (c *deleteCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *deleteCommand) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := c.c2seClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}
	defer client.Close()

	req := &pb.DeleteRuleRequest{
		RuleId: c.flags.RuleID,
	}

	resp, err := client.DeleteRule(ctx, req)
	if err != nil {
		return fmt.Errorf("cannot delete rule #%d: %s", c.flags.RuleID, err)
	}

	fmt.Printf("Rule #%d deleted!\n", resp.RuleId)

	return nil
}
