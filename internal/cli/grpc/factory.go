package grpc

import (
	"errors"

	"github.com/spf13/cobra"
	"gitlab.com/teserakt/c2se/internal/pb"
	"google.golang.org/grpc"
)

const (
	// EndpointFlag is the global flag name used to store the api endpoint url
	EndpointFlag = "endpoint"
)

var (
	// ErrFlagUndefined is returned when the endpoint flag cannot be found on given command
	ErrFlagUndefined = errors.New("cannot retrieve endpoint flag on given cobra command")
)

// ClientFactory allows to create pb.C2ScriptEngineClient instances
type ClientFactory interface {
	NewClient(cmd *cobra.Command) (pb.C2ScriptEngineClient, error)
}

type clientFactory struct {
}

var _ ClientFactory = &clientFactory{}

// NewClientFactory creates a new C2ScriptEngineClient factory
func NewClientFactory() ClientFactory {
	return &clientFactory{}
}

// NewClient creates a new ob.C2ScriptEngineClient instance connecting to given api endpoint
func (c *clientFactory) NewClient(cmd *cobra.Command) (pb.C2ScriptEngineClient, error) {
	flag := cmd.Flag(EndpointFlag)
	if flag == nil {
		return nil, ErrFlagUndefined
	}

	cnx, err := grpc.Dial(flag.Value.String(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewC2ScriptEngineClient(cnx), nil
}
