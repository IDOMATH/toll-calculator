package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/idomath/toll-calculator/types"
	"net/http"
)

type HttpClient struct {
	Endpoint string
}

func NewHttpClient(endpoint string) *HttpClient {
	return &HttpClient{
		Endpoint: endpoint,
	}
}

func (c *HttpClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with a non 200 status code: %d", resp.StatusCode)
	}
	return nil
}