package api

import (
	context "context"
	"reflect"
	"testing"

	"gitlab.com/teserakt/c2se/internal/models"

	"github.com/golang/mock/gomock"
	"gitlab.com/teserakt/c2se/internal/mocks"
	"gitlab.com/teserakt/c2se/internal/pb"
)

func TestServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConverter := mocks.NewMockConverter(mockCtrl)
	mockRuleService := mocks.NewMockRuleService(mockCtrl)

	server := NewServer(":0", mockRuleService, mockConverter)

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

		if reflect.DeepEqual(resp.Rule, pbRule) == false {
			t.Errorf("Expected rule to be %#v, got %#v", pbRule, resp.Rule)
		}
	})
}
