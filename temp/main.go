package main

import (
	"context"
	"github.com/idomath/toll-calculator/aggregator/client"
	"github.com/idomath/toll-calculator/types"
	"log"
	"time"
)

func main() {
	c, err := client.NewGrpcClient(":3001")
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuId: 1,
		Value: 234.12,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
}
