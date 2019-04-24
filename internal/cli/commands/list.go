package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/rvflash/elapsed"
	"github.com/spf13/cobra"

	"gitlab.com/teserakt/c2se/internal/cli"
	"gitlab.com/teserakt/c2se/internal/pb"
)

type listCommand struct {
	cobraCmd          *cobra.Command
	c2seClientFactory cli.APIClientFactory
}

var _ Command = &listCommand{}

// NewListCommand creates a new command to list all the rules
func NewListCommand(c2seClientFactory cli.APIClientFactory) Command {

	listCmd := &listCommand{
		c2seClientFactory: c2seClientFactory,
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: "List all rules",
		RunE:  listCmd.run,
	}

	listCmd.cobraCmd = cobraCmd

	return listCmd
}

func (c *listCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *listCommand) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := c.c2seClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}

	resp, err := client.ListRules(ctx, &pb.ListRulesRequest{})
	if err != nil {
		return fmt.Errorf("api client error: %s", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	if len(resp.Rules) == 0 {
		fmt.Fprintln(w, "No rules are defined yet.")

		return nil
	}

	fmt.Fprintln(w, " #ID\t Description\t Triggers\t Targets\t Last executed")
	fmt.Fprintln(w, " ---\t -----------\t --------\t -------\t -------------")

	for _, rule := range resp.Rules {

		t, err := ptypes.Timestamp(rule.LastExectued)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(
			w,
			" %d\t %s\t %d\t %d\t %s\n",
			rule.Id,
			rule.Description,
			len(rule.Triggers),
			len(rule.Targets),
			elapsed.Time(t),
		)
	}

	return nil
}
