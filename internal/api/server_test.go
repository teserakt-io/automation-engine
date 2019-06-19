package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"gitlab.com/teserakt/c2ae/internal/config"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
	"gitlab.com/teserakt/c2ae/internal/services"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func assertRulesModified(t *testing.T, rulesModifiedChan <-chan bool, expectedModified bool) {
	rulesModified := <-rulesModifiedChan
	if rulesModified != expectedModified {
		t.Errorf("Expected rulesModified to be %t, got %t", expectedModified, rulesModified)
	}
}

func getRootDir() string {
	_, filename, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(filename), "..", "..")
}

func TestServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockConverter := models.NewMockConverter(mockCtrl)
	mockRuleService := services.NewMockRuleService(mockCtrl)

	grpcLis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to obtain a free address: %v", err)
	}
	httpLis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to obtain a free address: %v", err)
	}

	certPath := filepath.Join(getRootDir(), "test/data/c2ae-cert.pem")
	keyPath := filepath.Join(getRootDir(), "test/data/c2ae-key.pem")

	serverCfg := config.ServerCfg{
		GRPCAddr: grpcLis.Addr().String(),
		GRPCCert: certPath,
		GRPCKey:  keyPath,
		HTTPAddr: httpLis.Addr().String(),
		HTTPCert: certPath,
		HTTPKey:  keyPath,
	}

	grpcLis.Close()
	httpLis.Close()

	server := NewServer(serverCfg, mockRuleService, mockConverter, log.NewNopLogger())

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

		mockRuleService.EXPECT().All(gomock.Any()).Times(1).Return(rules, nil)
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

		mockRuleService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1)

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

		mockRuleService.EXPECT().ByID(gomock.Any(), 1).Times(1).Return(ruleBefore, nil)

		mockConverter.EXPECT().PbToTriggers(pbTriggers).Times(1).Return(triggers, nil)
		mockConverter.EXPECT().PbToTargets(pbTargets).Times(1).Return(targets, nil)

		mockRuleService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1)

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

		mockRuleService.EXPECT().ByID(gomock.Any(), 1).Times(1).Return(rule, nil)
		mockRuleService.EXPECT().Delete(gomock.Any(), rule).Times(1)

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

		mockRuleService.EXPECT().ByID(gomock.Any(), 1).Times(1).Return(rule, nil)
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

	t.Run("ListenAndServe listen for grpc or http requests", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		errChan := make(chan error)

		go func() {
			errChan <- server.ListenAndServe(ctx)
		}()

		select {
		case err := <-errChan:
			t.Errorf("Expected no error, got %v", err)
		case <-time.After(1 * time.Second):
		}

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

		mockRuleService.EXPECT().All(gomock.Any()).AnyTimes().Return(rules, nil)
		mockConverter.EXPECT().RulesToPb(rules).AnyTimes().Return(pbRules, nil)

		// Test retrieve all rules with a GRPC client
		grpcClient := newGrpcClient(t, serverCfg.GRPCAddr, certPath)
		grpcResp, err := grpcClient.ListRules(context.Background(), &pb.ListRulesRequest{})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(grpcResp.Rules) != len(pbRules) {
			t.Errorf("Expected %d rules, got %d", len(pbRules), len(grpcResp.Rules))
		}

		for i, r := range grpcResp.Rules {
			if r.Id != pbRules[i].Id {
				t.Errorf("Expected ruleID to be %d, got %d", pbRules[i].Id, r.Id)
			}
		}

		// Test retrieve all rule with a HTTP request
		httpEndpoint := fmt.Sprintf("https://%s/rules", serverCfg.HTTPAddr)
		req, err := http.NewRequest("GET", httpEndpoint, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		httpResp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("HTTP request failed: %v", err)
		}
		d, err := ioutil.ReadAll(httpResp.Body)
		t.Log(string(d), err)
	})
}

func newGrpcClient(t *testing.T, addr string, certPath string) pb.C2AutomationEngineClient {
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		t.Fatalf("Failed to create TLS credentials: %v", err)
	}

	cnx, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	return pb.NewC2AutomationEngineClient(cnx)
}
