package cli

import (
	"errors"

	"github.com/spf13/cobra"
	"gitlab.com/teserakt/c2ae/internal/pb"
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

// C2AutomationEngineClient override the protobuf client definition to offer a Close method
// for the grpc connection
type C2AutomationEngineClient interface {
	pb.C2AutomationEngineClient
	Close() error
}

// APIClientFactory allows to create pb.C2AutomationEngineClient instances
type APIClientFactory interface {
	NewClient(cmd *cobra.Command) (C2AutomationEngineClient, error)
}

type apiClientFactory struct {
}

var _ APIClientFactory = &apiClientFactory{}

// NewAPIClientFactory creates a new C2AutomationEngineClient factory
func NewAPIClientFactory() APIClientFactory {
	return &apiClientFactory{}
}

type c2AutomationEngineClient struct {
	pb.C2AutomationEngineClient
	cnx *grpc.ClientConn
}

var _ C2AutomationEngineClient = &c2AutomationEngineClient{}
var _ pb.C2AutomationEngineClient = &c2AutomationEngineClient{}

// NewClient creates a new ob.C2AutomationEngineClient instance connecting to given api endpoint
func (c *apiClientFactory) NewClient(cmd *cobra.Command) (C2AutomationEngineClient, error) {
	flag := cmd.Flag(EndpointFlag)
	if flag == nil {
		return nil, ErrFlagUndefined
	}

	// TODO check https://godoc.org/google.golang.org/grpc#DialOption for available DialOptions
	cnx, err := grpc.Dial(flag.Value.String(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &c2AutomationEngineClient{
		C2AutomationEngineClient: pb.NewC2AutomationEngineClient(cnx),
		cnx:                      cnx,
	}, nil
}

func (c *c2AutomationEngineClient) Close() error {
	return c.cnx.Close()
}
