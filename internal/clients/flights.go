package clients

import (
	"context"
	"time"

	model "github.com/iramkrish/search-bff/internal/modal"
)

type FlightClient struct{}

func NewFlightClient() *FlightClient {
	return &FlightClient{}
}

func (c *FlightClient) Search(ctx context.Context, q string) ([]model.Flight, error) {
	select {
	case <-time.After(150 * time.Millisecond):
		return []model.Flight{
			{ID: "f1", Number: "AI-202"},
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
