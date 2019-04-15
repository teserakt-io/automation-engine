package api

import (
	context "context"
	"log"
	"net"

	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/pb"
	"gitlab.com/teserakt/c2se/internal/services"

	grpc "google.golang.org/grpc"
)

// Server interface
type Server interface {
	pb.C2ScriptEngineServer
	ListenAndServe() error
}

type apiServer struct {
	addr        string
	ruleService services.RuleService
	converter   models.Converter
}

var _ pb.C2ScriptEngineServer = &apiServer{}

// NewServer creates a new Server implementing the C2ScriptEngineServer interface
func NewServer(addr string, ruleService services.RuleService, converter models.Converter) Server {
	return &apiServer{
		addr:        addr,
		ruleService: ruleService,
		converter:   converter,
	}
}

func (s *apiServer) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterC2ScriptEngineServer(grpcServer, s)

	log.Printf("Starting api grpc server listening on %s\n", s.addr)
	return grpcServer.Serve(lis)
}

func (s *apiServer) ListRules(ctx context.Context, req *pb.ListRulesRequest) (*pb.RulesResponse, error) {
	rules, err := s.ruleService.All()
	if err != nil {
		return nil, err
	}

	pbRules, err := s.converter.RulesToPb(rules)
	if err != nil {
		return nil, err
	}

	return &pb.RulesResponse{
		Rules: pbRules,
	}, nil
}

func (s *apiServer) GetRule(ctx context.Context, req *pb.GetRuleRequest) (*pb.RuleResponse, error) {

	rule, err := s.ruleService.ByID(int(req.RuleId))
	if err != nil {
		return nil, err
	}

	pbRule, err := s.converter.RuleToPb(rule)
	if err != nil {
		return nil, err
	}

	return &pb.RuleResponse{
		Rule: pbRule,
	}, nil
}

func (s *apiServer) AddRule(ctx context.Context, req *pb.AddRuleRequest) (*pb.RuleResponse, error) {
	triggers, err := s.converter.PbToTriggers(req.Triggers)
	if err != nil {
		return nil, err
	}

	targets, err := s.converter.PbToTargets(req.Targets)
	if err != nil {
		return nil, err
	}

	rule := &models.Rule{
		Description: req.Description,
		ActionType:  req.Action,
		Triggers:    triggers,
		Targets:     targets,
	}

	err = s.ruleService.Save(rule)
	if err != nil {
		return nil, err
	}

	pbRule, err := s.converter.RuleToPb(*rule)
	if err != nil {
		return nil, err
	}

	return &pb.RuleResponse{
		Rule: pbRule,
	}, nil
}
func (s *apiServer) UpdateRule(ctx context.Context, req *pb.UpdateRuleRequest) (*pb.RuleResponse, error) {

	rule, err := s.ruleService.ByID(int(req.RuleId))
	if err != nil {
		return nil, err
	}

	triggers, err := s.converter.PbToTriggers(req.Triggers)
	if err != nil {
		return nil, err
	}

	targets, err := s.converter.PbToTargets(req.Targets)
	if err != nil {
		return nil, err
	}

	rule.Description = req.Description
	rule.ActionType = req.Action
	rule.Triggers = triggers
	rule.Targets = targets

	if err := s.ruleService.Save(&rule); err != nil {
		return nil, err
	}

	pbRule, err := s.converter.RuleToPb(rule)
	if err != nil {
		return nil, err
	}

	return &pb.RuleResponse{
		Rule: pbRule,
	}, nil
}

func (s *apiServer) DeleteRule(ctx context.Context, req *pb.DeleteRuleRequest) (*pb.DeleteRuleResponse, error) {
	rule, err := s.ruleService.ByID(int(req.RuleId))
	if err != nil {
		return nil, err
	}

	if err := s.ruleService.Delete(rule); err != nil {
		return nil, err
	}

	return &pb.DeleteRuleResponse{RuleId: int32(rule.ID)}, nil
}
