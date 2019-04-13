package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"

	"gitlab.com/teserakt/c2se/internal/cli/grpc"
	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/pb"
)

type addTriggerCommand struct {
	cobraCmd          *cobra.Command
	c2seClientFactory grpc.ClientFactory
	flags             addTriggerCommandFlags
}

type addTriggerCommandFlags struct {
	RuleID   int32
	Type     string
	Settings map[string]string
}

var _ Command = &addTriggerCommand{}

// NewAddTriggerCommand creates a new command to create a trigger on a rule
func NewAddTriggerCommand(c2seClientFactory grpc.ClientFactory) Command {
	addTriggerCmd := &addTriggerCommand{
		c2seClientFactory: c2seClientFactory,
	}

	cobraCmd := &cobra.Command{
		Use:   "add-trigger",
		Short: "Create a new trigger on a rule",
		RunE:  addTriggerCmd.run,
	}

	cobraCmd.Flags().Int32VarP(&addTriggerCmd.flags.RuleID, "rule", "r", 0, "The ruleID to add the trigger on")
	cobraCmd.Flags().StringVarP(&addTriggerCmd.flags.Type, "type", "t", "", "The trigger type")
	cobraCmd.Flags().StringToStringVarP(
		&addTriggerCmd.flags.Settings,
		"setting",
		"s",
		nil,
		"Used to set trigger settings",
	)

	cobraCmd.MarkFlagCustom("type", CompletionFuncNameTriggerType)

	cobraCmd.MarkFlagRequired("rule")
	cobraCmd.MarkFlagRequired("type")

	addTriggerCmd.cobraCmd = cobraCmd

	return addTriggerCmd
}

func (c *addTriggerCommand) CobraCmd() *cobra.Command {
	return c.cobraCmd
}

func (c *addTriggerCommand) run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	triggerType, ok := pb.TriggerType_value[c.flags.Type]
	if !ok {
		return fmt.Errorf("unknown trigger type %s", c.flags.Type)
	}

	client, err := c.c2seClientFactory.NewClient(cmd)
	if err != nil {
		return fmt.Errorf("cannot create api client: %s", err)
	}

	resp, err := client.GetRule(ctx, &pb.GetRuleRequest{RuleId: c.flags.RuleID})
	if err != nil {
		return fmt.Errorf("cannot retrieve rule #%d: %s", c.flags.RuleID, err)
	}

	triggerSettings, err := mapToTriggerSettings(c.flags.Settings, pb.TriggerType(triggerType))
	if err != nil {
		return err
	}

	if err := triggerSettings.Validate(); err != nil {
		return fmt.Errorf("trigger settings validation error: %s", err)
	}

	encodedSettings, err := triggerSettings.Encode()
	if err != nil {
		return err
	}

	newTrigger := &pb.Trigger{
		Type:     pb.TriggerType(triggerType),
		Settings: encodedSettings,
	}

	updateReq := &pb.UpdateRuleRequest{
		RuleId:      c.flags.RuleID,
		Action:      resp.Rule.Action,
		Description: resp.Rule.Description,
		Targets:     resp.Rule.Targets,
		Triggers:    append(resp.Rule.Triggers, newTrigger),
	}

	resp, err = client.UpdateRule(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("cannot update rule #%d: %s", c.flags.RuleID, err)
	}

	fmt.Printf("New trigger successfully added on rule #%d\n", c.flags.RuleID)

	return nil
}

func mapToTriggerSettings(userSettings map[string]string, triggerType pb.TriggerType) (models.TriggerSettings, error) {
	var decoderConfig *mapstructure.DecoderConfig

	switch triggerType {
	case pb.TriggerType_TIME_INTERVAL:
		decoderConfig = &mapstructure.DecoderConfig{
			Result: &models.TriggerSettingsTimeInterval{},
		}
	case pb.TriggerType_CLIENT_SUBSCRIBED:
		decoderConfig = &mapstructure.DecoderConfig{
			Result: &models.TriggerSettingsEvent{
				EventType: models.EventTypeClientSubscribed,
			},
		}
	case pb.TriggerType_CLIENT_UNSUBSCRIBED:
		decoderConfig = &mapstructure.DecoderConfig{
			Result: &models.TriggerSettingsEvent{
				EventType: models.EventTypeClientUnsubscribed,
			},
		}
	}

	decoderConfig.WeaklyTypedInput = true
	decoderConfig.Metadata = &mapstructure.Metadata{}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(userSettings); err != nil {
		return nil, err
	}

	for _, unused := range decoderConfig.Metadata.Unused {
		fmt.Printf("WARN: setting %s is provided, but was ignored.\n", unused)
	}

	return decoderConfig.Result.(models.TriggerSettings), nil
}
