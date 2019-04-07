package api

import (
	context "context"
	"log"
	"net"

	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/pb"
	"gitlab.com/teserakt/c2se/internal/services"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Server interface
type Server interface {
	ListenAndServe() error
}

type apiServer struct {
	addr        string
	ruleService services.RuleService
}

var _ pb.C2ScriptEngineServer = &apiServer{}

// NewServer creates a new Server implementing the C2ScriptEngineServer interface
func NewServer(addr string, ruleService services.RuleService) Server {
	return &apiServer{
		addr:        addr,
		ruleService: ruleService,
	}
}

func (s *apiServer) ListenAndServe() error {

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterC2ScriptEngineServer(grpcServer, s)

	return grpcServer.Serve(lis)
}

func (s *apiServer) ListRules(ctx context.Context, req *pb.ListRulesRequest) (*pb.RulesResponse, error) {
	rules, err := s.ruleService.All()
	if err != nil {
		return nil, err
	}

	pbRules, err := models.RulesToPb(rules)
	if err != nil {
		return nil, err
	}

	response := &pb.RulesResponse{
		Rules: pbRules,
	}

	return response, err
}
func (s *apiServer) AddRule(ctx context.Context, req *pb.AddRuleRequest) (*pb.RuleResponse, error) {

	triggers, err := models.PbToTriggers(req.Triggers)
	if err != nil {
		return nil, err
	}

	targets, err := models.PbToTargets(req.Targets)
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

	pbRule, err := models.RuleToPb(*rule)
	if err != nil {
		return nil, err
	}

	response := &pb.RuleResponse{
		Rule: pbRule,
	}

	return response, err
}
func (*apiServer) UpdateRule(ctx context.Context, req *pb.UpdateRuleRequest) (*pb.RuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRule not implemented")
}
func (*apiServer) DeleteRule(ctx context.Context, req *pb.DeleteRuleRequest) (*pb.DeleteRuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRule not implemented")
}
