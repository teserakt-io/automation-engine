package cli

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

// C2ScriptEngineClient override the protobuf client definition to offer a Close method
// for the grpc connection
type C2ScriptEngineClient interface {
	pb.C2ScriptEngineClient
	Close() error
}

// APIClientFactory allows to create pb.C2ScriptEngineClient instances
type APIClientFactory interface {
	NewClient(cmd *cobra.Command) (C2ScriptEngineClient, error)
}

type apiClientFactory struct {
}

var _ APIClientFactory = &apiClientFactory{}

// NewAPIClientFactory creates a new C2ScriptEngineClient factory
func NewAPIClientFactory() APIClientFactory {
	return &apiClientFactory{}
}

type c2scriptEngineClient struct {
	pb.C2ScriptEngineClient
	cnx *grpc.ClientConn
}

var _ C2ScriptEngineClient = &c2scriptEngineClient{}
var _ pb.C2ScriptEngineClient = &c2scriptEngineClient{}

// NewClient creates a new ob.C2ScriptEngineClient instance connecting to given api endpoint
func (c *apiClientFactory) NewClient(cmd *cobra.Command) (C2ScriptEngineClient, error) {
	flag := cmd.Flag(EndpointFlag)
	if flag == nil {
		return nil, ErrFlagUndefined
	}

	// TODO check https://godoc.org/google.golang.org/grpc#DialOption for available DialOptions
	cnx, err := grpc.Dial(flag.Value.String(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &c2scriptEngineClient{
		C2ScriptEngineClient: pb.NewC2ScriptEngineClient(cnx),
		cnx:                  cnx,
	}, nil
}

func (c *c2scriptEngineClient) Close() error {
	return c.cnx.Close()
}
