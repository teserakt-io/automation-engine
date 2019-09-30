package services

//go:generate mockgen -destination=c2client_mocks.go -package=services -self_package github.com/teserakt-io/automation-engine/internal/services github.com/teserakt-io/automation-engine/internal/services C2,C2EventStreamClient

import (
	"context"

	c2pb "github.com/teserakt-io/c2/pkg/pb"
	"go.opencensus.io/trace"

	"github.com/teserakt-io/automation-engine/internal/pb"
)

// C2EventStreamClient wrap around the c2pb.C2_SubscribeToEventStreamClient definition
type C2EventStreamClient interface {
	c2pb.C2_SubscribeToEventStreamClient
}

// C2 describes a C2 client service interface
type C2 interface {
	NewClientKey(ctx context.Context, clientName string) error
	NewTopicKey(ctx context.Context, topic string) error
	SubscribeToEventStream(ctx context.Context) (C2EventStreamClient, error)
}

type c2 struct {
	c2PbClientFactory pb.C2PbClientFactory
}

var _ C2 = &c2{}

// NewC2 creates a new C2 client service
func NewC2(c2pbClientFactory pb.C2PbClientFactory) C2 {
	return &c2{
		c2PbClientFactory: c2pbClientFactory,
	}
}

func (c *c2) NewClientKey(ctx context.Context, clientName string) error {
	ctx, span := trace.StartSpan(ctx, "C2Client.NewClientKey")
	defer span.End()

	client, err := c.c2PbClientFactory.Create()
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.NewClientKey(ctx, &c2pb.NewClientKeyRequest{Client: &c2pb.Client{Name: clientName}})

	return err
}

func (c *c2) NewTopicKey(ctx context.Context, topic string) error {
	ctx, span := trace.StartSpan(ctx, "C2Client.NewTopicKey")
	defer span.End()

	client, err := c.c2PbClientFactory.Create()
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.NewTopic(ctx, &c2pb.NewTopicRequest{Topic: topic})

	return err
}

func (c *c2) SubscribeToEventStream(ctx context.Context) (C2EventStreamClient, error) {
	ctx, span := trace.StartSpan(ctx, "C2Client.SubscribeToEventStream")
	defer span.End()

	client, err := c.c2PbClientFactory.Create()
	if err != nil {
		return nil, err
	}

	stream, err := client.SubscribeToEventStream(ctx, &c2pb.SubscribeToEventStreamRequest{})
	if err != nil {
		return nil, err
	}

	return stream, nil
}
