package api

import (
	context "context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.opencensus.io/trace"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/teserakt-io/automation-engine/internal/config"
	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/pb"
	"github.com/teserakt-io/automation-engine/internal/services"
)

// Server interface
type Server interface {
	pb.C2AutomationEngineServer
	ListenAndServe(ctx context.Context) error
	RulesModifiedChan() <-chan bool
}

type apiServer struct {
	cfg         config.ServerCfg
	ruleService services.RuleService
	converter   models.Converter
	logger      log.Logger

	rulesModified chan bool
}

var _ pb.C2AutomationEngineServer = &apiServer{}

// NewServer creates a new Server implementing the C2AutomationEngineServer interface
func NewServer(
	cfg config.ServerCfg,
	ruleService services.RuleService,
	converter models.Converter,
	logger log.Logger,
) Server {
	return &apiServer{
		cfg:         cfg,
		ruleService: ruleService,
		converter:   converter,
		logger:      logger,

		rulesModified: make(chan bool),
	}
}

func (s *apiServer) RulesModifiedChan() <-chan bool {
	return s.rulesModified
}

func (s *apiServer) ListenAndServe(ctx context.Context) error {
	var lc net.ListenConfig
	grpcLis, err := lc.Listen(ctx, "tcp", s.cfg.GRPCAddr)
	if err != nil {
		s.logger.Log("msg", "failed to listen", "addr", s.cfg.GRPCAddr, "error", err)
		return err
	}
	defer grpcLis.Close()

	httpLis, err := lc.Listen(ctx, "tcp", s.cfg.HTTPAddr)
	if err != nil {
		s.logger.Log("msg", "failed to listen", "addr", s.cfg.HTTPAddr, "error", err)
		return err
	}
	defer httpLis.Close()

	errChan := make(chan error)
	go func() {
		errChan <- s.listenAndServeGRPC(ctx, grpcLis)
	}()
	go func() {
		errChan <- s.listenAndServeHTTP(ctx, httpLis)
	}()

	s.logger.Log("msg", "api server ready to accept connexions")

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *apiServer) listenAndServeGRPC(ctx context.Context, lis net.Listener) error {
	creds, err := credentials.NewServerTLSFromFile(s.cfg.GRPCCert, s.cfg.GRPCKey)
	if err != nil {
		s.logger.Log("msg", "failed to get credentials", "cert", s.cfg.GRPCCert, "key", s.cfg.GRPCKey, "error", err)
		return err
	}

	s.logger.Log("msg", "using TLS for gRPC", "cert", s.cfg.GRPCCert, "key", s.cfg.GRPCKey)

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterC2AutomationEngineServer(grpcServer, s)

	s.logger.Log("msg", "starting grpc listener", "addr", lis.Addr().String())
	return grpcServer.Serve(lis)
}

func (s *apiServer) listenAndServeHTTP(ctx context.Context, lis net.Listener) error {
	creds, err := credentials.NewClientTLSFromFile(s.cfg.GRPCCert, "")
	if err != nil {
		return fmt.Errorf("failed to create TLS credentials from %v: %v", s.cfg.GRPCCert, err)
	}

	httpMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	err = pb.RegisterC2AutomationEngineHandlerFromEndpoint(ctx, httpMux, s.cfg.HTTPGRPCAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register http listener : %v", err)
	}

	s.logger.Log("msg", "starting http listener", "addr", lis.Addr().String())
	return http.ServeTLS(lis, httpMux, s.cfg.HTTPCert, s.cfg.HTTPKey)
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

	// Force creation of new triggers
	for i := 0; i < len(req.Triggers); i++ {
		req.Triggers[i].Id = 0
	}

	triggers, err := s.converter.PbToTriggers(req.Triggers)
	if err != nil {
		return nil, err
	}

	// Force creation of new targets
	for i := 0; i < len(req.Targets); i++ {
		req.Targets[i].Id = 0
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

	deletedTriggers := models.FilterNonExistingTriggers(rule.Triggers, triggers)
	if len(deletedTriggers) > 0 {
		s.logger.Log("msg", "deleting removed triggers", "count", len(deletedTriggers))
		err := s.ruleService.DeleteTriggers(ctx, deletedTriggers...)
		if err != nil {
			return nil, err
		}
	}

	deletedTargets := models.FilterNonExistingTargets(rule.Targets, targets)
	if len(deletedTargets) > 0 {
		s.logger.Log("msg", "deleting removed targets", "count", len(deletedTargets))
		err := s.ruleService.DeleteTargets(ctx, deletedTargets...)
		if err != nil {
			return nil, err
		}
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
