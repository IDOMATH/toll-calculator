package client

import (
	"context"
	"github.com/idomath/toll-calculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGrpcClient(endpoint string) (*GrpcClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	if err != nil {
		return nil, err
	}
	return &GrpcClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GrpcClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	return err
}
