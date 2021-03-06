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
	"regexp"
	"time"

	"github.com/spf13/cobra"

	"github.com/teserakt-io/automation-engine/internal/cli"
	"github.com/teserakt-io/automation-engine/internal/pb"
)

type addTargetCommand struct {
	cobraCmd          *cobra.Command
	c2aeClientFactory cli.APIClientFactory
	flags             addTargetCommandFlags
}

type addTargetCommandFlags struct {
	RuleID int32
	Type   string
	Expr   string
}

var _ Command = &addTargetCommand{}

// NewAddTargetCommand creates a new command to create a target on a rule
func NewAddTargetCommand(c2aeClientFactory cli.APIClientFactory) Command {
	addTargetCmd := &addTargetCommand{
		c2aeClientFactory: c2aeClientFactory,
	}

	cobraCmd := &cobra.Command{
		Use:   "add-target",
		Short: "Create a new target on a rule",
		RunE:  addTargetCmd.run,
	}

	cobraCmd.Flags().Int32Var(&addTargetCmd.flags.RuleID, "rule", 0, "The ruleID to add the target on")
	cobraCmd.Flags().StringVar(&addTargetCmd.flags.Type, "type", "", "The target type")
	cobraCmd.Flags().StringVar(
		&addTargetCmd.flags.Expr,
		"expr",
		"",
		"A regular expression used to match clients or topics",
	)

	cobraCmd.MarkFlagCustom("type", CompletionFuncNameTargetType)

	cobraCmd.MarkFlagRequired("rule")
	cobraCmd.MarkFlagRequired("type")
	cobraCmd.MarkFlagRequired("expr")

	addTargetCmd.cobraCmd = cobraCmd

	return addTargetCmd
}

func (c *addTargetCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *addTargetCommand) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	targetType, ok := pb.TargetType_value[c.flags.Type]
	if !ok {
		return fmt.Errorf("unknown trigger type %s", c.flags.Type)
	}

	_, err := regexp.Compile(c.flags.Expr)
	if err != nil {
		return fmt.Errorf("invalid expr: %s", err)
	}

	client, err := c.c2aeClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}
	defer client.Close()

	resp, err := client.GetRule(ctx, &pb.GetRuleRequest{RuleId: c.flags.RuleID})
	if err != nil {
		return fmt.Errorf("cannot retrieve rule #%d: %s", c.flags.RuleID, err)
	}

	target := &pb.Target{
		Type: pb.TargetType(targetType),
		Expr: c.flags.Expr,
	}

	updateReq := &pb.UpdateRuleRequest{
		RuleId:      c.flags.RuleID,
		Action:      resp.Rule.Action,
		Description: resp.Rule.Description,
		Targets:     append(resp.Rule.Targets, target),
		Triggers:    resp.Rule.Triggers,
	}

	_, err = client.UpdateRule(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("cannot update rule #%d: %s", c.flags.RuleID, err)
	}

	fmt.Printf("New target successfully added on rule #%d\n", c.flags.RuleID)

	return nil
}
