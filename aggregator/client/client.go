package client

import (
	"context"
	"github.com/idomath/toll-calculator/types"
)

type Client interface {
	Aggregate(ctx context.Context, request *types.AggregateRequest) error
}
