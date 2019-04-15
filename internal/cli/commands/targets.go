package commands

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/teserakt/c2se/internal/cli/grpc"
	"gitlab.com/teserakt/c2se/internal/pb"
)

type addTargetCommand struct {
	cobraCmd          *cobra.Command
	c2seClientFactory grpc.ClientFactory
	flags             addTargetCommandFlags
}

type addTargetCommandFlags struct {
	RuleID int32
	Type   string
	Expr   string
}

var _ Command = &addTargetCommand{}

// NewAddTargetCommand creates a new command to create a target on a rule
func NewAddTargetCommand(c2seClientFactory grpc.ClientFactory) Command {
	addTargetCmd := &addTargetCommand{
		c2seClientFactory: c2seClientFactory,
	}

	cobraCmd := &cobra.Command{
		Use:   "add-target",
		Short: "Create a new target on a rule",
		RunE:  addTargetCmd.run,
	}

	cobraCmd.Flags().Int32VarP(&addTargetCmd.flags.RuleID, "rule", "", 0, "The ruleID to add the target on")
	cobraCmd.Flags().StringVarP(&addTargetCmd.flags.Type, "type", "", "", "The target type")
	cobraCmd.Flags().StringVarP(
		&addTargetCmd.flags.Expr,
		"expr",
		"",
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
		return fmt.Errorf("Invalid expr: %s", err)
	}

	client, err := c.c2seClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}

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

	resp, err = client.UpdateRule(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("cannot update rule #%d: %s", c.flags.RuleID, err)
	}

	fmt.Printf("New target successfully added on rule #%d\n", c.flags.RuleID)

	return nil
}
