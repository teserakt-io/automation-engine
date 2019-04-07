package cli

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/teserakt/c2se/internal/pb"
)

type createCommand struct {
	cobraCmd   *cobra.Command
	c2seClient pb.C2ScriptEngineClient
}

var _ Command = &createCommand{}

// NewCreateCommand creates a new command to create a new rule
func NewCreateCommand(c2seClient pb.C2ScriptEngineClient) Command {

	createCmd := &createCommand{
		c2seClient: c2seClient,
	}

	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new rule",
		Run:   createCmd.run,
	}

	createCmd.cobraCmd = cobraCmd

	return createCmd
}

func (c *createCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *createCommand) Execute() error {
	return c.CobraCmd().Execute()
}

func (c *createCommand) run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.AddRuleRequest{
		Description: "Test command",
		Action:      pb.ActionType_KEY_ROTATION,
	}

	resp, err := c.c2seClient.AddRule(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rule #%d created!\n", resp.Rule.Id)
}
