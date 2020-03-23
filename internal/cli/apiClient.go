// Copyright 2020 Teserakt AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/credentials"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/teserakt-io/automation-engine/internal/pb"
)

const (
	// EndpointFlag is the global flag name used to store the api endpoint url
	EndpointFlag = "endpoint"
	// CertFlag is the global flag name used to store the api certificate path
	CertFlag = "cert"
)

var (
	// ErrEndpointFlagUndefined is returned when the endpoint flag cannot be found on given command
	ErrEndpointFlagUndefined = errors.New("cannot retrieve endpoint flag on given cobra command")
	// ErrCertFlagUndefined is returned when the cert flag cannot be found on given command
	ErrCertFlagUndefined = errors.New("cannot retrieve cert flag on given cobra command")
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
	endpointFlag := cmd.Flag(EndpointFlag)
	if endpointFlag == nil || len(endpointFlag.Value.String()) == 0 {
		return nil, ErrEndpointFlagUndefined
	}

	certFlag := cmd.Flag(CertFlag)
	if certFlag == nil || len(certFlag.Value.String()) == 0 {
		return nil, ErrCertFlagUndefined
	}

	creds, err := credentials.NewClientTLSFromFile(certFlag.Value.String(), "")
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS credentials from certificate %v: %v", certFlag.Value.String(), err)
	}

	cnx, err := grpc.Dial(endpointFlag.Value.String(), grpc.WithTransportCredentials(creds))
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
