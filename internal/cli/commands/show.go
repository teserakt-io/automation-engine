package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/teserakt-io/automation-engine/internal/cli"
	"github.com/teserakt-io/automation-engine/internal/pb"
)

type showCommand struct {
	cobraCmd          *cobra.Command
	c2aeClientFactory cli.APIClientFactory
	flags             showCommandFlags
}

type showCommandFlags struct {
	RuleID int32
}

var _ Command = &showCommand{}

// NewShowCommand creates a new command to show a given rule
func NewShowCommand(c2aeClientFactory cli.APIClientFactory) Command {
	showCmd := &showCommand{
		c2aeClientFactory: c2aeClientFactory,
	}

	cobraCmd := &cobra.Command{
		Use:   "show",
		Short: "Show a given rule",
		RunE:  showCmd.run,
	}

	cobraCmd.Flags().Int32Var(&showCmd.flags.RuleID, "rule", 0, "The ruleID to show")

	cobraCmd.MarkFlagRequired("rule")

	showCmd.cobraCmd = cobraCmd

	return showCmd
}

func (c *showCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *showCommand) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := c.c2aeClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}
	defer client.Close()

	resp, err := client.GetRule(ctx, &pb.GetRuleRequest{RuleId: c.flags.RuleID})
	if err != nil {
		return fmt.Errorf("cannot retrieve rule #%d: %s", c.flags.RuleID, err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(resp.Rule); err != nil {
		return fmt.Errorf("cannot json encode rule: %s", err)
	}

	return nil
}
