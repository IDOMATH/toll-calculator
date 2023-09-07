package main

import (
	"context"
	"github.com/idomath/toll-calculator/types"
)

type GrpcAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGrpcAggregatorServer(svc Aggregator) *GrpcAggregatorServer {
	return &GrpcAggregatorServer{
		svc: svc,
	}
}

func (s *GrpcAggregatorServer) Aggregate(ctx context.Context, request *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		ObuId: int(request.ObuId),
		Value: request.Value,
		Unix:  request.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}
