package pb

//go:generate mockgen -copyright_file ../../doc/COPYRIGHT_TEMPLATE.txt -destination=c2client_mocks.go -package=pb -self_package github.com/teserakt-io/automation-engine/internal/pb github.com/teserakt-io/automation-engine/internal/pb C2PbClient,C2PbClientFactory

import (
	c2pb "github.com/teserakt-io/c2/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// C2PbClient wrap e4.C2Client in a local interface and provide a way to close its connection
type C2PbClient interface {
	c2pb.C2Client
	Close() error
}

type c2PbClient struct {
	c2pb.C2Client
	cnx *grpc.ClientConn
}

var _ C2PbClient = &c2PbClient{}

func (c *c2PbClient) Close() error {
	return c.cnx.Close()
}

// C2PbClientFactory defines a factory creating C2PbClient
type C2PbClientFactory interface {
	Create() (C2PbClient, error)
}

type c2PbClientFactory struct {
	endpoint string
	creds    credentials.TransportCredentials
}

var _ C2PbClientFactory = &c2PbClientFactory{}

// NewC2PbClientFactory creates a new factory for C2 protobuf client
func NewC2PbClientFactory(endpoint string, certPath string) (C2PbClientFactory, error) {
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return nil, err
	}

	return &c2PbClientFactory{
		endpoint: endpoint,
		creds:    creds,
	}, nil
}

func (f *c2PbClientFactory) Create() (C2PbClient, error) {
	cnx, err := grpc.Dial(f.endpoint, grpc.WithTransportCredentials(f.creds))
	if err != nil {
		return nil, err
	}

	return &c2PbClient{
		C2Client: c2pb.NewC2Client(cnx),
		cnx:      cnx,
	}, nil
}
