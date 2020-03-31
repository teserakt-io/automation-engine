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
	"time"

	"github.com/spf13/cobra"

	"github.com/teserakt-io/automation-engine/internal/cli"
	"github.com/teserakt-io/automation-engine/internal/pb"
)

type createCommand struct {
	cobraCmd          *cobra.Command
	c2aeClientFactory cli.APIClientFactory
	flags             createCommandFlags
}

type createCommandFlags struct {
	Description string
	Action      string
}

var _ Command = &createCommand{}

// NewCreateCommand creates a new command to create a new rule
func NewCreateCommand(c2aeClientFactory cli.APIClientFactory) Command {
	createCmd := &createCommand{
		c2aeClientFactory: c2aeClientFactory,
	}

	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new rule",
		RunE:  createCmd.run,
	}

	cobraCmd.Flags().StringVar(&createCmd.flags.Description, "description", "", "short description of the rule")
	cobraCmd.Flags().StringVar(&createCmd.flags.Action, "action", "", "action to be performed when the rule will trigger")

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

	client, err := c.c2aeClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}
	defer client.Close()

	resp, err := client.AddRule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add rule: %s", err)
	}

	fmt.Printf("Rule #%d created!\n", resp.Rule.Id)

	return nil
}
