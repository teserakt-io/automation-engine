package api

import (
	context "context"
	"net"

	"github.com/go-kit/kit/log"
	"go.opencensus.io/trace"
	grpc "google.golang.org/grpc"

	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
	"gitlab.com/teserakt/c2ae/internal/services"
)

// Server interface
type Server interface {
	pb.C2AutomationEngineServer
	ListenAndServe(ctx context.Context, errorChan chan<- error)
	RulesModifiedChan() <-chan bool
}

type apiServer struct {
	addr        string
	ruleService services.RuleService
	converter   models.Converter
	logger      log.Logger

	rulesModified chan bool
}

var _ pb.C2AutomationEngineServer = &apiServer{}

// NewServer creates a new Server implementing the C2AutomationEngineServer interface
func NewServer(
	addr string,
	ruleService services.RuleService,
	converter models.Converter,
	logger log.Logger,
) Server {
	return &apiServer{
		addr:        addr,
		ruleService: ruleService,
		converter:   converter,
		logger:      logger,

		rulesModified: make(chan bool),
	}
}

func (s *apiServer) RulesModifiedChan() <-chan bool {
	return s.rulesModified
}

func (s *apiServer) ListenAndServe(ctx context.Context, errorChan chan<- error) {
	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", s.addr)

	if err != nil {
		s.logger.Log("msg", "failed to listen", "error", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterC2AutomationEngineServer(grpcServer, s)

	s.logger.Log("msg", "starting api grpc server", "addr", s.addr)
	errorChan <- grpcServer.Serve(lis)
}

func (s *apiServer) ListRules(ctx context.Context, req *pb.ListRulesRequest) (*pb.RulesResponse, error) {
	ctx, span := trace.StartSpan(ctx, "ListRules")
	defer span.End()

	rules, err := s.ruleService.All(ctx)
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
	ctx, span := trace.StartSpan(ctx, "GetRule")
	defer span.End()

	rule, err := s.ruleService.ByID(ctx, int(req.RuleId))
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
	ctx, span := trace.StartSpan(ctx, "AddRule")
	defer span.End()

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

	err = s.ruleService.Save(ctx, rule)
	if err != nil {
		return nil, err
	}

	pbRule, err := s.converter.RuleToPb(*rule)
	if err != nil {
		return nil, err
	}

	s.notifyRulesModified()

	return &pb.RuleResponse{
		Rule: pbRule,
	}, nil
}
func (s *apiServer) UpdateRule(ctx context.Context, req *pb.UpdateRuleRequest) (*pb.RuleResponse, error) {
	ctx, span := trace.StartSpan(ctx, "UpdateRule")
	defer span.End()

	rule, err := s.ruleService.ByID(ctx, int(req.RuleId))
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

	if err := s.ruleService.Save(ctx, &rule); err != nil {
		return nil, err
	}

	pbRule, err := s.converter.RuleToPb(rule)
	if err != nil {
		return nil, err
	}

	s.notifyRulesModified()

	return &pb.RuleResponse{
		Rule: pbRule,
	}, nil
}

func (s *apiServer) DeleteRule(ctx context.Context, req *pb.DeleteRuleRequest) (*pb.DeleteRuleResponse, error) {
	ctx, span := trace.StartSpan(ctx, "DeleteRule")
	defer span.End()

	rule, err := s.ruleService.ByID(ctx, int(req.RuleId))
	if err != nil {
		return nil, err
	}

	if err := s.ruleService.Delete(ctx, rule); err != nil {
		return nil, err
	}

	s.notifyRulesModified()

	return &pb.DeleteRuleResponse{RuleId: int32(rule.ID)}, nil
}

func (s *apiServer) notifyRulesModified() {
	select {
	case s.rulesModified <- true:
	default:
		s.logger.Log("msg", "skipped writting ruleModified event, channel is busy")
	}
}
