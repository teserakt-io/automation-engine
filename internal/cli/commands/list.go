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
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/rvflash/elapsed"
	"github.com/spf13/cobra"

	"github.com/teserakt-io/automation-engine/internal/cli"
	"github.com/teserakt-io/automation-engine/internal/pb"
)

type listCommand struct {
	cobraCmd          *cobra.Command
	c2aeClientFactory cli.APIClientFactory
}

var _ Command = &listCommand{}

// NewListCommand creates a new command to list all the rules
func NewListCommand(c2aeClientFactory cli.APIClientFactory) Command {
	listCmd := &listCommand{
		c2aeClientFactory: c2aeClientFactory,
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

	client, err := c.c2aeClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}
	defer client.Close()

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
		t, err := ptypes.Timestamp(rule.LastExecuted)
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
