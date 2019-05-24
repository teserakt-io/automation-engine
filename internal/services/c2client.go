package services

//go:generate mockgen -destination=c2client_mocks.go -package=services -self_package gitlab.com/teserakt/c2ae/internal/services gitlab.com/teserakt/c2ae/internal/services C2,C2Requester

import (
	"context"
	"errors"
	"time"

	"gitlab.com/teserakt/c2ae/internal/pb"
	e4 "gitlab.com/teserakt/e4common"
)

// C2 describes a C2 client service interface
type C2 interface {
	NewClientKey(clientID []byte) error
	NewTopicKey(topic string) error
}

// C2Requester defines a type able to make request to C2 backend
type C2Requester interface {
	C2Request(e4.C2Request) (e4.C2Response, error)
}

type c2Requester struct {
	c2PbClientFactory pb.C2PbClientFactory
}

var _ C2Requester = &c2Requester{}

// NewC2Requester creates a new C2Requester
func NewC2Requester(c2PbClientFactory pb.C2PbClientFactory) C2Requester {
	return &c2Requester{
		c2PbClientFactory: c2PbClientFactory,
	}
}

func (r *c2Requester) C2Request(in e4.C2Request) (e4.C2Response, error) {
	client, err := r.c2PbClientFactory.Create()
	if err != nil {
		return e4.C2Response{}, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.C2Command(ctx, &in)
	if err != nil {
		return e4.C2Response{}, err
	}

	if resp.Success == false {
		return e4.C2Response{}, errors.New(resp.Err)
	}

	return *resp, nil
}

type c2 struct {
	c2Requester C2Requester
}

var _ C2 = &c2{}

// NewC2 creates a new C2 client service
func NewC2(c2Requester C2Requester) C2 {
	return &c2{
		c2Requester: c2Requester,
	}
}

func (c *c2) NewClientKey(clientID []byte) error {
	request := e4.C2Request{
		Command: e4.C2Request_NEW_CLIENT_KEY,
		Id:      clientID,
	}

	_, err := c.c2Requester.C2Request(request)

	return err
}

func (c *c2) NewTopicKey(topic string) error {
	request := e4.C2Request{
		Command: e4.C2Request_NEW_TOPIC,
		Topic:   topic,
	}

	_, err := c.c2Requester.C2Request(request)

	return err
}
