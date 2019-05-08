package api

import (
	context "context"
	"reflect"
	"testing"
	"time"

	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/pb"
	"gitlab.com/teserakt/c2se/internal/services"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
)

func assertRulesModified(t *testing.T, rulesModifiedChan <-chan bool, expectedModified bool) {
	rulesModified := <-rulesModifiedChan
	if rulesModified != expectedModified {
		t.Errorf("Expected rulesModified to be %t, got %t", expectedModified, rulesModified)
	}
}

func TestServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConverter := models.NewMockConverter(mockCtrl)
	mockRuleService := services.NewMockRuleService(mockCtrl)

	server := NewServer(":0", mockRuleService, mockConverter, log.NewNopLogger())

	rulesModifiedChan := make(chan bool)
	go func() {
		for {
			select {
			case modified := <-server.RulesModifiedChan():
				rulesModifiedChan <- modified
			case <-time.After(100 * time.Millisecond):
				rulesModifiedChan <- false
			}
		}
	}()

	t.Run("ListRules returns all the rules", func(t *testing.T) {
		rules := []models.Rule{
			models.Rule{ID: 1},
			models.Rule{ID: 2},
			models.Rule{ID: 3},
		}

		pbRules := []*pb.Rule{
			&pb.Rule{Id: 1},
			&pb.Rule{Id: 2},
			&pb.Rule{Id: 3},
		}

		mockRuleService.EXPECT().All().Times(1).Return(rules, nil)
		mockConverter.EXPECT().RulesToPb(rules).Times(1).Return(pbRules, nil)

		resp, err := server.ListRules(context.Background(), &pb.ListRulesRequest{})
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		assertRulesModified(t, rulesModifiedChan, false)

		if reflect.DeepEqual(resp.Rules, pbRules) == false {
			t.Errorf("Expected rules to be %#v, got %#v", pbRules, resp.Rules)
		}
	})

	t.Run("AddRule creates a new rule", func(t *testing.T) {

		pbTargets := []*pb.Target{
			&pb.Target{Id: 1},
			&pb.Target{Id: 2},
		}

		pbTriggers := []*pb.Trigger{
			&pb.Trigger{Id: 1},
			&pb.Trigger{Id: 2},
		}

		req := &pb.AddRuleRequest{
			Action:      pb.ActionType_KEY_ROTATION,
			Description: "description",
			Targets:     pbTargets,
			Triggers:    pbTriggers,
		}

		mockConverter.EXPECT().PbToTriggers(pbTriggers).Times(1)
		mockConverter.EXPECT().PbToTargets(pbTargets).Times(1)

		mockRuleService.EXPECT().Save(gomock.Any()).Times(1)

		pbRule := &pb.Rule{Id: 1}
		mockConverter.EXPECT().RuleToPb(gomock.Any()).Times(1).Return(pbRule, nil)

		resp, err := server.AddRule(context.Background(), req)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		assertRulesModified(t, rulesModifiedChan, true)

		if reflect.DeepEqual(resp.Rule, pbRule) == false {
			t.Errorf("Expected rule to be %#v, got %#v", pbRule, resp.Rule)
		}
	})

	t.Run("UpdateRule properly updates rule", func(t *testing.T) {
		targets := []models.Target{
			models.Target{ID: 1},
			models.Target{ID: 2},
		}

		pbTargets := []*pb.Target{
			&pb.Target{Id: 1},
			&pb.Target{Id: 2},
		}

		triggers := []models.Trigger{
			models.Trigger{ID: 1},
			models.Trigger{ID: 2},
		}

		pbTriggers := []*pb.Trigger{
			&pb.Trigger{Id: 1},
			&pb.Trigger{Id: 2},
		}

		req := &pb.UpdateRuleRequest{
			RuleId:      1,
			Action:      pb.ActionType_KEY_ROTATION,
			Description: "new description",
			Targets:     pbTargets,
			Triggers:    pbTriggers,
		}

		ruleBefore := models.Rule{
			ID:          1,
			Description: "before",
		}

		updatedRule := models.Rule{
			ID:          1,
			ActionType:  pb.ActionType_KEY_ROTATION,
			Description: "new description",
			Triggers:    triggers,
			Targets:     targets,
		}

		updatedPbRule := &pb.Rule{
			Id: 2,
		}

		mockRuleService.EXPECT().ByID(1).Times(1).Return(ruleBefore, nil)

		mockConverter.EXPECT().PbToTriggers(pbTriggers).Times(1).Return(triggers, nil)
		mockConverter.EXPECT().PbToTargets(pbTargets).Times(1).Return(targets, nil)

		mockRuleService.EXPECT().Save(gomock.Any()).Times(1)

		mockConverter.EXPECT().RuleToPb(updatedRule).Times(1).Return(updatedPbRule, nil)

		resp, err := server.UpdateRule(context.Background(), req)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		assertRulesModified(t, rulesModifiedChan, true)

		if reflect.DeepEqual(updatedPbRule, resp.Rule) == false {
			t.Errorf("Expected rule to be %#v, got %#v", updatedPbRule, resp.Rule)
		}
	})

	t.Run("DeleteRule deletes given rule", func(t *testing.T) {

		req := &pb.DeleteRuleRequest{
			RuleId: 1,
		}

		rule := models.Rule{ID: 1}

		mockRuleService.EXPECT().ByID(1).Times(1).Return(rule, nil)
		mockRuleService.EXPECT().Delete(rule).Times(1)

		resp, err := server.DeleteRule(context.Background(), req)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		assertRulesModified(t, rulesModifiedChan, true)

		if resp.RuleId != req.RuleId {
			t.Errorf("Expected ruleId to be %d, got %d", req.RuleId, resp.RuleId)
		}
	})

	t.Run("GetRule returns expected rule", func(t *testing.T) {

		req := &pb.GetRuleRequest{
			RuleId: 1,
		}

		rule := models.Rule{ID: 1}
		pbRule := &pb.Rule{Id: 1}

		mockRuleService.EXPECT().ByID(1).Times(1).Return(rule, nil)
		mockConverter.EXPECT().RuleToPb(rule).Times(1).Return(pbRule, nil)

		resp, err := server.GetRule(context.Background(), req)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		assertRulesModified(t, rulesModifiedChan, false)

		if reflect.DeepEqual(resp.Rule, pbRule) == false {
			t.Errorf("Expected rule to be %#v, got %#v", pbRule, resp.Rule)
		}
	})
}