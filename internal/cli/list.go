package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/teserakt/c2se/internal/pb"
)

type listCommand struct {
	cobraCmd   *cobra.Command
	c2seClient pb.C2ScriptEngineClient
}

var _ Command = &listCommand{}

// NewListCommand creates a new command to list all the rules
func NewListCommand(c2seClient pb.C2ScriptEngineClient) Command {

	listCmd := &listCommand{
		c2seClient: c2seClient,
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: "List all rules",
		Run:   listCmd.run,
	}

	listCmd.cobraCmd = cobraCmd

	return listCmd
}

func (c *listCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *listCommand) Execute() error {
	return c.CobraCmd().Execute()
}

func (c *listCommand) run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.c2seClient.ListRules(ctx, &pb.ListRulesRequest{})
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	if len(resp.Rules) == 0 {
		fmt.Fprintln(w, "No rules are defined yet.")

		return
	}

	fmt.Fprintln(w, " #ID\t Description\t Triggers\t Targets\t Last executed")

	for _, rule := range resp.Rules {
		fmt.Fprintf(
			w,
			" %d\t %s\t %d\t %d\t %s\n",
			rule.Id,
			rule.Description,
			len(rule.Triggers),
			len(rule.Targets),
			rule.LastExectued,
		)
	}

}
